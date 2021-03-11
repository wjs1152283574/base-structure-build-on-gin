package appmysql

// MYSQL 数据库初始化

import (
	"fmt"

	"net"
	"time"

	"goweb/utils/parsecfg"

	"github.com/jinzhu/gorm"
)

// DB DB
var DB *gorm.DB

func init() {
	InitMysql()
}

// InitMysql 初始化数据库连接
func InitMysql() {
	// 数据链对象--mysql
	fmt.Println(parsecfg.GlobalConfig.UseDbType)
	hp := net.JoinHostPort(parsecfg.GlobalConfig.Mysql.Write.Host, parsecfg.GlobalConfig.Mysql.Write.Port) // 需要使用这个方法将host/port 拼接起来才能正常运行
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", parsecfg.GlobalConfig.Mysql.Write.User, parsecfg.GlobalConfig.Mysql.Write.PassWord, hp, parsecfg.GlobalConfig.Mysql.Write.DataBase, parsecfg.GlobalConfig.Mysql.Write.ChatSet)
	fmt.Println(str)

	db, err := gorm.Open(parsecfg.GlobalConfig.UseDbType, str)
	if err != nil {
		panic(err)
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(parsecfg.GlobalConfig.Mysql.Write.SetMaxIdleConns)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(parsecfg.GlobalConfig.Mysql.Write.SetMaxOpenConns)

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
