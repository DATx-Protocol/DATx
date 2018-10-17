package main

import (
	"datx/ListenServer/explorer"

	"github.com/gin-gonic/gin"
)

func init() {
	explorer.LoadConfig()
}

func main() {
	router := gin.Default()
	//router.LoadHTMLGlob("templates/**/*")
	router.POST("/token_balance", postTokenBalance)
	router.POST("/wallet_balance", postWalletBalance)
	router.POST("/wallet_trx_list", postWalletTrxList)
	router.POST("/resource", postDATXResource)
	router.POST("/new_account", postDATXSignup)
	router.POST("/address_map", postAddressMap)
	router.Run(":9101")
}
