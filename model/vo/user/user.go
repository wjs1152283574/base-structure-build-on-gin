package user

import (
	"time"

	"gorm.io/gorm"
)

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

// GetSendAll send all
type GetSendAll struct {
	Mobile string `json:"mobile"`
}

// AdminUserList 获取用户列表
type AdminUserList struct {
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
