package appuser

import (
	"goweb/dao/appmysql"
	"time"

	"github.com/jinzhu/gorm"
)

func init() {
	appmysql.DB.AutoMigrate(
		&User{},
		&Address{},
		&Email{},
		&Language{},
		&CreditCard{})
}

// User 用户模型
type User struct {
	gorm.Model
	Birthday *time.Time `json:"birthday"`
	Age      int        `gorm:"default:18" json:"age"`
	Name     string     `gorm:"not null;unique" json:"username" binding:"required"` // string默认长度为255, 使用这种tag重设。
	Gender   int        `gorm:"default:1" json:"gender"`                            // 自增
	Pwd      string     `gorm:"not null" json:"password" binding:"required"`
	Mobile   *string    `gorm:"size:12" json:"mobile"`

	CreditCard []CreditCard `json:"credit_cards"` // One-To-Many (拥有一个 - CreditCard表的UserID作外键)
	Emails     []Email      `json:"emails"`       // One-To-Many (拥有多个 - Email表的UserID作外键)

	BillingAddress   Address `gorm:"foreignkey:BillingAddressID" json:"billing_addr"` // One-To-One (属于 - 本表的BillingAddressID作外键)
	BillingAddressID *int

	ShippingAddress   Address `gorm:"foreignkey:ShippingAddressID" json:"shipping_addr"` // One-To-One (属于 - 本表的ShippingAddressID作外键)
	ShippingAddressID *int    // 使用指针: 不会存零值  不传就是null而不hi零

	IgnoreMe  int        `gorm:"-"`                                         // 忽略这个字段
	Languages []Language `gorm:"many2many:user_languages;" json:"language"` // Many-To-Many , 'user_languages'是连接表
}

// AfterCreate 创建 User 之后执行得钩子函数 --- 自动执行
func (u *User) AfterCreate(scope *gorm.Scope) error {
	return scope.DB().Model(u).Update("role", "admin").Error
}

// BeforeCreate 创建 User 之前执行得钩子函数 --- 自动执行
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	return nil
}

// CreateUser 创建用户
func (u *User) CreateUser() error {
	return appmysql.DB.Create(u).Error
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

// CheckUsr 检测用户是否存在,并返回基础信息
func (u *User) CheckUsr(res *ResUser) error {
	return appmysql.DB.Where("name = ?", u.Name).First(res).Error
}

// Email 邮箱
type Email struct {
	ID         int    `json:"email_id"`
	UserID     int    `gorm:"index" json:"-"`                              // 外键 (属于), tag `index`是为该列创建索引
	Email      string `gorm:"type:varchar(100);unique_index" json:"email"` // `type`设置sql类型, `unique_index` 为该列设置唯一索引
	Subscribed bool   `json:"subscrib"`
}

// Address 地址
type Address struct {
	ID       int
	Address1 string `gorm:"not null;type:varchar(100)" json:"addr1"` // 设置字段为非空并唯一
	Address2 string `gorm:"type:varchar(100)" json:"addr2"`
}

// Language 语种
type Language struct {
	ID   int
	Name string `gorm:"index:idx_name_code" json:"language_name"` // 创建索引并命名，如果找到其他相同名称的索引则创建组合索引
	Code string `gorm:"index:idx_name_code" json:"language_code"` // `unique_index` also works
}

// CreditCard 卡s
type CreditCard struct {
	gorm.Model
	UserID uint
	Number string `gorm:"not null;type:varchar(25);unique" json:"cnums"`
}
