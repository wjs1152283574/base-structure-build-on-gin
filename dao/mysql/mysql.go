/*
 * @Author: Casso-Wong
 * @Date: 2021-06-04 14:41:27
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-09-10 11:15:49
 */
package mysql

// MYSQL 数据库初始化

import (
	"fmt"
	"goweb/utils/parsecfg"

	"net"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB DB
var DB *gorm.DB

func init() {
	if parsecfg.GlobalConfig.Env == "dev" {
		InitMysql(parsecfg.GlobalConfig.Mysql.Write.Host, parsecfg.GlobalConfig.Mysql.Write.Port, parsecfg.GlobalConfig.Mysql.Write.User, parsecfg.GlobalConfig.Mysql.Write.PassWord, parsecfg.GlobalConfig.Mysql.Write.DataBase, parsecfg.GlobalConfig.Mysql.Write.Charset)
	}
	if parsecfg.GlobalConfig.Env == "test" {
		InitMysql(parsecfg.GlobalConfig.Mysql.Write.HostLive, parsecfg.GlobalConfig.Mysql.Write.PortLive, parsecfg.GlobalConfig.Mysql.Write.User, parsecfg.GlobalConfig.Mysql.Write.PassWord, parsecfg.GlobalConfig.Mysql.Write.DataBase, parsecfg.GlobalConfig.Mysql.Write.Charset)
	}
	// you can add other env here
	fmt.Println("所有被编译器发现的 init 函数都会安排在 main 函数之前执行 。init 函数用在设置包、初始化变量或其他要在程序运行前优先完成的引导工作。")
}

// InitMysql 初始化数据库连接
func InitMysql(host, port, user, pass, dbname, chaset string) {
	// 数据链对象--mysql
	hp := net.JoinHostPort(host, port)
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", user, pass, hp, dbname, chaset)
	db, err := gorm.Open(parsecfg.GlobalConfig.DbType, str)
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
