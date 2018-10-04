package delayqueue

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//RedisResponse ...
type RedisResponse struct {
	Code    uint64 `json:"code"`    //	1成功	0失败
	Message error  `json:"message"` //	系统自带的错误信息
	Val     string `json:"val"`     //	如果GET成功就是Key对应的Val
}

//RedisRequest ...
var RedisRequest = func(w http.ResponseWriter, r *http.Request) {
	var cmd RedisCmdInfo
	var resp RedisResponse
	//cmd := RedisCmdInfo{"SET", "lmx", "01XXXXXXXXX"}
	if r.Method == "POST" || r.Method == "GET" {
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &cmd); err != nil {
			resp = RedisResponse{0, err, ""}
		} else {
			if res, err := ExecRedisCmd(cmd); err != nil {
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
