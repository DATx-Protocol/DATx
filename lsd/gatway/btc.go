package gatway

import (
	"crypto/tls"
	"datx/lsd/chainlib"
	"datx/lsd/delayqueue"
	"datx/lsd/server"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	BTCRetrySeconds    int64 = 1
	BTCDelaySeconds    int64 = 60
	BTCIrreversibleCnt int64 = 6
)

type BTCTrxIn struct {
	Sequence int64  `json:"sequence"`
	Witness  string `json:"witness"`
	PrevOut  struct {
		Spent   bool   `json:"spent"`
		TxIndex int    `json:"tx_index"`
		Type    int    `json:"type"`
		Addr    string `json:"addr"`
		Value   int    `json:"value"`
		N       int    `json:"n"`
		Script  string `json:"script"`
	} `json:"prev_out"`
	Script string `json:"script"`
}

type BTCTrxOut struct {
	AddrTagLink string `json:"addr_tag_link,omitempty"`
	AddrTag     string `json:"addr_tag,omitempty"`
	Spent       bool   `json:"spent"`
	TxIndex     int    `json:"tx_index"`
	Type        int    `json:"type"`
	Addr        string `json:"addr"`
	Value       int    `json:"value"`
	N           int    `json:"n"`
	Script      string `json:"script"`
}

type BTCTransaction struct {
	Ver         int         `json:"ver"`
	Inputs      []BTCTrxIn  `json:"inputs"`
	Weight      int         `json:"weight"`
	BlockHeight int         `json:"block_height"`
	RelayedBy   string      `json:"relayed_by"`
	Out         []BTCTrxOut `json:"out"`
	LockTime    int         `json:"lock_time"`
	Result      int         `json:"result"`
	Size        int         `json:"size"`
	Time        int         `json:"time"`
	TxIndex     int         `json:"tx_index"`
	VinSz       int         `json:"vin_sz"`
	Hash        string      `json:"hash"`
	VoutSz      int         `json:"vout_sz"`
}

type BTCResponse struct {
	Hash160       string           `json:"hash160"`
	Address       string           `json:"address"`
	NTx           int              `json:"n_tx"`
	TotalReceived int64            `json:"total_received"`
	TotalSent     int64            `json:"total_sent"`
	FinalBalance  int              `json:"final_balance"`
	Txs           []BTCTransaction `json:"txs"`
}

type BTCLatestBlockNUm struct {
	Hash       string `json:"hash"`
	Time       int    `json:"time"`
	BlockIndex int    `json:"block_index"`
	Height     int    `json:"height"`
	TxIndexes  []int  `json:"txIndexes"`
}

type BTCBrowser struct {
	url    string
	apikey string

	tickAccount string

	handleHeight int64

	close chan bool

	tick *server.ChainServer
}

func NewBTCBrowser(accountAddr string, server *server.ChainServer) *BTCBrowser {
	return &BTCBrowser{
		url:          "https://testnet.blockchain.info",
		apikey:       "",
		tickAccount:  accountAddr,
		handleHeight: 0,
		close:        make(chan bool),
		tick:         server,
	}
}

//GetTrxs https://blockchain.info/rawaddr/$bitcoin_address
func (btc *BTCBrowser) GetTrxs(addr string) ([]chainlib.Transaction, error) {
	requrl := btc.url + "/rawaddr/" + addr
	// requrl := "https://testnet.blockchain.info/rawtx/6c711db1296824782a1206bdc27034c23710a1b4fc3c504b88f0e49583ae29cb"

	// log.Printf("BTC GetTrxs request url is : %s\n", requrl)

	res, err := http.Get(requrl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("BTC GetTrxs:%v %v\n", requrl, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp BTCResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	latestBlockNum, err := btc.GetLatestBlockNum()
	if err != nil {
		return nil, err
	}

	result := make([]chainlib.Transaction, 0)
	for _, v := range resp.Txs {

		if int64(v.BlockHeight) <= btc.handleHeight {
			continue
		}

		for _, out := range v.Out {
			if out.Addr != "" && out.Addr != addr {
				continue
			}

			var temp chainlib.Transaction
			temp.Category = "BTC"
			temp.BlockNum = int64(v.BlockHeight)

			infrom := make([]string, 0, len(v.Inputs))
			for _, in := range v.Inputs {
				infrom = append(infrom, in.PrevOut.Addr)
			}
			temp.From = strings.Join(infrom, ",")
			temp.To = out.Addr
			temp.Amount = float64(out.Value) / 100000000 //	监听单位是聪，转为比特币
			temp.TransactionID = v.Hash
			temp.IsIrrevisible = false
			temp.Time = time.Unix(int64(v.Time), 0)
			temp.Memo = out.Script

			if (latestBlockNum - temp.BlockNum) > BTCIrreversibleCnt {
				temp.IsIrrevisible = true
			}

			result = append(result, temp)
		}
	}

	return result, nil
}

//GetLatestBlockNum https://blockchain.info/latestblock
func (btc *BTCBrowser) GetLatestBlockNum() (int64, error) {
	requrl := btc.url + "/latestblock"

	// log.Printf("BTC GetLatestBlockNum request url is : %s\n", requrl)

	res, err := http.Get(requrl)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.Status != "200 OK" {
		return 0, fmt.Errorf("BTC GetLatestBlockNum:%v %v\n", requrl, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var resp BTCLatestBlockNUm
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}

	result := int64(resp.Height)
	return result, nil
}

func (btc *BTCBrowser) Irreversible(blocknum int64) (bool, error) {
	latest, err := btc.GetLatestBlockNum()
	if err != nil {
		return false, err
	}

	sub := latest - blocknum
	if sub > BTCIrreversibleCnt {
		return true, nil
	}

	return false, fmt.Errorf("%d not irreversible,need %d confirmed.", blocknum, sub)
}

func (btc *BTCBrowser) SetTickAccountAddr(account string) {
	btc.tickAccount = account
}

func (btc *BTCBrowser) Tick() {
	trxs, err := btc.GetTrxs(btc.tickAccount)
	if err != nil {
		log.Printf("Get trxs err: %v\n", err)
		return
	}

	for _, trx := range trxs {
		if trx.IsIrrevisible {
			// log.Printf("BTC trx is irreversible on Tick(): %v\n", trx.TransactionID)
			//exec push action

			if btc.handleHeight < trx.BlockNum {
				btc.handleHeight = trx.BlockNum
			}

			var result error
			if trx.To != "" {
				result = chainlib.PushCharge(trx)
			} else {
				result = btc.pushExtract(trx)
			}

			log.Printf("[Tick] BTC push action result:%v\n", result)
		} else {
			jobid := trx.Category + "_" + trx.TransactionID
			if job, _ := delayqueue.Get(jobid); job != nil {
				continue
			}

			log.Printf("BTC add btc to delay task on Tick(): %v  %v\n", trx.TransactionID, time.Now().Unix())
			btc.tick.AddTask(trx, BTCDelaySeconds)
		}
	}
}

//ReTry ...
func (btc *BTCBrowser) ReTry(trx chainlib.Transaction) bool {
	log.Printf("BTC trx on ReTry(): %v\n", trx.TransactionID)
	blockNum := trx.BlockNum

	sta, err := btc.Irreversible(blockNum)
	if err != nil || !sta {
		btc.tick.AddTask(trx, BTCDelaySeconds)
		return false
	}

	if trx.To != "" {
		chainlib.PushCharge(trx)
	} else {
		btc.pushExtract(trx)
	}

	return sta
}

//Close ...
func (btc *BTCBrowser) Close() {
	btc.close <- true
}

func (btc *BTCBrowser) pushExtract(trx chainlib.Transaction) error {
	//get trxID by memo
	//eg: https://127.0.0.1:8080/btc/decodeMemo?script=6a0d626974636f696e6a732d6c6962
	url := fmt.Sprintf("https://127.0.0.1:8080/btc/decodeMemo?script=%s", trx.Memo)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("get btc url : %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("BTC decodeMemo Response error: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	trxid := string(body)
	log.Printf("get trx id from script:%v\n", trxid)

	trx.Memo = trxid
	return chainlib.PushExtract(trx)
}
