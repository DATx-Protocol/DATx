package delayqueue

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	// RedisPool RedisPool连接池实例
	RedisPool *redis.Pool
)

// 初始化连接池
func initRedisPool() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:      Setting.Redis.MaxIdle,
		MaxActive:    Setting.Redis.MaxActive,
		IdleTimeout:  300 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisTestOnBorrow,
		Wait:         true,
	}

	return pool
}

// 连接redis
func redisDial() (redis.Conn, error) {
	conn, err := redis.Dial(
		"tcp",
		Setting.Redis.Host,
		redis.DialConnectTimeout(time.Duration(Setting.Redis.ConnectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(Setting.Redis.ReadTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(Setting.Redis.WriteTimeout)*time.Millisecond),
	)
	if err != nil {
		log.Printf("连接redis失败#%s", err.Error())
		return nil, err
	}

	if Setting.Redis.Password != "" {
		if _, err := conn.Do("AUTH", Setting.Redis.Password); err != nil {
			conn.Close()
			log.Printf("redis认证失败#%s", err.Error())
			return nil, err
		}
	}

	_, err = conn.Do("SELECT", Setting.Redis.Db)
	if err != nil {
		conn.Close()
		log.Printf("redis选择数据库失败#%s", err.Error())
		return nil, err
	}

	return conn, nil
}

// 从池中取出连接后，判断连接是否有效
func redisTestOnBorrow(conn redis.Conn, t time.Time) error {
	_, err := conn.Do("PING")
	if err != nil {
		log.Printf("从redis连接池取出的连接无效#%s", err.Error())
	}

	return err
}

// 执行redis命令, 执行完成后连接自动放回连接池
func execRedisCommand(command string, args ...interface{}) (interface{}, error) {
	redis := RedisPool.Get()
	defer redis.Close()

	return redis.Do(command, args...)
}

// SET
func execRedisSet(key string, val []byte) (interface{}, error) {
	return execRedisCommand("SET", key, val)
}

// DEL
func execRedisDel(key string) (interface{}, error) {
	return execRedisCommand("DEL", key)
}

// GET
func execRedisGet(key string) ([]byte, error) {
	return redis.Bytes(execRedisCommand("GET", key))
}

// EXISTS
func execRedisEXISTS(key string) (bool, error) {
	return redis.Bool(execRedisCommand("EXISTS", key))
}
