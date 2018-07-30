package httpplugin

import (
	"crypto/md5"
	"datx_chain/chainlib/application"
	"datx_chain/chainlib/controller"
	"datx_chain/chainlib/types"
	"datx_chain/plugins/chain_plugin"
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
	"encoding/hex"
	"encoding/json"
	"fmt"
	template "html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var TransferHandler = func(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" || r.Method == "POST" {
		var response string
		r.ParseForm()
		log.Printf("ParseFrom: %v", r.Form)
		data := "20Awazu18"
		from := r.Form["From"][0]
		data += from
		to := r.Form["To"][0]
		data += to
		a := r.Form["Amount"][0]
		amount, _ := strconv.Atoi(r.Form["Amount"][0])
		data += a
		memo := r.Form["Memo"]
		time := r.Form["time"]
		token := r.Form["token"]

		if len(memo) == 0 || len(time) == 0 || len(token) == 0 {
			response = fmt.Sprintf("the token is wrong")
			return
		}
		data += memo[0]
		data += time[0]
		log.Printf("data  is  %v", data)
		d := []byte(data)
		t := md5.Sum(d)
		var x []byte
		for _, v := range t {
			x = append(x, v)
		}
		b := hex.EncodeToString(t[:])
		log.Printf(" first md5 %v ", b)
		c := md5.Sum([]byte(b))
		if hex.EncodeToString(c[:]) != token[0] {
			response = fmt.Sprintf("the token is wrong")
			log.Printf("mytoken is :%v,posttoken is %v", hex.EncodeToString(c[:]), token[0])
			return
		}
		// transaction, e := ioutil.ReadAll(r.Body)
		// if e != nil {
		// 	res = fmt.Sprintf("Marshal json err: %v", e)
		// 	w.Write([]byte(res))
		// }
		// if err := json.Unmarshal(transaction, &transfer); err == nil {
		// from := transfer.From
		// to := transfer.To
		// amount := transfer.Amount
		var wait sync.WaitGroup
		wait.Add(1)

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

//TransferHandle promote a transfer operation
var TransferHandle = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//get exec file current path
		transfer := controller.Transfer{}
		transfer.From = "avazudemo"
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
		var res string
		transactionList, e := FindTransaction()

		if e != nil {
			res = fmt.Sprintf("FindTransaction with error :%v ", e)
			w.Write([]byte(res))

		}
		trxs, err := json.Marshal(transactionList)
		if err != nil {
			res = fmt.Sprintf("Marshal json err: %v", err)
			w.Write([]byte(res))
		}
		w.Write(trxs)

	}
}

//GetTransactionByHashHandle promote handler of transaction with trx hash
var GetTransactionByHashHandle = func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "GET" {
		var res string
		r.ParseForm()
		log.Printf("ParseFrom: %v", r.Form)
		TrxHash := r.Form["TrxHash"]
		if len(TrxHash) == 0 {
			return
		}

		transaction, e := json.Marshal(QueryTransactionById(common.HexToHash(TrxHash[0])))
		w.Write(transaction)
		if e != nil {
			res = fmt.Sprintf("Marshal json err: %v", e)
			w.Write([]byte(res))

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
