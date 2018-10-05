package server

import (
	"datx/ListenServer/chainlib"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//btc multisig
func BTCMultiSig(trx chainlib.Transaction) (string, error) {
	url := "https://localhost:8080/btc/withdraw?isTestnet=1&to=" + trx.To + "&value=" + fmt.Sprintf("%f", trx.Amount) + "&fee=100000&trxid=" + trx.TransactionID

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("BTCMultiSig Response error: %v\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	trxid := string(body)
	log.Printf("BTCMultiSig return trx id:%v\n", trxid)
	return trxid, nil
}

//eth multisig
func ETHMultiSig(trx chainlib.Transaction) (string, error) {
	url := "https://localhost:8080/eth/withdraw?to=" + trx.To + "&value=" + fmt.Sprintf("%f", trx.Amount) + "&data=" + trx.Memo + "&trxid=" + trx.TransactionID

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ETHMultiSig Response error: %v\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	trxid := string(body)
	log.Printf("ETHMultiSig return trx id:%v\n", trxid)
	return trxid, nil
}

//eos multisig
func EOSMultiSig(trx chainlib.Transaction) (string, error) {
	url := "https://localhost:8080/eos/withdraw?to=" + trx.To + "&value=" + fmt.Sprintf("%f EOS", trx.Amount) + "&trxid=" + trx.TransactionID

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("EOSMultiSig Response error: %v\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	trxid := string(body)
	log.Printf("EOSMultiSig return trx id:%v\n", trxid)
	return trxid, nil
}
