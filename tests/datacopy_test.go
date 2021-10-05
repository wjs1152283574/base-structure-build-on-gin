/*
 * @PackageName: tests
 * @Description: 数据复制单元测试
 * @Author: Casso-Wong
 * @Date: 2021-10-05 14:32:15
 * @Last Modified by: Casso-Wong
 * @Last Modified time: 2021-10-05 14:32:15
 */
package tests

import (
	"goweb/utils/datacopy"
	"testing"
)

type Per struct {
	Name   string
	Age    int
	Gender int
}

type Stu struct {
	Name   string
	Age    int
	Gender int
	Test   string // 这个字段经过copy后应该还是零值，即copy只覆盖相同字段，才符合需求
}

func TestDataCopy(t *testing.T) {
	p := Per{
		Name:   "casso",
		Age:    23,
		Gender: 1,
	}
	s := Stu{}
	if err := datacopy.DataCopy(&p, &s); err != nil || s.Name == "" || s.Test != "" {
		t.Error(err)
	}
}
