package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"goweb/dao/appredis"
	"goweb/utils/alimsg"
	"goweb/utils/customerjwt"
	"goweb/utils/tool"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

func init() {
	go WebsocketManager.Start()
}

// TranstMsg 客户端 > 服务端
type TranstMsg struct {
	From        string     `json:"from"`
	To          []string   `json:"to"`
	Group       string     `json:"group"`
	GroupName   string     `json:"group_name"`
	GroupIcon   string     `json:"group_icon"` // 群图标
	GroupType   int        `json:"group_type"`
	Msg         string     `json:"msg"`
	Time        *time.Time `json:"time"`
	ContentType string     `json:"content_type"`
	MsgType     int64      `json:"msg_type"`
}

// TranstMsgReverse 服务端 > 客户端
type TranstMsgReverse struct {
	From        map[string]string `json:"from"`
	To          []string          `json:"to"`
	Group       string            `json:"group"`
	GroupName   string            `json:"group_name"`
	GroupIcon   string            `json:"group_icon"` // 群图标
	GroupType   int               `json:"group_type"`
	Msg         string            `json:"msg"`
	Time        *time.Time        `json:"time"`
	ContentType string            `json:"content_type"`
	MsgType     int64             `json:"msg_type"`
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

// Read 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		log.Printf("Read:client [%s] disconnect", c.ID)
		if err := c.Socket.Close(); err != nil {
			// log.Printf("client [%s] 服务端断开连接出错 READ: %s", c.ID, err)
		}
	}()
	for {
		messageType, message, err := c.Socket.ReadMessage() // 从客户端那边发送过来的数据再写入Message通道(后台消息分发是程序直接往Message通道里面直接写入数据)
		if err != nil || messageType == websocket.CloseMessage {
			if err != nil {
				fmt.Println(err, "客户端 主动断开连接")
			} else {
				fmt.Println(err, "服务端**** 主动断开连接")
			}
			return
		}
		if len(WebsocketManager.Group) > 0 {
			var s []string
			for k, _ := range WebsocketManager.Group {
				s = append(s, k)
			}
			log.Printf("在线人数:%d---存活链接:%d--%#v", WebsocketManager.groupCount, WebsocketManager.groupCount, s)
		}
		log.Printf("client [%s] receive message: %s", c.ID, string(message))
		c.Message <- message // 写入channel,等待写入对应长连接通道
	}
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
}

// NewGroupsDedials XXX
func NewGroupsDedials(to []string, nums int) GroupsDedials {
	return GroupsDedials{
		NotTalk:    0,
		CanConn:    0,
		Limit:      50,
		NotName:    0,
		Status:     1,
		NeedAccess: 0,
		Nums:       nums,
		To:         to,
	}
}

// NewTranstMsg 消息体构造函数
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

// Groups 所有用户组
var Groups = make(map[string]GroupsDedials)

// GetInfoByMobile 根据电话获取昵称/头像
func GetInfoByMobile(mobile []string, flag int) (res []map[string]string) {
	item := make(map[string]string)
	if flag == 1 { // from
		nickName, _ := appredis.Get(mobile[0] + "nick_name")
		icon, _ := appredis.Get(mobile[0] + "icon")
		item["nick_name"] = string(nickName)
		item["icon"] = string(icon)
		item["mobile"] = mobile[0]
		res = append(res, item)
	} else { // to
		for _, v := range mobile {
			nickName, _ := appredis.Get(v + "nick_name")
			icon, _ := appredis.Get(v + "icon")
			item["nick_name"] = string(nickName)
			item["icon"] = string(icon)
			item["mobile"] = v
			res = append(res, item)
		}
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

// CheckInlineOutlingSend 功能: 检测在线/离线; 发送数据给指定用户  传入: 需要发送用户切片/原始数据(需要改造为指定 msg_type再传入)
func CheckInlineOutlingSend(mobiles []string, msg TranstMsg, flag bool) {
	conn := appredis.RedisDefaultPool.Get() // 获取redis链接
	defer conn.Close()
	if flag { // 使用一个布尔值来判断是否要存入数据
		message, _ := json.Marshal(msg)
		_ = storeMsgToRedis(message) // 消息存储
	}
	res, er := json.Marshal(TranslateMessage(msg, GetInfoByMobile([]string{msg.From}, 1)[0])) // 改变from 类型 并返回最新数据体格式
	fmt.Print("发送前序列化出错:", er, "发送至列表:", mobiles, "\n")
	for _, v := range mobiles {
		if _, ok := WebsocketManager.Group[v]; ok { // 在线
			msg.To = []string{}                                                              // 暂时先这样, 后期会去掉to 属性 , 因为人多的话很长
			ert := WebsocketManager.Group[v].Socket.WriteMessage(websocket.TextMessage, res) // 直接 发送到指定用户连接通道
			log.Printf("用户%s在线，发送消息%s,err:%#v", v, msg.Msg, ert)
		} else { // 离线  : 按原样存储  又或者: 目标用户登录的事另外一台集群的服务器(则需要把消息存入redis中转并发给目标服务器处理)
			msg.To = []string{v} // 因为数据分发都是发送给 msg.To
			res, _ := json.Marshal(msg)
			log.Printf("用户%s离线，存储消息%s", v, msg.Msg)
			if _, err := conn.Do("RPUSH", v+":history", res); err != nil { // 历史消息格式
				fmt.Println(err, "store msg to redis fail--casso")
			}
		}
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() { // 解析客户端发来的数据,在此分发数据到对应的长连接通道
	defer func() {
		log.Printf("Write:client [%s] disconnect", c.ID)
		// WebsocketManager.UnRegister <- c
		// if err := c.Socket.Close(); err != nil {
		// 	// log.Printf("client [%s] 服务端断开连接出错 WRITE: %s", c.ID, err)
		// }
	}()
	// 优化思路:1. 代码封装拆分(主流业务逻辑清晰) ; 2. 增加Groups存储信息(禁言/审核等flag) 3. 数据接入统一改变frmo ; 4. 数据统一存入redis,方便集群; 5. 定义清晰的msg_type区分
	for {
		select {
		case message, ok := <-c.Message:
			if !ok { // 通道 关闭且无值
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte("write---error close"))
				return
			}
			var msg TranstMsg
			if err := json.Unmarshal(message, &msg); err == nil {
				if msg.Group == "" { // 认为是刚发起聊天 || 添加好友
					if msg.MsgType == 4 { // 为添加好友请求,需将改请求推送给目标用户
						msg.MsgType = 204                          // 修改为与前端约定好的状态码: 204 为添加好友请求
						CheckInlineOutlingSend(msg.To, msg, false) // 发送消息 不包括自己
					} else if msg.MsgType == 1 { // 正常邀请聊天 : 这里看判断一下 本次一次性邀请是否超过 群组会员数量限制
						if len(msg.To) > 50 {
							msg.MsgType = 500 // 修改为与前端约定好的状态码: 500 为超出群聊人数上限
							msg.Msg = "超出群聊人数上限"
							CheckInlineOutlingSend([]string{msg.From}, msg, false)
						} else { // 目前仅 会邀请私聊, 创建群聊会使用http 创建
							gn := msg.From + ":chat" + alimsg.Code() // 生成随机组名
							msg.To = append(msg.To, msg.From)        // 组所有成员: 包括自己
							msg.Group = gn
							_ = StoreCurrenGroups(gn, NewGroupsDedials(msg.To, 2)) // 改造Groups 并存入redis
							_ = storeGroups(msg, false, false)                     // 存储新建分组
							msg.MsgType = 101                                      // 修改为与前端约定好的状态码: 101 为邀请聊天
							// 并且需要把组名发给相关人员-- 之后聊天用组名
							CheckInlineOutlingSend(msg.To, msg, false) // 发送消息 包括自己
						}
					}
				} else {
					if msg.MsgType == 1 { // 正常消息通讯
						if details, err := GetGroupDetailsFromRedis(msg); err == nil { // 从redis中取出本组所有群员
							msg.MsgType = 200
							if details.NotTalk == 1 {
								msg.MsgType = 500 // 修改为与前端约定好的状态码: 500 不存在 群组
								msg.Msg = "全体禁言"
								CheckInlineOutlingSend([]string{msg.From}, msg, false)
							} else if details.Status == 2 {
								msg.MsgType = 500 // 修改为与前端约定好的状态码: 500 不存在 群组
								msg.Msg = "群已被锁定"
								CheckInlineOutlingSend([]string{msg.From}, msg, false)
							} else {
								msg.GroupType = 3                             // 3 是私聊
								CheckInlineOutlingSend(details.To, msg, true) // 发送消息 包括自己:所有群员
							}
						}
					} else {
						// 除了 1 之外都是其他的讯息类型,前端可以自定义(以获得对方是否已阅读/是否正在输入等讯息),建议100以内
						// 基本信息都需要传: groupid groupname msgtype 等
						if details, err := GetGroupDetailsFromRedis(msg); err == nil {
							CheckInlineOutlingSend(details.To, msg, false) // 通知目标
						}
					}
				}
			} else {
				// 数据无法反序列化: 关闭连接
				c.Socket.WriteMessage(websocket.TextMessage, []byte("数据格式有误,服务端主动断开连接"))
				fmt.Println(err, "数据无法反序列化: 关闭连接--casso")
				return
			}
		}
	}
}

// GetGroupDetailsFromRedis 从redis  返回分组详细并根据详细进行一些判断
func GetGroupDetailsFromRedis(msg TranstMsg) (results GroupsDedials, err error) {
	if ok := appredis.Exists(msg.Group); !ok {
		msg.MsgType = 500 // 修改为与前端约定好的状态码: 500 不存在 群组
		msg.Msg = "不存在群组"
		CheckInlineOutlingSend([]string{msg.From}, msg, false)
		err = errors.New("不存在群组")
	} else {
		if res, r := appredis.Get(msg.Group); r != nil {
			err = r
		} else {
			_ = json.Unmarshal(res, &results)
		}
	}
	return
}

// StoreCurrenGroups 将程序运行时的Groups 存入redis中: 取Groups时也需要从redis中取出
func StoreCurrenGroups(gn string, g GroupsDedials) error {
	conn := appredis.RedisDefaultPool.Get() // 获取redis链接
	defer conn.Close()
	fmt.Println(tool.ParseSlicess(g.To))
	g.To = tool.ParseSlicess(g.To)
	res, _ := json.Marshal(g)
	_, err := conn.Do("set", gn, res) // 将新分组存入redis
	// fmt.Println(err, "casocasocaoscoasc")
	return err
}

// 新建分组存入redis
func storeGroups(msg TranstMsg, flag, flag2 bool) error {
	conn := appredis.RedisDefaultPool.Get() // 获取redis链接
	defer conn.Close()
	var stores = make(map[string]interface{})
	stores["group"] = msg.Group
	stores["belong"] = msg.To
	stores["type"] = msg.GroupType
	stores["group_name"] = msg.GroupName // 组名
	stores["group_icon"] = msg.GroupIcon // 图标
	stores["add"] = flag
	stores["del"] = flag2
	r, _ := json.Marshal(stores)
	_, err := conn.Do("LPUSH", "groups", r) // 将新建组存入redis,等待入库mysql
	return err
}

// 消息存入redis: 可进行额外操作
func storeMsgToRedis(msg []byte) error { // 函数执行完毕返回即释放redis连接,如果在Write方法里面创建连接的话 要用户退出长连接才会释放redis连接
	conn := appredis.RedisDefaultPool.Get() // 获取redis链接
	defer conn.Close()
	_, err := conn.Do("LPUSH", "msg", msg)
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
				ert := WebsocketManager.Group[client.ID].Socket.WriteMessage(websocket.TextMessage, []byte("new connect,force disconnected!"))
				client.Socket.Close() //  服务端主动关闭连接
				fmt.Println(ert, "注册时已存在在线相同账号，强制离线！！！")
			}
			manager.Lock.Unlock() // 解锁
			appredis.SetHash("onlineUser", fmt.Sprint(manager.groupCount))
		// 注销
		case client := <-manager.UnRegister:
			log.Printf("client: [%s] **注销**************disconnect****************\n", client.ID)
			manager.Lock.Lock() // 加锁
			appredis.SetArr("underline", client.ID)
			if _, ok := manager.Group[client.ID]; ok {
				close(client.Message) // 在此关闭通道, 接可以同步结束Write 协程 (退出死循环监听)
				// manager.Group[client.ID].Socket.Close()                        // 服务端关闭链接
				delete(manager.Group, client.ID)                               // 删除manager
				manager.groupCount--                                           // 在线人数 --
				appredis.SetHash("onlineUser", fmt.Sprint(manager.groupCount)) // 写入下线人员
			} else {
				// fmt.Println("注销失败----clientID 查找失败")
			}
			manager.Lock.Unlock()
		}
	}
}

// Send 向指定的 client 发送数据
func (manager *Manager) Send(ID string, group string, message []byte) {
	data := &MessageData{
		ID:      ID,
		Group:   group,
		Message: message,
	}
	manager.Message <- data
}

// SendGroup 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group:   group,
		Message: message,
	}
	manager.GroupMessage <- data
}

// SendAll 广播
func (manager *Manager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message: message,
	}
	manager.BroadCastMessage <- data
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
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}

// WebsocketManager Manager 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group:            make(map[string]*Client),
	Register:         make(chan *Client, 128),
	UnRegister:       make(chan *Client, 128),
	GroupMessage:     make(chan *GroupMessageData, 128),
	Message:          make(chan *MessageData, 128),
	BroadCastMessage: make(chan *BroadCastMessageData, 128),
	groupCount:       0,
	clientCount:      0,
}

// WsClient gin 处理 websocket handler
func (manager *Manager) WsClient(ctx *gin.Context) {
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}
	var jwt = customerjwt.NewJWT()
	tokenInfo, tokenErr := jwt.ParseToken([]string{ctx.GetHeader("Sec-WebSocket-Protocol")}[0])
	if tokenErr != nil {
		log.Println("token invaliad")
		return
	}
	conn, conErr := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if conErr != nil {
		log.Printf("websocket connect error: %s", tokenInfo.Name)
		return
	}
	client := &Client{
		ID:      tokenInfo.Name,          // 电话作为唯一的 ID
		Socket:  conn,                    // socket链接
		Message: make(chan []byte, 1024), // 通道内消息缓存大小
	}

	manager.RegisterClient(client)
	go client.Read() // 每次注册 另起协程处理
	go client.Write()
	// time.Sleep(time.Second * 5)
	// 以下需处理用户刚上线时,需要接收离线消息: 先去redis查看是否有自己的离线消息,然后Send
	redisConn := appredis.RedisDefaultPool.Get() // 获取redis链接
	defer redisConn.Close()
	histo, _ := redis.Strings(redisConn.Do("LRANGE", client.ID+":history", 0, -1))
	var msglist []TranstMsg
	for _, value := range histo {
		tempMovie := TranstMsg{}
		err := json.Unmarshal([]byte(value), &tempMovie)
		if err == nil {
			msglist = append(msglist, tempMovie)
		}
	}
	if len(msglist) > 0 { // 该用户存在离线消息
		// fmt.Println(msglist[0].From, msglist[0].Msg, msglist)
		// 遍历每条都发送给用户,并删除这个 redis list
		// 记得处理一下 content-type="history"
		for _, v := range msglist {
			CheckInlineOutlingSend([]string{client.ID}, v, false)
		}
		redisConn.Do("del", client.ID+":history") // 发送完毕,删除
	} else {
		fmt.Println("无历史消息!!!!")
	}
	redisConn.Do("srem", "underline", client.ID) // 上线 移除离线名单
}

// TestSendGroup 测试组广播
func TestSendGroup() {
	for {
		time.Sleep(time.Second * 20)
		WebsocketManager.SendGroup("leffss", []byte("SendGroup message ----"+time.Now().Format("2006-01-02 15:04:05")))
	}
}

// TestSendAll 测试广播
func TestSendAll() {
	for {
		time.Sleep(time.Second * 25)
		res := make(map[string]interface{})
		res["action"] = "send_group_msg"
		res2 := make(map[string]interface{})
		res["params"] = res2
		res2["group_id"] = 340893636
		res2["message"] = "app star chat!!"
		json01, _ := json.Marshal(res)
		WebsocketManager.SendAll([]byte(json01))
		fmt.Println(WebsocketManager.Info())
	}
}
