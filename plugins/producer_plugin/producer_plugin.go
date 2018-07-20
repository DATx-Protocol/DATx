package producerplugin

import (
	"DATx/chainlib/application"
	"DATx/utils/common"
	"DATx/utils/crypto"
	"DATx/utils/helper"
	"crypto/ecdsa"
	"datx_chain/chainlib/controller"
	"datx_chain/chainlib/types"

	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"gopkg.in/yaml.v2"
)

//enum start block result
const (
	succeeded uint8 = iota
	failed
	exhausted
)

//enum pending block status
const (
	producing uint8 = iota
	speculating
)

//ProducerConfig struct
type ProducerConfig struct {
	//account name
	Producers    []string `yaml:"Producers"`
	PrivateKey   string   `yaml:"PrivateKey"`
	maxTrxTimeMs int64    `yaml:"MaxTrxTimeMs"`
}

//ProducerPlugin struct
type ProducerPlugin struct {
	producers []string //account name

	privateKey *ecdsa.PrivateKey

	timer *time.Timer

	config ProducerConfig

	persistentTransaction map[common.Hash]uint64

	chain *controller.Controller

	pendingBlockMode   uint8
	producerWaterMarks map[string]uint32 //the pairs of accountname/blocknum
	productionEnabled  bool
}

//NewProducerPlugin new
func NewProducerPlugin() *ProducerPlugin {
	return &ProducerPlugin{
		productionEnabled:     false,
		persistentTransaction: make(map[common.Hash]uint64, 0),
	}
}

//Init interface relization
func (pp *ProducerPlugin) Init() error {
	err, data := helper.GetFileHelper("producer_config.yaml", application.App().GetConfigFolder())
	if err != nil {
		log.Printf("producer_plugin init with producer config error={%v}", err)
	}

	var config ProducerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Printf("producer_plugin init unmarshal config  error={%v}", err)
		return err
	}

	log.Printf("producer init config=%v", config)
	pp.producers = config.Producers
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		log.Printf("convert private key error,%v", err.Error())
		return err
	}
	pp.privateKey = privateKey
	pp.config = config
	return nil
}

//Open interface relization
func (pp *ProducerPlugin) Open() error {
	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
		return err
	}
	chainplugin := plugin.(*chainplugin.ChainPlugin)
	pp.chain = chainplugin.Chain()

	//listen sync block
	pp.chain.SyncBlockChan = make(chan *types.Block, 10)
	go func() {
		for {
			select {
			case v := <-pp.chain.SyncBlockChan:
				log.Printf("on incoming block: %v", v)
				pp.OnIncomingBlock(v)
			default:
			}
		}
	}()

	// pp.scheduleProductionLoop()
	//start timer

	timeout := time.Duration(types.BlockIntervalMs) * time.Millisecond
	pp.timer = time.NewTimer(timeout)
	go func() {
		for {
			select {
			case <-pp.timer.C:
				pp.scheduleProductionLoop()
			}
		}
	}()

	//init

	return nil
}

//Close interface relization
func (pp *ProducerPlugin) Close() {
	pp.timer.Stop()
}

func (pp *ProducerPlugin) resetTimer(d time.Duration) {
	if !pp.timer.Stop() {
		select {
		case <-pp.timer.C:
		default:
		}
	}

	pp.timer.Reset(d)
}

func (pp *ProducerPlugin) scheduleProductionLoop() {
	result := pp.startBlock()

	if result == failed {
		// log.Print("Failed to start a pending block, will try again later")
		delaytime := types.BlockIntervalMs / 10
		pp.resetTimer(time.Duration(delaytime) * time.Millisecond)
	} else if pp.pendingBlockMode == producing {
		log.Print("successed\n")
		pp.MaybeProduceBlock()
	} else if pp.pendingBlockMode == speculating && len(pp.producers) > 0 {
		timeout := time.Duration(types.BlockIntervalMs) * time.Millisecond

		pp.resetTimer(timeout)
	}
}

//start block return result
func (pp *ProducerPlugin) startBlock() uint8 {
	headstate := pp.chain.HeadBlockState()

	now := time.Now()
	headtime := pp.chain.HeadBlockTime()
	var base time.Time
	if now.Unix() > headtime.Unix() {
		base = now
	} else {
		base = headtime
	}

	mintimenextblock := types.BlockIntervalMs - (uint64(base.Unix()*1000) % types.BlockIntervalMs)
	blocktime := base.Add(time.Duration(mintimenextblock) * time.Millisecond)

	// If we would wait less than 50ms (1/10 of block_interval), wait for the whole block interval.
	interv := int64(types.BlockIntervalMs) * int64(time.Millisecond) / 10
	if blocktime.Sub(now).Nanoseconds() < interv {
		blocktime = blocktime.Add(time.Duration(types.BlockIntervalUs) * time.Microsecond)
	}

	pp.pendingBlockMode = producing
	b := types.NewBlockTime(blocktime)
	b.SetTime(blocktime)
	scheduledProducer := headstate.GetScheduledProducer(*b)
	currWaterMark, ok := pp.producerWaterMarks[scheduledProducer.ProducerName]

	if !pp.productionEnabled {
		// pp.pendingBlockMode = speculating
	} else if !ok {
		pp.pendingBlockMode = speculating
	}

	if pp.pendingBlockMode == producing {
		if ok && currWaterMark >= (headstate.BlockNum+1) {
			pp.pendingBlockMode = speculating
		}
	}

	//exception handle
	helper.CatchException(nil, func() {
		panic(errors.New("producerplugin start block panic"))
	})

	blocksToConfirm := 0
	if pp.pendingBlockMode == producing && ok {
		if currWaterMark < headstate.BlockNum {
			cm := uint16(headstate.BlockNum - currWaterMark)
			if cm < math.MaxUint16 {
				blocksToConfirm = int(cm)
			} else {
				blocksToConfirm = math.MaxUint16
			}

		}
	}

	pp.chain.AbortBlock()
	pp.chain.StartBlock(b, uint16(blocksToConfirm), types.Incomplete)

	pbs := pp.chain.PendingBlockState()
	if pbs != nil {
		if pp.pendingBlockMode == producing && pbs.BlockSigningKey != scheduledProducer.SigningKey {
			pp.pendingBlockMode = speculating
		}

		bexhausted := false
		unappliedTrxs := pp.chain.GetUnappliedTransaction()

		//remove all persisted transactions that have now expired
		headpoint := pbs.Header.TimeStamp.Time.Unix()
		for k, v := range pp.persistentTransaction {
			value := int64(v)
			if value <= headpoint {
				delete(pp.persistentTransaction, k)
			}
		}

		//push transaction
		for i := 0; i < len(unappliedTrxs); i++ {
			trx := unappliedTrxs[i]
			if _, ok := pp.persistentTransaction[trx.ID]; ok {
				pp.chain.PushTransaction(trx, types.MaxTime(), false, 0)
				trx = nil
			}
		}

		if pp.pendingBlockMode == producing && len(unappliedTrxs) > 0 {
			for _, v := range unappliedTrxs {
				if bexhausted {
					break
				}

				if v == nil {
					continue
				}

				if v.PackedTrx.Expiration() < uint64(pbs.Header.TimeStamp.Time.Unix()) {
					pp.chain.DropUnappliedTransaction(v)
					continue
				}

				deadline := time.Now().Add(5 * time.Second)
				if blocktime.Unix() < deadline.Unix() {
					deadline = blocktime
				}

				trace := pp.chain.PushTransaction(v, deadline, false, 0)
				if trace.Except != nil {
					bexhausted = true
				}
			}
			return succeeded
		}

	}

	return failed
}

func (pp *ProducerPlugin) calculateNextBlockTime(producerName string) time.Time {
	return time.Now().Add(time.Second)
}

//ProduceBlock produce a block
func (pp *ProducerPlugin) produceBlock() {
	pbs := pp.chain.PendingBlockState()
	// hbs := pp.chain.HeadBlockState()
	if pbs == nil {
		log.Print("pending_block_state does not exist but it should, another plugin may have corrupted it")
		return
	}

	pp.chain.FinalizeBlock()
	pp.chain.SignBlock(pbs.Block.Signature, true)
	pp.chain.CommitBlock(true)

	// newhbs := pp.chain.HeadBlockState()

	// if hbs.ActiveSchedule.Version != newhbs.ActiveSchedule.Version {
	// 	newProducers := make(map[string]struct{}, len(newhbs.ActiveSchedule.Producers))
	// 	for _, v := range newhbs.ActiveSchedule.Producers {
	// 		newProducers[v.ProducerName] = struct{}{}
	// 	}

	// 	for _, p := range hbs.ActiveSchedule.Producers {
	// 		delete(newProducers, p.ProducerName)
	// 	}

	// 	for value := range newProducers {
	// 		pp.producerWaterMarks[value] = pp.chain.HeadBlockNum()
	// 	}
	// }
	// pp.producerWaterMarks[newhbs.Header.Producer] = pp.chain.HeadBlockNum()
}

//MaybeProduceBlock call produceblock
func (pp *ProducerPlugin) MaybeProduceBlock() (res bool) {
	// timeout := time.Duration(types.BlockIntervalMs) * time.Millisecond
	defer pp.scheduleProductionLoop()

	helper.CatchException(nil, func() {
		log.Print("ProducerPlugin MaybeProduceBlock panic")
		pp.chain.AbortBlock()
		res = false
	})

	pp.produceBlock()
	return true
}

//OnIncomingBlock handle accept a block
func (pp *ProducerPlugin) OnIncomingBlock(block *types.Block) {
	timeout := time.Duration(types.BlockIntervalMs) * time.Millisecond
	defer pp.resetTimer(timeout)

	id := block.Hash()
	existing := pp.chain.FetchBlockByID(id)
	if existing == nil {
		return
	}

	pp.chain.AbortBlock()

	helper.CatchException(nil, func() {

	})

	pp.chain.PushBlock(existing, types.Complete)

	pp.productionEnabled = true
}

//OnIncomingTransactionAsync handle incoming transaction
func (pp *ProducerPlugin) OnIncomingTransactionAsync(trx *types.PackedTransaction, persistuntilexpired bool, next func(inerr error, trace *types.TransactionTrace)) {
	blocktime := pp.chain.PendingBlockTime()

	sendResponse := func(err error, trx *types.PackedTransaction, pt *types.TransactionTrace) {
		next(err, pt)
		if err != nil {
			pp.chain.TransactionAckChan <- &types.TrxTrace{Err: err, Trx: trx}
		} else {
			pp.chain.TransactionAckChan <- &types.TrxTrace{Err: nil, Trx: trx}
		}
	}

	id := trx.ID()

	// expire := int64(trx.Expiration())
	// btexpire := blocktime.Unix()

	// if expire < btexpire {
	// 	errs := fmt.Errorf("expired transaction id={%v}", id)
	// 	sendResponse(errs, trx, nil)
	// 	return
	// }

	if pp.chain.IsKnownUnexpiredTransaction(id) {
		errs := fmt.Errorf("duplicate transaction id={%v}", id)
		sendResponse(errs, trx, nil)
		return
	}

	deadline := time.Now().Add(time.Duration(pp.config.maxTrxTimeMs) * time.Millisecond)
	if blocktime.Unix() < deadline.Unix() {
		deadline = blocktime
	}

	helper.CatchException(nil, func() {
		panic("OnIncomingTransactionAsync panic")
	})

	meta := types.NewMetaDataByPackedTrx(trx)
	trace := pp.chain.PushTransaction(meta, deadline, false, 0)
	if trace.Except != nil {
		sendResponse(trace.Except, trx, nil)
	} else {
		if persistuntilexpired {
			pp.persistentTransaction[trx.ID()] = trx.Expiration()
		}

		sendResponse(nil, trx, trace)
	}
}
