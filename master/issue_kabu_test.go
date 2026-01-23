package master

import (
	"encoding/json"
	"testing"
)

func TestIssueMstKabuUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMIssueMstKabu",
		"sIssueCode":"6501",
		"sIssueName":"KYOKUYO",
		"sGyousyuCode":"0050",
		"sUpdateNumber":"1"
	}`)

	var info IssueMstKabu
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal issue kabu: %v", err)
	}
	if info.IssueCode() != "6501" {
		t.Fatalf("unexpected issue code: %s", info.IssueCode())
	}
	if info.IssueName() != "KYOKUYO" {
		t.Fatalf("unexpected issue name: %s", info.IssueName())
	}
	if got := info.Fields.Value(IssueKabuFieldGyousyuCode); got != "0050" {
		t.Fatalf("unexpected gyousyu code: %s", got)
	}

	key, ok := MasterKey(MasterIssueMstKabu, info.Fields)
	if !ok || key != "6501" {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
