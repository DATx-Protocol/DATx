package chainlib

import (
	"bytes"
	"crypto/sha256"
	"datx/lsd/common"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

type CommonFunc interface {
	SetTickAccountAddr(account string)

	Tick()

	ReTry(trx Transaction) bool

	Close()
}

//TransactionInfo ...
type TransactionInfo struct {
	ID  string `json:"id"`
	Trx struct {
		Trx struct {
			Actions []ActionInfo `json:"actions"`
		} `json:"trx"`
	} `json:"trx"`
}

//ActionInfo ...
type ActionInfo struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	Data    struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Quantity string `json:"quantity"`
		Category string `json:"category"`
		Memo     string `json:"memo"`
	} `json:"data"`
}

//TransferInfo ...
type TransferInfo struct {
	Hash     string `json:"hash"`
	From     string `form:"from" json:"from" binding:"required"`
	To       string `form:"to" json:"to" bdinding:"required"`
	Quantity string `form:"quantity" json:"quantity" bdinding:"required"`
	Memo     string `form:"memo" json:"memo"`
}

//ChargeInfo ...
type ChargeInfo struct {
	BPName   string `json:"bpname"`
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	BlockNum int64  `json:"blocknum"`
	Quantity string `json:"quantity"`
	Category string `json:"category"`
	Memo     string `json:"memo"`
}

//ExtractInfo info
type ExtractInfo struct {
	TrxID    string `json:"trxid"`
	Producer string `json:"producer"`
}

// RandomSha256 ...
func RandomSha256() string {
	randStr := strconv.FormatUint(rand.Uint64(), 10)
	hash := sha256.New()
	hash.Write([]byte(randStr))
	return hex.EncodeToString(hash.Sum(nil))
}

// ParseTrxID ...
func ParseTrxID(inStr string) (string, error) {
	js, err := simplejson.NewJson([]byte(inStr))
	if err != nil {
		return "", fmt.Errorf("simplejson error: %v\n", inStr)
	}
	return js.Get("transaction_id").MustString(), nil
}

// ExecShell ...
func ExecShell(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr := string(stdout.Bytes())
	errStr := string(stderr.Bytes())
	if err != nil {
		return errStr, fmt.Errorf("Execute command error: %v\n", errStr)
	}
	return outStr, nil
}

// ClWalletUnlock ...
func ClWalletUnlock(password string) (string, error) {
	command := "cldatx wallet unlock " + " --password " + password
	return ExecShell(command)
}

// ClPushAction ...
func ClPushAction(account, action, data, permission string) (string, error) {
	actionstr := fmt.Sprintf("cldatx push action %s %s '%s' -j -f -p %s", account, action, data, permission)
	log.Printf("\n****************************************\n%s\n****************************************\n", actionstr)
	result, err := ExecShell(actionstr)
	if err != nil {
		if strings.Contains(err.Error(), "Locked wallet") {
			//Ensure that your wallet is unlocked before using it!
			wname, wpassword := common.GetWalletNameAndPassword()
			if len(wpassword) == 0 {
				log.Print("Push Action before setup your wallet name and password in your ~/datxos-wallet/wallet_password.ini.\n")
				return "", fmt.Errorf("~/datxos-wallet/wallet_password.ini not found")
			}
			unlockwallet := fmt.Sprintf("cldatx wallet unlock -n %s --password %s", wname, wpassword)
			_, err := ExecShell(unlockwallet)
			if err != nil {
				log.Printf("Push Action unlock wallet: %v\n", err)
				return "", nil
			}

			result, err = ExecShell(actionstr)
		}
	}

	return result, err
}

// ClPushTransfer ...
// string	返回TrxID
func ClPushTransfer(account string, action string, trans TransferInfo) (string, error) {
	//Ensure that your wallet is unlocked before using it!
	js, _ := json.Marshal(trans)
	transStr := string(js)
	outStr, err := ClPushAction(account, action, transStr, account)
	if err != nil {
		return "", err
	}
	TrxID, err := ParseTrxID(outStr)
	if err != nil {
		return "", err
	}
	return TrxID, nil
}

// ClPushCharge ...
func ClPushCharge(account string, action string, charge ChargeInfo) (string, error) {
	//Ensure that your wallet is unlocked before using it!
	js, _ := json.Marshal(charge)
	chargeStr := string(js)
	outStr, err := ClPushAction(account, action, chargeStr, charge.BPName)
	if err != nil {
		return "", err
	}
	TrxID, err := ParseTrxID(outStr)
	if err != nil {
		return "", err
	}
	return TrxID, nil
}

// ClGetTrxBlockNum ...
func ClGetTrxBlockNum(trxID string) (int64, error) {
	command := "cldatx get transaction " + trxID
	outStr, err := ExecShell(command)
	if err != nil {
		return 0, err
	}
	js, err := simplejson.NewJson([]byte(outStr))
	if err != nil {
		return 0, fmt.Errorf("simplejson error: %v\n", outStr)
	}
	return js.Get("block_num").MustInt64(), nil
}

// ClGetTrxInfo ...
func ClGetTrxInfo(trxID string) (*TransactionInfo, error) {
	command := "cldatx get transaction " + trxID
	outStr, err := ExecShell(command)
	if err != nil {
		return nil, err
	}
	var trx TransactionInfo
	if err := json.Unmarshal([]byte(outStr), &trx); err != nil {
		return nil, fmt.Errorf("ClGetTrxInfo Unmarshal error: %v\n", outStr)
	}
	return &trx, nil
}

// ClGetLIBNum ...
func ClGetLIBNum() (int64, error) {
	command := "cldatx get info"
	outStr, err := ExecShell(command)
	if err != nil {
		return 0, err
	}
	js, err := simplejson.NewJson([]byte(outStr))
	if err != nil {
		return 0, fmt.Errorf("simplejson error: %v\n", outStr)
	}
	return js.Get("last_irreversible_block_num").MustInt64(), nil
}

//PushCharge push charge action to blockchain
func PushCharge(trx Transaction) error {
	var charge ChargeInfo
	charge.BPName = common.GetCfgProducerName()
	charge.Hash = trx.TransactionID
	charge.From = trx.From
	charge.To = trx.To
	charge.BlockNum = trx.BlockNum
	amount := int64(10000 * trx.Amount)
	if amount == 0 {
		return fmt.Errorf("PushCharge trx amount is zero")
	}
	charge.Quantity = strconv.FormatInt(amount, 10)
	charge.Category = trx.Category
	charge.Memo = trx.Memo

	_, err := ClPushCharge("datxos.charg", "charge", charge)
	if err != nil {
		log.Printf("[PushCharge] failed: %v\n %v\n:", trx, err)
		return err
	}

	return nil
}

//PushExtract push extract action to blockchain
func PushExtract(trx Transaction) error {
	var extract ExtractInfo
	extract.TrxID = trx.Memo
	extract.Producer = common.GetCfgProducerName()

	bytes, err := json.Marshal(extract)
	if err != nil {
		log.Printf("PushExtract marshal failed:%v %v\n", trx, err)
		return err
	}
	_, err = ClPushAction("datxos.extra", "setsuccess", string(bytes), extract.Producer)
	if err != nil {
		log.Printf("[PushExtract] failed:  %v\n %v\n", trx, err)
		return err
	}

	return nil
}
