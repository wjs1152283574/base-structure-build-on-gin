package user

import (
	"fmt"
	dto "goweb/model/dto/user"
	entity "goweb/model/entity/user"
	vo "goweb/model/vo/user"
	"goweb/utils/contxtverify"
	"goweb/utils/customerjwt"
	"goweb/utils/passmd5"
	"goweb/utils/response"
	"goweb/utils/statuscode"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SignUp  新建用户
func SignUp(c *gin.Context) {
	var postData entity.SignUpReq
	if err := c.ShouldBind(&postData); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.InvalidParam.Code, statuscode.InvalidParam.Msg, nil)
		return
	}
	var me dto.User
	me.Mobile = postData.Mobile
	if err := me.Check(); err == nil { // 通过 `First` API  查找, 不存在侧会报错
		response.ReturnJSON(c, http.StatusOK, statuscode.AlreadyExit.Code, statuscode.AlreadyExit.Msg, nil)
		return
	}
	me.Name = postData.Name
	me.Pwd = passmd5.Base64Md5(postData.Mobile)
	var res vo.ResUser
	if err := me.Create(&res); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}
	var payLoad = customerjwt.CustomClaims{TimeStr: time.Now().Format("2006-01-02 15:04:05"), Name: res.Name, Password: postData.Pwd}
	token, err := customerjwt.NewJWT().CreateToken(payLoad)
	if err != nil {
		fmt.Println("生成Token失败,可调用登录接口")
	}
	response.ReturnJSON(c, http.StatusOK, statuscode.Suucess.Code, statuscode.Suucess.Msg, map[string]interface{}{
		"infos": res,
		"token": token,
	})
}

// GetUser  获取用户信息
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.InvalidParam.Code, statuscode.InvalidParam.Msg, err)
		return
	}
	var me dto.User
	var res vo.ResUser
	me.ID = uint(id)
	if err := me.Get(&res); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}
	response.ReturnJSON(c, http.StatusOK, statuscode.Suucess.Code, statuscode.Suucess.Msg, res)
}

// SignInReq 登录请求结构
type SignInReq struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

// SignIn 用户登陆 签发token
func SignIn(c *gin.Context) {
	var postData entity.SignUpReq
	if err := c.ShouldBind(&postData); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}
	var payLoad = customerjwt.CustomClaims{TimeStr: time.Now().Format("2006-01-02 15:04:05"), Name: postData.Name, Password: postData.Pwd}
	token, err := customerjwt.NewJWT().CreateToken(payLoad)
	if err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}
	response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, map[string]interface{}{
		"token": token,
	})
}

// UserList 用户列表
func UserList(c *gin.Context) {
	page, limit, err := contxtverify.CheckPageLimit(c)
	if err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.InvalidParam.Code, statuscode.InvalidParam.Msg, err)
		return
	}
	var u dto.User
	var res []vo.AdminUserList
	mobile := c.GetString("mobile")
	if err := contxtverify.CheckAdmin(mobile); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.PermitionDenid.Code, statuscode.PermitionDenid.Msg, err)
		return
	}
	count, err := u.AdminGetList(page, limit, &res)
	if err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}
	response.ReturnJSONPage(c, http.StatusOK, statuscode.PermitionDenid.Code, statuscode.PermitionDenid.Msg, count, res)
}
