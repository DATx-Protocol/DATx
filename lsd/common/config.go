package common

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
)

//GetConfig ...
func GetConfig() (*ini.File, error) {
	var cfg *ini.File
	var err error
	if runtime.GOOS == "linux" {
		// log.Println("Unix/Linux type OS detected")
		cfg, err = ini.Load(os.Getenv("HOME") + "/.local/share/datxos/noddatx/config/config.ini")
	} else if runtime.GOOS == "darwin" {
		// log.Println("Mac OS detected")
		cfg, err = ini.Load(os.Getenv("HOME") + "/Library/Application Support/datxos/noddatx/config/config.ini")
	} else {
		return nil, fmt.Errorf("%s detected,not support", runtime.GOOS)
	}

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

//GetCfgProducerName ...
func GetCfgProducerName() string {
	var result string
	cfg, err := GetConfig()
	if err != nil {
		log.Printf("Get config err:%v\n", err)
		return result
	}

	result = cfg.Section("").Key("producer-name").String()
	return result
}

//GetCfgProducerKey Get node key
func GetCfgProducerKey() []string {
	cfg, err := GetConfig()
	if err != nil {
		log.Printf("Get config err:%v\n", err)
		return nil
	}

	sig := cfg.Section("").Key("signature-provider").String()
	result := strings.Split(sig, "=KEY:")
	return result
}

//GetTrusteeAccount get account addr by name
func GetTrusteeAccount(name string) string {
	var result string
	cfg, err := GetConfig()
	if err != nil {
		log.Printf("GetETHTrusteeAccount Get config err:%v\n", err)
		return result
	}

	result = cfg.Section("").Key(name).String()
	return result
}

//GetWalletNameAndPassword ...
func GetWalletNameAndPassword() (string, string) {
	var cfg *ini.File
	var err error
	walletPath := os.Getenv("HOME") + "/datxos-wallet/wallet_password.ini"

	cfg, err = ini.Load(walletPath)
	if err != nil {
		log.Printf("GetWalletPassword: %v\n", err)
		return "", ""
	}

	name := cfg.Section("").Key("wallet-namer").String()
	password := cfg.Section("").Key("wallet-password").String()
	return name, password
}
