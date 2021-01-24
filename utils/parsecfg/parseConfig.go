package parsecfg

import (
	"sync"

	"github.com/spf13/viper"
)

// 解析配置文件

// ConfigParams 作为全局使用的配置参数 使用sync.Map (开箱即用并发安全锁)
var ConfigParams sync.Map

// ParserConfig  解析yaml配置文件 传入文件路径/文件路径/文件类型
func ParserConfig(filePath, fileName, fileType string) {
	config := viper.New()
	config.AddConfigPath(filePath) //设置读取的文件路径
	config.SetConfigName(fileName) //设置读取的文件名
	config.SetConfigType(fileType) //设置文件的类型

	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	// 读取成功 --> 将参数存入 ConfigParams
	ConfigParams.Store("AppDebug", config.GetBool("AppDebug"))
	ConfigParams.Store("UseDbType", config.GetString("UseDbType"))
	ConfigParams.Store("ApiPort", config.GetString("HttpServer.Api.Port"))
	ConfigParams.Store("WebPort", config.GetString("HttpServer.Web.Port"))
	ConfigParams.Store("AllowCrossDomain", config.GetBool("HttpServer.AllowCrossDomain"))

	// mysql 数据库连接参数
	ConfigParams.Store("MysqlHost", config.GetString("Mysql.Write.Host"))
	ConfigParams.Store("MysqlDataBase", config.GetString("Mysql.Write.DataBase"))
	ConfigParams.Store("MysqlPort", config.GetString("Mysql.Write.Port"))
	ConfigParams.Store("MysqlPrefix", config.GetString("Mysql.Write.Prefix"))
	ConfigParams.Store("MysqlUser", config.GetString("Mysql.Write.User"))
	ConfigParams.Store("MysqlPassword", config.GetString("Mysql.Write.Password"))
	ConfigParams.Store("MysqlCharset", config.GetString("Mysql.Write.Charset"))

	ConfigParams.Store("MysqlSetMaxIdleConns", config.GetString("Mysql.Write.SetMaxIdleConns"))
	ConfigParams.Store("MysqlSetMaxOpenConns", config.GetString("Mysql.Write.SetMaxOpenConns"))
	ConfigParams.Store("MysqlSetConnMaxLifetime", config.GetString("Mysql.Write.SetConnMaxLifetime"))

	ConfigParams.Store("MysqlReConnectInterval", config.GetString("Mysql.Write.ReConnectInterval"))
	ConfigParams.Store("MysqlPingFailRetryTimes", config.GetString("Mysql.Write.PingFailRetryTimes"))
	ConfigParams.Store("MysqlIsOpenReadDb", config.GetString("Mysql.IsOpenReadDb"))

	// 日志
	ConfigParams.Store("GinLogPath", config.GetString("Logs.GinLogName"))

}
