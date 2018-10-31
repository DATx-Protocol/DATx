package explorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

// GetDATXTrxBlockNum ...
func GetDATXTrxBlockNum(trxID string) (int64, error) {
	formData := make(map[string]interface{})
	formData["id"] = trxID
	bytesData, err := json.Marshal(formData)
	if err != nil {
		return 0, fmt.Errorf("datx get_transaction parameter error %v", formData)
	}

	URL := WalletConfig.DatxIP + "/v1/history/get_transaction"
	request, err := http.NewRequest("POST", URL, bytes.NewReader(bytesData))
	if err != nil {
		return 0, fmt.Errorf("datx get_transaction request error %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("datx get_transaction resp error %v", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("datx get_transaction not found %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("datx get_transaction body error %v", err)
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return 0, fmt.Errorf("datx get_transaction simplejson error %v", string(body))
	}
	return js.Get("block_num").MustInt64(), nil
}

// GetDATXLIBNum ...
func GetDATXLIBNum() (int64, error) {
	URL := WalletConfig.DatxIP + "/v1/chain/get_info"
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer([]byte("")))
	if err != nil {
		return 0, fmt.Errorf("datx get_info request error %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("datx get_info resp error %v", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("datx get_info not found %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("datx get_info body error %v", err)
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return 0, fmt.Errorf("datx get_info simplejson error %v", string(body))
	}
	return js.Get("last_irreversible_block_num").MustInt64(), nil
}

// WaitIrreversible ...
func WaitIrreversible(trxID string) error {
	period := 0
	for {
		blockNum, err := GetDATXTrxBlockNum(trxID)
		if err != nil { // 报错退出循环
			return err
		}
		libNum, err := GetDATXLIBNum()
		if err != nil { // 报错退出循环
			return err
		}
		if libNum > blockNum { // 达到不可逆
			log.Printf("libNum : %d\tblockNum : %d\n", libNum, blockNum)
			return nil
		}
		if period > 60 { // 无法达到不可逆
			log.Printf("libNum : %d\tblockNum : %d\n", libNum, blockNum)
			return fmt.Errorf("datx transaction %v cannot get irreversible", trxID)
		}
		time.Sleep(time.Duration(1) * time.Second) // 否则等待1秒
		period++
	}
}
