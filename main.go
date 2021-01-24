package main

import (
	"goweb/dao/daomysql"
	"goweb/database/mysqlmodel"
	"goweb/routers"
	"goweb/utils/parsecfg"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var basePath string // 项目根路径

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	basePath = path
}

func main() {
	//获取项目的执行路径
	
	path := basePath + "/config" // 项目执行目录的config目录下存放配置文件 `pro_cfg.yaml`
	parsecfg.ParserConfig(path, "pro_cfg", "yaml")

	dataType, _ := parsecfg.ConfigParams.Load("UseDbType")
	dbUser, _ := parsecfg.ConfigParams.Load("MysqlUser")
	dbPwd, _ := parsecfg.ConfigParams.Load("MysqlPassword")
	dbName, _ := parsecfg.ConfigParams.Load("MysqlDataBase")
	dbCharset, _ := parsecfg.ConfigParams.Load("MysqlCharset")
	daomysql.InitMysql(dataType.(string), dbUser.(string), dbPwd.(string), dbName.(string), dbCharset.(string)) // 初始化数据库连接

	defer daomysql.DB.Close() // 预操作: 关闭数据库连接

	mysqlmodel.MyMigrate() // 根据给定的结构体进行迁移 迁移的结构体需要在dbmodel中定义

	r := routers.InitAPIRouter(basePath) // 初始化路由监听

	APIPort, _ := parsecfg.ConfigParams.Load("ApiPort")
	r.Run(APIPort.(string)) // 在 0.0.0.0:8080 上监听并服务
}
