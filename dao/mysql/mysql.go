/*
 * @Author: Casso-Wong
 * @Date: 2021-06-04 14:41:27
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-09-17 00:45:35
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
		InitMysql(parsecfg.GlobalConfig.Mysql.Dev.Host,
			parsecfg.GlobalConfig.Mysql.Dev.Port,
			parsecfg.GlobalConfig.Mysql.Dev.User,
			parsecfg.GlobalConfig.Mysql.Dev.PassWord,
			parsecfg.GlobalConfig.Mysql.Dev.DataBase,
			parsecfg.GlobalConfig.Mysql.Dev.Charset,
			parsecfg.GlobalConfig.Mysql.Dev.SetMaxIdleConns,
			parsecfg.GlobalConfig.Mysql.Dev.SetMaxOpenConns)
	}
	if parsecfg.GlobalConfig.Env == "test" {
		InitMysql(parsecfg.GlobalConfig.Mysql.Prod.Host,
			parsecfg.GlobalConfig.Mysql.Prod.Port,
			parsecfg.GlobalConfig.Mysql.Prod.User,
			parsecfg.GlobalConfig.Mysql.Prod.PassWord,
			parsecfg.GlobalConfig.Mysql.Prod.DataBase,
			parsecfg.GlobalConfig.Mysql.Prod.Charset,
			parsecfg.GlobalConfig.Mysql.Prod.SetMaxIdleConns,
			parsecfg.GlobalConfig.Mysql.Prod.SetMaxOpenConns)
	}
	if parsecfg.GlobalConfig.Env == "stage" {
		InitMysql(parsecfg.GlobalConfig.Mysql.Stage.Host,
			parsecfg.GlobalConfig.Mysql.Stage.Port,
			parsecfg.GlobalConfig.Mysql.Stage.User,
			parsecfg.GlobalConfig.Mysql.Stage.PassWord,
			parsecfg.GlobalConfig.Mysql.Stage.DataBase,
			parsecfg.GlobalConfig.Mysql.Stage.Charset,
			parsecfg.GlobalConfig.Mysql.Stage.SetMaxIdleConns,
			parsecfg.GlobalConfig.Mysql.Stage.SetMaxOpenConns)
	}
	// you can add other env here
	fmt.Println("所有被编译器发现的 init 函数都会安排在 main 函数之前执行 。init 函数用在设置包、初始化变量或其他要在程序运行前优先完成的引导工作。")
}

// InitMysql 初始化数据库连接
func InitMysql(host, port, user, pass, dbname, chaset string, maxidle, maxopen int) {
	// 数据链对象--mysql
	hp := net.JoinHostPort(host, port)
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", user, pass, hp, dbname, chaset)
	db, err := gorm.Open(parsecfg.GlobalConfig.DbType, str)
	if err != nil {
		panic(err)
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(maxidle)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(maxopen)

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
