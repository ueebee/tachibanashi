package master

import (
	"encoding/json"
	"testing"
)

func TestIssueMstSakUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMIssueMstSak",
		"sIssueCode":"160060018",
		"sIssueName":"NK225 2506",
		"sGensisanCode":"101",
		"sUpdateNumber":"4"
	}`)

	var info IssueMstSak
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal issue sak: %v", err)
	}
	if info.IssueCode() != "160060018" {
		t.Fatalf("unexpected issue code: %s", info.IssueCode())
	}
	if got := info.Fields.Value(IssueSakFieldGensisanCode); got != "101" {
		t.Fatalf("unexpected gensisan code: %s", got)
	}

	key, ok := MasterKey(MasterIssueMstSak, info.Fields)
	if !ok || key != "160060018" {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
