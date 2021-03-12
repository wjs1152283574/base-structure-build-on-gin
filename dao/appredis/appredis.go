package appredis

import (
	"encoding/json"
	"fmt"
	"goweb/utils/parsecfg"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisDefaultPool redis 连接池
var RedisDefaultPool *redis.Pool

func init() {
	RedisDefaultPool = NewPool(parsecfg.GlobalConfig.Redis.Host, parsecfg.GlobalConfig.Redis.Port, parsecfg.GlobalConfig.Redis.MaxIdle, parsecfg.GlobalConfig.Redis.MaxActive)
}

// NewPool 项目运行初始化redis连接池
func NewPool(addr, port string, max, maxactive int) *redis.Pool { // 传入 ip:port 最大闲置连接数 最大活跃连接数
	str := addr + ":" + port
	fmt.Println(str)
	return &redis.Pool{
		MaxIdle:     max,
		MaxActive:   maxactive,
		IdleTimeout: 240 * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", str)
		},
	}
}

// SetJSON 设置key
func SetJSON(key string, data interface{}, time int) error {
	conn := RedisDefaultPool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// SetHash 设置基础类型
func SetHash(key, val string) error {
	conn := RedisDefaultPool.Get()
	defer conn.Close()

	_, err := conn.Do("set", key, val)
	return err
}

// Exists 是否存在key
func Exists(key string) bool {
	conn := RedisDefaultPool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get 返回key值
func Get(key string) ([]byte, error) {
	conn := RedisDefaultPool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Delete 删除key
func Delete(key string) (bool, error) {
	conn := RedisDefaultPool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes 删除相似key
func LikeDeletes(key string) error {
	conn := RedisDefaultPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

// Rpop 返回key值
func Rpop(key string) ([]byte, error) {
	conn := RedisDefaultPool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("rpop", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}
