package main

import (
	"goweb/dao/mysql"
	"goweb/routers"
	"goweb/utils/parsecfg"
	"goweb/utils/timer"
	tripledes "goweb/utils/tripleDES"
)

var Kafka = "haha"

func init() {
	// TODO: 增加一层bootstrap项目启动库
	// 为了方便测试，只有在项目启动时才应该初始化各项资源。也就是还需要一个bootstrap项目启动库
	// 单元测试就需要引入这些工具包，而这些包包含了一些初始化函数并从配置文件获取数据，项目没运行时配置不会被初始化
	tripledes.GlobalTripleDES.Key = parsecfg.GlobalConfig.TripleDes.Key
	tripledes.GlobalTripleDES.Iv = parsecfg.GlobalConfig.TripleDes.Iv
}

func main() {
	defer mysql.DB.Close()                                  // 预操作: 关闭 数据库连接
	defer timer.Conrs.Stop()                                // 预操作: 关闭定时器任务
	routers.InitAPIRouter().Run(parsecfg.GlobalConfig.Port) // 在 main中阻塞监听
}
