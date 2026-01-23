package master

import (
	"encoding/json"
	"testing"
)

func TestHosyoukinMstUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMHosyoukinMst",
		"sSystemKouzaKubun":"102",
		"sIssueCode":"6501",
		"sZyouzyouSizyou":"00",
		"sHenkouDay":"20230110",
		"sDaiyoHosyokinRitu":"60.000000",
		"sGenkinHosyokinRitu":"0.000000"
	}`)

	var info HosyoukinMst
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal hosyoukin mst: %v", err)
	}
	if info.IssueCode() != "6501" {
		t.Fatalf("unexpected issue code: %s", info.IssueCode())
	}
	if info.MarketCode() != "00" {
		t.Fatalf("unexpected market code: %s", info.MarketCode())
	}
	if got := info.Fields.Value(HosyoukinMstFieldDaiyoRitu); got != "60.000000" {
		t.Fatalf("unexpected daiyo ritu: %s", got)
	}

	key, ok := MasterKey(MasterHosyoukinMst, info.Fields)
	want := JoinIndex("102", "6501", "00", "20230110")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
