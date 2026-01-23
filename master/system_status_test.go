package master

import (
	"encoding/json"
	"testing"
)

func TestSystemStatusUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMSystemStatus",
		"sSystemStatusKey":"001",
		"sLoginKyokaKubun":"1",
		"sSystemStatus":"1",
		"sCreateTime":"",
		"sUpdateTime":"",
		"sUpdateNumber":"",
		"sDeleteFlag":"",
		"sDeleteTime":""
	}`)

	var status SystemStatus
	if err := json.Unmarshal(raw, &status); err != nil {
		t.Fatalf("unmarshal system status: %v", err)
	}
	if status.Key() != "001" {
		t.Fatalf("unexpected key: %s", status.Key())
	}
	if got := status.Fields.Value(SystemStatusFieldLoginKyoka); got != "1" {
		t.Fatalf("unexpected login flag: %s", got)
	}

	key, ok := MasterKey(MasterSystemStatus, status.Fields)
	if !ok || key != "001" {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
