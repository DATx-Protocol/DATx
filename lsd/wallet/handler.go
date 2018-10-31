package main

import (
	"datx/lsd/chainlib"
	"datx/lsd/explorer"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TrxResult ...
type TrxResult struct {
	Hash string `json:"hash"`
}

func jsonBindingError(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":    400,
		"message": "Json Binding Error",
		"data":    nil,
	})
}

func explorerError(ctx *gin.Context, err error) {
	if strings.Contains(err.Error(), "not found") {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
	}
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

func postWalletTrxList(ctx *gin.Context) {
	var request explorer.WalletTrxRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		jsonBindingError(ctx)
		return
	}
	trxList, err := explorer.GetWalletTrxList(request.Category, request.Address, request.Limit)
	if err != nil {
		explorerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"data":    trxList,
	})
}

func postDATXResource(ctx *gin.Context) {
	var request explorer.DATXResourceRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		jsonBindingError(ctx)
		return
	}
	trxList, err := explorer.GetDATXResource(request.Account)
	if err != nil {
		explorerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"data":    trxList,
	})
}

func postDATXSignup(ctx *gin.Context) {
	var request explorer.SignupAccountRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		jsonBindingError(ctx)
		return
	}

	trxList, err := explorer.GetSignupTrxList(request.SysAccount, 0, 1000)
	if err != nil {
		explorerError(ctx, err)
		return
	}
	signupAcc, err := explorer.MatchSignupAccount(request, trxList)
	if err != nil {
		explorerError(ctx, err)
		return
	}
	outStr, err := explorer.ClSystemNewaccount(signupAcc)
	if err != nil {
		explorerError(ctx, err)
		return
	}
	trxID, err := chainlib.ParseTrxID(outStr)
	if err != nil {
		explorerError(ctx, err)
		return
	}
	err = explorer.WaitIrreversible(trxID)
	if err != nil {
		explorerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"data":    TrxResult{trxID},
	})
}

func postAddressMap(ctx *gin.Context) {
	var request explorer.AddressMapRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		jsonBindingError(ctx)
		return
	}

	outStr, err := explorer.ClRecordUser(request.DatxAddress, request.Address)
	if err != nil {
		explorerError(ctx, err)
		return
	}
	trxID, err := chainlib.ParseTrxID(outStr)
	if err != nil {
		explorerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"data":    TrxResult{trxID},
	})
}

// AccountsResult ...
type AccountsResult struct {
	Names []string `json:"account_names"`
}

func postGetAccounts(ctx *gin.Context) {
	var request explorer.SignupPKRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		jsonBindingError(ctx)
		return
	}

	outStr, err := explorer.ClGetAccounts(request.PublicKey)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid public key") {
			explorerError(ctx, fmt.Errorf("public_key not found"))
		} else {
			explorerError(ctx, err)
		}
		return
	}

	names := &AccountsResult{}
	if err := json.Unmarshal([]byte(outStr), &names); err != nil {
		explorerError(ctx, fmt.Errorf("datx get_actions unmarshal error %v", outStr))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "OK",
		"data":    names,
	})
}
