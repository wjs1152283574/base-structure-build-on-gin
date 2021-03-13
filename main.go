package main

import (
	"goweb/dao/appmysql"
	"goweb/routers"
	"goweb/utils/parsecfg"
	"goweb/utils/timer"
)

func main() {
	defer appmysql.DB.Close()                               // 预操作: 关闭数据库连接
	defer timer.Conrs.Stop()                                // 预操作: 关闭定时器任务
	routers.InitAPIRouter().Run(parsecfg.GlobalConfig.Port) // 在 main中阻塞监听
}
