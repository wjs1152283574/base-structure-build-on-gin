package main

import (
	"fmt"
	"goweb/dao/appmysql"
	"goweb/utils/parsecfg"
	_ "goweb/utils/parsecfg"
	_ "goweb/utils/ws"

	"github.com/robfig/cron"
)

func main() {
	fmt.Println(parsecfg.GlobalConfig.Mysql.Write.Host)
	defer appmysql.DB.Close() // 预操作: 关闭数据库连接
	// appuser.MyMigrate()       // 根据给定的结构体进行迁移 迁移的结构体需要在dbmodel中定义
	corns := cron.New() // 定时任务
	corns.Start()
	defer corns.Stop()

}
