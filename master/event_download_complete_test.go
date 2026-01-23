package master

import (
	"encoding/json"
	"testing"
)

func TestEventDownloadCompleteUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMEventDownloadComplete"
	}`)

	var info EventDownloadComplete
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal event download complete: %v", err)
	}
	if got := info.Fields.Value("sCLMID"); got != "CLMEventDownloadComplete" {
		t.Fatalf("unexpected CLMID: %s", got)
	}
}
