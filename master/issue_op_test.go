package master

import (
	"encoding/json"
	"testing"
)

func TestIssueMstOpUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMIssueMstOp",
		"sIssueCode":"130030005",
		"sIssueName":"TPX 2503P2000",
		"sKousiPrice":"2000.000000",
		"sUpdateNumber":"5"
	}`)

	var info IssueMstOp
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal issue op: %v", err)
	}
	if info.IssueCode() != "130030005" {
		t.Fatalf("unexpected issue code: %s", info.IssueCode())
	}
	if got := info.Fields.Value(IssueOpFieldKousiPrice); got != "2000.000000" {
		t.Fatalf("unexpected kousi price: %s", got)
	}

	key, ok := MasterKey(MasterIssueMstOp, info.Fields)
	if !ok || key != "130030005" {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
