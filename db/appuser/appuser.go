package appuser

import (
	"goweb/dao/appmysql"
	"time"

	"github.com/jinzhu/gorm"
)

func init() {
	appmysql.DB.AutoMigrate(
		&User{},
		&Msg{},
	)
}

// User 用户模型
type User struct {
	gorm.Model
	Birthday *time.Time `json:"birthday"`
	Age      int        `gorm:"default:18" json:"age"`
	Name     string     `gorm:"not null;unique" json:"username" binding:"required"`
	Gender   int        `gorm:"default:1" json:"gender"`
	Pwd      string     `gorm:"not null" json:"password" binding:"required"`
	Mobile   *string    `gorm:"size:12" json:"mobile" binding:"required"`
}

// Msg 消息
type Msg struct {
	gorm.Model
	UID int    `json:"u_id"` // 用户ID
	Msg string `json:"msg"`
}

// AfterCreate 创建 User 之后执行得钩子函数 --- 自动执行
func (u *User) AfterCreate(scope *gorm.Scope) error {
	return scope.DB().Model(u).Update("role", "admin").Error
}

// BeforeCreate 创建 User 之前执行得钩子函数 --- 自动执行
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	return nil
}

// Create 创建用户
func (u *User) Create(res *ResUser) error {
	return appmysql.DB.Create(u).First(res).Error
}

// Create 新建消息
func (m *Msg) Create() error {
	return appmysql.DB.Create(m).Error
}

// ResUser query response
type ResUser struct {
	ID        int        `json:"usr_id"`
	Age       int        `json:"age"`
	Name      string     `json:"username"`
	Gender    int        `json:"sex"`
	Mobile    string     `json:"phone"`
	Birthday  *time.Time `json:"birthday"`
	CreatedAt time.Time  `json:"join_day"`
}

// Get 返回基础信息
func (u *User) Get(res *ResUser) error {
	return appmysql.DB.First(res).Error
}

// Check 检测用户电话是否存在
func (u *User) Check() error {
	return appmysql.DB.Where("mobile = ?", u.Name).First(u).Error
}
