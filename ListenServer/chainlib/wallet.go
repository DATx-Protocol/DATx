package chainlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// TokenInfo ...
type TokenInfo struct {
	Symbol  string  `json:"symbol"`  // 代币，DATX，DBTC，DETH，DEOS，BTC，ETH，EOS
	Balance float64 `json:"balance"` // 数量
	Price   float64 `json:"price"`   // 单价（USD）
}

// GetTokenPrice ...
func GetTokenPrice(symbol string) (float64, error) {
	httpClient := &http.Client{}
	var url string
	if symbol == "DATX" {
		url = "https://api.coinmarketcap.com/v2/ticker/2567/"
	} else if symbol == "BTC" {
		url = "https://api.coinmarketcap.com/v2/ticker/1/"
	} else if symbol == "ETH" {
		url = "https://api.coinmarketcap.com/v2/ticker/1027/"
	} else if symbol == "EOS" {
		url = "https://api.coinmarketcap.com/v2/ticker/1765/"
	} else {
		return 0, errors.New("symbol not found")
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
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
func GetDATXBalance(symbol string, account string) (float64, error) {
	formData := make(map[string]interface{})
	formData["code"] = "datxio." + strings.ToLower(symbol)
	formData["account"] = account
	bytesData, err := json.Marshal(formData)
	if err != nil {
		return 0, err
	}

	URL := "http://127.0.0.1:8888/v1/chain/get_currency_balance"
	request, err := http.NewRequest("POST", URL, bytes.NewReader(bytesData))
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	// ["1216.5698 DBTC"] 要做处理
	quantity := string(body)
	quantity = quantity[2 : len(quantity)-2]
	symbolpos := strings.Index(quantity, " "+strings.ToUpper(symbol))
	quantity = quantity[:symbolpos]
	balance, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// GetETHBalance ...
func GetETHBalance(symbol string, account string) (float64, error) {
	// https: //api.etherscan.io/api?module=account&action=balance&address=0xf2a0132b971f6a4b980d7d3fd9555a42988c6062&tag=latest&apikey=8FT3VZVAS94PIHPYKEWPWHC4ZICB71RFSM
	URL, err := url.Parse("https://api.etherscan.io/api")
	if err != nil {
		return 0, err
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
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		return 0, err
	}
	result := js.Get("result").MustString()
	balance, err := strconv.ParseFloat(result, 64)
	if err != nil {
		return 0, err
	}
	return balance / 1000000000000000000, nil
}

// GetBTCBalance ...
func GetBTCBalance(symbol string, account string) (float64, error) {
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

// GetTokenBalance ... DATX，DBTC，DETH，DEOS，BTC，ETH
func GetTokenBalance(symbol string, account string) (float64, error) {
	if symbol == "DATX" || symbol == "DBTC" || symbol == "DETH" || symbol == "DEOS" {
		return GetDATXBalance(symbol, account)
	} else if symbol == "BTC" {
		return GetBTCBalance(symbol, account)
	} else if symbol == "ETH" {
		return GetETHBalance(symbol, account)
	}
	return 0, errors.New("symbol not found")
}
