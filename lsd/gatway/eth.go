package gatway

import (
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
	ETHREtrySeconds    int64 = 2
	ETHDelaySeconds    int64 = 60
	ETHIrreversibleCnt int64 = 12
)

type ETHTransaction struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

type ETHResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Result  []ETHTransaction `json:"result"`
}

type ETHBlockResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		BlockNumber string `json:"blockNumber"`
		TimeStamp   string `json:"timeStamp"`
		BlockMiner  string `json:"blockMiner"`
		BlockReward string `json:"blockReward"`
		Uncles      []struct {
			Miner         string `json:"miner"`
			UnclePosition string `json:"unclePosition"`
			Blockreward   string `json:"blockreward"`
		} `json:"uncles"`
		UncleInclusionReward string `json:"uncleInclusionReward"`
	} `json:"result"`
}

type ETHLatestBlockNum struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type ETHBrowser struct {
	url string

	apikey string

	tickAddress string

	handleHeight int64

	tick *server.ChainServer

	close chan bool
}

func NewETHBrowser(accountAddr string, server *server.ChainServer) *ETHBrowser {
	return &ETHBrowser{
		url:          "https://api.etherscan.io/api",
		apikey:       "8FT3VZVAS94PIHPYKEWPWHC4ZICB71RFSM",
		tickAddress:  accountAddr,
		handleHeight: 0,
		tick:         server,
		close:        make(chan bool),
	}
}

func (eth *ETHBrowser) GetTrxs(account string, startpos, endpos int64) ([]chainlib.Transaction, error) {
	req, err := http.NewRequest("GET", eth.url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("module", "account")
	q.Add("action", "txlist")
	q.Add("address", account)
	q.Add("startblock", string(startpos))
	q.Add("endblock", string(endpos))
	q.Add("sort", "asc")
	q.Add("apikey", eth.apikey)
	req.URL.RawQuery = q.Encode()

	//log.Printf("allurl: %s", req.URL.String())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, fmt.Errorf("ETH GetTrxs:%v\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rest ETHResponse
	if err := json.Unmarshal(body, &rest); err != nil {
		return nil, err
	}

	if rest.Status != "1" {
		return nil, fmt.Errorf("error : %s", rest.Message)
	}

	if len(rest.Result) == 0 {
		return nil, fmt.Errorf("error : trx not found")
	}

	// log.Printf("\n%v\n", rest.Result)
	result := make([]chainlib.Transaction, 0, len(rest.Result))

	for _, v := range rest.Result {
		var temp chainlib.Transaction
		var err error

		temp.Category = "ETH"
		temp.BlockNum, err = strconv.ParseInt(v.BlockNumber, 10, 64)
		if err != nil {
			return nil, err
		}

		temp.From = v.From
		temp.To = v.To

		temp.Amount, err = strconv.ParseFloat(v.Value, 64)
		temp.Amount /= 1000000000000000000
		if err != nil {
			return nil, err
		}
		if temp.Amount == 0 {
			continue
		}

		var numtime int64
		numtime, err = strconv.ParseInt(v.TimeStamp, 10, 64)
		if err != nil {
			return nil, err
		}
		temp.Time = time.Unix(numtime, 0)

		temp.TransactionID = v.Hash
		temp.IsIrrevisible = false

		numconfirm, _ := strconv.ParseInt(v.Confirmations, 10, 64)
		if numconfirm >= ETHIrreversibleCnt {
			temp.IsIrrevisible = true
		}

		result = append(result, temp)
	}

	return result, nil
}

func (eth *ETHBrowser) GetBlock(num int64) (*chainlib.Block, error) {
	req, err := http.NewRequest("GET", eth.url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("module", "block")
	q.Add("action", "getblockreward")
	q.Add("blockno", string(num))
	q.Add("apikey", eth.apikey)
	req.URL.RawQuery = q.Encode()

	log.Printf("ETH allurl: %s", req.URL.String())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, fmt.Errorf("ETH GetBlock:%v\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rest ETHBlockResp
	if err := json.Unmarshal(body, &rest); err != nil {
		return nil, err
	}

	if rest.Status != "1" {
		return nil, fmt.Errorf("ETH error : %s", rest.Message)
	}

	var result chainlib.Block
	result.BlockNum, err = strconv.ParseInt(rest.Result.BlockNumber, 10, 64)
	if err != nil {
		return nil, err
	}

	result.BlockID = ""
	var latest int64
	latest, err = eth.GetLatestBlockNum()
	if err != nil {
		return nil, err
	}

	result.Irreversible = false
	sub := latest - result.BlockNum
	if sub >= ETHIrreversibleCnt {
		result.Irreversible = true
	}

	return &result, nil
}

func (eth *ETHBrowser) Irreversible(blocknum int64) (bool, error) {
	latest, err := eth.GetLatestBlockNum()
	if err != nil {
		return false, err
	}

	sub := latest - blocknum
	if sub >= ETHIrreversibleCnt {
		return true, nil
	}

	return false, fmt.Errorf("ETH %d not irreversible,need %d confirmed.", blocknum, sub)
}

func (eth *ETHBrowser) GetLatestBlockNum() (int64, error) {
	req, err := http.NewRequest("GET", eth.url, nil)
	if err != nil {
		return 0, err
	}

	q := req.URL.Query()
	q.Add("module", "proxy")
	q.Add("action", "eth_blockNumber")
	q.Add("apikey", eth.apikey)
	req.URL.RawQuery = q.Encode()

	// log.Printf("allurl: %s", req.URL.String())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("ETH GetLatestBlockNum:%v\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var rest ETHLatestBlockNum
	if err := json.Unmarshal(body, &rest); err != nil {
		return 0, err
	}

	str := rest.Result

	num, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		return 0, err
	}

	log.Printf("\nETH latest blocknum: %v\n", num)
	return num, nil
}

func (eth *ETHBrowser) SetTickAccountAddr(account string) {
	eth.tickAddress = account
}

func (eth *ETHBrowser) Tick() {
	trxs, err := eth.GetTrxs(eth.tickAddress, eth.handleHeight, 99999999)
	if err != nil {
		log.Printf("ETH Get trxs on tick %v\n", err)
		return
	}

	for _, trx := range trxs {
		if trx.IsIrrevisible {
			// log.Printf("trx is irreversible: %v\n", trx.TransactionID)
			//exec push action

			if eth.handleHeight < trx.BlockNum {
				eth.handleHeight = trx.BlockNum
				log.Printf("ETH trx irreversible from: %v  %v\n", trx.TransactionID, eth.handleHeight)
			}

			if strings.EqualFold(trx.To, eth.tickAddress) {
				chainlib.PushCharge(trx)
			}

		} else {
			jobid := trx.Category + "_" + trx.TransactionID
			if job, _ := delayqueue.Get(jobid); job != nil {
				log.Printf("ETH tick trx is existed: %v %v\n", trx.TransactionID, time.Now().Unix())
				continue
			}

			log.Printf("ETH Add eth task on tick: %v  %v\n", trx.TransactionID, time.Now().Unix())
			eth.tick.AddTask(trx, ETHDelaySeconds)
		}
	}
}

func (eth *ETHBrowser) ReTry(trx chainlib.Transaction) bool {
	blockNum := trx.BlockNum

	log.Printf("ETH ReTry eth on tick: %v  %v\n", trx.TransactionID, time.Now().Unix())

	sta, err := eth.Irreversible(blockNum)
	if err != nil || !sta {
		eth.tick.AddTask(trx, ETHDelaySeconds)
		return false
	}

	if strings.EqualFold(trx.To, eth.tickAddress) {
		chainlib.PushCharge(trx)
	} else {
		return false
	}

	return sta
}

func (eth *ETHBrowser) Close() {
	eth.close <- true
}
