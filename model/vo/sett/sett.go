package sett

// AlimsgNeed 发送短信 配置获取
type AlimsgNeed struct {
	MsgLocation string `json:"msg_location"`  // 短信服务所在地
	MsgKey      string `json:"msg_key"`       // key
	MsgScretKey string `json:"msg_scret_key"` // scrit key
	MsgScheme   string `json:"msg_scheme"`    // https
	MsgSignName string `json:"msg_sign_name"` // 签名
	MsgTemplate string `json:"msg_template"`  // 模板
}

// SetRes 总设置返回结构
type SetRes struct {
	ConBannr      int            `gorm:"default:1" json:"con_banner"` // 1 默认显示  2 不显示  联系人banner
	ConTag        int            `gorm:"default:1" json:"con_tag"`    // 1 默认显示  2 不显示  联系人标签
	Wild          int            `gorm:"default:1" json:"wild"`       // 1 默认显示  2 不显示  心情广场
	WildBanner    string         `json:"wild_banner"`                 // 心情广场banner
	Video         int            `gorm:"default:1" json:"video"`      // 1 默认显示  2 不显示  视频
	VideoIcon     string         `json:"video_icon"`                  // 视频图标
	VideoTitle    string         `json:"video_title"`                 // 视频标题
	VideiSubTitle string         `json:"video_sub_title"`             // 视频副标题
	Live          int            `gorm:"default:1" json:"live"`       // 1 默认显示  2 不显示  直播
	LiveIcon      string         `json:"live_icon"`                   // 直播图标
	LiveTitle     string         `json:"live_title"`                  // 直播标题
	LiveSubTitle  string         `json:"live_sub_title"`              // 直播副标题
	MyApp         int            `gorm:"default:1" json:"my_app"`     // 1 默认显示  2 不显示  我的应用
	Recommend     int            `gorm:"default:1" json:"reccommend"` // 1 默认显示  2 不显示  官方推荐
	Bag           int            `gorm:"default:1" json:"bag"`        // 1 默认显示  2 不显示  钱包
	WebLink       int            `gorm:"default:1" json:"web_link"`   // 1 默认显示  2 不显示  星聊网页版链接
	ConBannerRes  []ConBannerRes `json:"con_banners"`                 // 联系人banner列表
	ConTagRes     []ConTagRes    `json:"con_tags"`                    // 联系人标签列表
	MyAppRes      []MyAppRes     `json:"my_apps"`                     // 我的应用列表
	BagRes        []BagRes       `json:"bags"`                        // 钱包列表
	NavFunRes     []NavFunRes    `json:"nav_funs"`                    // 左侧菜单功能列表

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
	MypageItem  string `json:"mypage_item"`                 // xx网页版
	AbuotLogo   string `json:"about_logo"`                  // 关于--logo
	AboutTtile  string `json:"about_title"`                 // 关于--标题
}

// ConBannerRes x
type ConBannerRes struct {
	URL   string `json:"url"`   // banner 链接
	Image string `json:"image"` //  banner
}

// ConTagRes x
type ConTagRes struct {
	Route string `json:"route"` // 前端路由
	Name  string `json:"name"`  // 标签名称
	Icon  string `json:"icon"`  // 标签图标
}

// MyAppRes x
type MyAppRes struct {
	Route string `json:"route"` // 前端路由
	Name  string `json:"name"`
	Icon  string `json:"icon"`
}

// BagRes x
type BagRes struct {
	Route string `json:"route"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
}

// NavFunRes x
type NavFunRes struct {
	Route string `json:"route"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Show  int    `gorm:"default:1" json:"show"` // 1 默认显示  2 不显示
}

// SetResAdmin 总设置返回结构
type SetResAdmin struct {
	UpdLinkIos  string `json:"upd_link_ios"`                      // 更新包链接 IOS
	VersionsIos string `gorm:"default:'2.0'" json:"versions_ios"` // 版本号  IOS
	UpdDescIos  string `json:"upd_desc_ios"`                      // 更新描述  IOS

	Team int `gorm:"default:1" json:"team"` // 1 显示 2 不显示
}

type Codes struct {
	ID      int    `json:"id"`       // set id
	AppCode string `json:"app_code"` // appcode
	AppName string `json:"app_name"`
}
