/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:14:18
 * @Last Modified by:   Casso-Wong
 * @Last Modified time: 2021-06-05 10:14:18
 */
package sendall

// 统一发送给全部用户，后期可能会根据传入app_code 区别对待用户，所以放在这里统一处理

// NotifyAll 通知所有移动端用户更新配置
// func NotifyAll(code int, codemsg string) {
// 	var u dto.User
// 	var res []vo.GetSendAll
// 	u.GetFrontU(&res)
// 	var msg ws.TranstMsg
// 	msg.MsgType = int64(code)
// 	msg.Msg = codemsg
// 	for _, v := range res {
// 		ws.CheckInlineOutlingSend([]string{v.Mobile}, msg, false)
// 	}
// }
