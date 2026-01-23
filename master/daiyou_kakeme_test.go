package master

import (
	"encoding/json"
	"testing"
)

func TestDaiyouKakemeUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMDaiyouKakeme",
		"sSystemKouzaKubun":"102",
		"sIssueCode":"6501",
		"sTekiyouDay":"20220422",
		"sHosyokinDaiyoKakeme":"80.000000"
	}`)

	var info DaiyouKakeme
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal daiyou kakeme: %v", err)
	}
	if info.IssueCode() != "6501" {
		t.Fatalf("unexpected issue code: %s", info.IssueCode())
	}
	if info.TekiyouDay() != "20220422" {
		t.Fatalf("unexpected tekiyou day: %s", info.TekiyouDay())
	}
	if got := info.Fields.Value(DaiyouKakemeFieldHosyokinKakeme); got != "80.000000" {
		t.Fatalf("unexpected kakeme: %s", got)
	}

	key, ok := MasterKey(MasterDaiyouKakeme, info.Fields)
	want := JoinIndex("102", "6501", "20220422")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
