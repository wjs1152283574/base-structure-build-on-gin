package logfmt

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

// Log 自定义日志结构体
type Log struct {
	AppID        string                 `json:"app_id"`        // 预留，分布式判定节点
	Time         string                 `json:"time"`          // 时间
	Method       string                 `json:"method"`        // HTTP METHOD
	Agent        string                 `json:"agent"`         // 客户端代理
	ClientIP     string                 `json:"ip"`            // 来访IP
	Path         string                 `json:"path"`          // 访问路径
	ErrorMessage string                 `json:"error_message"` // 错误信息
	Size         int                    `json:"size"`          // response size
	Keys         map[string]interface{} `json:"key"`
}

// LogFmt  自定义gin框架日志输出
func LogFmt(param gin.LogFormatterParams) string {
	var logs Log
	logs.AppID = "goweb"
	logs.Time = time.Now().Local().Format("2006-01-02 15:04:05")
	logs.Method = param.Method
	logs.Agent = param.Request.UserAgent()
	logs.ClientIP = param.ClientIP
	logs.Path = param.Path
	logs.ErrorMessage = param.ErrorMessage
	logs.Size = param.BodySize
	logs.Keys = param.Keys
	res, _ := json.Marshal(logs)
	return string(res) + "\n"
}
