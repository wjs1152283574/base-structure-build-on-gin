/*
 * @Author: Casso-Wong
 * @Date: 2021-06-07 11:31:58
 * @Last Modified by:   Casso-Wong
 * @Last Modified time: 2021-06-07 11:31:58
 */
package contxtverify

import (
	"errors"
	dto "goweb/model/dto/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 统一验证gin handler 的 常用操作，例如参数验证等

// CheckPageLimit check page and limit query
func CheckPageLimit(c *gin.Context) (page, limit int, err error) {
	page, _ = strconv.Atoi(c.Query("page"))
	limit, _ = strconv.Atoi(c.Query("limit"))
	if page == 0 || limit == 0 {
		return 0, 0, errors.New("InvalidParam")
	}
	return
}

// CheckAdmin is or not; can be use for admin user checkking
func CheckAdmin(userID uint) error {
	var u dto.User
	u.ID = uint(userID)
	if err := u.Check(); err != nil {
		return err
	}
	if u.Status != 1 { // 异常账号
		return errors.New("locked user")
	}
	if u.Type != 3 && u.Type != 4 {
		return errors.New("PermitionDenid")
	}
	return nil
}

func CheckFront(mobile string) error {
	var u dto.User
	u.Mobile = mobile
	if err := u.Check(); err != nil {
		return err
	}
	if u.Status != 1 { // 异常账号
		return errors.New("locked user")
	}
	return nil
}
