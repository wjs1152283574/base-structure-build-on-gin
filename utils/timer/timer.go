/*
 * @Description:定时任务
 * @Author: Casso-Wong
 */

package timer

// StoreMsg StoreMsg from redis to mysql
type StoreMsg struct{}

// Run insert msg records to mysql form redis
// func (s StoreMsg) Run() {
// 	conn := appredis.RedisDefaultPool.Get()
// 	defer conn.Close()
// 	if res, err := appredis.Rpop("msg"); err == nil {
// 		var msg ws.TranstMsg
// 		_ = json.Unmarshal(res, &msg)
// 		// 先写入用户消息记录表
// 		var msgModel userlog.UserMsg
// 		msgModel.ContentType = msg.ContentType
// 		msgModel.From = msg.From
// 		msgModel.Group = msg.Group
// 		msgModel.GroupName = msg.GroupName
// 		msgModel.GroupType = msg.GroupType
// 		msgModel.Msg = msg.Msg
// 		if len(msg.To) > 2 {
// 			msgModel.AcceptTyp = 2
// 		} else {
// 			msgModel.AcceptTyp = 1
// 		}
// 		var strTo = ""
// 		for i, v := range msg.To {
// 			if i != len(msg.To)-1 {
// 				strTo += v + ","
// 			} else {
// 				strTo += v
// 			}
// 		}
// 		msgModel.To = strTo
// 		msgModel.Time = msg.Time
// 		msgModel.MsgType = int(msg.MsgType)
// 		if err := msgModel.CreateUserMsg(); err == nil {
// 			fmt.Printf("来自%s的消息入库成功\n", msg.From)
// 		}
// 	}
// }

// StoreGroup StoreMsg from redis to mysql
type StoreGroup struct{}

// ForStoreGroup ...
type ForStoreGroup struct {
	Type      int      `json:"type"`       // 组类型  1 群  2 团队
	Add       bool     `json:"add"`        // 为true : 需要根据组名去更新数据库里面得的belong信息
	Belong    []string `json:"belong"`     // 组员 电话
	Group     string   `json:"group"`      // 组编号
	GroupName string   `json:"group_name"` // 组编号
}

// Run insert groups records to mysql form redis
// func (s StoreGroup) Run() {
// 	conn := appredis.RedisDefaultPool.Get()
// 	defer conn.Close()
// 	if res, err := appredis.Rpop("groups"); err == nil {
// 		var groups group.Group
// 		var chatGroups group.ChatGroup
// 		var groupItem ForStoreGroup
// 		_ = json.Unmarshal(res, &groupItem)
// 		if groupItem.Type == 2 { // 团队
// 			groups.GroupID = groupItem.Group
// 			var strTo = ""
// 			for i, v := range groupItem.Belong {
// 				if i != len(groupItem.Belong)-1 {
// 					strTo += v + ","
// 				} else {
// 					strTo += v
// 				}
// 			}
// 			groups.Belong = strTo
// 			groups.Manager = strings.Split(groupItem.Group, ":")[0]
// 			groups.Name = groupItem.GroupName
// 			groups.CreateGroup()
// 		} else { // 正常聊天 -- 写入聊天组表
// 			chatGroups.ChatGroupID = groupItem.Group
// 			var strTo = ""
// 			for i, v := range groupItem.Belong {
// 				if i != len(groupItem.Belong)-1 {
// 					strTo += v + ","
// 				} else {
// 					strTo += v
// 				}
// 			}
// 			if len(groupItem.Belong) > 2 {
// 				chatGroups.GroupType = 1 // 群
// 			} else {
// 				chatGroups.GroupType = 0 // 私
// 			}
// 			chatGroups.ChatGroupName = groupItem.GroupName
// 			chatGroups.Belong = strTo
// 			chatGroups.Manager = strings.Split(groupItem.Group, ":")[0]
// 			chatGroups.CreateChatGroup()
// 		}

// 		// group.GroupID = groupItem.group
// 	}
// }
