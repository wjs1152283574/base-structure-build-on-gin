/*
 * @Author: Casso-Wong
 * @Date: 2021-06-04 14:41:32
 * @Last Modified by:   Casso-Wong
 * @Last Modified time: 2021-06-04 14:41:32
 */
package appredis

import (
	"encoding/json"
	"goweb/utils/parsecfg"
	"time"

	"github.com/gomodule/redigo/redis"
)

// redis 数据库初始化 以及具体配置

// RedisDefaultPool redis 连接池
var RedisDefaultPool *redis.Pool

func init() {
	if parsecfg.GlobalConfig.Env == "dev" {
		RedisDefaultPool = NewPool(parsecfg.GlobalConfig.Redis.Host, parsecfg.GlobalConfig.Redis.Port, parsecfg.GlobalConfig.Redis.Auth, parsecfg.GlobalConfig.Redis.MaxIdle, parsecfg.GlobalConfig.Redis.MaxActive)
	}
	if parsecfg.GlobalConfig.Env == "test" {
		RedisDefaultPool = NewPool(parsecfg.GlobalConfig.Redis.HostLive, parsecfg.GlobalConfig.Redis.PortLive, parsecfg.GlobalConfig.Redis.Auth, parsecfg.GlobalConfig.Redis.MaxIdle, parsecfg.GlobalConfig.Redis.MaxActive)
	}
}

// NewPool 项目运行初始化redis连接池
func NewPool(addr, port, auth string, max, maxactive int) *redis.Pool { // 传入 ip:port 最大闲置连接数 最大活跃连接数
	str := addr + ":" + port
	return &redis.Pool{
		MaxIdle:     max,
		MaxActive:   maxactive,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", str)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", auth); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", 0); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
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

// SetExpire 是这基础类型+ TTL
func SetExpire(key string, data interface{}, time int) error {
	conn := RedisDefaultPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, data)
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

// SetHash 设置基础类型
func SetArr(key, val string) error {
	conn := RedisDefaultPool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, val)
	return err
}

// DelArr 移除list指定值
func DelArr(key, val string) error {
	conn := RedisDefaultPool.Get()
	defer conn.Close()
	_, err := conn.Do("LREM", key, 0, val)
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

// GetLIke 获取相似Key
func GetLIke(key string) (res [][]byte, err error) {
	conn := RedisDefaultPool.Get()
	defer conn.Close()
	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return [][]byte{}, err
	}
	var r [][]byte
	for _, v := range keys {
		val, _ := Get(v)
		r = append(r, val)
	}
	return r, nil
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

// Rpush 插入列表尾部
func Rpush(key string, val []byte) ([]byte, error) {
	conn := RedisDefaultPool.Get()
	defer conn.Close()
	reply, err := redis.Bytes(conn.Do("RPUSH", key, val))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// GetList 获取列表
func GetList(key string) (res []string, err error) {
	conn := RedisDefaultPool.Get()
	defer conn.Close()
	reply, err := redis.Strings(conn.Do("LRANGE", key, 0, -1))
	if err != nil {
		return nil, err
	}
	return reply, nil
}
