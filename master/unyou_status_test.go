package master

import (
	"encoding/json"
	"testing"
)

func TestUnyouStatusUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMUnyouStatus",
		"sSystemKouzaKubun":"102",
		"sUnyouCategory":"01",
		"sUnyouUnit":"0101",
		"sEigyouDayC":"0",
		"sUnyouStatus":"001",
		"sTaisyouGyoumu":"04",
		"sGyoumuZyoutai":"001",
		"sCreateTime":"",
		"sUpdateTime":"",
		"sUpdateNumber":"",
		"sDeleteFlag":"",
		"sDeleteTime":"",
		"sEventName":"",
		"sMeyasuTime":""
	}`)

	var status UnyouStatus
	if err := json.Unmarshal(raw, &status); err != nil {
		t.Fatalf("unmarshal unyou status: %v", err)
	}
	if got := status.Fields.Value(UnyouStatusFieldUnyouStatus); got != "001" {
		t.Fatalf("unexpected status: %s", got)
	}

	key, ok := MasterKey(MasterUnyouStatus, status.Fields)
	want := JoinIndex("102", "01", "0101", "0", "001", "04")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}

func TestUnyouStatusKabuUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMUnyouStatusKabu",
		"sSystemKouzaKubun":"102",
		"sZyouzyouSizyou":"00",
		"sUnyouCategory":"01",
		"sUnyouUnit":"0101",
		"sEigyouDayC":"0",
		"sUnyouStatus":"001",
		"sCreateTime":"",
		"sUpdateTime":"",
		"sUpdateNumber":"",
		"sDeleteFlag":"",
		"sDeleteTime":""
	}`)

	var status UnyouStatusKabu
	if err := json.Unmarshal(raw, &status); err != nil {
		t.Fatalf("unmarshal unyou status kabu: %v", err)
	}
	if got := status.Fields.Value(UnyouStatusFieldZyouzyouSizyou); got != "00" {
		t.Fatalf("unexpected market: %s", got)
	}

	key, ok := MasterKey(MasterUnyouStatusKabu, status.Fields)
	want := JoinIndex("102", "00", "01", "0101", "0")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}

func TestUnyouStatusHaseiUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMUnyouStatusHasei",
		"sSystemKouzaKubun":"102",
		"sZyouzyouSizyou":"01",
		"sGensisanCode":"101",
		"sSyouhinType":"03",
		"sUnyouCategory":"02",
		"sUnyouUnit":"0201",
		"sEigyouDayC":"0",
		"sUnyouStatus":"001",
		"sCreateTime":"",
		"sUpdateTime":"",
		"sUpdateNumber":"",
		"sDeleteFlag":"",
		"sDeleteTime":""
	}`)

	var status UnyouStatusHasei
	if err := json.Unmarshal(raw, &status); err != nil {
		t.Fatalf("unmarshal unyou status hasei: %v", err)
	}
	if got := status.Fields.Value(UnyouStatusFieldGensisanCode); got != "101" {
		t.Fatalf("unexpected gensisan code: %s", got)
	}

	key, ok := MasterKey(MasterUnyouStatusHasei, status.Fields)
	want := JoinIndex("102", "01", "101", "03", "02", "0201", "0")
	if !ok || key != want {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
