/*
 * @Author: Casso-Wong
 * @Date: 2020-12-24 10:15:26
 * @Description: 可实现分页功能,已经在使用
 */
package response

// 基于gin 自定义返回格式
import (
	"github.com/gin-gonic/gin"
)

// ReturnJSON 参数: http状态码 自定义状态码 提示信息字符串  json数据
func ReturnJSON(Context *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {
	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	// 返回 json数据
	Context.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

// ReturnJSONPage 参数: http状态码 自定义状态码 提示信息字符串  json数据 添加总条数: 可根据前端limit返回对应数据
func ReturnJSONPage(Context *gin.Context, httpCode int, dataCode int, msg string, totals int, data interface{}) {
	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	// 返回 json数据
	Context.JSON(httpCode, gin.H{
		"code":  dataCode,
		"msg":   msg,
		"data":  data,
		"total": totals, // 根据所传limit以及数据库内总条数生成这个pageNums:便于前端分页
	})
}

// ReturnJSONFromString 将json字符窜以标准json格式返回（例如，从redis读取json、格式的字符串，返回给浏览器json格式）
func ReturnJSONFromString(Context *gin.Context, httpCode int, jsonStr string) {
	Context.Header("Content-Type", "application/json; charset=utf-8")
	Context.String(httpCode, jsonStr)
}
