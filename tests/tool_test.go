package tests

import (
	"goweb/utils/tool"
	"reflect"
	"testing"
)

func TestParseSlicess(t *testing.T) {
	// type test struct {
	// 	source []string
	// 	res    []string
	// }
	// tests := []test{
	// 	{source: []string{"aa", "bb", "cc"}, res: []string{"gg", "qq"}},
	// 	{source: []string{"ww", "qq"}, res: []string{"tt", "jj"}},
	// }
	// for _, v := range tests {
	// 	final := GetMaxint(v.source)
	// 	if !reflect.DeepEqual(v.res, final) {
	// 		t.Errorf("want:%#v,got:%#v", v.res, final)
	// 	}
	// }
	source := []string{"aa", "bb", "cc", "cc"}
	res := []string{"aa", "bb", "cc"}
	res2 := tool.ParseSlicess(source)
	if reflect.DeepEqual(res, res2) {
		t.Errorf("ss")
	}
}

func TestDeledeSlice(t *testing.T) {
	type test struct {
		source []string
		target string
		res    []string
	}
	tests := []test{
		{source: []string{"张三", "李四", "王五"}, target: "张三", res: []string{"李四", "王五"}},
		{source: []string{"17620439807", "17620439808", "17620439809"}, target: "17620439807", res: []string{"17620439808", "17620439809"}},
	}
	for _, v := range tests {
		final := tool.DeledeSlice(v.source, v.target)
		if !reflect.DeepEqual(v.res, final) {
			t.Errorf("want:%#v,got:%#v", v.res, final)
		}
	}
}

func TestGetMax(t *testing.T) {
	type test struct {
		source []int
		res    []int
	}
	tests := []test{
		{source: []int{1, 1, 1, 2, 2, 2, 3, 5, 9}, res: []int{1, 2}},
		{source: []int{1, 1, 1, 2, 5, 6}, res: []int{1}},
	}
	for _, v := range tests {
		final := tool.GetMaxint(v.source)
		if !reflect.DeepEqual(v.res, final) {
			t.Errorf("want:%#v,got:%#v", v.res, final)
		}
	}
}
