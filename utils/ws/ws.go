/*
 * @Description:聊天业务相关 websocket模块
 * 实现功能:
	1. 发起/加入 群聊/单聊
* 记录:
	1. 相关数据直接存入redis,方便后期集群部署;
	2. 每一条消息都需要进行序列化(分发)/反序列化操作(解析处理),目前没有想到可行的 替代方案;
	3. 分发实现细节封装,保持主线流程清晰易读;
	4. 集群时如果用户连接的不是同一台服务器则,在线也不能互相发送消息: 需要用到redis作为中转站把消息发送至对应服务器进行处理
 * @Author: Casso-Wong
*/

package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	appredis "goweb/dao/daoredis"
	"goweb/utils/customerjwt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

// TranstMsg 客户端 > 服务端
type TranstMsg struct {
	From        string     `json:"from"`
	To          []string   `json:"to"`
	Group       string     `json:"group"`
	GroupName   string     `json:"group_name"`
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
	GroupType   int               `json:"group_type"`
	Msg         string            `json:"msg"`
	Time        *time.Time        `json:"time"`
	ContentType string            `json:"content_type"`
	MsgType     int64             `json:"msg_type"`
}

// Manager 所有 websocket 信息
type Manager struct {
	Group                   map[string]map[string]*Client
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *Client
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
}

// Client 单个 websocket 信息
type Client struct {
	ID, Group string
	Socket    *websocket.Conn
	Message   chan []byte
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
		log.Printf("client [%s] disconnect", c.ID)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.ID, err)
		}
	}()
	for {
		messageType, message, err := c.Socket.ReadMessage() // 从客户端那边发送过来的数据再写入Message通道(后台消息分发是程序直接往Message通道里面直接写入数据)
		if err != nil || messageType == websocket.CloseMessage {
			break
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
	To         []string `json:"to"`          // 群员电话列表
}

// newGroupsDedials 每个分组的详细信息构造函数
func newGroupsDedials(to []string) GroupsDedials {
	return GroupsDedials{
		NotTalk:    0,
		CanConn:    0,
		Limit:      50,
		NotName:    0,
		Status:     1,
		NeedAccess: 0,
		To:         to,
	}
}

// Groups 所有用户组: 优化后已经把分组信息存入redis,已经不用本地存储Groups了,如果确定不需要集群的话使用本地存储会更快(不需要连接redis取值)
// var Groups = make(map[string]GroupsDedials)

// GetInfoByMobile 根据电话获取昵称/头像(在http接口中,用户更改/上传昵称/头像时会入库并存入redis,所以此处直接从redis中查找昵称跟头像)
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
	res.From = from
	return
}

// CheckInlineOutlingSend 功能: 检测在线/离线; 发送数据给指定用户  传入: 需要发送用户切片/原始数据(需要改造为指定 msg_type再传入)
func CheckInlineOutlingSend(mobiles []string, msg TranstMsg) {
	conn := appredis.RedisDefaultPool.Get() // 获取redis链接
	defer conn.Close()
	for _, v := range mobiles {
		if _, ok := WebsocketManager.Group[v]; ok { // 在线
			res, _ := json.Marshal(TranslateMessage(msg, GetInfoByMobile([]string{msg.From}, 1)[0])) // 改变from 类型 并返回最新数据体格式
			_ = WebsocketManager.Group[v][v].Socket.WriteMessage(websocket.TextMessage, res)         // 直接 发送到指定用户连接通道
		} else { // 离线  : 按原样存储  又或者: 目标用户登录的事另外一台集群的服务器(则需要把消息存入redis中转并发给目标服务器处理)
			msg.To = []string{v} // 因为数据分发都是发送给 msg.To
			res, _ := json.Marshal(msg)
			_, err := conn.Do("RPUSH", v+":history", res) // 历史消息格式
			fmt.Println(err, "store history msg to redis fail--casso")
		}
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() { // 解析客户端发来的数据,在此分发数据到对应的长连接通道
	defer func() {
		log.Printf("client [%s] disconnect", c.ID)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.ID, err)
		}
	}()
	// 优化思路:1. 代码封装拆分(主流业务逻辑清晰) ; 2. 增加Groups存储信息(禁言/审核等flag) 3. 数据接入统一改变frmo(数据发送给用户时from需要有昵称跟头像信息,而客户端发送的数据的from只有电话号码所以需要转换) ; 4. 数据统一存入redis,方便集群; 5. 定义清晰的msg_type区分
	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// 注意: 消息确认分发后会改变form的格式,不再是string类型; 此时再走到 Write方法里面就会报:反序列化失败,也正好不走里面的逻辑,所以就留下了这个错误
			var msg TranstMsg
			if err := json.Unmarshal(message, &msg); err == nil {
				_ = storeMsgToRedis(message) // 消息存储
				// 这里可以根据接收到 的消息进行正确的消息分发: 需要用到的方法已经封装好
			} else {
				// 数据无法反序列化: 关闭连接
				err = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
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
		CheckInlineOutlingSend([]string{msg.From}, msg)
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
	res, _ := json.Marshal(g)
	_, err := conn.Do("set", gn, res) // 将新分组存入redis
	return err
}

// 新建分组存入redis
func storeGroups(msg TranstMsg, flag bool) error {
	conn := appredis.RedisDefaultPool.Get() // 获取redis链接
	defer conn.Close()
	var stores = make(map[string]interface{})
	stores["group"] = msg.Group
	stores["belong"] = msg.To
	stores["type"] = msg.GroupType
	stores["group_name"] = msg.GroupName
	stores["add"] = flag
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
			log.Printf("client [%s] connect", client.ID)
			// log.Printf("register client [%s] to group [%s]", client.ID, client.Group)
			manager.Lock.Lock()                     // 加锁
			if manager.Group[client.Group] == nil { // 如果是新的分组,则创建该新分组
				manager.Group[client.Group] = make(map[string]*Client) // 创建一个总管(manager)级别的分组, 但此时分组内部还是没有对应数据,只是make了内存地址
				manager.groupCount++                                   // 分组数自增
			}
			// 至此, 已经没有新分组的可能
			manager.Group[client.Group][client.ID] = client // {"12":{"17620439807":client,"17620439808":client,"17620439809":client}}
			manager.clientCount++                           // 长链接数自增-- 即 在线人数
			manager.Lock.Unlock()                           // 解锁
			appredis.SetHash("onlineUser", fmt.Sprint(manager.clientCount))
		// 注销
		case client := <-manager.UnRegister:
			log.Printf("unregister client [%s] from group [%s]", client.ID, client.Group)
			conn := appredis.RedisDefaultPool.Get() // 获取redis链接
			defer conn.Close()
			conn.Do("SADD", "underline", client.ID) // 下线人员记录: 上线时查找相应离线 消息
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Group]; ok {
				if _, ok := manager.Group[client.Group][client.ID]; ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.ID)
					manager.clientCount--
					if len(manager.Group[client.Group]) == 0 {
						//log.Printf("delete empty group [%s]", client.Group)
						delete(manager.Group, client.Group)
						manager.groupCount--
						appredis.SetHash("onlineUser", fmt.Sprint(manager.clientCount))
					}

				}
			}
			manager.Lock.Unlock()
		}
	}
}

// SendService 处理单个 client 发送数据 : 共用Write 方法来给通道写入数据
func (manager *Manager) SendService() {
	for {
		select {
		case data := <-manager.Message:
			if groupMap, ok := manager.Group[data.Group]; ok { // 在线 才将消息加入队列
				if conn, ok := groupMap[data.ID]; ok { // 存在 该链接才 将消息加入队列
					conn.Message <- data.Message
				}
			}
		}
	}
}

// SendGroupService 处理 group 广播数据
func (manager *Manager) SendGroupService() {
	for {
		select {
		// 发送广播数据到某个组的 channel 变量 Send 中
		case data := <-manager.GroupMessage:
			if groupMap, ok := manager.Group[data.Group]; ok {
				for _, conn := range groupMap {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// SendAllService 处理广播数据
func (manager *Manager) SendAllService() {
	for {
		select {
		case data := <-manager.BroadCastMessage:
			for _, v := range manager.Group {
				for _, conn := range v {
					conn.Message <- data.Message
				}
			}
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
	Group:            make(map[string]map[string]*Client),
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
	tokenInfo, err := jwt.ParseToken([]string{ctx.GetHeader("Sec-WebSocket-Protocol")}[0])
	if err != nil {
		log.Println("token invaliad")
		return
	}
	// fmt.Println(tokenInfo.Name) // token解析出来的认证用户的信息
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("websocket connect error: %s", tokenInfo.Name)
		return
	}

	client := &Client{
		ID:      tokenInfo.Name, // 电话作为唯一的 ID
		Group:   tokenInfo.Name, // connectors 组: 所有链接用户都在此处
		Socket:  conn,
		Message: make(chan []byte, 1024),
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
			r, _ := json.Marshal(v)
			_ = client.Socket.WriteMessage(websocket.TextMessage, r) // 直接 发送到用户连接通道
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
