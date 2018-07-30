package httpplugin

import (
	"datx_chain/chainlib/application"
	"datx_chain/chainlib/controller"
	"datx_chain/chainlib/types"
	"datx_chain/plugins/chain_plugin"
	"datx_chain/plugins/producer_plugin"
	"datx_chain/utils/common"
	"encoding/json"

	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

//CreateTransfer create a action transfer
func CreateTransfer(account, from, to string, amount uint16) types.Action {
	var result types.Action

	result.Account = account
	result.ActionName = "transfer"

	data := controller.Transfer{
		From:   from,
		To:     to,
		Amount: amount,
	}

	if v, err := json.Marshal(&data); err == nil {
		result.Data = v
	} else {
		log.Printf("Marshal hson err : %v", err)
	}

	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
	}
	chainplugin := plugin.(*chainplugin.ChainPlugin)
	chain := chainplugin.Chain()

	if _, ok := chain.TestAccounts[from]; !ok {
		var account controller.UserAccount
		account.Name = from
		account.Amount = 100

		chain.TestAccounts[from] = &account
	}

	if _, ok := chain.TestAccounts[to]; !ok {
		var account controller.UserAccount
		account.Name = to
		account.Amount = 100

		chain.TestAccounts[to] = &account
	}

	return result
}

//PushActions push action to transaction
func PushActions(actions ...types.Action) {
	var trx types.SignedTransaction

	trx.Actions = append(trx.Actions, actions...)

	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
		return
	}
	chainplugin := plugin.(*chainplugin.ChainPlugin)
	chain := chainplugin.Chain()

	trx.Expiration = uint64(chain.HeadBlockTime().Unix()) + uint64(30*time.Second)
	trx.RefBlockNum = chain.LastIrreversibleBlockNum()
	trx.RefBlockPerfix = chain.LastIrreversibleBlockID()

	metaTrx := types.NewTrxMetaData(&trx)
	chain.PushTransaction(metaTrx, types.MaxTime(), false, 0)
}

//CreatePackedTransaction create packed trx
func CreatePackedTransaction(actions ...types.Action) *types.PackedTransaction {
	var trx types.SignedTransaction

	trx.Actions = append(trx.Actions, actions...)

	res := types.NewPackedTransaction(&trx, types.None)

	return res
}

//PushTransaction http handler
func PushTransaction(pack *types.PackedTransaction, next func(inerr error, trace *types.TransactionTrace)) {
	go func() {
		//init chain
		plugin, err := application.App().Find("producer")
		if err != nil {
			log.Print("you do not add producerplugin to app before init the producerplugin")
			return
		}
		producer := plugin.(*producerplugin.ProducerPlugin)
		if producer != nil {
			producer.OnIncomingTransactionAsync(pack, true, next)
		}
	}()
}

//FindTransaction all transaction from forkdb
func FindTransaction() ([]*TransactionDetail, error) {
	result := &TransactionDetail{}
	tdList := make([]*TransactionDetail, 0)
	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
	}

	chainplugin := plugin.(*chainplugin.ChainPlugin)
	index := chainplugin.Chain().ForkDB.GetIndex()
	result.BlockHeight = 0

	result.Pending = "True"
	result.TrxType = "Transfer"
	result.Amount = 0

	transfer := &controller.Transfer{}
	//result.BlockHeight = int64(head.BlockHeaderState.BlockNum)
	for _, val := range index {
		blockheight := int64(val.BlockHeaderState.BlockNum)
		if result.BlockHeight < blockheight {
			result.BlockHeight = blockheight
		}
		t := val.BlockHeaderState.Header.TimeStamp.Time.Int64()
		result.TimeStamp = time.Unix(t, 0).String()

		for _, v := range val.Trxs {
			for _, action := range v.Trx.Transaction.Actions {
				data := action.Data
				if err := json.Unmarshal(data, &transfer); err == nil {
					result.Amount += float64(transfer.Amount)
					result.TrxHash = v.ID.String()
					tdList = append(tdList, result)
					log.Printf("*TransactionMetaData:%v,blockheight is :%v", result.TrxHash, result.BlockHeight)
				} else {
					return nil, err
				}

			}

		}

	}

	return tdList, nil

}

//QueryTransactionById query transactiondetail by trx hash
func QueryTransactionById(id common.Hash) *TransactionDetail {
	result := &TransactionDetail{}
	transfer := &controller.Transfer{}
	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
	}
	chainplugin := plugin.(*chainplugin.ChainPlugin)
	index := chainplugin.Chain().ForkDB.GetIndex()
	result.TrxType = "Transfer"

	result.Amount = 0
	for _, val := range index {
		blockheight := int64(val.BlockHeaderState.BlockNum)
		if result.BlockHeight < blockheight {
			result.BlockHeight = blockheight
		}

		t := val.BlockHeaderState.Header.TimeStamp.Time.Int64()
		result.TimeStamp = time.Unix(t, 0).String()

		for _, v := range val.Trxs {
			log.Printf("hash is -------:%v", id)
			if v.ID == id {
				result.TrxHash = v.ID.String()
				for _, action := range v.Trx.Actions {
					data := action.Data
					if err := json.Unmarshal(data, &transfer); err == nil {
						result.Amount += float64(transfer.Amount)
						result.AccountFrom = transfer.From

					} else {
						log.Printf("Unmarshal json err : %v", err)
					}
				}

			}
		}

	}

	log.Printf("QueryTransactionById is  %v", result)
	return result
}

//BlockByIDResult return json of Block
type BlockByIDResult struct {
	BlockHeight int64
	BlockID     string
	TranxCnt    int64
	Amount      float64
	Rewards     float64
	Producer    string
	TimeStamp   string
	Pending     string
}

//QueryBlockByID query by BlockID
func QueryBlockByID(id common.Hash) *BlockByIDResult {

	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
	}
	chainplugin := plugin.(*chainplugin.ChainPlugin)
	chain := chainplugin.Chain()

	var blockstate = chain.ForkDB.GetBlock(id)

	res := &BlockByIDResult{}
	res.BlockHeight = int64(blockstate.BlockHeaderState.BlockNum)
	res.BlockID = blockstate.BlockHeaderState.ID.String()
	res.Producer = blockstate.BlockHeaderState.Header.Producer
	ts := blockstate.BlockHeaderState.Header.TimeStamp.Time.Int64()
	res.TimeStamp = time.Unix(ts, 0).String()
	res.Pending = "True"
	res.TranxCnt = int64(len(blockstate.Trxs))
	res.Amount = 0
	res.Rewards = 50

	transfer := controller.Transfer{}

	for _, t := range blockstate.Trxs {
		for _, a := range t.Trx.Transaction.Actions {
			if err := json.Unmarshal(a.Data, &transfer); err == nil {
				res.Amount += float64(transfer.Amount)
			} else {
				return nil
			}
		}
	}
	return res
}

//BlockSlice used in sort
type BlockSlice []*BlockByIDResult

func (a BlockSlice) Len() int {
	return len(a)
}
func (a BlockSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a BlockSlice) Less(i, j int) bool {
	return a[j].BlockHeight < a[i].BlockHeight
}

//QueryBlocks return json of BlockList
func QueryBlocks() []*BlockByIDResult {

	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
	}
	chainplugin := plugin.(*chainplugin.ChainPlugin)
	chain := chainplugin.Chain()

	var index = chain.ForkDB.GetIndex()

	res := make([]*BlockByIDResult, 0)
	for key := range index {
		res = append(res, QueryBlockByID(key))
	}

	sort.Sort(BlockSlice(res))

	blockcnt := len(res)
	if blockcnt >= 20 {
		blockcnt = 20
	}
	return res[0:blockcnt]
}

//GeneralInfo return json of GeneralInfo
type GeneralInfo struct {
	Price       float64 //https://api.coinmarketcap.com/v2/ticker/2567/
	BlockHeight int64
	LibNum      int64
	TranxCnt    int64
	AccountCnt  int64
	AgentCnt    int64
}

//QueryGeneralInfo GeneralInfo
func QueryGeneralInfo() *GeneralInfo {

	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
	}
	chainplugin := plugin.(*chainplugin.ChainPlugin)
	chain := chainplugin.Chain()

	var index = chain.ForkDB.GetIndex()

	res := &GeneralInfo{}
	res.Price = 0
	res.BlockHeight = 0
	res.LibNum = 0
	res.TranxCnt = 0
	res.AccountCnt = 45
	res.AgentCnt = 2

	httpClient := &http.Client{}
	request, err := http.NewRequest("GET", "https://api.coinmarketcap.com/v2/ticker/2567/", nil)
	if err != nil {
		res.Price = 0
	}
	response, err := httpClient.Do(request)
	if err != nil {
		res.Price = 0
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		res.Price = 0
	}
	//fmt.Println(string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		res.Price = 0
	} else {
		res.Price = js.Get("data").Get("quotes").Get("USD").Get("price").MustFloat64()
	}

	for _, val := range index {
		blockheight := int64(val.BlockHeaderState.BlockNum)
		if res.BlockHeight < blockheight {
			res.BlockHeight = blockheight
		}
		res.TranxCnt += int64(len(val.Trxs))
	}

	if res.BlockHeight > 50 {
		res.LibNum = res.BlockHeight - 50
	}

	return res
}
