package master

import (
	"encoding/json"
	"testing"
)

func TestIssueDetailResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMMfdsGetIssueDetail",
		"aCLMMfdsIssueDetail":[{"sIssueCode":"6501","pBPSB":"6155.38"}]
	}`)

	var resp IssueDetailResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal issue detail: %v", err)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("unexpected entry count: %d", len(resp.Entries))
	}
	entry := resp.Entries[0]
	if entry.IssueCode != "6501" {
		t.Fatalf("unexpected issue code: %s", entry.IssueCode)
	}
	if got := entry.Fields.Value("pBPSB"); got != "6155.38" {
		t.Fatalf("unexpected field: %s", got)
	}
}

func TestSyoukinZanResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMMfdsGetSyoukinZan",
		"aCLMMfdsSyoukinZan":[{"sIssueCode":"6501","pSFD":"2024/12/30"}]
	}`)

	var resp SyoukinZanResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal syoukin zan: %v", err)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("unexpected entry count: %d", len(resp.Entries))
	}
	entry := resp.Entries[0]
	if entry.IssueCode != "6501" {
		t.Fatalf("unexpected issue code: %s", entry.IssueCode)
	}
	if got := entry.Fields.Value("pSFD"); got != "2024/12/30" {
		t.Fatalf("unexpected field: %s", got)
	}
}

func TestShinyouZanResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMMfdsGetShinyouZan",
		"aCLMMfdsShinyouZan":[{"sIssueCode":"6501","pMBD":"2024/12/20"}]
	}`)

	var resp ShinyouZanResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal shinyou zan: %v", err)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("unexpected entry count: %d", len(resp.Entries))
	}
	entry := resp.Entries[0]
	if entry.IssueCode != "6501" {
		t.Fatalf("unexpected issue code: %s", entry.IssueCode)
	}
	if got := entry.Fields.Value("pMBD"); got != "2024/12/20" {
		t.Fatalf("unexpected field: %s", got)
	}
}

func TestHibuInfoResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"sCLMID":"CLMMfdsGetHibuInfo",
		"aCLMMfdsHibuInfo":[{"sIssueCode":"6501","pBWRQ":"0.05"}]
	}`)

	var resp HibuInfoResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal hibu info: %v", err)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("unexpected entry count: %d", len(resp.Entries))
	}
	entry := resp.Entries[0]
	if entry.IssueCode != "6501" {
		t.Fatalf("unexpected issue code: %s", entry.IssueCode)
	}
	if got := entry.Fields.Value("pBWRQ"); got != "0.05" {
		t.Fatalf("unexpected field: %s", got)
	}
}
