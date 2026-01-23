package master

import (
	"encoding/json"
	"testing"
)

func TestIssueSizyouKiseiHaseiUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMIssueSizyouKiseiHasei",
		"sSystemKouzaKubun":"102",
		"sIssueCode":"160060018",
		"sZyouzyouSizyou":"01",
		"sTeisiKubun":"1",
		"sUpdateNumber":"6"
	}`)

	var info IssueSizyouKiseiHasei
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal issue sizyou kisei hasei: %v", err)
	}
	if info.IssueCode() != "160060018" {
		t.Fatalf("unexpected issue code: %s", info.IssueCode())
	}
	if info.MarketCode() != "01" {
		t.Fatalf("unexpected market code: %s", info.MarketCode())
	}
	if got := info.Fields.Value(IssueKiseiHaseiFieldTeisiKubun); got != "1" {
		t.Fatalf("unexpected teisi kubun: %s", got)
	}

	key, ok := MasterKey(MasterIssueSizyouKiseiHasei, info.Fields)
	want := JoinIndex("102", "160060018", "01")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
