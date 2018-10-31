package explorer

import (
	"bytes"
	"datx/lsd/chainlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

// AddressMapInfo ...
type AddressMapInfo struct {
	DatxAddress string `json:"datxaddress"`
	Address     string `json:"address"`
	BPName      string `json:"bpname"`
}

// AddressMapRequest ...
type AddressMapRequest struct {
	DatxAddress string `json:"datxaddress"`
	Address     string `json:"address"`
}

// GetCurrentBP ...
func GetCurrentBP() (string, error) {
	URL := WalletConfig.DatxIP + "/v1/chain/get_info"
	request, err := http.NewRequest("POST", URL, bytes.NewBuffer([]byte("")))
	if err != nil {
		return "", fmt.Errorf("datx get_info request error %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("datx get_info resp error %v", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return "", fmt.Errorf("datx get_info not found %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("datx get_info body error %v", err)
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		return "", fmt.Errorf("datx get_info simplejson error %v", err)
	}
	return js.Get("head_block_producer").MustString(), nil
}

// ClRecordUser ...
func ClRecordUser(datxaddress string, address string) (string, error) {
	var addressMap AddressMapInfo
	addressMap.DatxAddress = datxaddress
	addressMap.Address = address
	addressMap.BPName = WalletConfig.Name
	js, _ := json.Marshal(addressMap)
	addressStr := "'" + string(js) + "'"
	command := "cldatx -u " + WalletConfig.DatxIP + " push action datxos.charg recorduser " +
		addressStr + " -j " + " -f " + " -p " + WalletConfig.Name
	log.Println(command)
	chainlib.ClWalletUnlock(WalletConfig.PassWord)
	return chainlib.ExecShell(command)
}
