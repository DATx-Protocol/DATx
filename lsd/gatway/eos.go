package gatway

import (
	"datx/lsd/chainlib"
	"datx/lsd/delayqueue"
	"datx/lsd/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	EOSDelaySeconds int64 = 2
)

type EOSBlocks struct {
	BlockNum                         int    `json:"block_num"`
	BlockID                          string `json:"block_id"`
	PreviousBlockID                  string `json:"previous_block_id"`
	Timestamp                        int    `json:"timestamp"`
	ActionMroot                      string `json:"action_mroot"`
	TransactionMroot                 string `json:"transaction_mroot"`
	Producer                         string `json:"producer"`
	ScheduleVersion                  int    `json:"schedule_version"`
	ProducerSignature                string `json:"producer_signature"`
	BlockSigningKey                  string `json:"block_signing_key"`
	DposProposedIrreversibleBlocknum int    `json:"dpos_proposed_irreversible_blocknum"`
	DposIrreversibleBlocknum         int    `json:"dpos_irreversible_blocknum"`
	BftIrreversibleBlocknum          int    `json:"bft_irreversible_blocknum"`
	PendingScheduleLibNum            int    `json:"pending_schedule_lib_num"`
	PendingScheduleHash              string `json:"pending_schedule_hash"`
	Irreversible                     bool   `json:"irreversible"`
	TrxCount                         int    `json:"trx_count"`
}

type EOSTransactions struct {
	TransactionID    string   `json:"transaction_id"`
	ActionsCount     int      `json:"actions_count"`
	RefBlockNum      int      `json:"ref_block_num"`
	RefBlockPrefix   int64    `json:"ref_block_prefix"`
	MaxNetUsageWords int      `json:"max_net_usage_words"`
	MaxCPUUsageMs    int      `json:"max_cpu_usage_ms"`
	DelaySec         int      `json:"delay_sec"`
	Signatures       []string `json:"signatures"`
	BlockNum         int      `json:"block_num"`
	Expiration       int64    `json:"expiration"`
	Irreversible     bool     `json:"irreversible"`
	// BlockID	没有这个字段
	// Expiration	类型修正
	// Signatures	类型修正
	// SigningKeys	没有这个字段
	// TransactionExtensions	不清楚什么类型，舍弃
}

type EOSAction struct {
	Name           string `json:"name"`
	TrxID          string `json:"trx_id"`
	ActionNum      int    `json:"action_num"`
	HandlerAccount string `json:"handler_account"`
	Authorization  []struct {
		Actor      string `json:"actor"`
		Permission string `json:"permission"`
	} `json:"authorization"`
	Expiration int `json:"expiration"`
	Data       struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Quantity string `json:"quantity"`
		Memo     string `json:"memo"`
	} `json:"data"`
}

type EOSAccountAction struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Page    int         `json:"page"`
		PerPage int         `json:"per_page"`
		Actions []EOSAction `json:"data"`
		Total   int         `json:"total"`
		HasNext bool        `json:"has_next"`
		HasPrev bool        `json:"has_prev"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

type EOSAccounts struct {
	Name        string `json:"name"`
	Permissions []struct {
		PermName     string `json:"perm_name"`
		Parent       string `json:"parent"`
		RequiredAuth struct {
			Threshold int `json:"threshold"`
			Keys      []struct {
				Key    string `json:"key"`
				Weight int    `json:"weight"`
			} `json:"keys"`
			Accounts []interface{} `json:"accounts"`
			Waits    []interface{} `json:"waits"`
		} `json:"required_auth"`
	} `json:"permissions"`
	Privileged     bool `json:"privileged"`
	LastCodeUpdate int  `json:"last_code_update"`
	Created        int  `json:"created"`
	RAMQuota       int  `json:"ram_quota"`
	RAMUsage       int  `json:"ram_usage"`
	NetWeight      int  `json:"net_weight"`
	CPUWeight      int  `json:"cpu_weight"`
	NetLimit       struct {
		Used      int `json:"used"`
		Available int `json:"available"`
		Max       int `json:"max"`
	} `json:"net_limit"`
	CPULimit struct {
		Used      int `json:"used"`
		Available int `json:"available"`
		Max       int `json:"max"`
	} `json:"cpu_limit"`
	EosBalance     float64 `json:"eos_balance"`
	TotalResources struct {
		Owner     string `json:"owner"`
		NetWeight string `json:"net_weight"`
		CPUWeight string `json:"cpu_weight"`
		RAMBytes  int    `json:"ram_bytes"`
	} `json:"total_resources"`
	DelegatedBandwidth struct {
		From      string `json:"from"`
		To        string `json:"to"`
		NetWeight string `json:"net_weight"`
		CPUWeight string `json:"cpu_weight"`
	} `json:"delegated_bandwidth"`
	Creator string `json:"creator"`
}

type EOSResponseBlocks struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    EOSBlocks `json:"data"`

	Error interface{} `json:"error"`
}

type EOSResponseTransactions struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    EOSTransactions `json:"data"`

	Error interface{} `json:"error"`
}

type EOSResponseAccounts struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    EOSAccounts `json:"data"`

	Error interface{} `json:"error"`
}

type EOSBrowser struct {
	url string

	tickAccount string

	currentPage int64

	close chan bool

	tick *server.ChainServer
}

func NewEOSBrowser(account string, server *server.ChainServer) *EOSBrowser {
	return &EOSBrowser{
		url:         "https://api.eosmonitor.io/v1/",
		tickAccount: account,
		currentPage: 1,
		tick:        server,
		close:       make(chan bool),
	}
}

func (eos *EOSBrowser) GetBlocks(blocknum int64) (*chainlib.Block, error) {
	var allurl = eos.url + "blocks/" + string(blocknum)

	// fmt.Printf("GetBlocks request url is : %s\n", allurl)

	res, err := http.Get(allurl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("EOS GetBlocks:%v %v", allurl, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp EOSResponseBlocks
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("err: %s", resp.Message)
	}

	var result chainlib.Block
	result.BlockNum = int64(resp.Data.BlockNum)
	result.BlockID = resp.Data.BlockID
	result.Irreversible = resp.Data.Irreversible

	return &result, nil
}

func (eos *EOSBrowser) Irreversible(blocknum int64) (bool, error) {
	block, err := eos.GetBlocks(blocknum)
	if err != nil {
		return false, err
	}

	return block.Irreversible, nil
}

//GetTransaction GET https://api.eosmonitor.io/v1/transactions/<transaction_id>
func (eos *EOSBrowser) GetTransaction(trxid string) (*EOSTransactions, error) {
	var allurl = eos.url + "transactions/" + trxid

	fmt.Printf("GetTransaction request url is : %s\n", allurl)

	res, err := http.Get(allurl)
	if err != nil {
		fmt.Printf("geturl:%v\n", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("GetTransaction:%v %v\n", allurl, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp EOSResponseTransactions
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("response status err: %s", resp.Message)
	}

	return &resp.Data, nil
}

//GetAccountActions https://api.eosmonitor.io/v1/actions?account=eostea111111&name=transfer&page=1&per_page=30
func (eos *EOSBrowser) GetAccountActions(accountAddr string) ([]chainlib.Transaction, error) {
	var allurl = eos.url + "actions?account=" + accountAddr + "&name=transfer&page=" + fmt.Sprintf("%d", eos.currentPage) + "&per_page=10"

	fmt.Printf("EOS GetAccountActions request url is : %s\n", allurl)

	res, err := http.Get(allurl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("EOS GetAccountActions:%v %v", allurl, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp EOSAccountAction
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	total := int64(resp.Data.Total)
	remain := total - 10*eos.currentPage
	t := eos.currentPage + 1
	if remain > 0 {
		eos.currentPage = t
	}

	result := make([]chainlib.Transaction, 0)
	for _, v := range resp.Data.Actions {
		if v.Name != "transfer" {
			continue
		}

		eostrx, err := eos.GetTransaction(v.TrxID)
		if err != nil {
			fmt.Printf("GetAccountActions :%v %v", accountAddr, err)
			continue
		}

		var temp chainlib.Transaction
		temp.TransactionID = v.TrxID
		temp.Category = "EOS"
		temp.BlockNum = int64(eostrx.BlockNum)
		temp.From = v.Data.From
		temp.To = v.Data.To
		amountpos := strings.Index(v.Data.Quantity, " ") // EOS的Quantity是金额+空格+币种
		amountstr := v.Data.Quantity[:amountpos]         // 只需要空格前面的金额
		temp.Amount, _ = strconv.ParseFloat(amountstr, 64)
		temp.Memo = v.Data.Memo

		temp.Time = time.Now()
		temp.IsIrrevisible = eostrx.Irreversible

		result = append(result, temp)
	}

	return result, nil
}

func (eos *EOSBrowser) GetAccounts(args string) (*EOSResponseAccounts, error) {
	var allurl = eos.url + "accounts/" + args

	fmt.Printf("GetAccounts request url is : %s\n", allurl)

	res, err := http.Get(allurl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("EOS GetAccounts:%v %v\n", allurl, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp EOSResponseAccounts
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (eos *EOSBrowser) IsIrreversible(trxid string) (bool, error) {
	fmt.Printf("trsid len:%d\n", len(trxid))

	trx, err := eos.GetTransaction(trxid)
	if err != nil {
		return false, err
	}

	blockid := int64(trx.BlockNum)
	block, err := eos.GetBlocks(blockid)
	if err != nil {
		return false, err
	}

	return block.Irreversible, nil
}

//SetTickAccountAddr set account
func (eos *EOSBrowser) SetTickAccountAddr(account string) {
	eos.tickAccount = account
}

//Tick execute per second
func (eos *EOSBrowser) Tick() {
	trxs, err := eos.GetAccountActions(eos.tickAccount)
	if err != nil {
		fmt.Printf("Get trxs err: %v\n", err)
		return
	}

	for _, trx := range trxs {
		if trx.IsIrrevisible {
			// fmt.Printf("eos trx is irreversible: %v\n", trx.TransactionID)

			var result error
			if trx.To == eos.tickAccount {
				result = chainlib.PushCharge(trx)
			} else if trx.From == eos.tickAccount {
				result = chainlib.PushExtract(trx)
			}

			if result != nil {
				jobid := trx.Category + "_" + trx.TransactionID
				if job, _ := delayqueue.Get(jobid); job != nil {
					continue
				}

				fmt.Printf("BTC push action faile and retry on Tick(): %v  %v\n", trx.TransactionID, time.Now().Unix())
				eos.tick.AddTask(trx, EOSDelaySeconds)
			}

		} else {
			jobid := trx.Category + "_" + trx.TransactionID
			if job, _ := delayqueue.Get(jobid); job != nil {
				fmt.Printf("eos trx is existed: %v\n", trx.TransactionID)
				continue
			}

			fmt.Printf("add eos task: %v  %v\n", trx.TransactionID, time.Now().Unix())
			eos.tick.AddTask(trx, EOSDelaySeconds)
		}
	}
}

//ReTry ...
func (eos *EOSBrowser) ReTry(trx chainlib.Transaction) bool {
	blockNum := trx.BlockNum

	sta, err := eos.Irreversible(blockNum)
	if err != nil || !sta {
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

	return sta
}

//Close ...
func (eos *EOSBrowser) Close() {
	eos.close <- true
}
