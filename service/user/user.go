package user

import (
	dto "goweb/model/dto/user"
	entity "goweb/model/entity/user"
	"goweb/utils/contxtverify"
	"goweb/utils/customerjwt"
	"goweb/utils/datacopy"
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
	datacopy.DataCopy(&postData, &me)

	if err := me.Check(); err == nil { // 通过 `First` API  查找, 数据不存在会报错
		response.ReturnJSON(c, http.StatusOK, statuscode.AlreadyExit.Code, statuscode.AlreadyExit.Msg, nil)
		return
	}

	me.Pwd = passmd5.Base64Md5(postData.Mobile)
	res, err := me.Create()
	if err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}

	var payLoad = customerjwt.CustomClaims{TimeStr: time.Now().Format("2006-01-02 15:04:05"), Name: res.Name, Password: postData.Pwd}
	token, err := customerjwt.NewJWT().CreateToken(payLoad)
	if err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.FailToken.Code, statuscode.FailToken.Msg, err)
		return
	}

	response.ReturnJSON(c, http.StatusOK, statuscode.Success.Code, statuscode.Success.Msg, map[string]interface{}{
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
	me.ID = uint(id)
	res, err := me.Get()
	if err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}

	response.ReturnJSON(c, http.StatusOK, statuscode.Success.Code, statuscode.Success.Msg, res)
}

// SignIn 用户登陆 签发token
func SignIn(c *gin.Context) {
	var postData entity.SignInReq
	if err := c.ShouldBind(&postData); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}

	var payLoad = customerjwt.CustomClaims{TimeStr: time.Now().Format("2006-01-02 15:04:05"), Name: postData.Name, Password: postData.Pass}
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

	mobile := c.GetString("mobile")
	if err := contxtverify.CheckAdmin(mobile); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.PermitionDenid.Code, statuscode.PermitionDenid.Msg, err)
		return
	}

	var u dto.User
	res, count, err := u.AdminGetList(page, limit)
	if err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}

	response.ReturnJSONPage(c, http.StatusOK, statuscode.PermitionDenid.Code, statuscode.PermitionDenid.Msg, count, res)
}

// UserUpd 用户编辑
func UserUpd(c *gin.Context) {
	var postData entity.UserUpdReq
	if err := c.ShouldBind(&postData); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}

	var me dto.User
	datacopy.DataCopy(&postData, &me)

	if err := me.Update(); err != nil {
		response.ReturnJSON(c, http.StatusOK, statuscode.Faillure.Code, statuscode.Faillure.Msg, err)
		return
	}

	response.ReturnJSON(c, http.StatusOK, statuscode.Success.Code, statuscode.Success.Msg, nil)
}
