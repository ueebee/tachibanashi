package master

import (
	"encoding/json"
	"testing"
)

func TestIssueSizyouMstKabuUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMIssueSizyouMstKabu",
		"sIssueCode":"6501",
		"sZyouzyouSizyou":"00",
		"sYobineTaniNumber":"101",
		"sUpdateNumber":"2"
	}`)

	var info IssueSizyouMstKabu
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal issue sizyou kabu: %v", err)
	}
	if info.IssueCode() != "6501" {
		t.Fatalf("unexpected issue code: %s", info.IssueCode())
	}
	if info.MarketCode() != "00" {
		t.Fatalf("unexpected market code: %s", info.MarketCode())
	}
	if got := info.Fields.Value(IssueSizyouKabuFieldYobineTaniNumber); got != "101" {
		t.Fatalf("unexpected yobine tani: %s", got)
	}

	key, ok := MasterKey(MasterIssueSizyouMstKabu, info.Fields)
	want := JoinIndex("6501", "00")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
