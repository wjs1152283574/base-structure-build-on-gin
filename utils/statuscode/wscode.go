/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:14:28
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-07-22 12:33:47
 */
package statuscode

// websocket 状态码

// WsResponse 成功返回
var WsResponse = comment{
	Code: 200,
	Msg:  "返回成功！", // 返回成功，一般不用msg
}

// WsAddCon 添加好友请求！
var WsAddCon = comment{
	Code: 204,
	Msg:  "添加好友请求！",
}

// WsNewWilds 有新的心情更新！
var WsNewWilds = comment{
	Code: 205,
	Msg:  "有新的心情更新！",
}

// WsAddConRes 对方同意添加好友！
var WsAddConRes = comment{
	Code: 206,
	Msg:  "我同意了你的好友请求，现在我们可以开始聊天了",
}

// WsChangeGroupLImit 群主修改群员上限！
var WsChangeGroupLImit = comment{
	Code: 207,
	Msg:  "群主修改群员上限！",
}

// WsGroupAir 群公告更新！
var WsGroupAir = comment{
	Code: 208,
	Msg:  "群公告更新！",
}

// WsGroupNotTal 全体禁言！
var WsGroupNotTal = comment{
	Code: 209,
	Msg:  "全体禁言！",
}

// WsGroupCancelNotTal 解除全体禁言！
var WsGroupCancelNotTal = comment{
	Code: 210,
	Msg:  "解除全体禁言！",
}

// WsNeddPromote 进群审核开启！
var WsNeddPromote = comment{
	Code: 211,
	Msg:  "进群审核开启！",
}

// WsNotNeddPromote 进群审核关闭！
var WsNotNeddPromote = comment{
	Code: 212,
	Msg:  "进群审核关闭！",
}

// WsDelGroup 群主已解散群聊！
var WsDelGroup = comment{
	Code: 213,
	Msg:  "群主已解散群聊！",
}

// WsLockGroup 群已被锁定！
var WsLockGroup = comment{
	Code: 214,
	Msg:  "群已被锁定！",
}

// WsCancenLockGroup 群锁定解除！
var WsCancenLockGroup = comment{
	Code: 215,
	Msg:  "群锁定解除！",
}

// WsLockUser 账号已被锁定！
var WsLockUser = comment{
	Code: 216,
	Msg:  "账号已被锁定！",
}

// WsDeleted 已被拉黑！
var WsDeleted = comment{
	Code: 217,
	Msg:  "已被拉黑！",
}

// WsChatGroupCreated 群组创建成功,快邀请小伙伴聊天吧！
var WsChatGroupCreated = comment{
	Code: 218,
	Msg:  "群组创建成功,快邀请小伙伴聊天吧！",
}

// WsChangedChatGroupName 群名称修改！
var WsChangedChatGroupName = comment{
	Code: 219,
	Msg:  "群名称修改！",
}

// WsChangedChatGroupIcon 群图标修改！
var WsChangedChatGroupIcon = comment{
	Code: 220,
	Msg:  "群图标修改！",
}

// WsChangedChatGroupRemark 群内昵称修改！
var WsChangedChatGroupRemark = comment{
	Code: 221,
	Msg:  "群内昵称修改！",
}

// WsJoinChatGroup 入群通知！
var WsJoinChatGroup = comment{
	Code: 222,
	Msg:  "入群通知！",
}

// WsWantJoinChatGroup 入群申请通知！
var WsWantJoinChatGroup = comment{
	Code: 223,
	Msg:  "入群申请通知！",
}

// WshatGroupManagerChange 群管理员变更！
var WshatGroupManagerChange = comment{
	Code: 224,
	Msg:  "群管理员变更！",
}

// WsCancelTalkSome 群员解除禁言！
var WsCancelTalkSome = comment{
	Code: 225,
	Msg:  "群员解除禁言！",
}

// WsCancelTalk 解除指定群员禁言！
var WsCancelTalk = comment{
	Code: 226,
	Msg:  "解除指定群员禁言！",
}

// WsNotTalkSome 群员禁言！
var WsNotTalkSome = comment{
	Code: 227,
	Msg:  "群员禁言！",
}

// WsNotTalk 指定群员禁言！
var WsNotTalk = comment{
	Code: 228,
	Msg:  "指定群员禁言！",
}

// WsNotTalkEach 关闭群员私聊！
var WsNotTalkEach = comment{
	Code: 229,
	Msg:  "关闭群员私聊！",
}

// WsTalkEach 开启群员私聊！
var WsTalkEach = comment{
	Code: 230,
	Msg:  "开启群员私聊！",
}

// WsNotInvi 关闭群员邀请！
var WsNotInvi = comment{
	Code: 231,
	Msg:  "关闭群员邀请！",
}

// WsCanInvi 开启群员邀请！
var WsCanInvi = comment{
	Code: 232,
	Msg:  "开启群员邀请！",
}

// WsChatGroupNotUpload 关闭群员上传！
var WsChatGroupNotUpload = comment{
	Code: 233,
	Msg:  "关闭群员上传！",
}

// WsChatGroupCanUpload 开启群员上传！
var WsChatGroupCanUpload = comment{
	Code: 234,
	Msg:  "开启群员上传！",
}

// WsChatGroupNotNotice 关闭退群通知！
var WsChatGroupNotNotice = comment{
	Code: 235,
	Msg:  "关闭退群通知！",
}

// WsChatGroupNotice 开启退群通知！
var WsChatGroupNotice = comment{
	Code: 236,
	Msg:  "开启退群通知！",
}

// WsLeaveChatGroup 群员退群通知！
var WsLeaveChatGroup = comment{
	Code: 237,
	Msg:  "群员退群通知！",
}

// WsDelChatItem 移除群员通知！
var WsDelChatItem = comment{
	Code: 238,
	Msg:  "移除群员通知！",
}

// WsDeletedChatItem 已被移除群聊！
var WsDeletedChatItem = comment{
	Code: 239,
	Msg:  "已被移除群聊！",
}

// WsDomChanged 配置修改通知！
var WsDomChanged = comment{
	Code: 240,
	Msg:  "配置修改通知！",
}

// WsManagerSending 后台管理员消息！
var WsManagerSending = comment{
	Code: 241,
	Msg:  "后台管理员消息！",
}

// WsOutOfLimit 超出群员限制！
var WsOutOfLimit = comment{
	Code: 242,
	Msg:  "超出群员限制！",
}

// WsCantTalkNow 已开启全体禁言,不可发言！
var WsCantTalkNow = comment{
	Code: 242,
	Msg:  "已开启全体禁言,不可发言！",
}

// WsCantTalkNow2 群已被锁定,不可发言！
var WsCantTalkNow2 = comment{
	Code: 243,
	Msg:  "群已被锁定,不可发言！",
}

// WsChatGroupNotExit 群组不存在！
var WsChatGroupNotExit = comment{
	Code: 244,
	Msg:  "群组不存在！",
}

// WsMuiltyLogin 账号已在其他设备登陆！
var WsMuiltyLogin = comment{
	Code: 245,
	Msg:  "账号已在其他设备登陆！",
}

// WsSendAll 一键群发！
var WsSendAll = comment{
	Code: 246,
	Msg:  "一键群发！",
}

// WsGroupExit 已经存在聊天组，请使用返回组ID直接进行聊天！
var WsGroupExit = comment{
	Code: 248,
	Msg:  "已经存在聊天组，请使用返回组ID直接进行聊天！",
}

// WsInBlackList 已被对方拉黑，不可发送消息！
var WsInBlackList = comment{
	Code: 249,
	Msg:  "已被对方拉黑，不可发送消息！",
}

// WsDoIt 已被对方拉黑，不可发送消息！
var WsDoIt = comment{
	Code: 250,
	Msg:  "您的投诉/举报已处理，被举报人已受到冻结账号处理！",
}

// WsComplaintUselee 已被对方拉黑，不可发送消息！
var WsComplaintUselee = comment{
	Code: 251,
	Msg:  "经证实，您的投诉/举报内容无效，请知悉！",
}

// WsInvaliData 不可识别的消息类型！
var WsInvaliData = comment{
	Code: 252,
	Msg:  "不可识别的消息类型！",
}

// SysNotify 有新的系统通告！
var SysNotify = comment{
	Code: 253,
	Msg:  "有新的系统通告！",
}

// HearBeatResponse 服务端心跳响应！
var HearBeatResponse = comment{
	Code: 254,
	Msg:  "服务端心跳响应！",
}

// BurnAfterRead 阅后即焚消息！
var BurnAfterRead = comment{
	Code: 255,
	Msg:  "阅后即焚消息！",
}

// BurnAfterReadChecked 阅后即焚消息已读！
var BurnAfterReadChecked = comment{
	Code: 256,
	Msg:  "阅后即焚消息已读！",
}

// VideoChat 视频聊天！
var VideoChat = comment{
	Code: 257,
	Msg:  "视频聊天！",
}

// VoiceChat 语音聊天！
var VoiceChat = comment{
	Code: 258,
	Msg:  "语音聊天！",
}

// AcceptVideoChat 接受视频聊天！
var AcceptVideoChat = comment{
	Code: 259,
	Msg:  "接受视频聊天！",
}

// AcceptVoiceChat 接受语音聊天！
var AcceptVoiceChat = comment{
	Code: 260,
	Msg:  "接受语音聊天！",
}

// CutVideoChat 挂断视频聊天！
var CutVideoChat = comment{
	Code: 261,
	Msg:  "挂断视频聊天！",
}

// CutVoiceChat 挂断语音聊天！
var CutVoiceChat = comment{
	Code: 262,
	Msg:  "挂断语音聊天！",
}

// InvalidGroup 聊天组ID格式有误！
var InvalidGroup = comment{
	Code: 263,
	Msg:  "聊天组ID格式有误！",
}

// AssitantSend 官方小助手发送消息！
var AssitantSend = comment{
	Code: 264,
	Msg:  "官方小助手发送消息！",
}

// GroupNotExit 群组不存在！
var GroupNotExit = comment{
	Code: 265,
	Msg:  "群组不存在！",
}

// WsCome 邀请聊天！
var WsCome = comment{
	Code: 503,
	Msg:  "邀请聊天！",
}

// WrongMsg 消息格式有误！
var WrongMsg = comment{
	Code: 500,
	Msg:  "消息格式有误！",
}
