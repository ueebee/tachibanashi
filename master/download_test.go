package master

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ueebee/tachibanashi/auth"
)

type streamClient struct {
	data string
}

func (c *streamClient) DoJSON(ctx context.Context, method, path string, req, resp any) error {
	return nil
}

func (c *streamClient) DoStream(ctx context.Context, method, path string, req any) (*http.Response, io.Reader, error) {
	reader := strings.NewReader(c.data)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(reader),
		Header:     make(http.Header),
	}
	return resp, reader, nil
}

func (c *streamClient) VirtualURLs() auth.VirtualURLs {
	return auth.VirtualURLs{Master: "https://example.invalid/"}
}

func TestDownloadStoresRecords(t *testing.T) {
	stream := strings.Join([]string{
		`{"sCLMID":"CLMDateZyouhou","sDayKey":"001","sTheDay":"20240101","sUpdateTime":"20240101120000","sUpdateNumber":"2","sDeleteFlag":"0"}`,
		`{"sCLMID":"CLMEventDownloadComplete","sResultCode":"0","sResultText":""}`,
		"",
	}, "\n")

	service := NewService(&streamClient{data: stream})
	store := NewMemoryStore()

	if err := service.Download(context.Background(), store, nil); err != nil {
		t.Fatalf("Download() error = %v", err)
	}

	record, ok := store.Get(MasterDateZyouhou, "001")
	if !ok {
		t.Fatalf("record not stored")
	}
	if got := record.Fields.Value(DateInfoFieldTheDay); got != "20240101" {
		t.Fatalf("record field mismatch: %v", got)
	}
	if record.Meta.Serial != 2 {
		t.Fatalf("record meta serial mismatch: %v", record.Meta.Serial)
	}
}

func TestDownloadRequiresDestination(t *testing.T) {
	service := NewService(&streamClient{data: `{}`})
	err := service.Download(context.Background(), nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDownloadKeepsFields(t *testing.T) {
	message, err := parseDownloadMessage(map[string]json.RawMessage{
		"sCLMID":        json.RawMessage(`"CLMDateZyouhou"`),
		"sDayKey":       json.RawMessage(`"001"`),
		"sTheDay":       json.RawMessage(`"20240101"`),
		"sUpdateTime":   json.RawMessage(`"00000000000000"`),
		"sUpdateNumber": json.RawMessage(`"0"`),
		"sDeleteFlag":   json.RawMessage(`"0"`),
	})
	if err != nil {
		t.Fatalf("parseDownloadMessage() error = %v", err)
	}
	if message.Meta.UpdatedAt != "" {
		t.Fatalf("expected empty UpdatedAt, got %v", message.Meta.UpdatedAt)
	}
	if message.Fields.Value(DateInfoFieldTheDay) != "20240101" {
		t.Fatalf("expected field value")
	}
	if message.Key != "001" {
		t.Fatalf("expected key 001")
	}
}

func TestDownloadStreamContinuesAfterComplete(t *testing.T) {
	stream := strings.Join([]string{
		`{"sCLMID":"CLMDateZyouhou","sDayKey":"001","sTheDay":"20240101","sUpdateTime":"20240101120000","sUpdateNumber":"1","sDeleteFlag":"0"}`,
		`{"sCLMID":"CLMEventDownloadComplete","sResultCode":"0","sResultText":""}`,
		`{"sCLMID":"CLMDateZyouhou","sDayKey":"001","sTheDay":"20240102","sUpdateTime":"20240102120000","sUpdateNumber":"2","sDeleteFlag":"0"}`,
		"",
	}, "\n")

	service := NewService(&streamClient{data: stream})
	store := NewMemoryStore()

	if err := service.DownloadStream(context.Background(), store, nil); err != nil {
		t.Fatalf("DownloadStream() error = %v", err)
	}

	record, ok := store.Get(MasterDateZyouhou, "001")
	if !ok {
		t.Fatalf("record not stored")
	}
	if got := record.Fields.Value(DateInfoFieldTheDay); got != "20240102" {
		t.Fatalf("record field mismatch: %v", got)
	}
	if record.Meta.Serial != 2 {
		t.Fatalf("record meta serial mismatch: %v", record.Meta.Serial)
	}
}
