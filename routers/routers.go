package routers

import (
	"goweb/service"
	"goweb/utils/cors"
	"goweb/utils/customerjwt"
	"goweb/utils/parsecfg"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// 路由模块

// InitAPIRouter 路由初始化
func InitAPIRouter(basePath string) *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	appDebug, _ := parsecfg.ConfigParams.Load("AppDebug")
	if appDebug.(bool) {
		gin.DisableConsoleColor()
		fileObj, _ := os.OpenFile(".gin.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		gin.DefaultWriter = io.MultiWriter(fileObj)
	}
	// if yml_config.CreateYamlFactory().GetBool("AppDebug") == false {
	// 	//1.将日志写入日志文件
	// 	gin.DisableConsoleColor()
	// 	f, _ := os.Create(variable.BasePath + yml_config.CreateYamlFactory().GetString("Logs.GinLogName"))
	// 	gin.DefaultWriter = io.MultiWriter(f)
	// 	// 2.如果是有nginx前置做代理，基本不需要gin框架记录访问日志，开启下面一行代码，屏蔽上面的三行代码，性能提升 5%
	// 	//gin.SetMode(gin.ReleaseMode)

	// 	router = gin.Default()
	// } else {
	// 	// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
	// 	router = gin.Default()
	// 	pprof.Register(router)
	// }

	//  初始化
	router = gin.Default()
	//根据配置进行设置跨域
	dbCharset, _ := parsecfg.ConfigParams.Load("AllowCrossDomain")
	if dbCharset.(bool) {
		router.Use(cors.Next()) // 启用跨域中间件
	}

	// router.GET("/", func(context *gin.Context) {
	// 	context.String(http.StatusOK, "Api 模块接口 hello word！")
	// })

	// //处理静态资源（不建议gin框架处理静态资源，参见 Public/readme.md 说明 ）
	// router.Static("/public", "./public")             //  定义静态资源路由与实际目录映射关系
	// router.StaticFile("/abcd", "./public/readme.md") // 可以根据文件名绑定需要返回的文件名

	v1 := router.Group("/v1")
	{
		v12 := v1.Group("/auth")
		v12.Use(customerjwt.JWTAuth())
		{
			v12.GET("/user", service.GetUserByName) // 查看用户信息
		}
		v1.POST("/user", service.CreateUser) // 用户注册
		v1.POST("/login", service.UserLogin) // 用户登录
		// v1.GET("/goroutine", func(c *gin.Context) {
		// 	cC := c.Copy() // 在这里使用goroutine需要用Context副本
		// 	// 如果没有固定数量得线程池得话  有可能创建得线程会超出最大承受数量
		// 	go service.GetUserByName(cC)
		// })
	}
	v2 := router.Group("/v2")
	{
		v2.POST("/user", service.CreateUser)
		v2.GET("/user", service.GetUserByName)
	}
	return router
}
