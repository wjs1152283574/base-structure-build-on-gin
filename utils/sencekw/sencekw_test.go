package sencekw

import "testing"

func TestGetSencitiveResponse(t *testing.T) {
	type test struct {
		source string
		target string
		res    string
	}
	tests := []test{
		{source: "阿时间打军事基地你，。，。妈。。。。的", target: "阿时间打军事基地******。。。。的"},
		{source: "asdasdaFuc*k", target: "asdasda*****"},
	}
	for _, v := range tests {
		final := GetSencitiveResponse(v.source)
		if v.target != final {
			t.Errorf("want:%#v,got:%#v", v.res, final)
		}
	}
}
