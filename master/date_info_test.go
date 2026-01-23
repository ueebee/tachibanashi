package master

import (
	"encoding/json"
	"testing"
)

func TestDateInfoUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMDateZyouhou",
		"sDayKey":"001",
		"sMaeEigyouDay_1":"20231031",
		"sMaeEigyouDay_2":"20231030",
		"sMaeEigyouDay_3":"20231027",
		"sTheDay":"20231101",
		"sYokuEigyouDay_1":"20231102",
		"sYokuEigyouDay_2":"20231106",
		"sYokuEigyouDay_3":"20231107",
		"sYokuEigyouDay_4":"20231108",
		"sYokuEigyouDay_5":"20231109",
		"sYokuEigyouDay_6":"20231110",
		"sYokuEigyouDay_7":"20231113",
		"sYokuEigyouDay_8":"20231114",
		"sYokuEigyouDay_9":"20231115",
		"sYokuEigyouDay_10":"20231116",
		"sKabuUkewatasiDay":"20231106",
		"sKabuKariUkewatasiDay":"20231107",
		"sBondUkewatasiDay":"20231106"
	}`)

	var info DateInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal date info: %v", err)
	}
	if info.DayKey() != "001" {
		t.Fatalf("unexpected day key: %s", info.DayKey())
	}
	if info.TheDay() != "20231101" {
		t.Fatalf("unexpected the day: %s", info.TheDay())
	}
	if got := info.Fields.Value(DateInfoFieldYokuEigyouDay10); got != "20231116" {
		t.Fatalf("unexpected next day 10: %s", got)
	}

	key, ok := MasterKey(MasterDateZyouhou, info.Fields)
	if !ok || key != "001" {
		t.Fatalf("unexpected master key: %v %s", ok, key)
	}
}
