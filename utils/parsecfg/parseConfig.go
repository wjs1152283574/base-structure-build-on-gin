/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:13:59
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-09-11 15:06:51
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
	Kafka            Kafka
	OpenApi          Open
	TripleDes        Triple
}

// MysqlCfg mysql配置
type MysqlCfg struct {
	Write MysqlWriteCfg
	Read  MysqlReadCfg
}

// MysqlWriteCfg mysql配置
type MysqlWriteCfg struct {
	Host            string
	HostLive        string
	DataBase        string
	Port            string
	PortLive        string
	HostStage       string
	PortStage       string
	PreFix          string
	User            string
	PassWord        string
	SetMaxIdleConns int
	SetMaxOpenConns int
	Charset         string
}

// MysqlReadCfg mysql配置
type MysqlReadCfg struct {
	Host            string
	DataBase        string
	Port            string
	PreFix          string
	User            string
	PassWord        string
	SetMaxIdleConns int
	SetMaxOpenConns int
	ChatSet         string
}

// RedisCfg redis配置
type RedisCfg struct {
	Host        string
	HostLive    string
	Port        string
	PortLive    string
	Auth        string
	HostStage   string
	PortStage   string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

// TimerCfg 定时器配置
type TimerCfg struct {
	Store string
}

// Kafka
type Kafka struct {
	Host  string
	Port  string
	Auth  string
	Topic string
}

// Open
type Open struct {
	Key    string
	Domian string
}

// TripleDes
type Triple struct {
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
