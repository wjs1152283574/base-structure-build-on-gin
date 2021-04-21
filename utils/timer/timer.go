/*
 * @Description:定时任务
 * @Author: Casso-Wong
 */

package timer

import (
	"fmt"
	"goweb/utils/parsecfg"

	"github.com/robfig/cron"
)

// StoreMsg StoreMsg from redis to mysql
type StoreMsg struct{}

// StoreGroup StoreMsg from redis to mysql
type StoreGroup struct{}

// Run 固定方法,指定接收者可添加任务(在main.go中可见)
func (s *StoreMsg) Run() {
	fmt.Println("定制任务操作逻辑")
}

// Run 固定方法,指定接收者可添加任务(在main.go中可见)
func (s StoreGroup) Run() {
	fmt.Println("定制任务操作逻辑")
}

// Conrs 定时器
var Conrs *cron.Cron

func init() {
	Conrs = cron.New() // 定时任务
	Conrs.Start()
	Conrs.AddJob(parsecfg.GlobalConfig.Timer.Store, &StoreMsg{})
}
