/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:14:24
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-07-15 16:40:09
 */
package statuscode

// 同一定义全局返回code; 其他开发人员如有需要可自行添加.

type comment struct {
	Code int
	Msg  string
}

// Success 操作成功
var Success = comment{
	Code: 20000,
	Msg:  "操作成功!",
}

// Faillure 操作失败
var Faillure = comment{
	Code: -1,
	Msg:  "操作失败!",
}

// FailToken 无效otken
var FailToken = comment{
	Code: -2,
	Msg:  "无效otken!",
}

// ExprieToken token过期
var ExprieToken = comment{
	Code: -3,
	Msg:  "token过期!",
}

// InvalidParam 操作失败
var InvalidParam = comment{
	Code: 40000,
	Msg:  "版本过期,请及时更新版本!",
}

// InvalidMsgCode 操作失败
var InvalidMsgCode = comment{
	Code: 40001,
	Msg:  "新用户,请输入正确验证码!",
}

// WrongMsgCode 操作失败
var WrongMsgCode = comment{
	Code: 40002,
	Msg:  "验证码错误!",
}

// NotExit 操作失败
var NotExit = comment{
	Code: 40003,
	Msg:  "用户名或密码错误!",
}

// WrongMoblie 操作失败
var WrongMoblie = comment{
	Code: 40004,
	Msg:  "请输入正确手机号!",
}

// TooSoon 操作失败
var TooSoon = comment{
	Code: 40004,
	Msg:  "操作过于频繁!",
}

// UserNotExit 操作失败
var UserNotExit = comment{
	Code: 40005,
	Msg:  "用户不存在!",
}

// PassInvalid 操作失败
var PassInvalid = comment{
	Code: 40006,
	Msg:  "密码不正确!",
}

// UnlimitFileSize 操作失败
var UnlimitFileSize = comment{
	Code: 40007,
	Msg:  "文件过大!",
}

// UnlimitAdd 操作失败
var UnlimitAdd = comment{
	Code: 40008,
	Msg:  "好友添加数量限制!",
}

// UnlimitAddHis 操作失败
var UnlimitAddHis = comment{
	Code: 40009,
	Msg:  "对方添加好友数量限制!",
}

// AllreadyExit 好友请求已发出
var AllreadyExit = comment{
	Code: 40010,
	Msg:  "添加好友请求已发出!",
}

// NotRecords 好友请求已发出
var NotRecords = comment{
	Code: 40011,
	Msg:  "暂无数据!",
}

// PermitionDenid 好友请求已发出
var PermitionDenid = comment{
	Code: 40012,
	Msg:  "权限不足!",
}

// GroupOutLimit 好友请求已发出
var GroupOutLimit = comment{
	Code: 40013,
	Msg:  "超出群员限制!",
}

// LockedUser 好友请求已发出
var LockedUser = comment{
	Code: 40014,
	Msg:  "已被锁定,无法登陆!",
}

// AlreadyExit 数据已存在
var AlreadyExit = comment{
	Code: 40015,
	Msg:  "数据已存在!",
}

// WrongCaptcha 图片验证失败
var WrongCaptcha = comment{
	Code: 40016,
	Msg:  "图片验证失败!",
}

// AllreadyExit2 好友请求已生成，前往新朋友界面添加
var AllreadyExit2 = comment{
	Code: 40017,
	Msg:  "好友请求已生成，前往新朋友界面添加!",
}

// SendMSgFail 发送短信失败！
var SendMSgFail = comment{
	Code: 40017,
	Msg:  "发送短信失败！",
}

// SendMSgFailNoConfig 发送短信失败,未配置短信！
var SendMSgFailNoConfig = comment{
	Code: 40018,
	Msg:  "发送短信失败,未配置短信！",
}

// BlackyourSelf 不可拉黑自身！
var BlackyourSelf = comment{
	Code: 40019,
	Msg:  "不可拉黑自身！",
}

// IncludingSencitive 包含敏感词语！
var IncludingSencitive = comment{
	Code: 40020,
	Msg:  "包含敏感词语！",
}

// ConnectReful 长连接出错！
var ConnectReful = comment{
	Code: 40021,
	Msg:  "长连接出错！",
}

// WrongAppCode 无法访问！
var WrongAppCode = comment{
	Code: 40022,
	Msg:  "appcode不存在,无法访问！",
}

// PhoneRegistered 号码已被注册
var PhoneRegistered = comment{
	Code: 40023,
	Msg:  "号码已被注册",
}

// IllegalUser 非法用户！
var IllegalUser = comment{
	Code: 40024,
	Msg:  "非法用户！",
}

// CandoThat 解散群组请另外调用解散群组接口
var CandoThat = comment{
	Code: 40025,
	Msg:  "解散群组请另外调用解散群组接口",
}

// CanTChatWild 不可评论
var CanTChatWild = comment{
	Code: 40026,
	Msg:  "不可评论",
}

// GetUserInfoFailed 查询用户详细完整数据失败
var GetUserInfoFailed = comment{
	Code: 40027,
	Msg:  "查询用户详细完整数据失败",
}

// ===================OpenAPI专用状态码====================
// OpenApiURLValuesInvalid URL参数不合法
var OpenApiURLValuesInvalid = comment{
	Code: 40127,
	Msg:  "URL参数不合法",
}

// OpenApiCreateAccountFailed 创建星聊账号失败
var OpenApiCreateAccountFailed = comment{
	Code: 40128,
	Msg:  "创建星聊账号失败",
}

// OpenApiUpdPassFailed 屌用更新登录密码方法失败
var OpenApiUpdPassFailed = comment{
	Code: 40129,
	Msg:  "屌用更新登录密码方法失败",
}

// MakeAccountFailed 屌用个人开户方法失败
var OpenApiMakeAccountFailed = comment{
	Code: 40130,
	Msg:  "调用个人开户方法失败",
}

// OpenApiCompanyAccountFailed 调用企业开户方法失败
var OpenApiCompanyAccountFailed = comment{
	Code: 40131,
	Msg:  "调用企业开户方法失败",
}

// OpenApiBankListFailed 调用查询银行编码方法失败
var OpenApiBankListFailed = comment{
	Code: 40132,
	Msg:  "调用查询银行编码方法失败",
}

// OpenApiMsmCodeFailed 调用获取手机验证码方法失败
var OpenApiMsmCodeFailed = comment{
	Code: 40133,
	Msg:  "调用获取手机验证码方法失败",
}

// OpenApiUploadFailed 调用上传图片方法失败
var OpenApiUploadFailed = comment{
	Code: 40134,
	Msg:  "调用上传图片方法失败",
}

// OpenApiMybankList 调用查询银行账户方法失败
var OpenApiMybankListFailed = comment{
	Code: 40135,
	Msg:  "调用查询银行账户方法失败",
}

// OpenApiBalanceFailed 调用查询账户余额方法失败
var OpenApiBalanceFailed = comment{
	Code: 40136,
	Msg:  "调用查询账户余额方法失败",
}

// OpenApiChargeFailed 调用充值方法失败
var OpenApiChargeFailed = comment{
	Code: 40137,
	Msg:  "调用充值方法失败",
}

// OpenApiTakeCashFailed 调用提现方法失败
var OpenApiTakeCashFailed = comment{
	Code: 40138,
	Msg:  "调用提现方法失败",
}

// OpenApiUpdPayPassFailed 调用修改支付密码方法失败
var OpenApiUpdPayPassFailed = comment{
	Code: 40139,
	Msg:  "调用修改支付密码方法失败",
}

// OpenApiBillsFailed 调用查询账户流水方法失败
var OpenApiBillsFailed = comment{
	Code: 40140,
	Msg:  "调用查询账户流水方法失败",
}

// OpenApiSetKeyFailed 调用设置交易所安全码方法失败
var OpenApiSetKeyFailed = comment{
	Code: 40141,
	Msg:  "调用设置交易所安全码方法失败",
}

// OpenApiVerifyKeyFailed 调用校验交易所安全码方法失败
var OpenApiVerifyKeyFailed = comment{
	Code: 40142,
	Msg:  "调用校验交易所安全码方法失败",
}

// OpenApiRateBalanceFailed 调用查询通用积分余额方法失败
var OpenApiRateBalanceFailed = comment{
	Code: 40143,
	Msg:  "调用查询通用积分余额方法失败",
}

// OpenApiGetUrlFailed 调用获取积分交易中心访问URL方法失败
var OpenApiGetUrlFailed = comment{
	Code: 40144,
	Msg:  "调用获取积分交易中心访问URL方法失败",
}

//
// OpenApiRateBalanceListFailed 调用查询积分余额列表方法失败
var OpenApiRateBalanceListFailed = comment{
	Code: 40145,
	Msg:  "调用查询积分余额列表方法失败",
}

//
// OpenApiMakeOrderFailed 屌用积分充值下单方法失败 19
var OpenApiMakeOrderFailed = comment{
	Code: 40146,
	Msg:  "屌用积分充值下单方法失败",
}

//调用积分充值提交方法失败
// OpenApiChargeSubmitFailed 调用积分充值提交方法失败 20
var OpenApiChargeSubmitFailed = comment{
	Code: 40147,
	Msg:  "调用积分充值提交方法失败",
}

// OpenApiRateCashFailed 调用积分提现方法失败 21
var OpenApiRateCashFailed = comment{
	Code: 40148,
	Msg:  "调用积分提现方法失败",
}

// OpenApiRateTransferOrderFailed 调用积分转账下单方法失败 22
var OpenApiRateTransferOrderFailed = comment{
	Code: 40149,
	Msg:  "调用积分转账下单方法失败",
}

// OpenApiRateTransferOrderSubmitFailed 调用积分转账提交方法失败
var OpenApiRateTransferOrderSubmitFailed = comment{
	Code: 40150,
	Msg:  "调用积分转账提交方法失败",
}

// OpenApiRateBillList 调用积分流水列表方法失败
var OpenApiRateBillList = comment{
	Code: 40151,
	Msg:  "调用积分流水列表方法失败",
}

// OpenApiAccountExist 已开户,不能再次开户
var OpenApiAccountExist = comment{
	Code: 40152,
	Msg:  "已开户，不能再次开户",
}
