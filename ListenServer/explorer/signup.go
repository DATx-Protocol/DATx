package explorer

import (
	"bytes"
	"datx/ListenServer/chainlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// SignupTrxInfo ...
type SignupTrxInfo struct {
	Trxid    string `json:"trxid"`
	Quantity string `json:"quantity"`
	Blocknum int64  `json:"blocknum"`
	From     string `json:"from"`
	To       string `json:"to"`
	Time     string `json:"time"`
	Memo     string `json:"memo"`
}

// SignupTrxRequest ...
type SignupTrxRequest struct {
	Account string `form:"account" json:"account" binding:"required"`
	Limit   int64  `form:"limit" json:"limit" binding:"required"`
}

// SignupAccountRequest ...
type SignupAccountRequest struct {
	SysAccount string `form:"to" json:"to" binding:"required"`
	NewAccount string `form:"account" json:"account" binding:"required"`
	Quantity   string `form:"quantity" json:"quantity" binding:"required"`
	Memo       string `form:"memo" json:"memo" binding:"required"`
}

// SignupAccountInfo ...
type SignupAccountInfo struct {
	SysAccount string `json:"sys_account"`
	NewAccount string `json:"new_account"`
	PublicKey  string `json:"public_key"`
	Quantity   string `json:"quantity"`
}

// GetSignupTrxList ... DATX
func GetSignupTrxList(account string, pos int64, offset int64) ([]*SignupTrxInfo, error) {
	formData := DATXGetActionsFormData{account, pos, offset}
	bytesData, err := json.Marshal(formData)
	if err != nil {
		return nil, fmt.Errorf("datx get_actions parameter error %v", formData)
	}

	URL := WalletConfig.DatxIP + "/v1/history/get_actions"
	request, err := http.NewRequest("POST", URL, bytes.NewReader(bytesData))
	if err != nil {
		return nil, fmt.Errorf("datx get_actions request error %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("datx get_actions resp error %v", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, fmt.Errorf("datx get_actions not found %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("datx get_actions body error %v", err)
	}

	actions := &DATXActions{}
	if err := json.Unmarshal([]byte(body), &actions); err != nil {
		return nil, fmt.Errorf("datx get_actions unmarshal error %v", string(body))
	}

	trxList := make([]*SignupTrxInfo, 0)
	for _, act := range actions.Actions {
		if act.ActionTrace.Act.Token != "datxio.token" {
			continue // 注册转账都调datxio.token合约的transfer
		}
		if act.ActionTrace.Act.Data.To != account {
			continue // 收款账号一致
		}
		var trx SignupTrxInfo
		trx.Trxid = act.ActionTrace.Trxid
		trx.Quantity = act.ActionTrace.Act.Data.Quantity
		trx.Blocknum = act.Blocknum
		trx.From = act.ActionTrace.Act.Data.From
		trx.To = act.ActionTrace.Act.Data.To
		trx.Time = act.Time
		trx.Memo = act.ActionTrace.Act.Data.Memo

		trxList = append(trxList, &trx)
	}
	return trxList, nil
}

// GetSignupAccountList ...
func GetSignupAccountList(account string, trxList []*SignupTrxInfo) ([]*SignupAccountInfo, error) {
	signList := make([]*SignupAccountInfo, 0)
	for _, trx := range trxList {
		separatorPos := strings.Index(trx.Memo, " ")
		if separatorPos <= 0 {
			separatorPos = strings.Index(trx.Memo, "-")
		}
		if separatorPos <= 0 {
			continue // 非法Memo值
		}
		var sign SignupAccountInfo
		sign.SysAccount = account
		sign.NewAccount = trx.Memo[:separatorPos]
		sign.PublicKey = trx.Memo[separatorPos+1 : len(trx.Memo)]
		sign.Quantity = trx.Quantity

		signList = append(signList, &sign)
	}
	return signList, nil
}

// MatchSignupAccount ...
func MatchSignupAccount(request SignupAccountRequest, trxList []*SignupTrxInfo) (*SignupAccountInfo, error) {
	var sign SignupAccountInfo
	for _, trx := range trxList {
		if request.Memo != trx.Memo {
			continue
		}
		separatorPos := strings.Index(trx.Memo, " ")
		if separatorPos <= 0 {
			separatorPos = strings.Index(trx.Memo, "-")
		}
		if separatorPos <= 0 {
			return nil, fmt.Errorf("datx signup transfer memo not correct %v", trx.Memo)
		}
		if QuantityToAmount(request.Quantity) > QuantityToAmount(trx.Quantity) {
			return nil, fmt.Errorf("datx signup quantity cannot be larger than transfer quantity %v", trx.Quantity)
		}
		sign.SysAccount = request.SysAccount
		sign.NewAccount = trx.Memo[:separatorPos]
		sign.PublicKey = trx.Memo[separatorPos+1 : len(trx.Memo)]
		sign.Quantity = trx.Quantity

		return &sign, nil
	}
	return nil, fmt.Errorf("datx signup transfer record not found %v", request.Memo)
}

// ClSystemNewaccount ...
func ClSystemNewaccount(sign *SignupAccountInfo) (string, error) {
	ramAmount := strconv.FormatFloat(QuantityToAmount(sign.Quantity)-0.2, 'f', 4, 64)
	command := "cldatx -u " + WalletConfig.DatxIP + " system newaccount " +
		sign.SysAccount + " " + sign.NewAccount + " " + sign.PublicKey +
		" --stake-net '0.1 DATX' --stake-cpu '0.1 DATX' --buy-ram '" + ramAmount + " DATX'" +
		" -j " + " -f " + " -p " + sign.SysAccount
	fmt.Println(command)
	chainlib.ClWalletUnlock(WalletConfig.PassWord)
	return chainlib.ExecShell(command)
}

// Config ...
type Config struct {
	DatxIP   string `json:"datxip"`
	PassWord string `json:"password"`
}

// WalletConfig ...
var (
	WalletConfig = &Config{}
)

// LoadConfig ...
func LoadConfig() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("load config error: ", err)
	}
	err = json.Unmarshal(file, &WalletConfig)
	if err != nil {
		fmt.Println("para config failed: ", err)
	}
	fmt.Println(WalletConfig)
}
