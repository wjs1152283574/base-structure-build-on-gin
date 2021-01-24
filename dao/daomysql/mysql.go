package daomysql

// MYSQL 数据库初始化

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// DB DB
var DB *gorm.DB

// InitMysql 初始化数据库连接
func InitMysql(datatype, user, pwd, dbname, charset string) {
	// 数据链对象--mysql
	str := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True&loc=Local", user, pwd, dbname, charset)
	// fmt.Println(str)
	db, err := gorm.Open(datatype, str)

	if err != nil {
		fmt.Printf("open db faild, %#v", err)
		return
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)

	DB = db
}

// AmountGreaterThan1000 ...
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("amount > ?", 1000)
}

// PaidWithCreditCard ...
func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "C")
}

// PaidWithCod ...
func PaidWithCod(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "C")
}

// OrderStatus ...
func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(AmountGreaterThan1000).Where("status IN (?)", status)
	}
}

// 使用:
// db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)
