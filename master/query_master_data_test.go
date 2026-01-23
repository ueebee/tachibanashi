package master

import (
	"encoding/json"
	"testing"
)

func TestMasterDataResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMMfdsGetMasterData",
		"CLMIssueMstKabu":[{"sIssueCode":"6501","sIssueName":"Test"}],
		"CLMOrderErrReason":[{"sErrReasonCode":"-110007","sErrReasonText":"error"}]
	}`)

	var resp MasterDataResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal master data: %v", err)
	}
	issues := resp.Data[MasterIssueMstKabu]
	if len(issues) != 1 {
		t.Fatalf("unexpected issue count: %d", len(issues))
	}
	if got := issues[0].Value("sIssueName"); got != "Test" {
		t.Fatalf("unexpected issue name: %s", got)
	}
	reasons := resp.Data[MasterOrderErrReason]
	if len(reasons) != 1 {
		t.Fatalf("unexpected reason count: %d", len(reasons))
	}
	if got := reasons[0].Value("sErrReasonText"); got != "error" {
		t.Fatalf("unexpected reason text: %s", got)
	}
}
