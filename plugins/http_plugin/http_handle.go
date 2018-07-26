package httpplugin

import (
	"datx_chain/chainlib/application"
	"datx_chain/chainlib/controller"
	"datx_chain/chainlib/types"
	"datx_chain/plugins/chain_plugin"
	"datx_chain/utils/helper"
	"encoding/json"
	"fmt"
	template "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

//TransferHandle promote a transfer operation
var TransferHandle = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		//get exec file current path
		CurrentPath := helper.GetCurrentPath()
		in := strings.Index(CurrentPath, "datx_chain")
		if in < 1 {
			return
		}

		configpath := helper.MakePath(CurrentPath[:in-1], "datx_chain", application.App().GetConfigFolder(), "transfer.html")
		t, err := template.ParseFiles(configpath)
		if err != nil {
			log.Printf("parse err: %v", err)
			return
		}

		t.Execute(w, nil)
	} else {
		r.ParseForm()
		log.Printf("ParseFrom: %v", r.Form)

		from := r.Form["from"][0]
		to := r.Form["to"][0]
		amount, _ := strconv.Atoi(r.Form["amount"][0])

		var wait sync.WaitGroup
		wait.Add(1)

		var response string
		pkg := CreatePackedTransaction(CreateTransfer("datx", from, to, uint16(amount)))
		PushTransaction(pkg, func(inerr error, trace *types.TransactionTrace) {
			if inerr != nil {
				response = fmt.Sprintf("%v", inerr)
			} else {
				if trace != nil {
					id := trace.ID.String()

					//init chain
					var info string
					plugin, err := application.App().Find("chain")
					if err != nil {
						log.Print("you do not add producerplugin to app before init the producerplugin")
						return
					}
					chain := plugin.(*chainplugin.ChainPlugin)
					if chain != nil {
						for k, v := range chain.Chain().TestAccounts {
							info += fmt.Sprintf("%s %d;", k, v.Amount)
						}
					}

					response = fmt.Sprintf("transaction:%v\n trace:%v\n details: %v", id, trace, info)
				} else {
					response = fmt.Sprint("Unknow error")
				}
			}

			wait.Done()
		})

		wait.Wait()

		w.Write([]byte(response))
	}
}

//GetTransactionListHandle promote handler of TransactionList
var GetTransactionListHandle = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "GET" {
		transactionList, err := json.Marshal(FindTransaction())
		log.Printf("TransactionList: %v", transactionList)
		if err != nil {
			log.Printf("Marshal json err: %v", err)
		}
		w.Write(transactionList)

	}
}

//GetTransactionByHashHandle promote handler of transaction with trx hash
var GetTransactionByHashHandle = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "GET" {
		tra := &controller.TransactionDetail{}
		trx, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(trx, &tra); err == nil {
			id := helper.RLPHash(tra.TrxHash)
			transaction, e := json.Marshal(QueryTransactionById(id))
			log.Printf("Transaction of the hash is :%v ", transaction)
			w.Write(transaction)
			if e != nil {
				log.Printf("Marshal json err: %v", e)
			}
		} else {
			log.Printf("Unmarshal json err : %v", err)
		}

	}
}

//GetBlockListHandle 区块列表
var GetBlockListHandle = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res := QueryBlocks()
		respjson, _ := json.Marshal(res)
		response := string(respjson)
		w.Write([]byte(response))
	} else {
		res := QueryBlocks()
		respjson, _ := json.Marshal(res)
		response := string(respjson)
		w.Write([]byte(response))
	}
}

//GetGeneralInfo 总体信息
var GetGeneralInfo = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		res := QueryGeneralInfo()
		respjson, _ := json.Marshal(res)
		response := string(respjson)
		w.Write([]byte(response))
	} else {
		res := QueryGeneralInfo()
		respjson, _ := json.Marshal(res)
		response := string(respjson)
		w.Write([]byte(response))
	}
}
