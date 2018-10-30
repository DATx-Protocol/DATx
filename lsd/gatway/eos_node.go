package gatway

import (
	"datx/lsd/chainlib"
	"datx/lsd/delayqueue"
	"datx/lsd/server"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	EOSRetrySeconds int64 = 2
)

//Action ...
type Action struct {
	GlobalActionSeq  int    `json:"global_action_seq"`
	AccountActionSeq int    `json:"account_action_seq"`
	BlockNum         int    `json:"block_num"`
	BlockTime        string `json:"block_time"`
	ActionTrace      struct {
		Receipt struct {
			Receiver       string          `json:"receiver"`
			ActDigest      string          `json:"act_digest"`
			GlobalSequence int             `json:"global_sequence"`
			RecvSequence   int             `json:"recv_sequence"`
			AuthSequence   [][]interface{} `json:"auth_sequence"`
			CodeSequence   int             `json:"code_sequence"`
			AbiSequence    int             `json:"abi_sequence"`
		} `json:"receipt"`
		Act struct {
			Account       string `json:"account"`
			Name          string `json:"name"`
			Authorization []struct {
				Actor      string `json:"actor"`
				Permission string `json:"permission"`
			} `json:"authorization"`
			Data struct {
				From     string `json:"from"`
				To       string `json:"to"`
				Quantity string `json:"quantity"`
				Memo     string `json:"memo"`
			} `json:"data"`
			HexData string `json:"hex_data"`
		} `json:"act"`
		Elapsed       int           `json:"elapsed"`
		CPUUsage      int           `json:"cpu_usage"`
		Console       string        `json:"console"`
		TotalCPUUsage int           `json:"total_cpu_usage"`
		TrxID         string        `json:"trx_id"`
		InlineTraces  []interface{} `json:"inline_traces"`
	} `json:"action_trace"`
}

//AccountActions ...
type AccountActions struct {
	Actions               []Action `json:"actions"`
	LastIrreversibleBlock int      `json:"last_irreversible_block"`
}

type EOSNode struct {
	url string

	tickAccount string

	pos int64

	lastIrreversibleBlockNum int64

	close chan bool

	tick *server.ChainServer
}

//NewEOSNode ...
func NewEOSNode(url, account string, server *server.ChainServer) *EOSNode {
	return &EOSNode{
		url:                      url,
		tickAccount:              account,
		pos:                      0,
		lastIrreversibleBlockNum: 0,
		tick:                     server,
		close:                    make(chan bool),
	}
}

//GetAccountActions https://api.eosmonitor.io/v1/actions?account=eostea111111&name=transfer&page=1&per_page=30
func (eos *EOSNode) GetAccountActions(accountAddr string) ([]chainlib.Transaction, error) {
	cmdStr := fmt.Sprintf("cldatx -u %s get actions %s -j %d %d", eos.url, eos.tickAccount, eos.pos, 10)

	log.Printf("\n****************************************\n%s\n****************************************\n", cmdStr)
	res, err := chainlib.ExecShell(cmdStr)
	if err != nil {
		log.Printf("\nGetAccountActions get actions return: %s\n", err)
		return nil, err
	}

	var resp AccountActions
	if err := json.Unmarshal([]byte(res), &resp); err != nil {
		log.Printf("[EOSNode] GetAccountActions unmarsh :%v\n", err)
		return nil, err
	}

	if len(resp.Actions) > 0 {
		eos.pos = eos.pos + 10
	}

	eos.lastIrreversibleBlockNum = int64(resp.LastIrreversibleBlock)

	result := make([]chainlib.Transaction, 0)
	for _, v := range resp.Actions {
		if v.ActionTrace.Act.Name != "transfer" {
			continue
		}

		var temp chainlib.Transaction
		temp.TransactionID = v.ActionTrace.TrxID
		temp.Category = "EOS"
		temp.BlockNum = int64(v.BlockNum)
		temp.From = v.ActionTrace.Act.Data.From
		temp.To = v.ActionTrace.Act.Data.To
		amountpos := strings.Index(v.ActionTrace.Act.Data.Quantity, " ")
		amountstr := v.ActionTrace.Act.Data.Quantity[:amountpos]
		temp.Amount, _ = strconv.ParseFloat(amountstr, 64)
		temp.Memo = v.ActionTrace.Act.Data.Memo

		temp.Time = time.Now()
		temp.IsIrrevisible = false
		if eos.lastIrreversibleBlockNum >= temp.BlockNum {
			temp.IsIrrevisible = true
		}

		result = append(result, temp)
	}

	return result, nil
}

//SetTickAccountAddr set account
func (eos *EOSNode) SetTickAccountAddr(account string) {
	eos.tickAccount = account
}

//Tick execute per second
func (eos *EOSNode) Tick() {
	trxs, err := eos.GetAccountActions(eos.tickAccount)
	if err != nil {
		log.Printf("Get trxs err: %v\n", err)
		return
	}

	for _, trx := range trxs {
		if trx.IsIrrevisible {
			// log.Printf("eos trx is irreversible: %v\n", trx.TransactionID)

			var result error
			if trx.To == eos.tickAccount {
				result = chainlib.PushCharge(trx)
			} else if trx.From == eos.tickAccount {
				result = chainlib.PushExtract(trx)
			}

			log.Printf("[Tick] EOS push action result:%v\n", result)
		} else {
			jobid := trx.Category + "_" + trx.TransactionID
			if job, _ := delayqueue.Get(jobid); job != nil {
				log.Printf("eos trx is existed: %v\n", trx.TransactionID)
				continue
			}

			log.Printf("add eos task: %v  %v\n", trx.TransactionID, time.Now().Unix())
			eos.tick.AddTask(trx, EOSDelaySeconds)
		}
	}
}

//ReTry ...
func (eos *EOSNode) ReTry(trx chainlib.Transaction) bool {
	if eos.lastIrreversibleBlockNum < trx.BlockNum {
		eos.tick.AddTask(trx, EOSDelaySeconds)
		return false
	}

	if trx.To == eos.tickAccount {
		chainlib.PushCharge(trx)
	} else if trx.From == eos.tickAccount {
		chainlib.PushExtract(trx)
	} else {
		return false
	}

	return true
}

//Close ...
func (eos *EOSNode) Close() {
	eos.close <- true
}
