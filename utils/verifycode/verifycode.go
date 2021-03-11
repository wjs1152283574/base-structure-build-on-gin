/*
 * @Author: Casso
 * @Date: 2021-03-03 18:27:32
 * @LastEditTime: 2021-03-04 17:15:54
 * @LastEditors: Please set LastEditors
 * @Description: 图片验证码
 * @FilePath: /githubStarChat/starChat/utils/verifycode/verifycode.go
 */
package verifycode

import (
	"fmt"

	"github.com/mojocn/base64Captcha"
)

// Store xx
var Store = base64Captcha.DefaultMemStore

// GetCaptcha 获取验证码
func GetCaptcha() (string, string) {
	// 生成默认数字
	driver := base64Captcha.DefaultDriverDigit
	// 生成base64图片
	c := base64Captcha.NewCaptcha(driver, Store)

	// 获取
	id, b64s, err := c.Generate()
	if err != nil {
		fmt.Println("Register GetCaptchaPhoto get base64Captcha has err:", err)
		return "", ""
	}
	return id, b64s
}

// Verify 验证验证码
func Verify(id string, val string) bool {
	if id == "" || val == "" {
		return false
	}
	// 同时在内存清理掉这个图片
	return Store.Verify(id, val, true)
}
