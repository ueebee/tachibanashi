package event

import (
	"net/url"
	"strings"
	"testing"
)

func TestBuildWSURLDefaults(t *testing.T) {
	got, err := BuildWSURL("wss://example.invalid/ws", Params{})
	if err != nil {
		t.Fatalf("BuildWSURL() error = %v", err)
	}

	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}

	query := parsed.Query()
	if query.Get("p_rid") != "0" {
		t.Fatalf("p_rid = %s", query.Get("p_rid"))
	}
	if query.Get("p_board_no") != "1000" {
		t.Fatalf("p_board_no = %s", query.Get("p_board_no"))
	}
	if query.Get("p_eno") != "0" {
		t.Fatalf("p_eno = %s", query.Get("p_eno"))
	}
	if query.Get("p_evt_cmd") != "ST,KP,EC,SS,US" {
		t.Fatalf("p_evt_cmd = %s", query.Get("p_evt_cmd"))
	}
}

func TestBuildWSURLRID22(t *testing.T) {
	params := Params{
		RID:         22,
		BoardNo:     1000,
		Rows:        []int{1, 2},
		IssueCodes:  []string{"6501", "6502"},
		MarketCodes: []string{"00", "00"},
		Cmds:        []Command{CommandFD},
	}
	got, err := BuildWSURL("wss://example.invalid/ws", params)
	if err != nil {
		t.Fatalf("BuildWSURL() error = %v", err)
	}
	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	query := parsed.Query()
	if query.Get("p_rid") != "22" {
		t.Fatalf("p_rid = %s", query.Get("p_rid"))
	}
	if query.Get("p_gyou_no") != "1,2" {
		t.Fatalf("p_gyou_no = %s", query.Get("p_gyou_no"))
	}
	if query.Get("p_issue_code") != "6501,6502" {
		t.Fatalf("p_issue_code = %s", query.Get("p_issue_code"))
	}
	if query.Get("p_mkt_code") != "00,00" {
		t.Fatalf("p_mkt_code = %s", query.Get("p_mkt_code"))
	}
	if query.Get("p_evt_cmd") != "FD" {
		t.Fatalf("p_evt_cmd = %s", query.Get("p_evt_cmd"))
	}

	if strings.Contains(parsed.RawQuery, "%2C") {
		t.Fatalf("unexpected encoded comma: %s", parsed.RawQuery)
	}
}

func TestValidateRejectFDWithRID0(t *testing.T) {
	_, err := BuildWSURL("wss://example.invalid/ws", Params{
		Cmds: []Command{CommandFD},
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestValidateRID21(t *testing.T) {
	_, err := BuildWSURL("wss://example.invalid/ws", Params{
		RID:         21,
		BoardNo:     0,
		Rows:        []int{1},
		IssueCodes:  []string{"6501"},
		MarketCodes: []string{"00"},
		Cmds:        []Command{CommandFD},
	})
	if err != nil {
		t.Fatalf("BuildWSURL() error = %v", err)
	}
}

func TestValidateMismatch(t *testing.T) {
	_, err := BuildWSURL("wss://example.invalid/ws", Params{
		RID:        22,
		BoardNo:    1000,
		Rows:       []int{1},
		IssueCodes: []string{"6501"},
		Cmds:       []Command{CommandFD},
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
