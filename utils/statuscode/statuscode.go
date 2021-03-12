package statuscode

// 同一定义全局返回code; 其他开发人员如有需要可自行添加.

type comment struct {
	Code int
	Msg  string
}

// Suucess 操作成功
var Suucess = comment{
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

// Faillure 操作失败
var InvalidParam = comment{
	Code: 40000,
	Msg:  "参数不可用!",
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
