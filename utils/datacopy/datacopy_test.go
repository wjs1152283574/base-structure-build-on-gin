package datacopy

import (
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
	if err := DataCopy(&p, &s); err != nil || s.Name == "" || s.Test != "" {
		t.Error(err)
	}
}
