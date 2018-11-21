package gatway

import (
	"bytes"
	"datx/lsd/chainlib"
	"datx/lsd/delayqueue"
	"datx/lsd/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	EOSRetrySeconds int64 = 2
)

//ActionsPara...
type ActionsPara struct {
	Pos         int64  `json:"pos"`
	Offset      int64  `json:"offset"`
	AccountName string `json:"account_name"`
}

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
	reqpara := ActionsPara{
		Pos:         eos.pos,
		Offset:      10,
		AccountName: eos.tickAccount,
	}

	req, err := json.Marshal(reqpara)
	if err != nil {
		return nil, err
	}

	para := bytes.NewBuffer([]byte(req))
	request, err := http.NewRequest("POST", eos.url+"/v1/history/get_actions", para)
	if err != nil {
		log.Printf("[EOSNode] GetAccountActions new http request: %v\n", err)
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, errr := client.Do(request)
	if errr != nil {
		log.Printf("[EOSNode] GetAccountActions do http request: %v\n", errr)
		return nil, errr
	}
	if response.StatusCode != 200 {
		log.Printf("[EOSNode] GetAccountActions http reponse: %v\n", response.Body)
		return nil, fmt.Errorf("[EOSNode] GetAccountActions Response error: %v", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("[EOSNode] GetAccountActions read body :%v\n", err)
		return nil, err
	}

	var resp AccountActions
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Printf("[EOSNode] GetAccountActions unmarsh :%v\n", err)
		return nil, err
	}

	eos.pos = eos.pos + int64(len(resp.Actions)) + 1

	eos.lastIrreversibleBlockNum = int64(resp.LastIrreversibleBlock)

	if len(resp.Actions) == 0 {
		return nil, fmt.Errorf("the account %s no trxs", eos.tickAccount)
	}

	result := make([]chainlib.Transaction, 0)
	for _, v := range resp.Actions {
		if v.ActionTrace.Act.Name != "transfer" {
			continue
		}

		if strings.Contains(v.ActionTrace.Act.Data.From, "eosio.") || strings.Contains(v.ActionTrace.Act.Data.To, "eosio.") {
			log.Printf("[EOSNode] GetAccountActions: from=%v to=%v\n", v.ActionTrace.Act.Data.From, v.ActionTrace.Act.Data.To)
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
		log.Printf("[EOSNode] GetAccountActions : %v\n", err)
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
			eos.tick.AddTask(trx, EOSRetrySeconds)
		}
	}
}

//ReTry ...
func (eos *EOSNode) ReTry(trx chainlib.Transaction) bool {
	if eos.lastIrreversibleBlockNum < trx.BlockNum {
		eos.tick.AddTask(trx, EOSRetrySeconds)
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
