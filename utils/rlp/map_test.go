package rlp

import (
	"log"
	"testing"
)

func Test_Map(t *testing.T) {
	TeMap := make(map[string]interface{})
	TeMap["simon"] = []byte("23")
	TeMap["tom"] = []byte("567")

	val, err := EncodeToBytes(TeMap)
	if err != nil {
		t.Errorf("Encode to bytes err={%v}", err)
	}

	t.Logf("After encode data={%v}", val)

	afde := make(map[string]interface{})
	if err := DecodeBytes(val, &afde); err != nil {
		t.Errorf("Decode to bytes err={%v}", err)
	}

	t.Logf("After decode data={%v}", afde)
}

func Test_StructMap(t *testing.T) {
	type testruct struct {
		Value uint64
		Str   string
		Save  map[string]interface{}
	}

	var res testruct
	res.Value = 345678
	res.Str = "datx"

	TeMap := make(map[string]interface{})
	TeMap["simon"] = []byte("23")
	TeMap["tom"] = []byte("567")
	TeMap["go"] = []byte("gh")

	res.Save = TeMap

	val, err := EncodeToBytes(res)
	if err != nil {
		t.Errorf("Encode to bytes err={%v}", err)
	}

	t.Logf("After encode data={%v}", val)

	afde := new(testruct)
	afde.Save = make(map[string]interface{}) //very important
	if err := DecodeBytes(val, afde); err != nil {
		t.Errorf("Decode to bytes err={%v}", err)
	}

	log.Printf("%v\n", res)
	t.Logf("After decode data={%v}", afde)
}
