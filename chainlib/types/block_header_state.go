package types

import (
	"bytes"
	"datx_chain/chainlib/chainobject"
	"datx_chain/utils/common"
	"datx_chain/utils/crypto"
	"datx_chain/utils/helper"
	"log"
	"sort"
	"strings"
)

//BlockHeaderState struct
type BlockHeaderState struct {
	ID       common.Hash
	BlockNum uint32
	Header   BlockHeader

	DposProposedIrreversibleBlockNum uint32
	DposIrreversibleBlockNum         uint32
	BftIrreversibleBlockNum          uint32
	PendingScheduleLibNum            uint32

	PendingSchedleHash common.Hash

	PendingSchedule chainobject.ProducerSchedule
	ActiveSchedule  chainobject.ProducerSchedule

	ProducerToLastProduced   map[string]interface{} //the pairs of account name and block num
	ProducerToLastImpliedIrb map[string]interface{}

	BlockSigningKey common.Hash
	Count           []uint8
	Confirmations   []HeaderConfirmation
}

//NewBlockHeaderState new
func NewBlockHeaderState(b Block) BlockHeaderState {
	return BlockHeaderState{
		ID:                       b.ID,
		BlockNum:                 b.BlockNum,
		Header:                   b.BlockHeader,
		ProducerToLastProduced:   make(map[string]interface{}, 0),
		ProducerToLastImpliedIrb: make(map[string]interface{}, 0),
	}
}

//MakeBlockHerderState make
func MakeBlockHerderState() BlockHeaderState {
	return BlockHeaderState{
		ProducerToLastProduced:   make(map[string]interface{}, 0),
		ProducerToLastImpliedIrb: make(map[string]interface{}, 0),
	}
}

//GenerateNext Generate Next BlockHeaderState
func (bhs *BlockHeaderState) GenerateNext(when *BlockTime) *BlockHeaderState {
	result := MakeBlockHerderState()

	helper.CatchException(nil, func() {
		log.Print("GenerateNext panic")
	})

	// if when.Less(bhs.Header.TimeStamp) {
	// 	log.Printf("next block must be in the future:%v %v", when.Time, bhs.Header.TimeStamp.Time)
	// 	return nil
	// }

	// when.Plus()
	result.Header.TimeStamp = when
	result.Header.Previous = bhs.ID
	result.Header.BlockNum = bhs.BlockNum + 1
	result.Header.ScheduleVersion = bhs.ActiveSchedule.Version

	prokey := bhs.GetScheduledProducer(*when)
	result.BlockSigningKey = prokey.SigningKey
	result.Header.Producer = prokey.ProducerName

	result.BlockNum = bhs.BlockNum + 1
	result.ProducerToLastProduced = bhs.ProducerToLastProduced
	result.ProducerToLastImpliedIrb = bhs.ProducerToLastImpliedIrb
	result.ProducerToLastProduced[prokey.ProducerName] = result.BlockNum

	result.ActiveSchedule = bhs.ActiveSchedule
	result.PendingSchedule = bhs.PendingSchedule
	result.DposProposedIrreversibleBlockNum = bhs.DposProposedIrreversibleBlockNum
	result.BftIrreversibleBlockNum = bhs.BftIrreversibleBlockNum
	result.ProducerToLastImpliedIrb[prokey.ProducerName] = result.DposProposedIrreversibleBlockNum

	//calculate
	result.DposIrreversibleBlockNum = bhs.calcDposLastIrreversible()

	numActiveProducers := len(bhs.ActiveSchedule.Producers)
	requiredConfs := (uint8)(numActiveProducers*2/3 + 1)

	size := len(bhs.Count)

	if size < MaxTrackedDposConfirmations {
		result.Count = make([]uint8, size+1)
		result.Count = append(bhs.Count, requiredConfs)
	} else {
		result.Count = make([]uint8, size)
		result.Count = append(bhs.Count[1:], requiredConfs)
	}

	return &result
}

//Next BlockHeaderState
func (bhs *BlockHeaderState) Next(b *BlockHeader, trust bool) *BlockHeaderState {
	// if b.TimeStamp.Less(bhs.Header.TimeStamp) {
	// 	log.Print("Block must be later in time.\n")
	// 	return nil
	// }

	log.Printf("Slot Next: %v", b.TimeStamp.Slot)
	result := bhs.GenerateNext(b.TimeStamp)
	if strings.Compare(result.Header.Producer, b.Producer) != 0 {
		log.Print("Wrong producer specified.\n")
		return nil
	}

	if result.Header.ScheduleVersion != b.ScheduleVersion {
		log.Print("Schedule version om signed block is corrupted.\n")
		return nil
	}

	if value, ok := bhs.ProducerToLastProduced[b.Producer]; ok {
		if value.(uint32) >= (result.BlockNum - (uint32)(b.Confirmed)) {
			log.Printf("Producer %s double-confirming know range", b.Producer)
			// return nil
		}
	}

	result.SetConfirm(b.Confirmed)
	wasPendingPromoted := bhs.MaybePromotePending()
	if len(b.NewProducers.Producers) != 0 {
		if wasPendingPromoted {
			log.Print("cannot set pending producer schedule in the same block in whtch pending was promoted to active")
			return nil
		}
		result.SetNewProducers(b.NewProducers)
	}

	copy(result.Header.Signature, b.Signature)

	result.ID = result.Header.Hash()

	return result
}

//MaybePromotePending bool
func (bhs *BlockHeaderState) MaybePromotePending() bool {
	if len(bhs.PendingSchedule.Producers) > 0 && bhs.DposIrreversibleBlockNum >= bhs.PendingScheduleLibNum {
		bhs.ActiveSchedule = bhs.PendingSchedule

		newProducerToLastProdeced := make(map[string]interface{})
		for _, v := range bhs.ActiveSchedule.Producers {
			value, ok := bhs.ProducerToLastProduced[v.ProducerName]
			if ok {
				newProducerToLastProdeced[v.ProducerName] = value
			} else {
				newProducerToLastProdeced[v.ProducerName] = bhs.DposIrreversibleBlockNum
			}
		}

		newProducerToLastImpliedIrb := make(map[string]interface{})
		for _, v := range bhs.ActiveSchedule.Producers {
			value, ok := bhs.ProducerToLastImpliedIrb[v.ProducerName]
			if ok {
				newProducerToLastImpliedIrb[v.ProducerName] = value
			} else {
				newProducerToLastImpliedIrb[v.ProducerName] = bhs.DposIrreversibleBlockNum
			}
		}

		bhs.ProducerToLastProduced = newProducerToLastProdeced
		bhs.ProducerToLastImpliedIrb = newProducerToLastImpliedIrb
		bhs.ProducerToLastProduced[bhs.Header.Producer] = bhs.BlockNum

		return true
	}

	return false
}

//SetNewProducers set new producers
func (bhs *BlockHeaderState) SetNewProducers(pending chainobject.ProducerSchedule) {
	if pending.Version != (bhs.ActiveSchedule.Version + 1) {
		log.Print("Wrong producer schedule version specified.\n")
		return
	}

	if len(bhs.PendingSchedule.Producers) != 0 {
		log.Print("Cannot set new pending producers until last pending is confirmed.\n")
		return
	}

	bhs.Header.NewProducers = pending
	bhs.PendingSchedleHash = helper.RLPHash(bhs.Header.NewProducers)
	bhs.PendingSchedule = bhs.Header.NewProducers
	bhs.PendingScheduleLibNum = bhs.BlockNum

	produs := []chainobject.ProducerKey{{ProducerName: "alice", SigningKey: common.Hash{}}, {ProducerName: "bob", SigningKey: common.Hash{}}}
	bhs.ActiveSchedule = chainobject.ProducerSchedule{Version: 0, Producers: produs}

}

//SetActiveProducers set active producers
func (bhs *BlockHeaderState) SetActiveProducers() {
	// if active.Version != (bhs.ActiveSchedule.Version + 1) {
	// 	log.Print("Wrong producer schedule version specified.\n")
	// 	return
	// }

	// if len(bhs.ActiveSchedule.Producers) != 0 {
	// 	log.Print("Cannot set new pending producers until last pending is confirmed.\n")
	// 	return
	// }

	// bhs.ActiveSchedule = active
	produs := []chainobject.ProducerKey{{ProducerName: "alice", SigningKey: common.Hash{}}, {ProducerName: "bob", SigningKey: common.Hash{}}}
	bhs.ActiveSchedule = chainobject.ProducerSchedule{Version: 0, Producers: produs}
}

//SetConfirm set confirm
func (bhs *BlockHeaderState) SetConfirm(numPrev uint16) {
	bhs.Header.Confirmed = numPrev

	i := len(bhs.Count) - 1
	blocksToConfirm := numPrev + 1

	for i >= 0 && blocksToConfirm > 0 {
		bhs.Count[i]--

		if bhs.Count[i] == 0 {
			numi := bhs.BlockNum - (uint32)(len(bhs.Count)-1-i)
			bhs.DposProposedIrreversibleBlockNum = numi

			if i == (len(bhs.Count) - 1) {
				bhs.Count = make([]uint8, 0)
			} else {
				bhs.Count = bhs.Count[i+1 : len(bhs.Count)-i-1]
			}

			return
		}

		i--
		blocksToConfirm--
	}
}

//AddConfirm add confirmation
func (bhs *BlockHeaderState) AddConfirm(c *HeaderConfirmation) {
	for _, v := range bhs.Confirmations {
		if v.Producer != c.Producer {
			log.Print("block already confirmed by this producer")
			return
		}
	}

	key := bhs.ActiveSchedule.GetProducerKey(c.Producer)
	empt := common.Hash{}
	if key == empt {
		log.Print("producer not in current schedule")
		return
	}

	//

	bhs.Confirmations = append(bhs.Confirmations, *c)
}

//calcDposLastIrreversible calc dpos last irreversible block
func (bhs *BlockHeaderState) calcDposLastIrreversible() uint32 {
	blocknums := make([]uint32, len(bhs.ProducerToLastImpliedIrb))

	for _, v := range bhs.ProducerToLastImpliedIrb {
		blocknums = append(blocknums, v.(uint32))
	}

	if len(blocknums) == 0 {
		return 0
	}

	// 2/3 must be greater, so if i go 1/3 into the list sorted from low to high. then 2/3 are greater
	sort.Slice(blocknums, func(i, j int) bool {
		return blocknums[i] < blocknums[j]
	})

	index := (len(blocknums) - 1) / 3

	return blocknums[index]
}

//GetScheduledProducer get schedule producer
func (bhs *BlockHeaderState) GetScheduledProducer(t BlockTime) chainobject.ProducerKey {
	produs := []chainobject.ProducerKey{{ProducerName: "alice", SigningKey: common.Hash{}}, {ProducerName: "bob", SigningKey: common.Hash{}}}
	bhs.ActiveSchedule = chainobject.ProducerSchedule{Version: 0, Producers: produs}

	allBlocks := len(bhs.ActiveSchedule.Producers) * ProducerRepetitions
	index := (int)(t.Slot) % allBlocks

	index /= ProducerRepetitions
	return bhs.ActiveSchedule.Producers[index]
}

//SignDigest is to sign with BlockHeaderState
func (bhs *BlockHeaderState) SignDigest() []byte {

	//var headMroot = crypto.Keccak256Hash(bhs.Header.Hash().Bytes(), GetMroot().Bytes())
	//TODO should SignDiest with blockroot_merkle
	return crypto.Keccak256(bhs.Header.Hash().Bytes(), bhs.PendingSchedleHash.Bytes())
}

//Sign is to sign for block when trust
func (bhs *BlockHeaderState) Sign(sig []byte, trust bool) {
	var d = bhs.SignDigest()
	var buffer bytes.Buffer
	buffer.Write(bhs.Header.Signature)
	buffer.Write(d)
	bhs.Header.Signature = buffer.Bytes()
	if true == trust {
		bsKey, e := crypto.SigToPub(d, sig)
		if e != nil {
			log.Printf("Sign the block with error,%v %d", e.Error())
		}
		bhs.BlockSigningKey = crypto.Keccak256Hash(crypto.CompressPubkey(bsKey))
	}
	// sig, err := crypto.Sign(bhs.Header.Signature, prv)
	// if err != nil {
	// 	log.Printf("sign block with error,%v %d", err.Error())
	// 	return err
	// }
	// buffer.Write(sig)
	// bhs.Header.Signature = buffer.Bytes()

}
