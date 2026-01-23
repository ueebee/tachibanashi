package master

import (
	"encoding/json"
	"testing"
)

func TestOrderErrReasonUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMOrderErrReason",
		"sErrReasonCode":"-110007",
		"sErrReasonText":"error"
	}`)

	var info OrderErrReason
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal order err reason: %v", err)
	}
	if info.Code() != "-110007" {
		t.Fatalf("unexpected code: %s", info.Code())
	}
	if info.Text() != "error" {
		t.Fatalf("unexpected text: %s", info.Text())
	}

	key, ok := MasterKey(MasterOrderErrReason, info.Fields)
	if !ok || key != "-110007" {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
