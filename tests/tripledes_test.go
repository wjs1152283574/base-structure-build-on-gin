package tests

import (
	tripledes "goweb/utils/tripleDES"
	"testing"
)

var tdes = tripledes.TripleDES{
	Key: "hahaasdasdashahaasdasdas",
	Iv:  "qakmvrex",
}

func TestEncrypt(t *testing.T) {
	sources := "casso"

	if res, err := tdes.Encrypt(sources); err != nil {
		t.Error(res, err) // RkF1WyQ= <nil>
	}
}

func TestDecrypt(t *testing.T) {
	sources := "RkF1WyQ="
	if res, err := tdes.Decrypt(sources); err != nil {
		t.Error(res, err) // casso <nil>
	}
}
