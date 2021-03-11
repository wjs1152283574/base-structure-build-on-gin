/*
 * @Description:
 * @Author: Casso-Wong
 */

package alimsg

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// 封装发送短信功能:
// 1. 自动生成 6 位数随机码
// 2. 发送到阿里短信api
// 3. 存入redis

// Code 生成6位数随机码-- string
func Code() (res string) {
	rand.Seed(time.Now().Unix())
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	res = fmt.Sprintf("%06v", rnd.Intn(1000000))
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
