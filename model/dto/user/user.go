package user

import (
	"goweb/dao/mysql"
	"goweb/model/vo/user"
	"time"

	"github.com/jinzhu/gorm"
)

func init() {
	mysql.DB.AutoMigrate(
		&User{},
		&Msg{},
	)
}

// User 用户模型
type User struct {
	gorm.Model
	Birthday *time.Time `json:"birthday"`
	Status   int        `gorm:"default:1" json:"status"`
	Type     int        `gorm:"default:1" json:"type"`
	Age      int        `gorm:"default:18" json:"age"`
	Name     string     `gorm:"not null;unique" json:"username" binding:"required"`
	Gender   int        `gorm:"default:1" json:"gender"`
	Pwd      string     `gorm:"not null" json:"password" binding:"required"`
	Mobile   string     `gorm:"size:12" json:"mobile" binding:"required"`
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
func (u *User) Create(res *user.ResUser) error {
	return mysql.DB.Create(u).Scan(res).Error
}

// Create 新建消息
func (m *Msg) Create() error {
	return mysql.DB.Create(m).Error
}

// Get 返回基础信息
func (u *User) Get(res *user.ResUser) error {
	return mysql.DB.First(res).Error
}

// Check 检测用户电话是否存在
func (u *User) Check() error {
	return mysql.DB.Where("mobile = ?", u.Mobile).First(u).Error
}

// AdminGetList admin list
func (u *User) AdminGetList(page, limit int, res *[]user.AdminUserList) (count int, err error) {
	mysql.DB.Model(u).Count(&count)
	return count, mysql.DB.Model(u).Limit(limit).Offset((page - 1) * limit).Scan(res).Error
}

// GetFrontU get for send all
func (u *User) GetFrontU(res *[]user.GetSendAll) error {
	return mysql.DB.Model(u).Where("deleted_at is null").Scan(res).Error
}
