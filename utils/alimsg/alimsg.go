/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:13:40
 * @Last Modified by:   Casso-Wong
 * @Last Modified time: 2021-06-05 10:13:40
 */
package alimsg

import (
	"crypto/rand"
	"fmt"
	"goweb/db/sett"
	"io"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// 封装发送短信功能:
// 1. 自动生成 6 位数随机码
// 2. 发送到阿里短信api
// 3. 存入redis
var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

var table2 = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', '_', '@', '#'}

// Code 生成6位数随机码-- string
func Code() (res string) {
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	res = string(b)
	if res[0] == '0' {
		res = Code()
	}
	return
}

// OrderCode 生成15位UID
func OrderCode() (res string) {
	b := make([]byte, 11)
	n, err := io.ReadAtLeast(rand.Reader, b, 11)
	if n != 11 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	res = string(b)
	return
}

// SendMsg 发送短信
func SendMsg(moblie, Code string, need sett.AlimsgNeed) (code int) {
	client, err := dysmsapi.NewClientWithAccessKey(need.MsgLocation, need.MsgKey, need.MsgScretKey)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = need.MsgScheme
	request.PhoneNumbers = moblie
	request.SignName = need.MsgSignName
	// request.SignName = "海南生活"
	request.TemplateCode = need.MsgTemplate
	request.TemplateParam = fmt.Sprintf("{code:%s}", Code)
	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("ali sendmsg response code is %#v\n", response.BaseResponse.GetHttpStatus())
	return response.BaseResponse.GetHttpStatus()
}
