package service

import (
	"goweb/database/mysqlmodel"
	"goweb/utils/customerjwt"
	"goweb/utils/passmd5"
	"goweb/utils/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateUser  新建用户
func CreateUser(c *gin.Context) {
	var addUser mysqlmodel.User
	if err := c.ShouldBindJSON(&addUser); err != nil {
		response.ReturnJSON(c, http.StatusOK, 2001, "invalid params!", nil)
		return
	}
	addUser.Pwd = passmd5.Base64Md5(addUser.Pwd)
	if err := addUser.CreateUser(); err != nil {
		response.ReturnJSON(c, http.StatusOK, 2001, "invalid params!", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "user created!", "username": addUser.Name})
}

// GetUserByName  获取用户信息
func GetUserByName(c *gin.Context) {
	// 接口需要认证token --> 从token中解析出name  --> 根据name查看对应用户数据
	var usr mysqlmodel.User
	// 获取token并解析获得对应用户数据
	usrname := c.GetString("usrname")
	var queryStr = []string{"id", "name", "age", "gender", "mobile", "birthday", "created_at"}
	if res, err := usr.GetUser(usrname, queryStr); err != nil {
		response.ReturnJSON(c, http.StatusNotFound, 2004, "用户信息缺失", nil)
	} else {
		response.ReturnJSON(c, http.StatusOK, 0, "success", res)
	}
}

// UserLogin 用户登陆 签发token
func UserLogin(c *gin.Context) {
	name := c.PostForm("name")
	pwd := c.PostForm("pwd")
	if name != "" && pwd != "" {
		var payLoad = customerjwt.CustomClaims{TimeStr: time.Now().Format("2006-01-02 15:04:05"), Name: name, Password: pwd}
		if token, err := customerjwt.NewJWT().CreateToken(payLoad); err != nil {
			response.ReturnJSON(c, http.StatusBadRequest, 2004, "参数不全或不可用", nil)
		} else {
			data := map[string]string{"token": token}
			response.ReturnJSON(c, http.StatusOK, 1, "success!", data)
		}
	} else {
		response.ReturnJSON(c, http.StatusBadRequest, 2003, "参数不全或不可用!", nil)
	}

}
