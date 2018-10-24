package explorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// TokenValue ...
type TokenValue struct {
	Token   string `json:"token"`   // 代币，DATX，DBTC，DETH，DEOS，BTC，ETH，EOS
	Balance string `json:"balance"` // 数量
	Value   string `json:"value"`   // 价值（USD）
}

// TokenValueRequest ...
type TokenValueRequest struct {
	Token   string `form:"token" json:"token" binding:"required"`
	Address string `form:"address" json:"address" binding:"required"`
}

// WalletValueRequest ...
type WalletValueRequest struct {
	Category string `form:"category" json:"category" binding:"required"`
	Address  string `form:"address" json:"address" binding:"required"`
}

// QuantityToAmount ... 不做错误处理
func QuantityToAmount(quantity string) float64 {
	value := quantity[:strings.Index(quantity, " ")]
	amount, _ := strconv.ParseFloat(value, 64)
	return amount
}

// GetWalletValue ... DATX，BTC，ETH
func GetWalletValue(category string, account string) ([]*TokenValue, error) {
	tokens := make([]*TokenValue, 0)
	if category == "BTC" || category == "ETH" {
		tokenValue, err := GetTokenValue(category, account)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tokenValue)
	}
	if category == "DATX" {
		if tokenValue, err := GetTokenValue("DATX", account); err == nil {
			tokens = append(tokens, tokenValue)
		} else if tokenValue, err := GetTokenValue("DBTC", account); err == nil {
			tokens = append(tokens, tokenValue)
		} else if tokenValue, err := GetTokenValue("DETH", account); err == nil {
			tokens = append(tokens, tokenValue)
		} else if tokenValue, err := GetTokenValue("DEOS", account); err == nil {
			tokens = append(tokens, tokenValue)
		} else {
			return nil, fmt.Errorf("wallet not found")
		}
	}
	return tokens, nil
}

// GetTokenValue ... DATX，DBTC，DETH，DEOS，BTC，ETH
func GetTokenValue(token string, account string) (*TokenValue, error) {
	balance, err := GetTokenBalance(token, account)
	if err != nil {
		return nil, err
	}
	price, err := GetTokenPrice(token)
	if err != nil {
		return nil, err
	}
	quantity := strconv.FormatFloat(balance, 'f', 4, 64) + " " + token
	value := "$" + strconv.FormatFloat(balance*price, 'f', 2, 64)
	return &TokenValue{token, quantity, value}, nil
}

// GetTokenBalance ... DATX，DBTC，DETH，DEOS，BTC，ETH
func GetTokenBalance(token string, account string) (float64, error) {
	if token == "DATX" || token == "DBTC" || token == "DETH" || token == "DEOS" {
		return GetDATXBalance(token, account)
	} else if token == "BTC" {
		return GetBTCBalance(token, account)
	} else if token == "ETH" {
		return GetETHBalance(token, account)
	}
	return 0, fmt.Errorf("balance not found")
}

// GetTokenPrice ...
func GetTokenPrice(token string) (float64, error) {
	httpClient := &http.Client{}
	var url string
	if token == "DATX" {
		url = "https://api.coinmarketcap.com/v2/ticker/2567/"
	} else if token == "BTC" || token == "DBTC" {
		url = "https://api.coinmarketcap.com/v2/ticker/1/"
	} else if token == "ETH" || token == "DETH" {
		url = "https://api.coinmarketcap.com/v2/ticker/1027/"
	} else if token == "EOS" || token == "DEOS" {
		url = "https://api.coinmarketcap.com/v2/ticker/1765/"
	} else {
		return 0, fmt.Errorf("price not found")
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("price not found")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		return 0, err
	}
	return js.Get("data").Get("quotes").Get("USD").Get("price").MustFloat64(), nil
}

// GetDATXBalance ... DATX，DBTC，DETH，DEOS
func GetDATXBalance(token string, account string) (float64, error) {
	formData := make(map[string]interface{})
	// 转换成发币合约
	if token == "DATX" {
		formData["code"] = "datxos.token"
	} else {
		formData["code"] = "datxos." + strings.ToLower(token)
	}
	formData["account"] = account
	bytesData, err := json.Marshal(formData)
	if err != nil {
		return 0, err
	}

	URL := WalletConfig.DatxIP + "/v1/chain/get_currency_balance"
	request, err := http.NewRequest("POST", URL, bytes.NewReader(bytesData))
	if err != nil {
		return 0, fmt.Errorf("datx get_currency_balance request error %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("datx get_currency_balance resp error %v", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("datx get_currency_balance resp status error %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("datx get_currency_balance body error %v", err)
	}
	// [] 要做处理
	if len(string(body)) <= 2 {
		return 0, fmt.Errorf("datx get_currency_balance balance not found %v", string(body))
	}
	// ["1216.5698 DBTC"] 要做处理
	quantity := string(body)
	quantity = quantity[2 : len(quantity)-2]
	return QuantityToAmount(quantity), nil
}

// GetETHBalance ...
func GetETHBalance(token string, account string) (float64, error) {
	// http://api.etherscan.io/api?module=account&action=balance&address=0xf2a0132b971f6a4b980d7d3fd9555a42988c6062&tag=latest&apikey=8FT3VZVAS94PIHPYKEWPWHC4ZICB71RFSM
	URL, err := url.Parse("http://api.etherscan.io/api")
	if err != nil {
		return 0, fmt.Errorf("eth account balance url error %v", err)
	}

	params := url.Values{}
	params.Set("module", "account")
	params.Set("action", "balance")
	params.Set("address", account)
	params.Set("tag", "latest")
	params.Set("apikey", "bala8FT3VZVAS94PIHPYKEWPWHC4ZICB71RFSMnce")
	URL.RawQuery = params.Encode()

	resp, err := http.Get(URL.String())
	if err != nil {
		return 0, fmt.Errorf("eth account balance resp error %v", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("eth account balance resp status error %v", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("eth account balance body error %v", err)
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		return 0, fmt.Errorf("eth account balance simplejson error %v", err)
	}
	status := js.Get("status").MustString()
	result := js.Get("result").MustString()
	if status != "1" {
		return 0, fmt.Errorf("eth account balance result error %v", result)
	}
	balance, err := strconv.ParseFloat(result, 64)
	if err != nil || balance == 0 {
		return 0, fmt.Errorf("eth account balance balance not found %v", err)
	}
	return balance / 1000000000000000000, nil
}

// GetBTCBalance ...
func GetBTCBalance(token string, account string) (float64, error) {
	// https://blockchain.info/balance?active=32nxUrjexffHzwN4xVVwLivU2mpStQRtfz
	URL, err := url.Parse("https://blockchain.info/balance")
	if err != nil {
		return 0, err
	}

	params := url.Values{}
	params.Set("active", account)
	URL.RawQuery = params.Encode()

	resp, err := http.Get(URL.String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return 0, fmt.Errorf("balance not found")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		return 0, err
	}
	return js.Get(account).Get("final_balance").MustFloat64() / 100000000, nil
}
