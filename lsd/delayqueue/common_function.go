package delayqueue

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

//RedisCmdInfo ...
type RedisCmdInfo struct {
	Cmd string `json:"cmd"`
	Key string `json:"key"`
	Val string `json:"val"`
}

//ExecRedisCmd ...
func ExecRedisCmd(cmd RedisCmdInfo) (interface{}, error) {
	switch cmd.Cmd {
	case "SET":
		return redis.String(execRedisCommand("SET", cmd.Key, cmd.Val))
	case "DEL":
		return redis.Uint64(execRedisCommand("DEL", cmd.Key))
	case "GET":
		return redis.String(execRedisCommand("GET", cmd.Key))
	default:
		return redis.Uint64(execRedisCommand("EXISTS", cmd.Key))
	}
}

//PingRedisPool ping redis server
func PingRedisPool() bool {
	_, err := execRedisCommand("PING")
	if err != nil {
		log.Printf("Redis ping error: %v\n", err)
		return false
	}

	return true
}
