package tool

import (
	"errors"
	"log"
	"os"
	"time"
)

// ParseSlicess 字符串切片去重
func ParseSlicess(slc []string) []string {
	var res []string
	temMap := make(map[string]bool, len(slc))
	for _, v := range slc {
		if v != "" {
			if temMap[v] {
				temMap[v] = true
				res = append(res, v)
			}
		}
	}
	return res
}

// DeledeSlice 根据传入string删除相同的切片元素: 我的应用场景是知道切片中必然有且仅有一个相同的元素,想删掉
func DeledeSlice(slc []string, target string) []string {
	var index = 0
	for i, v := range slc {
		if v == target {
			index = i
			break
		}
	}
	// ["aa","bb","cc"]  例: index = 2  则: 会超出下标界限  下面的append
	if index == len(slc)-1 && len(slc) != 0 { // 删除最后一项
		slc = slc[:index] // 左包含,右不包含
	} else {
		slc = append(slc[:index], slc[index+1:]...) // 左包含,右不包含
	}
	return slc
}

// StartToday 返回当天零点
func StartToday() time.Time {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	return tm2
}

// GetTimesDiffer 获取相差时间秒数
func GetTimesDiffer(stime, etime string) (res int64, errs error) {
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", stime, time.Local)
	t2, err2 := time.ParseInLocation("2006-01-02 15:04:05", etime, time.Local)
	if err == nil && err2 == nil && t1.Before(t2) {
		res = t2.Unix() - t1.Unix()
		errs = nil
	} else {
		res = 0
		errs = err
	}
	return
}

//PathExists 判断文件夹是否存在 不存在则创建
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Printf("mkdir failed![%v]\n", err)
		} else {
			return true, nil
		}
	}
	return false, err
}

// TurnToSlice 将切片转为特定字符串隔开的字符串(为了存入数据库..我也是福了)
func TurnToSlice(s []string, slc string) (res string) { // 效果类似 : 1,2,3,4  || 1:2:3:4
	for i, v := range s {
		if i != len(s)-1 {
			res += v + slc
		} else {
			res += v
		}
	}
	return
}

// InSlice 实现类似Python 中的 in
func InSlice(t string, s []string) (bool, error) {
	if len(s) == 0 {
		err := errors.New("传入切片不可为空")
		return false, err
	}
	for _, v := range s {
		if t == v {
			return true, nil
		}
	}
	return false, nil
}

// GetMaxint 返回最大数值,及出现次数
func GetMaxint(s []int) []int {
	var m = make(map[int]int, len(s))
	for _, v := range s {
		if m[v] == 0 {
			m[v] = 1
		} else {
			m[v]++
		}
	}
	var max int
	for _, v := range m {
		if max == 0 || v > max {
			max = v
		}
	}
	var res []int
	for k, v := range m {
		if v == max {
			res = append(res, k)
		}
	}
	return res
}
