package gatway

import (
	"datx/ListenServer/chainlib"
	"datx/ListenServer/delayqueue"
	"datx/ListenServer/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
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

	//fmt.Printf("allurl: %s", req.URL.String())

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

	// fmt.Printf("\n%v\n", rest.Result)
	result := make([]chainlib.Transaction, 0, len(rest.Result))

	for _, v := range rest.Result {
		var temp chainlib.Transaction
		var err error

		temp.Category = "ETH"
		temp.BlockNum, err = strconv.ParseInt(v.BlockNumber, 10, 64)
		if err != nil {
			return nil, err
		}
		if v.From[0:2] == "0x" {
			temp.From = v.From[2:] // 去掉0x
		} else {
			temp.From = v.From
		}
		if v.To[0:2] == "0x" {
			temp.To = v.To[2:] // 去掉0x
		} else {
			temp.To = v.To
		}
		temp.Amount, err = strconv.ParseFloat(v.Value, 64)
		temp.Amount /= 1000000000000000000 //	监听单位是wei，转为ether
		if err != nil {
			return nil, err
		}

		var numtime int64
		numtime, err = strconv.ParseInt(v.TimeStamp, 10, 64)
		if err != nil {
			return nil, err
		}
		temp.Time = time.Unix(numtime, 0)

		if v.Hash[0:2] == "0x" {
			temp.TransactionID = v.Hash[2:] // 去掉0x
		} else {
			temp.TransactionID = v.Hash
		}
		temp.IsIrrevisible = false

		numconfirm, _ := strconv.ParseInt(v.Confirmations, 10, 64)
		if numconfirm > 6 {
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

	//https://api.etherscan.io/api?module=block&action=getblockreward&blockno=2165403&apikey=YourApiKeyToken

	q := req.URL.Query()
	q.Add("module", "block")
	q.Add("action", "getblockreward")
	q.Add("blockno", string(num))
	q.Add("apikey", eth.apikey)
	req.URL.RawQuery = q.Encode()

	fmt.Printf("ETH allurl: %s", req.URL.String())

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
	if sub > ETHIrreversibleCnt {
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
	if sub > ETHIrreversibleCnt {
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

	// fmt.Printf("allurl: %s", req.URL.String())

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

	fmt.Printf("\nETH latest blocknum: %v\n", num)
	return num, nil
}

func (eth *ETHBrowser) SetTickAccountAddr(account string) {
	eth.tickAddress = account
}

func (eth *ETHBrowser) Tick() {
	fmt.Printf("get blocknum from ontick: %v\n", eth.handleHeight)
	trxs, err := eth.GetTrxs(eth.tickAddress, eth.handleHeight, 99999999)
	if err != nil {
		fmt.Printf("Get trxs on tick err: %v\n", err)
		return
	}

	for _, trx := range trxs {
		if trx.IsIrrevisible {
			// fmt.Printf("trx is irreversible: %v\n", trx.TransactionID)
			//exec push action

			if eth.handleHeight < trx.BlockNum {
				eth.handleHeight = trx.BlockNum
				fmt.Printf("trx irreversible from: %v  %v\n", trx.TransactionID, eth.handleHeight)
			}

			var charge chainlib.ChargeInfo
			charge.Hash = trx.TransactionID
			charge.From = trx.From
			charge.To = trx.To
			charge.BlockNum = trx.BlockNum
			charge.Quantity = strconv.FormatFloat(trx.Amount, 'f', 4, 64)
			charge.Category = trx.Category
			charge.Memo = trx.Memo

			// 需要钱包密码
			_, err := chainlib.ClWalletUnlock("PW5JHPpaGrS7bKhmQJ5Rb7rNSXhp3S3sXN2fGWaqQNzQufQaWrkUJ")
			// 需要合约权限
			chargeID, err := chainlib.ClPushCharge("datxio.charg", "charge", charge)
			if err != nil {
				fmt.Printf("ETH push charge err:", err)
				continue
			}
			blockNum, err := chainlib.ClGetTrxBlockNum(chargeID)
			if err != nil {
				fmt.Printf("blockNum parse err:", err)
				continue
			}

			trans := trx
			trans.TransactionID = chargeID
			trans.BlockNum = blockNum
			trans.To = "lmx"    // 需要地址映射和权限
			eth.tick.AddHash(trans)
		} else {
			jobid := trx.Category + "_" + trx.TransactionID
			if job, _ := delayqueue.Get(jobid); job != nil {
				fmt.Printf("tick trx is existed: %v %v\n", trx.TransactionID, time.Now().Unix())
				continue
			}

			fmt.Printf("Add eth task on tick: %v  %v\n", trx.TransactionID, time.Now().Unix())
			eth.tick.AddTask(trx, ETHDelaySeconds)
		}
	}
}

func (eth *ETHBrowser) ReTry(trx chainlib.Transaction) bool {
	blockNum := trx.BlockNum

	fmt.Printf("ReTry eth on tick: %v  %v\n", trx.TransactionID, time.Now().Unix())

	sta, err := eth.Irreversible(blockNum)
	if err != nil || !sta {
		eth.tick.AddTask(trx, ETHDelaySeconds/2)
		return false
	}

	return sta
}

func (eth *ETHBrowser) Close() {
	eth.close <- true
}
