package explorer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// TrxInfo ...
type TrxInfo struct {
	Token    string `json:"token"` // DATX，DBTC，DETH，DEOS，BTC，ETH，EOS
	Trxid    string `json:"trxid"`
	Quantity string `json:"quantity"` // "2 DATX"
	Blocknum int64  `json:"blocknum"`
	From     string `json:"from"`
	To       string `json:"to"`
	Time     string `json:"time"`  // "2018-9-4 15:30:56"
	Value    string `json:"value"` // "$12.50"
}

// DATXActData ...
type DATXActData struct {
	Token string `json:"account"` // "datxio.token"
	Name  string `json:"name"`    // "transfer"
	Data  struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Quantity string `json:"quantity"` // "2.0000 DATX"
		Memo     string `json:"memo"`
	} `json:"data"`
}

// DATXActions ...
type DATXActions struct {
	Actions []struct {
		Blocknum    int64  `json:"block_num"`  // 19240
		Time        string `json:"block_time"` // "2018-07-22T02:55:42.000"
		ActionTrace struct {
			Act   DATXActData `json:"act"`
			Trxid string      `json:"trx_id"`
		} `json:"action_trace"`
	} `json:"actions"`
}

// DATXGetActionsFormData ...
type DATXGetActionsFormData struct {
	Account string `json:"account_name"`
	Pos     int64  `json:"pos"`
	Offset  int64  `json:"offset"`
}

// DATXTrxList ... DATX，DBTC，DETH，DEOS
func DATXTrxList(account string, pos int64, offset int64) ([]*TrxInfo, error) {
	formData := DATXGetActionsFormData{account, pos, offset}
	bytesData, err := json.Marshal(formData)
	if err != nil {
		return nil, fmt.Errorf("datx get_actions parameter error %v", formData)
	}

	URL := "http://127.0.0.1:8888/v1/history/get_actions"
	request, err := http.NewRequest("POST", URL, bytes.NewReader(bytesData))
	if err != nil {
		return nil, errors.New("datx get_actions request error")
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("datx get_actions resp error %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("datx get_actions body error %v", resp.Body)
	}

	actions := &DATXActions{}
	if err := json.Unmarshal([]byte(body), &actions); err != nil {
		return nil, fmt.Errorf("datx get_actions unmarshal body error %v", body)
	}

	trxList := make([]*TrxInfo, 0)
	for _, act := range actions.Actions {
		var trx TrxInfo
		switch act.ActionTrace.Act.Token {
		case "datxio": // 不需要展现
			continue
		case "datxio.token":
			trx.Token = "DATX"
		case "datxio.dbtc":
			trx.Token = "DBTC"
		case "datxio.deth":
			trx.Token = "DETH"
		case "datxio.deos":
			trx.Token = "DEOS"
		default:
			trx.Token = act.ActionTrace.Act.Token
		}
		trx.Trxid = act.ActionTrace.Trxid
		trx.Quantity = act.ActionTrace.Act.Data.Quantity
		trx.Blocknum = act.Blocknum
		trx.From = act.ActionTrace.Act.Data.From
		trx.To = act.ActionTrace.Act.Data.To
		trx.Time = act.Time
		// DATX的Quantity是金额+空格+币种
		amount, _ := strconv.ParseFloat(trx.Quantity[:strings.Index(trx.Quantity, " ")], 64)
		price, _ := GetTokenPrice(trx.Token)
		trx.Value = "$" + strconv.FormatFloat(amount*price, 'f', 2, 64)
		trxList = append(trxList, &trx)
	}
	return trxList, nil
}

// ETHTrx ...
type ETHTrx struct {
	Trxid    string `json:"hash"`
	Quantity string `json:"value"`
	Blocknum string `json:"blockNumber"`
	From     string `json:"from"`
	To       string `json:"to"`
	Time     string `json:"timeStamp"`
}

// ETHTrxs ...
type ETHTrxs struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Trxs    []ETHTrx `json:"result"`
}

// ETHTrxList ...
func ETHTrxList(account string) ([]*TrxInfo, error) {
	// http://api.etherscan.io/api?module=account&action=txlist&address=0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a&startblock=0&endblock=99999999&sort=asc&apikey=YourApiKeyToken
	request, err := url.Parse("http://api.etherscan.io/api")
	if err != nil {
		return nil, errors.New("eth url parse error")
	}

	params := url.Values{}
	params.Set("module", "account")
	params.Set("action", "txlist")
	params.Set("address", account)
	params.Set("startblock", "0")
	params.Set("endblock", "99999999")
	params.Set("sort", "asc")
	params.Set("apikey", "8FT3VZVAS94PIHPYKEWPWHC4ZICB71RFSM")
	request.RawQuery = params.Encode()

	resp, err := http.Get(request.String())
	if err != nil {
		return nil, fmt.Errorf("eth txlist resp error %v", err)
	}
	if resp.Status != "200 OK" {
		return nil, fmt.Errorf("eth txlist account\t%v %v", account, resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("eth txlist body error %v", resp.Body)
	}

	var ethTrxs ETHTrxs
	if err := json.Unmarshal([]byte(body), &ethTrxs); err != nil {
		return nil, fmt.Errorf("eth txlist unmarshal body error %v", string(body))
	}
	if ethTrxs.Status != "1" {
		return nil, fmt.Errorf("eth txlist account\t%v %v", account, ethTrxs.Message)
	}
	trxList := make([]*TrxInfo, 0)
	for _, act := range ethTrxs.Trxs {
		var trx TrxInfo
		trx.Token = "ETH"
		trx.Trxid = act.Trxid
		amount, _ := strconv.ParseFloat(act.Quantity, 64)
		amount /= 1000000000000000000
		trx.Quantity = strconv.FormatFloat(amount, 'f', 4, 64) + " ETH"
		blocknum, _ := strconv.ParseInt(act.Blocknum, 10, 64)
		trx.Blocknum = blocknum
		trx.From = act.From
		trx.To = act.To
		timeint64, _ := strconv.ParseInt(act.Time, 10, 64)
		trx.Time = time.Unix(timeint64, 0).Format("2006-01-02 15:04:05")
		price, _ := GetTokenPrice(trx.Token)
		trx.Value = "$" + strconv.FormatFloat(amount*price, 'f', 2, 64)
		trxList = append(trxList, &trx)
	}
	return trxList, nil
}
