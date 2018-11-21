package http

import (
	"datx/lsd/chainlib"
	"datx/lsd/delayqueue"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//RedisResponse ...
type RedisResponse struct {
	Code    uint64 `json:"code"`    //	1成功	0失败
	Message error  `json:"message"` //	系统自带的错误信息
	Val     string `json:"val"`     //	如果GET成功就是Key对应的Val
}

//RedisRequest ...
var RedisRequest = func(w http.ResponseWriter, r *http.Request) {
	var cmd delayqueue.RedisCmdInfo
	var resp RedisResponse
	//cmd := RedisCmdInfo{"SET", "lmx", "01XXXXXXXXX"}
	if r.Method == "POST" || r.Method == "GET" {
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &cmd); err != nil {
			resp = RedisResponse{0, err, ""}
		} else {
			if res, err := delayqueue.ExecRedisCmd(cmd); err != nil {
				resp = RedisResponse{0, err, ""}
			} else {
				if cmd.Cmd == "DEL" || cmd.Cmd == "EXISTS" {
					resp = RedisResponse{res.(uint64), nil, ""}
				} else {
					resp = RedisResponse{1, nil, res.(string)}
				}
			}
		}
		respjson, _ := json.Marshal(resp)
		w.Write([]byte(string(respjson)))
	}
}

//ETHExtractHandler get transaction and push set success to the contract of DatxExtract ...
var ETHExtractHandler = func(w http.ResponseWriter, r *http.Request) {
	isSuccess, str := func() (bool, string) {
		var trx chainlib.Transaction

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return false, fmt.Sprintf("ETHExtractHandler read body data err: %v", err)
		}

		if err := json.Unmarshal(body, &trx); err != nil {
			return false, fmt.Sprintf("ETHExtractHandler unmarshal json data err: %v", err)
		}

		log.Printf("ETH HTTP get Trx: %v\n", trx)

		jobid := trx.Category + "_" + trx.TransactionID
		if job, _ := delayqueue.Get(jobid); job != nil {
			return false, fmt.Sprintf("ETHExtractHandler trx is existed: %v", trx.TransactionID)
		}

		var job delayqueue.Job
		job.Topic = trx.Category
		job.Id = trx.Category + "_" + trx.TransactionID
		job.Delay = time.Now().Unix()
		job.TTR = 60

		bytes, err := json.Marshal(trx)
		if err != nil {
			return false, fmt.Sprintf("ETHExtractHandler marshal json data err: %v", err)
		}
		job.Body = string(bytes)

		if err = delayqueue.Push(job); err != nil {
			log.Printf("Push queue failed.%v\n", err)
			return false, fmt.Sprintf("ETHExtractHandler Push queue err: %v", err)
		}

		return true, fmt.Sprintf("ETHExtractHandler Push queue success: %v", trx.TransactionID)

	}()

	log.Printf("ETH Http Response: %v\n", str)

	if isSuccess {
		w.Write([]byte(str))
	} else {
		http.Error(w, str, 500)
	}
}
