package price

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ueebee/tachibanashi/auth"
)

type mockClient struct {
	urls       auth.VirtualURLs
	lastMethod string
	lastPath   string
	lastReq    any
	response   []byte
}

func (m *mockClient) DoJSON(ctx context.Context, method, path string, req, resp any) error {
	m.lastMethod = method
	m.lastPath = path
	m.lastReq = req
	if m.response != nil {
		return json.Unmarshal(m.response, resp)
	}
	return nil
}

func (m *mockClient) VirtualURLs() auth.VirtualURLs {
	return m.urls
}

func TestSnapshotBuildsRequestAndParsesEntries(t *testing.T) {
	response := []byte(`{
		"p_no":"1",
		"p_sd_date":"2020.01.02-03:04:05.000",
		"sCLMID":"CLMMfdsGetMarketPrice",
		"sResultCode":"0",
		"sResultText":"",
		"aCLMMfdsMarketPrice":[
			{"sIssueCode":"6501","pDPP":"5179","tDPP:T":"13:59","pPRP":"5197"},
			{"sIssueCode":"6502","pDPP":"","tDPP:T":"","pPRP":""}
		]
	}`)

	client := &mockClient{
		urls:     auth.VirtualURLs{Price: "https://example.invalid/price"},
		response: response,
	}
	svc := NewService(client)

	resp, err := svc.Snapshot(context.Background(), []string{" 6501 ", "6502, 6503"}, []string{"pDPP", "tDPP:T, pPRP"})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}

	if client.lastMethod != http.MethodGet {
		t.Fatalf("method mismatch: %s", client.lastMethod)
	}
	if client.lastPath != "https://example.invalid/price" {
		t.Fatalf("path mismatch: %s", client.lastPath)
	}

	req, ok := client.lastReq.(*MarketPriceRequest)
	if !ok {
		t.Fatalf("request type mismatch")
	}
	if req.CLMID != clmMarketPrice {
		t.Fatalf("CLMID mismatch: %s", req.CLMID)
	}
	if req.TargetIssueCode != "6501,6502,6503" {
		t.Fatalf("TargetIssueCode mismatch: %s", req.TargetIssueCode)
	}
	if req.TargetColumn != "pDPP,tDPP:T,pPRP" {
		t.Fatalf("TargetColumn mismatch: %s", req.TargetColumn)
	}

	if len(resp.Prices) != 2 {
		t.Fatalf("prices length mismatch: %d", len(resp.Prices))
	}
	first := resp.Prices[0]
	if first.IssueCode != "6501" {
		t.Fatalf("issue code mismatch: %s", first.IssueCode)
	}
	if got := first.Fields.Value("pDPP"); got != "5179" {
		t.Fatalf("last price mismatch: %s", got)
	}
}

func TestQuoteSnapshotConverts(t *testing.T) {
	response := []byte(`{
		"sCLMID":"CLMMfdsGetMarketPrice",
		"aCLMMfdsMarketPrice":[
			{"sIssueCode":"6501","pDPP":"5179","pPRP":"5197"}
		]
	}`)

	client := &mockClient{
		urls:     auth.VirtualURLs{Price: "https://example.invalid/price"},
		response: response,
	}
	svc := NewService(client)

	snapshot, err := svc.QuoteSnapshot(context.Background(), []string{"6501"}, []string{"pDPP", "pPRP"})
	if err != nil {
		t.Fatalf("QuoteSnapshot() error = %v", err)
	}
	if snapshot.Raw == nil {
		t.Fatalf("raw response missing")
	}
	if len(snapshot.Quotes) != 1 {
		t.Fatalf("quotes length mismatch: %d", len(snapshot.Quotes))
	}
	if snapshot.Quotes[0].Symbol != "6501" {
		t.Fatalf("symbol mismatch: %s", snapshot.Quotes[0].Symbol)
	}
	if got := snapshot.Quotes[0].Fields.Value("pPRP"); got != "5197" {
		t.Fatalf("prev close mismatch: %s", got)
	}
}

func TestHistoryBuildsRequestAndParsesEntries(t *testing.T) {
	response := []byte(`{
		"sCLMID":"CLMMfdsGetMarketPriceHistory",
		"sIssueCode":"6501",
		"sSizyouC":"00",
		"aCLMMfdsGetMarketPriceHistory":[
			{"sDate":"20240101","pDOP":"100","pDHP":"110","pDLP":"90","pDPP":"105","pDV":"1000"}
		]
	}`)

	client := &mockClient{
		urls:     auth.VirtualURLs{Price: "https://example.invalid/price"},
		response: response,
	}
	svc := NewService(client)

	resp, err := svc.History(context.Background(), "6501", "00")
	if err != nil {
		t.Fatalf("History() error = %v", err)
	}

	req, ok := client.lastReq.(*MarketPriceHistoryRequest)
	if !ok {
		t.Fatalf("request type mismatch")
	}
	if req.CLMID != clmMarketPriceHistory {
		t.Fatalf("CLMID mismatch: %s", req.CLMID)
	}
	if req.IssueCode != "6501" {
		t.Fatalf("IssueCode mismatch: %s", req.IssueCode)
	}
	if req.MarketCode != "00" {
		t.Fatalf("MarketCode mismatch: %s", req.MarketCode)
	}

	if resp.IssueCode != "6501" {
		t.Fatalf("response issue code mismatch: %s", resp.IssueCode)
	}
	if resp.MarketCode != "00" {
		t.Fatalf("response market code mismatch: %s", resp.MarketCode)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("entries length mismatch: %d", len(resp.Entries))
	}
	entry := resp.Entries[0]
	if entry.Date != "20240101" {
		t.Fatalf("date mismatch: %s", entry.Date)
	}
	if got := entry.Fields.Value("pDPP"); got != "105" {
		t.Fatalf("close price mismatch: %s", got)
	}
}

func TestHistoryRejectsMultipleIssues(t *testing.T) {
	client := &mockClient{
		urls: auth.VirtualURLs{Price: "https://example.invalid/price"},
	}
	svc := NewService(client)

	if _, err := svc.History(context.Background(), "6501,6502", ""); err == nil {
		t.Fatalf("expected error for multiple issue codes")
	}
}
