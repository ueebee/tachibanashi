package master

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestNewsHeadResponseUnmarshal(t *testing.T) {
	headline := "Market headline"
	encoded := base64.StdEncoding.EncodeToString([]byte(headline))

	raw := []byte(`{
		"sCLMID":"CLMMfdsGetNewsHead",
		"p_REC_MAX":"1",
		"aCLMMfdsNewsHead":[{"p_ID":"id","p_DT":"20230512","p_TM":"1530","p_CGL":"120","p_GNL":"62104","p_ISL":"4838","p_HDL":"` + encoded + `"}]
	}`)

	var resp NewsHeadResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal news head: %v", err)
	}
	if resp.RecordMax != "1" {
		t.Fatalf("unexpected record max: %s", resp.RecordMax)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("unexpected entry count: %d", len(resp.Entries))
	}
	entry := resp.Entries[0]
	if entry.ID != "id" {
		t.Fatalf("unexpected id: %s", entry.ID)
	}
	if entry.Headline != headline {
		t.Fatalf("unexpected headline: %s", entry.Headline)
	}
	if got := entry.Fields.Value("p_HDL"); got != headline {
		t.Fatalf("unexpected headline field: %s", got)
	}
}

func TestNewsBodyResponseUnmarshal(t *testing.T) {
	headline := "Market headline"
	body := "Market body"
	encodedHeadline := base64.StdEncoding.EncodeToString([]byte(headline))
	encodedBody := base64.StdEncoding.EncodeToString([]byte(body))

	raw := []byte(`{
		"sCLMID":"CLMMfdsGetNewsBody",
		"p_ID":"id",
		"aCLMMfdsNewsBody":[{"p_ID":"id","p_DT":"20230512","p_TM":"1530","p_CGL":"120","p_GNL":"62104","p_ISL":"4838","p_HDL":"` + encodedHeadline + `","p_TX":"` + encodedBody + `"}]
	}`)

	var resp NewsBodyResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("unmarshal news body: %v", err)
	}
	if resp.RequestID != "id" {
		t.Fatalf("unexpected request id: %s", resp.RequestID)
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("unexpected entry count: %d", len(resp.Entries))
	}
	entry := resp.Entries[0]
	if entry.Body != body {
		t.Fatalf("unexpected body: %s", entry.Body)
	}
	if got := entry.Fields.Value("p_TX"); got != body {
		t.Fatalf("unexpected body field: %s", got)
	}
}
