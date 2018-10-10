package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//router.LoadHTMLGlob("templates/**/*")
	router.POST("/token_balance", postTokenBalance)
	router.POST("/wallet_balance", postWalletBalance)
	router.Run(":8081")
}
