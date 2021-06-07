package routers

import (
	userlogi "goweb/service/usrlogi"
	"goweb/utils/cors"
	"goweb/utils/customerjwt"
	"goweb/utils/parsecfg"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// InitAPIRouter 路由初始化
func InitAPIRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if parsecfg.GlobalConfig.Debug {
		gin.DisableConsoleColor()
		fileObj, _ := os.OpenFile(".gin.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		gin.DefaultWriter = io.MultiWriter(fileObj)
	}
	//  初始化
	router = gin.Default()
	//根据配置进行设置跨域
	if parsecfg.GlobalConfig.AllowCrossDomain {
		router.Use(cors.Next()) // 启用跨域中间件
	}
	v1 := router.Group("/v1")
	{
		v12 := v1.Group("/auth")
		v12.Use(customerjwt.JWTAuth())
		{
			v12.GET("/user/:id", userlogi.GetUser) // 查看用户信息
		}
		v1.POST("/user", userlogi.SignUp)  // 用户注册
		v1.POST("/login", userlogi.SignIn) // 用户登录
	}
	v2 := router.Group("/v2")
	{
		v2.POST("/login", userlogi.SignIn) // 用户登录
	}
	return router
}
