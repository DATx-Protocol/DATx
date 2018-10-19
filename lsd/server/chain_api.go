package server

import (
	"bytes"
	"datx/lsd/chainlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//HTTPQuery ...
func HTTPQuery(command string) (string, error) {
	url := "http://127.0.0.1:8888/v1/chain/" + command
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "request error", err
	}
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return "response error", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "body error", err
	}
	return string(body), nil
}

func GetOuterTrxTable(code, scope, table string) ([]byte, error) {
	reqpara := TableParams{
		Code:  code,
		Scope: scope,
		Table: table,
		JSON:  "true",
		// Lower: 1,
		// Upper: -1,
		// Limit: 10,
	}
	req, err := json.Marshal(reqpara)
	if err != nil {
		return nil, err
	}

	para := bytes.NewBuffer([]byte(req))
	request, err := http.NewRequest("POST", "http://127.0.0.1:8888/v1/chain/get_table_rows", para)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, errr := client.Do(request)
	if errr != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Response error: %v\n", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetProducerSchedule() ([]Producers, error) {
	reqpara := ProducerScheduleParams{
		Limit:      21,
		LowerBound: 0,
		JSON:       "true",
	}
	req, err := json.Marshal(reqpara)
	if err != nil {
		return nil, err
	}

	para := bytes.NewBuffer([]byte(req))
	request, err := http.NewRequest("POST", "http://127.0.0.1:8888/v1/chain/get_producer_schedule", para)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, errr := client.Do(request)
	if errr != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Response error: %v\n", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	var resut ProducerSchedule
	if err := json.Unmarshal(body, &resut); err != nil {
		return nil, err
	}

	fmt.Printf("data:%v\n", resut.Active.Producers)
	return resut.Active.Producers, nil
}

func GetInfo() (*ChainInfo, error) {
	para := bytes.NewBuffer([]byte(""))
	request, err := http.NewRequest("POST", "http://127.0.0.1:8888/v1/chain/get_info", para)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, errr := client.Do(request)
	if errr != nil {
		return nil, errr
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Response error: %v\n", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var resut ChainInfo
	if err := json.Unmarshal(body, &resut); err != nil {
		return nil, err
	}

	blockTime, err := time.Parse("2006-01-02T15:04:05", resut.HeadBlockTime)
	if err != nil {
		fmt.Printf("Parase time err: %v\n", err)
		return nil, err
	}
	fmt.Printf("GetInfo:%v %v %v\n", resut.HeadBlockProducer, blockTime.Unix(), time.Now().Unix())
	return &resut, nil
}

func IsCurrentProducer() bool {
	localProducerName := chainlib.GetCfgProducerName()

	info, err := GetInfo()
	if err != nil {
		fmt.Printf("Get info err:%v\n", err)
		return false
	}

	headProducerName := info.HeadBlockProducer

	// commondstr := fmt.Sprint("cldatx system listproducers -l 21")
	// _, err := chainlib.ExecShell(commondstr)
	// if err != nil {
	// 	return false
	// }

	if localProducerName == headProducerName {
		return true
	}

	return false
}

//GetExtractActions get transaction by escrow account(*dbtc,deth,deos) when extracting
func GetExtractActions(addr string, pos, offset int32) ([]chainlib.Transaction, error) {
	reqpara := ActionsParams{
		AccountName: addr,
		Pos:         pos,
		Offset:      offset,
	}
	req, err := json.Marshal(reqpara)
	if err != nil {
		return nil, err
	}

	para := bytes.NewBuffer([]byte(req))
	request, err := http.NewRequest("POST", "http://127.0.0.1:8888/v1/history/get_actions", para)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, errr := client.Do(request)
	if errr != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("ChainAPI Response error: %v\n", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var resp ExtractActions
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	lastIrreversibleBlock := int64(resp.LastIrreversibleBlock)
	result := make([]chainlib.Transaction, 0)
	for _, v := range resp.Actions {
		if v.ActionTrace.Act.Name != "transfer" {
			continue
		}

		if v.ActionTrace.Act.Data.To != addr {
			continue
		}

		var temp chainlib.Transaction
		temp.TransactionID = v.ActionTrace.TrxID

		temp.BlockNum = int64(v.BlockNum)
		temp.From = v.ActionTrace.Act.Data.From
		temp.To = v.ActionTrace.Act.Data.To
		amountpos := strings.Index(v.ActionTrace.Act.Data.Quantity, " ") // EOS的Quantity是金额+空格+币种
		amountstr := v.ActionTrace.Act.Data.Quantity[:amountpos]         // 只需要空格前面的金额
		temp.Category = v.ActionTrace.Act.Data.Quantity[amountpos+1:]
		temp.Amount, err = strconv.ParseFloat(amountstr, 64)
		temp.Memo = v.ActionTrace.Act.Data.Memo
		if err != nil {
			log.Printf("Extract parse amount err:%v\n", err)
			return nil, err
		}

		temp.Time, err = time.Parse("2006-01-02T15:04:05", v.BlockTime)
		if err != nil {
			log.Printf("Extract parse time err:%v\n", err)
			return nil, err
		}
		temp.IsIrrevisible = false
		if lastIrreversibleBlock >= temp.BlockNum {
			temp.IsIrrevisible = true
		}

		result = append(result, temp)
	}

	return result, nil
}

func GetExtractTransaction(trxid string) (*chainlib.Transaction, error) {
	reqpara := struct {
		ID string
	}{
		ID: trxid,
	}
	req, err := json.Marshal(reqpara)
	if err != nil {
		return nil, err
	}

	para := bytes.NewBuffer([]byte(req))
	request, err := http.NewRequest("POST", "http://127.0.0.1:8888/v1/history/get_transaction", para)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, errr := client.Do(request)
	if errr != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("get_transaction Response error: %v\n", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result chainlib.Transaction
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CheckIrreversible(trx chainlib.Transaction) bool {
	blocknum := trx.BlockNum

	info, err := GetInfo()
	if err != nil {
		return false
	}

	lastIrreversiblrBlockNum := int64(info.LastIrreversibleBlockNum)
	if blocknum <= lastIrreversiblrBlockNum {
		return true
	}

	return false
}
