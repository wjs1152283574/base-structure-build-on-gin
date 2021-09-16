/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:13:59
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-09-17 00:31:01
 */
package parsecfg

import (
	"os"

	"github.com/spf13/viper"
)

// 解析配置文件

// GlobalConfig 全局配置
var GlobalConfig EnvCfg

// EnvCfg app 基本设置
type EnvCfg struct {
	Debug            bool
	DbType           string
	Port             string
	SuperUser1       string
	SuperUser2       string
	SuperUser3       string
	Pass             string
	AllowCrossDomain bool
	Env              string
	Mysql            MysqlCfg
	Redis            RedisCfg
	Timer            TimerCfg
	RateLimit        int
	Kafka            KafkaCfg
	OpenApi          OpenCfg
	TripleDes        TripleCfg
}

// MysqlCfg mysql配置
type MysqlCfg struct {
	Dev   MysqlDevCfg
	Prod  MysqlProdCfg
	Stage MysqlStageCfg
}

// MysqlDevCfg mysql dev 配置
type MysqlDevCfg struct {
	Host            string
	DataBase        string
	Port            string
	PreFix          string
	User            string
	PassWord        string
	SetMaxIdleConns int
	SetMaxOpenConns int
	Charset         string
}

// MysqlProdCfg mysql prod 配置
type MysqlProdCfg struct {
	Host            string
	DataBase        string
	Port            string
	PreFix          string
	User            string
	PassWord        string
	SetMaxIdleConns int
	SetMaxOpenConns int
	Charset         string
}

// MysqlStageCfg mysql stage 配置
type MysqlStageCfg struct {
	Host            string
	DataBase        string
	Port            string
	PreFix          string
	User            string
	PassWord        string
	SetMaxIdleConns int
	SetMaxOpenConns int
	Charset         string
}

// RedisCfg redis配置
type RedisCfg struct {
	Dev   RedisDevCfg
	Prod  RedisProdCfg
	Stage RedisStageCfg
}

// RedisDevCfg redis dev 配置
type RedisDevCfg struct {
	Host        string
	Port        string
	Auth        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

// RedisProdCfg redis prod 配置
type RedisProdCfg struct {
	Host        string
	Port        string
	Auth        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

// RedisStageCfg redis stage 配置
type RedisStageCfg struct {
	Host        string
	Port        string
	Auth        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

// TimerCfg 定时器配置
type TimerCfg struct {
	Store string
}

// KafkaCfg kafka 配置
type KafkaCfg struct {
	Dev   KafkaDevCfg
	Prod  KafkaProdCfg
	Stage KafkaStageCfg
}

// KafkaDevCfg kafka dev 配置
type KafkaDevCfg struct {
	Host  string
	Port  string
	Auth  string
	Topic string
}

// KafkaProdCfg kafka prod 配置
type KafkaProdCfg struct {
	Host  string
	Port  string
	Auth  string
	Topic string
}

// KafkaStageCfg kafka stage 配置
type KafkaStageCfg struct {
	Host  string
	Port  string
	Auth  string
	Topic string
}

// Open
type OpenCfg struct {
	Key    string
	Domian string
}

// TripleDes
type TripleCfg struct {
	Key string
	Iv  string
}

func init() {
	path, _ := os.Getwd()
	cfg := viper.New()
	cfg.AddConfigPath(path + "/config")
	cfg.SetConfigName("cfg")
	cfg.SetConfigType("json")
	if err := cfg.ReadInConfig(); err != nil { // 必须 先 读取 `ReadInConfig`
		panic(err)
	}
	if err := cfg.Unmarshal(&GlobalConfig); err != nil { // 才能反序列化到 结构体里面
		panic("读取配置文件出错")
	}
	cfg.WatchConfig() // 配置文件生效 之后 才进行监听 : 写在最后
}
