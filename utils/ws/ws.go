/*
 * @Author: Casso-Wong
 * @Date: 2021-06-05 10:15:44
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-09-17 00:56:29
 */
package ws

import (
	"encoding/json"
	appredis "goweb/dao/redis"
	dto "goweb/model/dto/user"
	"goweb/utils/alimsg"
	"goweb/utils/customerjwt"
	"goweb/utils/kafka"
	"goweb/utils/response"
	"goweb/utils/statuscode"
	"log"
	"net/http"
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
	From        string     `json:"from"`         // 发送人clientID
	To          []string   `json:"to"`           // 接收人clientIDs
	Group       string     `json:"group"`        // 组ID
	GroupName   string     `json:"group_name"`   // 组名称
	GroupIcon   string     `json:"group_icon"`   // 群图标
	GroupType   int        `json:"group_type"`   // 组类型
	Msg         string     `json:"msg"`          // 消息内容
	Time        *time.Time `json:"time"`         // 发送时间
	ContentType string     `json:"content_type"` // 内容类型
	MsgType     int64      `json:"msg_type"`     // 消息类型
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

// Read 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c // 仅在此处触发注销
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] 服务端断开连接出错 READ: %s，表示socket在WRITE已被关闭\n", c.ID, err)
		}
	}()
	for {
		// 从客户端那边发送过来的数据再写入Message通道(后台消息分发是程序直接往Message通道里面直接写入数据)
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			return
		}

		// 心跳消息,需马上回复
		if messageType == websocket.PingMessage {
			var hb HeartBeatRes
			hb.Ping = int(time.Now().Unix())
			hbs, _ := json.Marshal(hb)
			c.Socket.WriteMessage(websocket.TextMessage, hbs)
			continue
		}

		// 消息正常发送
		log.Printf("client [%s] receive message: %s\n", c.ID, string(message))
		c.Message <- message // 写入channel,等待写入对应长连接通道
		log.Printf("在线人数:%d\n", WebsocketManager.LenGroup())
	}
}

// Write 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		c.Socket.WriteMessage(websocket.CloseMessage, []byte("socket closed by WRITE "))
	}()
	for {
		message, ok := <-c.Message
		if !ok { // 通道 关闭且无值
			return
		}

		// 解析客户端发来的数据,在此分发数据到对应的长连接通道，这里忽略
		c.Socket.WriteMessage(websocket.TextMessage, message)
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

// NewGroupID return group ID
func NewGroupID(mobile string) string {
	return mobile + ":chat" + alimsg.Code()
}

// Start 启动 websocket 管理器
func (manager *Manager) Start() {
	log.Printf("websocket manage start")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			manager.Lock.Lock()                         // 加锁
			if _, ok := manager.Group[client.ID]; !ok { // 如果是新的分组,则创建该新分组
				manager.Group[client.ID] = client // 创建一个总管(manager)级别的分组, 但此时分组内部还是没有对应数据,只是make了内存地址
				manager.groupCount++              // 分组即链接
			}
			manager.Lock.Unlock()                                       // 解锁
			appredis.SetHash("onlineUser:"+client.ID, kafka.KafkaTopic) // 维护在线用户，value为kafka topic
		// 注销
		case client := <-manager.UnRegister:
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
		log.Println("token 不合法，长连接失败返回！本次token:", []string{ctx.GetHeader("Sec-WebSocket-Protocol")}[0])
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

		// 检测用户合法性
		var me dto.User
		me.ID = uint(tokenInfo.ID)
		if err := me.Check(); err != nil {
			response.ReturnJSON(ctx, http.StatusOK, statuscode.UserNotExit.Code, statuscode.UserNotExit.Msg, nil)
			return
		}

		// 生成websocket client 并处理
		client := &Client{
			ID:      me.Mobile,               // 电话作为唯一的 ID
			Socket:  conn,                    // socket链接
			Message: make(chan []byte, 1024), // 通道内消息缓存大小
		}

		manager.RegisterClient(client) // 注册到localCahe,后期需要集群到不同节点，预测利用redis-cluster 来作为调度中心（注册时：用户标识+节点标识）
		go client.Read()               // 阻塞处理socket消息通道读取
		go client.Write()              // 阻塞处理socket消息通道写入
		appredis.Dellist("underline", client.ID)
	}
}
