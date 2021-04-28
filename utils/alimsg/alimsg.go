/*
 * @Description:
 * @Author: Casso-Wong
 */

package alimsg

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// 封装发送短信功能:
// 1. 自动生成 6 位数随机码
// 2. 发送到阿里短信api
var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// Code 生成6位数随机码-- string 不能简单的使用当前时间并且截取6位, 当操作的时间近似相同时,返回的验证码也是相同的
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
	return
}

// OrderCode 生成12位订单号
func OrderCode() (res string) {
	b := make([]byte, 12)
	n, err := io.ReadAtLeast(rand.Reader, b, 12)
	if n != 12 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	res = "zh" + string(b)
	if res[0] == '0' {
		res = Code()
	}
	return
}

// SendMsg 发送短信
func SendMsg(moblie, Code string) (code int) {
	client, err := dysmsapi.NewClientWithAccessKey("******", "********", "********")
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = moblie
	request.SignName = "*******" // 隐私
	request.TemplateCode = "******"
	request.TemplateParam = fmt.Sprintf("{code:%s}", Code)
	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("ali sendmsg response code is %#v\n", response.BaseResponse.GetHttpStatus())
	return response.BaseResponse.GetHttpStatus()
}
