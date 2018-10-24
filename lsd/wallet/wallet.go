package main

import (
	"datx/lsd/explorer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	explorer.LoadConfig()
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"PUT", "GET", "OPTIONS", "POST"}
	router.Use(cors.New(config))
	//router.LoadHTMLGlob("templates/**/*")
	router.POST("/token_balance", postTokenBalance)
	router.POST("/wallet_balance", postWalletBalance)
	router.POST("/wallet_trx_list", postWalletTrxList)
	router.POST("/resource", postDATXResource)
	router.POST("/new_account", postDATXSignup)
	router.POST("/address_map", postAddressMap)
	router.POST("/get_accounts", postGetAccounts)
	router.Run(":9101")
}
