/*
 * @Author: Casso-Wong
 * @Date: 2021-06-04 14:46:03
 * @Last Modified by:   Casso-Wong
 * @Last Modified time: 2021-06-04 14:46:03
 */
package sett

import (
	"goweb/dao/mysql"
	vo "goweb/model/vo/sett"

	"github.com/jinzhu/gorm"
)

// 配置

func init() {
	mysql.DB.AutoMigrate(
		&Sett{},
		&ConBanner{},
		&ConTag{},
		&MyApp{},
		&Bag{},
		&NavFun{},
	)
	// 直接生成
}

// Sett 总配置:显示/隐藏
type Sett struct {
	gorm.Model
	ConBannr      int    `gorm:"default:1" json:"con_banner"` // 1 默认显示  2 不显示  联系人banner
	ConTag        int    `gorm:"default:1" json:"con_tag"`    // 1 默认显示  2 不显示  联系人标签
	Wild          int    `gorm:"default:1" json:"wild"`       // 1 默认显示  2 不显示  心情广场
	WildBanner    string `json:"wild_banner"`                 // 心情广场banner
	Video         int    `gorm:"default:1" json:"video"`      // 1 默认显示  2 不显示  视频
	VideoIcon     string `json:"video_icon"`                  // 视频图标
	VideoTitle    string `json:"video_title"`                 // 视频标题
	VideiSubTitle string `json:"video_sub_title"`             // 视频副标题
	Live          int    `gorm:"default:1" json:"live"`       // 1 默认显示  2 不显示  直播
	LiveIcon      string `json:"live_icon"`                   // 直播图标
	LiveTitle     string `json:"live_title"`                  // 直播标题
	LiveSubTitle  string `json:"live_sub_title"`              // 直播副标题
	MyApp         int    `gorm:"default:1" json:"my_app"`     // 1 默认显示  2 不显示  我的应用
	Recommend     int    `gorm:"default:1" json:"reccommend"` // 1 默认显示  2 不显示  官方推荐
	Bag           int    `gorm:"default:1" json:"bag"`        // 1 默认显示  2 不显示  钱包
	WebLink       int    `gorm:"default:1" json:"web_link"`   // 1 默认显示  2 不显示  星聊网页版链接

	UpdLink  string `json:"upd_link"`                      // 更新包链接
	Versions string `gorm:"default:'2.0'" json:"versions"` // 版本号
	UpdDesc  string `json:"upd_desc"`                      // 更新描述

	UpdLinkIos  string `json:"upd_link_ios"`                      // 更新包链接 IOS
	VersionsIos string `gorm:"default:'2.0'" json:"versions_ios"` // 版本号  IOS
	UpdDescIos  string `json:"upd_desc_ios"`                      // 更新描述  IOS

	Team int `gorm:"default:1" json:"team"` // 1 显示 2 不显示

	LoginTitle  string `gorm:"size:100" json:"login_title"` // 登录界面大标题文字
	IndexBanner string `json:"index_banner"`                // 首页banner图
	LfetTitle   string `json:"lfet_title"`                  // 左侧弹窗大标题
	MypageItem  string `json:"mypage_item"`                 // xx网页版文字
	AbuotLogo   string `json:"about_logo"`                  // 关于--logo
	AboutTtile  string `json:"about_title"`                 // 关于--标题

	AppCode string `gorm:"not null;" json:"app_code"` // 不同配置标识 code
	AppName string `gorm:"not null;" json:"app_name"` // 不同配置标识 name

	MsgLocation string `json:"msg_location"`  // 短信服务所在地
	MsgKey      string `json:"msg_key"`       // key
	MsgScretKey string `json:"msg_scret_key"` // scrit key
	MsgScheme   string `json:"msg_scheme"`    // https
	MsgSignName string `json:"msg_sign_name"` // 签名
	MsgTemplate string `json:"msg_template"`  // 模板
}

// ConBanner 联系人banner
type ConBanner struct {
	gorm.Model
	SID   int    `json:"sid"`   //
	URL   string `json:"url"`   // banner 链接
	Image string `json:"image"` //  banner
	Sort  int    `json:"sort"`
}

// ConTag 联系人标签
type ConTag struct {
	gorm.Model
	SID   int    `json:"sid"`   //
	Route string `json:"route"` // 前端路由
	Name  string `json:"name"`  // 标签名称
	Icon  string `json:"icon"`  // 标签图标
	Sort  int    `json:"sort"`
}

// MyApp 我的应用
type MyApp struct {
	gorm.Model
	SID   int    `json:"sid"`   //
	Route string `json:"route"` // 前端路由
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Sort  int    `json:"sort"`
}

// Bag 钱包
type Bag struct {
	gorm.Model
	SID   int    `json:"sid"` //
	Route string `json:"route"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Sort  int    `json:"sort"`
}

// NavFun 左侧菜单功能
type NavFun struct {
	gorm.Model
	SID   int    `json:"sid"` //
	Route string `json:"route"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Sort  int    `json:"sort"`
	Show  int    `gorm:"default:1" json:"show"` // 1 默认显示  2 不显示
}

// GetForAli  获取对应APP的短信配置
func (s *Sett) GetForAli(res *vo.AlimsgNeed) error {
	return mysql.DB.Table("setts").Select("msg_location,msg_key,msg_scret_key,msg_scheme,msg_sign_name,msg_template").Where("app_code =?", s.AppCode).First(res).Error
}

// Create  生成新的配置
func (s *Sett) Create() error {
	return mysql.DB.Model(s).Create(s).First(s).Error
}

// Upd  更新总配置
func (s *Sett) Upd() error {
	return mysql.DB.Model(s).Updates(s).Where("id = ?", s.ID).Error
}

// Create  Create
func (c *ConBanner) Create() error {
	return mysql.DB.Create(c).Error
}

// Del delete
func (c *ConBanner) Del() error {
	return mysql.DB.Exec("delete from con_banners where id = ?", c.ID).Error
}

// Upd  update
func (c *ConBanner) Upd() error {
	return mysql.DB.Save(c).Error
}

// Create  Create
func (c *ConTag) Create() error {
	return mysql.DB.Create(c).Error
}

// Del delete
func (c *ConTag) Del() error {
	return mysql.DB.Exec("delete from con_tags where id = ?", c.ID).Error
}

// Upd  update
func (c *ConTag) Upd() error {
	return mysql.DB.Save(c).Error
}

// Create  Create
func (m *MyApp) Create() error {
	return mysql.DB.Create(m).Error
}

// Del delete
func (m *MyApp) Del() error {
	return mysql.DB.Exec("delete from my_apps where id = ?", m.ID).Error
}

// Upd  update
func (c *MyApp) Upd() error {
	return mysql.DB.Save(c).Error
}

// Create  Create
func (b *Bag) Create() error {
	return mysql.DB.Create(b).Error
}

// Del delete
func (b *Bag) Del() error {
	return mysql.DB.Exec("delete from bags where id = ?", b.ID).Error
}

// Upd  update
func (c *Bag) Upd() error {
	return mysql.DB.Save(c).Error
}

// Create  Create
func (n *NavFun) Create() error {
	return mysql.DB.Create(n).Error
}

// Del delete
func (n *NavFun) Del() error {
	return mysql.DB.Exec("delete from nav_funs where id = ?", n.ID).Error
}

// Upd  update
func (n *NavFun) Upd() error {
	return mysql.DB.Model(n).Save(n).Error
}

// Get 获取总设置
func (s *Sett) Get(res *vo.SetRes) error {
	return mysql.DB.Table("setts").Select("*").Where("app_code = ?", s.AppCode).First(res).Error
}

// Check
func (s *Sett) Check() error {
	return mysql.DB.Table("setts").Select("*").Where("app_code = ?", s.AppCode).Scan(s).Error
}

// GetAdmin 获取总设置
func (s *Sett) GetAdmin(res *vo.SetResAdmin) error {
	return mysql.DB.Table("setts").Select("*").Where("app_code = ?", s.AppCode).First(res).Error
}

// GetAdmin 获取总设置
func (s *Sett) GetCode(res *[]vo.Codes) error {
	return mysql.DB.Table("setts").Select("id,app_code,app_name").Scan(res).Error
}
