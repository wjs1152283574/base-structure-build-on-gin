/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:15:44
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-07-15 18:19:42
 */
package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"goweb/dao/appredis"
	"goweb/utils/alimsg"
	"goweb/utils/customerjwt"
	"goweb/utils/parsecfg"
	"goweb/utils/response"
	"goweb/utils/sencekw"
	"goweb/utils/statuscode"
	"goweb/utils/tool"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func init() {
	go WebsocketManager.Start()
}

// WebsocketManager Manager 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group:       make(map[string]*Client),
	Register:    make(chan *Client, 128),
	UnRegister:  make(chan *Client, 128),
	Message:     make(chan *MessageData, 128),
	groupCount:  0,
	clientCount: 0,
}

// TranstMsg 客户端 > 服务端
type TranstMsg struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Group       string            `json:"group"`
	GroupName   string            `json:"group_name"`
	GroupIcon   string            `json:"group_icon"` // 群图标
	GroupType   int               `json:"group_type"`
	Msg         string            `json:"msg"`
	Time        *time.Time        `json:"time"`
	ContentType string            `json:"content_type"`
	MsgType     int64             `json:"msg_type"`
	ToInfos     map[string]string `json:"to_infos"`
}

// TranstMsgReverse 服务端 > 客户端
type TranstMsgReverse struct {
	From        map[string]string `json:"from"`         // 发送者
	To          []string          `json:"to"`           // 接受者电话（ID）可能是一个，也可能有很多接收者
	Group       string            `json:"group"`        // 消息组
	GroupName   string            `json:"group_name"`   // 组名
	GroupIcon   string            `json:"group_icon"`   // 群图标
	GroupType   int               `json:"group_type"`   // 组类型 单聊群聊
	Msg         string            `json:"msg"`          // 消息内容
	Time        *time.Time        `json:"time"`         // 时间
	ContentType string            `json:"content_type"` // 内容类型 视频 音频 文字
	MsgType     int64             `json:"msg_type"`     // 消息类型 200 正常
	ToInfos     map[string]string `json:"to_infos"`     // 接受者详细信息
}

// Manager 所有 websocket 信息
type Manager struct {
	Group                   map[string]*Client // {xx:{cc:client,zz:client,ss:client}}
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *Client
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
}

// Client 单个 websocket 信息
type Client struct {
	ID      string
	Socket  *websocket.Conn
	Message chan []byte
}

// MessageData 单个发送数据信息
type MessageData struct {
	ID, Group string
	Message   []byte
}

// GroupMessageData 组广播数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// BroadCastMessageData 广播发送数据信息
type BroadCastMessageData struct {
	Message []byte
}

// HeartBeatReq 心跳请求
type HeartBeatReq struct {
	Pong int `json:"pong"`
}

// HeartBeatRes 心跳响应
type HeartBeatRes struct {
	Ping int `json:"ping"`
}

// GroupsDedials redis 中组的详细信息
type GroupsDedials struct {
	NotTalk    int      `json:"not_talk"`    // 是否全体禁言 0 否 1 禁
	CanConn    int      `json:"can_conn"`    // 0 私聊   1 不可私聊
	Limit      int      `json:"limit"`       // 群员上限:默认50
	NotName    int      `json:"not_name"`    // 0 允许匿名  1 不允许
	Status     int      `json:"status"`      // 1 正常  2 锁定
	NeedAccess int      `json:"need_access"` // 0 无需审核  1 需审核
	Nums       int      `json:"nums"`        // 群员总数
	To         []string `json:"to"`          // 群员电话列表
	GroupID    string   `json:"group_id"`    // 组ID
	GroupType  int      `json:"group_type"`  // 1 单  2 群
	Name       string   `json:"name"`        // 群组名称（群聊） 这些固定信息不应该每次发消息都传输
	Icon       string   `json:"icon"`        // 群组头像（群聊）
}

// Read 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c // 仅在此处触发注销
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] 服务端断开连接出错 READ: %s，表示socket在WRITE已被关闭\n", c.ID, err)
		}
	}()
	for {
		messageType, message, err := c.Socket.ReadMessage() // 从客户端那边发送过来的数据再写入Message通道(后台消息分发是程序直接往Message通道里面直接写入数据)
		if err != nil || messageType == websocket.CloseMessage {
			return
		}
		if messageType == websocket.PingMessage { // 心跳消息,需马上回复
			var hb HeartBeatRes
			hb.Ping = int(time.Now().Unix())
			hbs, _ := json.Marshal(hb)
			c.Socket.WriteMessage(websocket.TextMessage, hbs)
		} else {
			log.Printf("client [%s] receive message: %s\n", c.ID, string(message))
			c.Message <- message // 写入channel,等待写入对应长连接通道
		}
		log.Printf("在线人数:%d\n", WebsocketManager.LenGroup())
	}
}

// Write 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() { // 解析客户端发来的数据,在此分发数据到对应的长连接通道
	defer func() {
		_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte("socket closed by WRITE "))
	}()
	for {
		message, ok := <-c.Message
		if !ok { // 通道 关闭且无值
			return
		}
		var msg TranstMsg
		if err := json.Unmarshal(message, &msg); err == nil {
			if msg.Msg != "" { // 空消息不需处理-- 在这里统一进行敏感词过滤，可以添加全局开关控制是否过滤
				msg.Msg = sencekw.GetSencitiveResponse(msg.Msg)
			}
			if msg.Group == "" { // 认为是刚发起聊天
				gid, filtererr := FilterGroup(msg.From, msg.To[0])
				if filtererr != nil || gid == "" { // 生成新组
					msg.Group = NewGroupID(msg.From) // 生成随机组名
				} else {
					msg.Group = SplitGroup(gid) // 使用已存在组名  -- 这里需要切掉:to
				}
				if msg.MsgType == int64(statuscode.WsCome.Code) || msg.MsgType == int64(statuscode.WsResponse.Code) || msg.MsgType == int64(statuscode.VideoChat.Code) || msg.MsgType == int64(statuscode.VoiceChat.Code) {
					msg.To = append(msg.To, msg.From)                                                   // 组所有成员: 包括自己
					CheckInlineOutlingSend(msg.To, msg, false)                                          // 通知目标:只发给拉起聊天人(通知双方)
					_ = StoreCurrenGroups(msg.Group, NewGroupsDedials(msg.To, 2, 1, msg.Group, "", "")) // 改造Groups 并存入redis
				} else {
					res, _ := json.Marshal(NewSendMsg("system", statuscode.WsInvaliData.Msg, "", int64(statuscode.WsInvaliData.Code)))
					c.Socket.WriteMessage(websocket.TextMessage, res) // 发送反馈消息CheckInlineOutlingSend([]string{msg.From}, msg, false) // 通知目标
					continue
				}
			} else {
				if !VerifyGroup(msg.Group) { // 统一验证组格式
					res, _ := json.Marshal(NewSendMsg("system", statuscode.InvalidGroup.Msg, "", int64(statuscode.InvalidGroup.Code)))
					c.Socket.WriteMessage(websocket.TextMessage, res) // 发送反馈消息
					continue
				}
				if msg.GroupType == 1 { // 单聊分支:包含黑名单检测等
					SingleChat(msg)
				}
				if msg.GroupType == 2 { // 群聊分支：包含禁言等
					GroupChat(msg, c)
				}
			}
		} else {
			// 数据无法反序列化: 退出goroutine关闭连接
			res, _ := json.Marshal(NewSendMsg("system", statuscode.WrongMsg.Msg, "", int64(statuscode.WrongMsg.Code)))
			c.Socket.WriteMessage(websocket.TextMessage, res) // 发送反馈消息
			return
		}
	}
}

// CheckInlineOutlingSend 功能: 检测在线/离线; 发送数据给指定用户  传入: 需要发送用户切片/原始数据(需要改造为指定 msg_type再传入)
func CheckInlineOutlingSend(mobiles []string, msg TranstMsg, flag bool) {
	if flag && msg.Msg != "" { // 使用一个布尔值来判断是否要存入数据
		message, _ := json.Marshal(msg)
		appredis.Lpush("msg", message) // 由于客户端已存储本地数据，此处并不要求零丢失，暂时吞掉此错误，后期可优化
	}
	gtype, _ := appredis.Hget(msg.Group, "gtype")
	if gtype == "2" { // 群聊无需toInfos，只需要from具体信息
		gicon, _ := appredis.Hget(msg.Group, "icon")
		gname, _ := appredis.Hget(msg.Group, "name")
		msg.GroupName = string(gname)
		msg.GroupIcon = string(gicon)
		RealSend(msg, mobiles)
	} else if msg.GroupType == 3 { // 系统通知，本来是走到群聊逻辑里面的，现在分离出来
		RealSend(msg, mobiles)
	} else {
		SingelMsg(mobiles, msg)
	}
}

// SingleChat 处理单聊
func SingleChat(msg TranstMsg) {
	listMe, _ := appredis.GetList(msg.From + ":balck")
	list, _ := appredis.GetList(msg.To[0] + ":balck")
	list = append(list, listMe...)
	if ok, _ := tool.InSlice(msg.From, list); ok && len(msg.To) == 1 { // 黑名单检测(仅用于单聊)
		msg.MsgType = int64(statuscode.WsInBlackList.Code)
		msg.Msg = statuscode.WsInBlackList.Msg
		CheckInlineOutlingSend([]string{msg.From}, msg, false)
	} else { // 正常消息通讯
		if ok := appredis.Exists(msg.Group); ok { // 是否存在组
			toList, _ := appredis.Smembers(msg.Group + ":to")
			if msg.From == toList[0] && len(toList) == 1 { // 处理自己跟自己聊天时只发一次消息
				CheckInlineOutlingSend([]string{msg.From}, msg, true)
			} else {
				CheckInlineOutlingSend(toList, msg, true) // 发送消息 包括自己:所有群员
			}
		} else { // 群组不存在，则直接重新拉起群组(为了适应版本更新，可能会存在聊天组被覆盖的问题),需要区分单聊群聊（找不到群组大概率是群聊，去数据库找，目前先默认为单聊）
			CheckInlineOutlingSend([]string{msg.From}, msg, true)                               // 通知目标:只发给拉起聊天人
			msg.To = append(msg.To, msg.From)                                                   // 组所有成员: 包括自己
			_ = StoreCurrenGroups(msg.Group, NewGroupsDedials(msg.To, 2, 1, msg.Group, "", "")) // 改造Groups 并存入redis
		}
	}
}

// GroupChat 处理群聊
func GroupChat(msg TranstMsg, c *Client) {
	if ok := appredis.Exists(msg.Group); ok { // 是否存在组
		// msg.MsgType = int64(statuscode.WsResponse.Code)
		toList, _ := appredis.Smembers(msg.Group + ":to")
		notTalk, _ := appredis.Hget(msg.Group, "notTalk")
		status, _ := appredis.Hget(msg.Group, "status")
		if notTalk == "1" { // 禁言
			msg.MsgType = int64(statuscode.WsCantTalkNow.Code)
			msg.Msg = statuscode.WsCantTalkNow.Msg
			CheckInlineOutlingSend([]string{msg.From}, msg, false)
		} else if status == "2" { // 锁定
			msg.MsgType = int64(statuscode.WsCantTalkNow2.Code)
			msg.Msg = statuscode.WsCantTalkNow2.Msg
			CheckInlineOutlingSend([]string{msg.From}, msg, false)
		} else { // 正常发送
			CheckInlineOutlingSend(toList, msg, true)
		}
	} else { // 群组不存在
		res, _ := json.Marshal(NewSendMsg("system", statuscode.GroupNotExit.Msg, "", int64(statuscode.GroupNotExit.Code)))
		c.Socket.WriteMessage(websocket.TextMessage, res) // 发送反馈消息
	}
}

// VerifyGroup 验证组id格式,可以添加额外验证
func VerifyGroup(gid string) bool {
	ll := strings.Split(gid, ":")
	return len(ll) == 2
}

// SplitGroup 截掉:to
func SplitGroup(str string) string {
	ll := strings.Split(str, ":")
	if len(ll) != 3 {
		return ""
	}
	return strings.Join(ll[0:2], ":")
}

// NewGroupsDedials 生成新组，应该在组信息中区分单聊群聊，不应该在每条消息单独区分
func NewGroupsDedials(to []string, nums, gtype int, gid, name, icon string) GroupsDedials {
	return GroupsDedials{
		NotTalk:    0,
		CanConn:    0,
		Limit:      50,
		NotName:    0,
		Status:     1,
		NeedAccess: 0,
		GroupID:    gid,
		Nums:       nums,
		To:         to,
		GroupType:  gtype,
		Name:       name,
		Icon:       icon,
	}
}

// NewTranstMsg 消息体构造函数,正常通讯消息，后面考虑去掉
func NewTranstMsg(from, gid, gname, gicon, msg, ctype string, to []string, msgtype int) *TranstMsg {
	t := time.Now()
	return &TranstMsg{
		From:        from,
		To:          to,
		Group:       gid,
		GroupName:   gname,
		GroupIcon:   gicon,
		GroupType:   1,
		Msg:         msg,
		Time:        &t,
		ContentType: ctype,
		MsgType:     int64(msgtype),
	}
}

// NewSendMsg 统一生成发送消息格式,只要函数能通过编译，即不可能返回错误，无错误函数，由于外部也需要调用，大驼峰
func NewSendMsg(from, msg, gid string, msgtype int64) TranstMsg {
	t := time.Now()
	return TranstMsg{
		From:    from,
		Msg:     msg,
		Time:    &t,
		MsgType: msgtype,
		Group:   gid,
	}
}

// GetInfoByMobile 根据电话获取昵称/头像: 需要优化，使用MGET来减少大量redis连接，并把to分离出来
func GetInfoByMobile(mobile []string, flag int) (res []map[string]string) {
	item := make(map[string]string)
	if flag == 1 { // from
		var icon, nickName string
		appredis.Mget([]interface{}{mobile[0] + ":icon", mobile[0] + ":nick_name"}, &icon, &nickName)
		item["nick_name"] = nickName
		item["icon"] = icon
		item["mobile"] = mobile[0]
		res = append(res, item)
	}
	return
}

// TranslateMessage 数据格式转换
func TranslateMessage(msg TranstMsg, from map[string]string) (res TranstMsgReverse) {
	res.To = msg.To
	res.Group = msg.Group
	res.GroupName = msg.GroupName
	res.GroupType = msg.GroupType
	res.Msg = msg.Msg
	res.Time = msg.Time
	res.ContentType = msg.ContentType
	res.MsgType = msg.MsgType
	res.GroupIcon = msg.GroupIcon
	res.From = from
	return
}

// RealSend 最终发送,目前用于群聊以及系统通道，后期可能需要加以区分
func RealSend(msg TranstMsg, mobiles []string) {
	res, _ := json.Marshal(TranslateMessage(msg, GetInfoByMobile([]string{msg.From}, 1)[0])) // 改变from 类型 并返回最新数据体格式
	for _, v := range mobiles {
		if _, ok := WebsocketManager.Group[v]; ok { // 在线
			msg.To = []string{}                                                              // 暂时先这样, 后期会去掉to 属性 , 因为人多的话很长
			ert := WebsocketManager.Group[v].Socket.WriteMessage(websocket.TextMessage, res) // 直接 发送到指定用户连接通道
			log.Printf("用户%s在线，发送消息:%#v,err:%#v\n", v, msg, ert)
		} else { // 离线  : 按原样存储  又或者: 目标用户登录的事另外一台集群的服务器(则需要把消息存入redis中转并发给目标服务器处理)
			msg.To = []string{v} // 因为数据分发都是发送给 msg.To
			res, _ := json.Marshal(msg)
			KafkaProducer(v, &res)
		}
	}
}

// KafkaProducer kafka 发送消息
func KafkaProducer(mobile string, res *[]byte) {
	kt, err := appredis.Get("onlineUser:" + mobile)
	if err == nil && string(kt) != "" { // 存在于另外的节点
		appredis.Rpush(mobile+":transport", *res) // 在这里存入中转消息
		// 接下来往kt(Topic)中发送 v(mobile)
		// 另外节点消费消息，从redis取出（并删除）源消息并发送
	} else { // 真正离线
		log.Printf("用户%s离线，存储消息\n", mobile)
		appredis.Rpush(mobile+":history", *res)
	}
}

// SingelMsg 单聊信息发送专用:判断ToInfos始终为对方的信息
func SingelMsg(mobiles []string, msg TranstMsg) {
	msg2 := TranslateMessage(msg, GetInfoByMobile([]string{msg.From}, 1)[0])
	if len(mobiles) == 2 { // 跟别人聊
		toinfos1 := GetInfoByMobile([]string{mobiles[0]}, 1)[0]
		toinfos2 := GetInfoByMobile([]string{mobiles[1]}, 1)[0]
		for i, v := range mobiles {
			if msg.From == v && i == 0 {
				msg2.ToInfos = toinfos2
			} else if msg.From == v && i == 1 {
				msg2.ToInfos = toinfos1
			}
			res, _ := json.Marshal(msg2)
			if _, ok := WebsocketManager.Group[v]; ok {
				log.Println("单聊分发TOinfos:", v, "消息:", msg2)
				if err := WebsocketManager.Group[v].Socket.WriteMessage(websocket.TextMessage, res); err != nil { // 直接 发送到指定用户连接通道
					msg.To = []string{v} // 因为数据分发都是发送给 msg.To
					res, _ := json.Marshal(msg)
					log.Printf("用户%s离线-----发送错误，存储消息:%#v,err:%#v\n", v, msg2, err)
					appredis.Rpush(v+":history", res)
				}
			} else {
				msg.To = []string{v} // 因为数据分发都是发送给 msg.To
				res, _ := json.Marshal(msg)
				log.Printf("用户%s离线，存储消息:%#v\n", v, msg2)
				appredis.Rpush(v+":history", res)
			}
		}
	} else if len(mobiles) == 1 { // 跟自己聊：肯定是在线状态
		// 也有可能是历史消息，长度也是1
		if mobiles[0] == msg.From {
			msg2.ToInfos = GetInfoByMobile([]string{mobiles[0]}, 1)[0]
			res, _ := json.Marshal(msg2)
			log.Println("单聊一人份消息self:", msg2)
			WebsocketManager.Group[mobiles[0]].Socket.WriteMessage(websocket.TextMessage, res) // 直接 发送到指定用户连接通道
		} else { // 不是发给自己
			log.Println("单聊一人份消息other:", msg2)
			msg2.ToInfos = GetInfoByMobile([]string{msg.From}, 1)[0] // 发送人为toinfos
			res, _ := json.Marshal(msg2)
			WebsocketManager.Group[mobiles[0]].Socket.WriteMessage(websocket.TextMessage, res) // 直接 发送到指定用户连接通道
		}
	}
}

// NewGroupID return group ID
func NewGroupID(mobile string) string {
	return mobile + ":chat" + alimsg.Code()
}

// FilterGroup 查询是否两人已存在聊天组
func FilterGroup(from, to string) (gid string, err error) {
	// 1. KEYS *MOBILE:CHAT*:TO*  获取所有类似key(双方都需要查)
	// 2. 遍历 Smembers 获取每个组并判断组员是否跟这两人重合（根据本次遍历的key切掉:to得到组号并返回）
	fg, _ := appredis.GetLIkeTo(from + ":chat*:to") // from 相似组
	tg, _ := appredis.GetLIkeTo(to + ":chat*:to")   // to 相似组
	if gid, err = FilterGroupSingel(fg, from, to); err == nil {
		return
	}
	if gid, err = FilterGroupSingel(tg, from, to); err == nil {
		return
	}
	return "", errors.New("没有查到已存在群组")
}

// FilterGroupSingel 查询单人名下是否已存在组
func FilterGroupSingel(b []string, from, to string) (gid string, err error) {
	for _, v := range b {
		reply, errs := appredis.Smembers(v)
		if errs != nil {
			fmt.Println("查重时失败,", v)
		}
		if len(reply) == 2 {
			flag1, _ := tool.InSlice(from, reply)
			flag2, _ := tool.InSlice(to, reply)
			if flag1 && flag2 {
				return v, nil
			}
		} else {
			fmt.Println("组员数量不对", reply)
		}
	}
	return "", errors.New("没有查到已存在群组")
}

// StoreCurrenGroups 将程序运行时的Groups 存入redis中: 取Groups时也需要从redis中取出,由于group内数据经常改动，改用hash类型
func StoreCurrenGroups(gn string, g GroupsDedials) error {
	appredis.PipeLineHset(gn, map[string]interface{}{ // 存入基本信息
		"notTalk":    g.NotTalk,
		"canConn":    g.CanConn,
		"limit":      g.Limit,
		"notName":    g.NotName,
		"status":     g.Status,
		"needAccess": g.NeedAccess,
		"nums":       g.Nums,
		"groupID":    g.GroupID,
		"gtype":      g.GroupType,
		"name":       g.Name,
		"icon":       g.Icon,
	})
	var final = []string{gn + ":to"}
	final = append(final, g.To...)
	err := appredis.Sadd(final...) // 组员使用集合存储（去重）
	fmt.Println(err, "SADD")
	return err
}

// Start 启动 websocket 管理器
func (manager *Manager) Start() {
	log.Printf("websocket manage start")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			log.Printf("client: [%s] **注册********************connect****************\n", client.ID)
			manager.Lock.Lock()                         // 加锁
			if _, ok := manager.Group[client.ID]; !ok { // 如果是新的分组,则创建该新分组
				manager.Group[client.ID] = client // 创建一个总管(manager)级别的分组, 但此时分组内部还是没有对应数据,只是make了内存地址
				manager.groupCount++              // 分组即链接
			} else { // 已有链接存在
				res, _ := json.Marshal(NewSendMsg("system", statuscode.WsMuiltyLogin.Msg, "", int64(statuscode.WsMuiltyLogin.Code)))
				ert := WebsocketManager.Group[client.ID].Socket.WriteMessage(websocket.TextMessage, res)
				client.Socket.Close() //  服务端主动关闭连接 ,前端停止重新连接请求
				log.Println(ert, "注册时已存在在线相同账号，强制离线！！！")
			}
			manager.Lock.Unlock()                                                        // 解锁
			appredis.SetHash("onlineUser:"+client.ID, parsecfg.GlobalConfig.Kafka.Topic) // 维护在线用户，value为kafka topic
		// 注销
		case client := <-manager.UnRegister:
			log.Printf("client: [%s] **注销*******disconnect********\n", client.ID)

			if _, ok := manager.Group[client.ID]; ok {
				close(client.Message)            // 在此关闭通道, 可以同步结束Write 协程 (退出死循环监听)
				manager.Lock.Lock()              // 加锁
				delete(manager.Group, client.ID) // 删除manager
				manager.groupCount--             // 在线人数 --
				manager.Lock.Unlock()
			}
			appredis.SetArr("underline", client.ID)
			appredis.Delete("onlineUser:" + client.ID) // 下线--移除
		}
	}
}

// RegisterClient 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// UnRegisterClient 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// LenGroup 当前组个数
func (manager *Manager) LenGroup() uint {
	return manager.groupCount
}

// LenClient 当前连接个数
func (manager *Manager) LenClient() uint {
	return manager.clientCount
}

// Info 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	return managerInfo
}

// WsClient gin 处理 websocket handler
func (manager *Manager) WsClient(ctx *gin.Context) {
	// 1、拿到token
	var jwt = customerjwt.NewJWT()
	tokenInfo, tokenErr := jwt.ParseToken([]string{ctx.GetHeader("Sec-WebSocket-Protocol")}[0])
	if tokenErr != nil || []string{ctx.GetHeader("Sec-WebSocket-Protocol")}[0] == "" { // 如果错误不为空或者拿到的token是空字符串
		fmt.Println("token 不合法，长连接失败返回！本次token:", []string{ctx.GetHeader("Sec-WebSocket-Protocol")}[0])
	} else {
		upGrader := websocket.Upgrader{
			// cross origin domain
			// 检验请求头信息：The application is responsible for checking the Origin header before calling the Upgrade function
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			// 处理 Sec-WebSocket-Protocol Header
			Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
		}
		conn, conErr := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if conErr != nil {
			response.ReturnJSON(ctx, http.StatusOK, statuscode.ConnectReful.Code, statuscode.ConnectReful.Msg, nil)
			return
		}
		client := &Client{
			ID:      tokenInfo.Name,          // 电话作为唯一的 ID
			Socket:  conn,                    // socket链接
			Message: make(chan []byte, 1024), // 通道内消息缓存大小
		}
		manager.RegisterClient(client) // 注册到localCahe,后期需要集群到不同节点，预测利用redis-cluster 来作为调度中心（注册时：用户标识+节点标识）
		go client.Read()               // 阻塞处理socket消息通道读取
		go client.Write()              // 阻塞处理socket消息通道写入
		go HistoryMsgs(client.ID)      // 异步处理历史消息
		appredis.Dellist("underline", client.ID)
	}
}

// HistoryMsgs deal with history messages
func HistoryMsgs(ClientID string) {
	history, redisErr := appredis.GetList(ClientID + ":history")
	if redisErr != nil {
		fmt.Println("历史消息获取出错!!!!") // 尝试panic ，输入日志（日志是顶层统一处理）
		return
	}
	var msglist []TranstMsg
	for _, value := range history {
		tempMovie := TranstMsg{}
		err := json.Unmarshal([]byte(value), &tempMovie)
		if err == nil {
			msglist = append(msglist, tempMovie)
		}
	}
	if len(msglist) == 0 {
		fmt.Println("无历史消息!!!!")
		return
	}
	// 剔除系统配置更改消息：只返回一条即可 meg_type:240
	var sendList []TranstMsg
	var systemNotify []TranstMsg
	for _, v := range msglist {
		if v.MsgType == 240 { // 此处过滤掉系统配置更改通知，此通知没有特殊含义，只需一次
			systemNotify = append(systemNotify, v)
		} else {
			sendList = append(sendList, v)
		}
	}
	if len(systemNotify) > 0 {
		sendList = append(sendList, systemNotify[0]) // 只发一条，其余忽略
	}
	for _, v := range sendList {
		CheckInlineOutlingSend([]string{ClientID}, v, false) // 无需存储
	}
	appredis.Delete(ClientID + ":history") // 处理完毕，删除相关历史消息
}
