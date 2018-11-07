package server

import (
	"crypto/tls"
	"datx/lsd/chainlib"
	"datx/lsd/common"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//btc multisig
func BTCMultiSig(trx chainlib.Transaction) (string, error) {
	url := fmt.Sprintf("%s/btc/withdraw?isTestnet=1&to=%s&value=%s&fee=2000&trxid=%s&nodeName=%s", common.GetCrossChainEnds(), trx.Memo, fmt.Sprintf("%.4f", trx.Amount), trx.TransactionID, common.GetCfgProducerName())
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
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
	url := fmt.Sprintf("%s/eth/withdraw?to=%s&value=%s&data=%s&trxid=%s&nodeName=%s", common.GetCrossChainEnds(), trx.To, fmt.Sprintf("%.4f", trx.Amount), trx.Memo, trx.TransactionID, common.GetCfgProducerName())
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
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
	url := fmt.Sprintf("%s/eos/withdraw?to=%s&value=%s&trxid=%s&nodeName=%s", common.GetCrossChainEnds(), trx.To, fmt.Sprintf("%.4f", trx.Amount), trx.TransactionID, common.GetCfgProducerName())
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
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
