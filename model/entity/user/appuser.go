package user

// SignUpReq 新建用户请求数据
type SignUpReq struct {
	Name   string `json:"username" binding:"required"`
	Pwd    string `json:"password" binding:"required"`
	Mobile string `json:"mobile" binding:"required"`
}
