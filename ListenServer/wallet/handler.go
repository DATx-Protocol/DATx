package main

import (
	"datx/ListenServer/explorer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func jsonBindingError(ctx *gin.Context) {
	ctx.JSON(400, gin.H{
		"code":    400,
		"message": "Json Binding Error",
		"data":    nil,
	})
}

func explorerError(ctx *gin.Context, err error) {
	ctx.JSON(400, gin.H{
		"code":    500,
		"message": err.Error(),
		"data":    nil,
	})
}

func postTokenBalance(ctx *gin.Context) {
	var request explorer.TokenValueRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		jsonBindingError(ctx)
		return
	}
	tokenValue, err := explorer.GetTokenValue(request.Token, request.Address)
	if err != nil {
		explorerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"data":    tokenValue,
	})
}

func postWalletBalance(ctx *gin.Context) {
	var request explorer.WalletValueRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		jsonBindingError(ctx)
		return
	}
	walletValue, err := explorer.GetWalletValue(request.Category, request.Address)
	if err != nil {
		explorerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"data":    walletValue,
	})
}
