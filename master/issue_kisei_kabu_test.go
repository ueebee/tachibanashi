package master

import (
	"encoding/json"
	"testing"
)

func TestIssueSizyouKiseiKabuUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMIssueSizyouKiseiKabu",
		"sSystemKouzaKubun":"102",
		"sIssueCode":"6501",
		"sZyouzyouSizyou":"00",
		"sTeisiKubun":"1",
		"sUpdateNumber":"3"
	}`)

	var info IssueSizyouKiseiKabu
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal issue sizyou kisei kabu: %v", err)
	}
	if info.IssueCode() != "6501" {
		t.Fatalf("unexpected issue code: %s", info.IssueCode())
	}
	if info.MarketCode() != "00" {
		t.Fatalf("unexpected market code: %s", info.MarketCode())
	}
	if got := info.Fields.Value(IssueKiseiKabuFieldTeisiKubun); got != "1" {
		t.Fatalf("unexpected teisi kubun: %s", got)
	}

	key, ok := MasterKey(MasterIssueSizyouKiseiKabu, info.Fields)
	want := JoinIndex("102", "6501", "00")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
