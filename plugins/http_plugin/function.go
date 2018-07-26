package httpplugin

import (
	"datx_chain/chainlib/application"
	"datx_chain/chainlib/controller"
	"datx_chain/chainlib/types"
	"datx_chain/plugins/chain_plugin"
	"datx_chain/plugins/producer_plugin"
	"datx_chain/utils/common"
	"datx_chain/utils/rlp"
	"encoding/json"
	"log"
	"time"
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
func FindTransaction() *controller.TransactionDetails {
	var result *controller.TransactionDetail
	var tdList *controller.TransactionDetails

	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
	}

	chainplugin := plugin.(*chainplugin.ChainPlugin)
	fdb := chainplugin.Chain().ForkDB
	head := fdb.Head().Block.BlockHeader
	trxList := fdb.Head().Trxs
	result.BlockHeight = int64(head.BlockNum)
	result.TimeStamp = head.TimeStamp.Time.String()
	result.Pending = "True"
	result.TrxType = "Transfer"
	result.Amount = 0

	transfer := &controller.Transfer{}
	for _, v := range trxList {
		hash, err := rlp.EncodeToBytes(v.ID)
		if err != nil {
			log.Printf("get rlp err={%v}", err)
		}
		result.TrxHash = rlp.ByteString(hash)

		for _, action := range v.Trx.Actions {
			data := action.Data
			if err := json.Unmarshal(data, &transfer); err == nil {
				result.Amount += float64(transfer.Amount)
			} else {
				log.Printf("Unmarshal json err : %v", err)
			}
		}
		tdList.TrxDetailList = append(tdList.TrxDetailList, result)
	}
	log.Printf("Transactionlist is  %v", tdList)
	return tdList

}

//QueryTransactionById query transactiondetail by trx hash
func QueryTransactionById(id common.Hash) *controller.TransactionDetail {
	var result *controller.TransactionDetail

	//init chain
	plugin, err := application.App().Find("chain")
	if err != nil {
		log.Print("you do not add chainplugin to app before init the producerplugin")
	}
	chainplugin := plugin.(*chainplugin.ChainPlugin)
	chain := chainplugin.Chain()

	trx, t := chain.ForkDB.GetTrx(id)
	result.TrxHash = trx.TransactionHash.String()
	result.TrxType = "Transfer"
	result.TimeStamp = t.Time.String()
	result.Amount = 0

	transfer := &controller.Transfer{}
	for _, action := range trx.Actions {
		data := action.Data
		if err := json.Unmarshal(data, &transfer); err == nil {
			result.Amount += float64(transfer.Amount)
			result.AccountFrom = transfer.From
		} else {
			log.Printf("Unmarshal json err : %v", err)
		}
	}

	log.Printf("QueryTransactionById is  %v", result)
	return result
}

//BlockByIDResult 返回的json格式
type BlockByIDResult struct {
	BlockHeight int64   //区块高度
	BlockID     string  //区块ID，从common.Hash转成string格式
	TranxCnt    int64   //交易量
	Amount      float64 //总金额
	Rewards     float64 //奖励数量
	Producer    string  //打块节点账户名
	TimeStamp   string  //时间戳，从BlockTime转成时间格式
	Pending     string  //Pending，统一填True
}

//QueryBlockByID 按BlockID查询
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
	res.TimeStamp = blockstate.BlockHeaderState.Header.TimeStamp.Time.String()
	res.Pending = "True"
	res.TranxCnt = 0
	res.Amount = 0
	res.Rewards = 50

	transfer := controller.Transfer{}

	for _, t := range blockstate.Trxs {
		for _, a := range t.Trx.Transaction.Actions {
			if err := json.Unmarshal(a.Data, &transfer); err == nil {
				res.TranxCnt++
				res.Amount += float64(transfer.Amount)
			} else {
				log.Printf("Unmarshal json err : %v", err)
			}
		}
	}
	return res
}

//QueryBlocks 区块列表
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

	return res
}

//GeneralInfo 返回的json格式
type GeneralInfo struct {
	Price       float64 //DATx 价格。https://api.coinmarketcap.com/v2/ticker/2567/
	BlockHeight int64   //区块高度
	LibNum      int64   //LIB NUM
	TranxCnt    int64   //交易总量。实时统计
	AccountCnt  int64   //账号数。实时统计
	AgentCnt    int64   //受托人总数
}

//QueryGeneralInfo 查询总体信息
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

	transfer := controller.Transfer{}

	for _, val := range index {
		blockheight := int64(val.BlockHeaderState.BlockNum)
		if res.BlockHeight < blockheight {
			res.BlockHeight = blockheight
		}
		for _, t := range val.Trxs {
			for _, a := range t.Trx.Transaction.Actions {
				if err := json.Unmarshal(a.Data, &transfer); err == nil {
					res.TranxCnt++
				} else {
					log.Printf("Unmarshal json err : %v", err)
				}
			}
		}
	}

	res.LibNum = res.BlockHeight - 50

	return res
}
