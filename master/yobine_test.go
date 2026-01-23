package master

import (
	"encoding/json"
	"testing"
)

func TestYobineUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMYobine",
		"sYobineTaniNumber":"101",
		"sTekiyouDay":"20140101",
		"sKizunPrice_1":"3000.000000",
		"sYobineTanka_1":"1.000000",
		"sDecimal_1":"0"
	}`)

	var info Yobine
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal yobine: %v", err)
	}
	if info.TaniNumber() != "101" {
		t.Fatalf("unexpected tani number: %s", info.TaniNumber())
	}
	if info.TekiyouDay() != "20140101" {
		t.Fatalf("unexpected tekiyou day: %s", info.TekiyouDay())
	}
	if got := info.Fields.Value(YobineFieldKizunPrice1); got != "3000.000000" {
		t.Fatalf("unexpected kizun price: %s", got)
	}

	key, ok := MasterKey(MasterYobine, info.Fields)
	want := JoinIndex("101", "20140101")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
