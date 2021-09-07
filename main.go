package main

import (
	"goweb/dao/mysql"
	"goweb/routers"
	"goweb/utils/parsecfg"
	"goweb/utils/timer"
)

func main() {
	defer mysql.DB.Close()                                  // 预操作: 关闭 数据库连接
	defer timer.Conrs.Stop()                                // 预操作: 关闭定时器任务
	routers.InitAPIRouter().Run(parsecfg.GlobalConfig.Port) // 在 main中阻塞监听
}
