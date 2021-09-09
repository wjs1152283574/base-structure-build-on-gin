package routers

import (
	"goweb/service/user"
	"goweb/utils/cors"
	"goweb/utils/customerjwt"
	logfmt "goweb/utils/log"
	"goweb/utils/parsecfg"
	"goweb/utils/ratelimit"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// InitAPIRouter 路由初始化
func InitAPIRouter() *gin.Engine {
	//  初始化
	router := gin.New()
	gin.SetMode("release")

	// 非调试模式（生产模式） 日志写到日志文件
	if parsecfg.GlobalConfig.Debug {
		// gin.DisableConsoleColor()
		fileObj, _ := os.OpenFile("gin.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		fileObjErr, _ := os.OpenFile("gin_err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		gin.DefaultWriter = io.MultiWriter(fileObj)         // 指定writer
		gin.DefaultErrorWriter = io.MultiWriter(fileObjErr) // 指定writer
		router.Use(gin.LoggerWithFormatter(logfmt.LogFmt))  // 自定义日志格式
	}
	router.Use(ratelimit.RateLimit()) // 限流中间件
	router.Use(gin.Recovery())
	//根据配置进行设置跨域
	if parsecfg.GlobalConfig.AllowCrossDomain {
		router.Use(cors.Next()) // 启用跨域中间件
	}

	// 前端
	app := router.Group("/app")
	{
		appAuth := app.Group("/auth")
		appAuth.Use(customerjwt.JWTAuth())
		{
			appAuth.GET("/user/:id", user.GetUser) // 查看用户信息
			appAuth.PUT("/user", user.UserUpd)     // 用户编辑
		}
		app.POST("/user", user.SignUp)  // 用户注册
		app.POST("/login", user.SignIn) // 用户登录
	}

	// 后台管理
	admin := router.Group("/admin")
	{
		admin.POST("/login", user.SignIn) // 用户登录
	}
	return router
}
