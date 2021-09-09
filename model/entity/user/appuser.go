package user

import "time"

// SignUpReq 新建用户请求数据
type SignUpReq struct {
	Name   string `json:"username" binding:"required"`
	Pwd    string `json:"password" binding:"required"`
	Mobile string `json:"mobile" binding:"required"`
}

// UserUpdReq 编辑用户请求数据
type UserUpdReq struct {
	Birthday *time.Time `json:"birthday"`
	Status   int        `json:"status"`
	Type     int        `json:"type"`
	Age      int        `json:"age"`
	Name     string     `json:"username" binding:"required"`
	Gender   int        `json:"gender"`
	Mobile   string     `json:"mobile" binding:"required"`
}
