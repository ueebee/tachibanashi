package event

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ueebee/tachibanashi/model"
)

func TestDecodeEventFD(t *testing.T) {
	hexValue := hex.EncodeToString([]byte("TPM"))
	raw := fmt.Sprintf("p_no\x021\x01p_date\x022018.12.03-13:11:22.122\x01p_cmd\x02FD\x01p_1_DPP\x026129\x01t_1_DPP:T\x0214:10\x01x_2_LISS\x02%s", hexValue)

	event, err := DecodeEvent([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}
	fd, ok := event.(FD)
	if !ok {
		t.Fatalf("event type mismatch")
	}
	if len(fd.Rows) != 2 {
		t.Fatalf("rows = %d", len(fd.Rows))
	}
	if got := fd.Rows[0].Fields.Value("pDPP"); got != "6129" {
		t.Fatalf("pDPP = %s", got)
	}
	if got := fd.Rows[0].Fields.Value("tDPP:T"); got != "14:10" {
		t.Fatalf("tDPP:T = %s", got)
	}
	if got := fd.Rows[1].Fields.Value("xLISS"); got != "TPM" {
		t.Fatalf("xLISS = %s", got)
	}
}

func TestQuoteBookApply(t *testing.T) {
	book := NewQuoteBook()
	first := FD{
		Rows: []FDRow{
			{Row: 1, Fields: model.Attributes{"pDPP": "100", "tDPP:T": "10:00"}},
		},
	}
	updated := book.Apply(first)
	if len(updated) != 1 {
		t.Fatalf("updated size = %d", len(updated))
	}
	if got := updated[0].Value("pDPP"); got != "100" {
		t.Fatalf("pDPP = %s", got)
	}

	second := FD{
		Rows: []FDRow{
			{Row: 1, Fields: model.Attributes{"pDV": "200"}},
		},
	}
	updated = book.Apply(second)
	if got := updated[0].Value("pDPP"); got != "100" {
		t.Fatalf("pDPP = %s", got)
	}
	if got := updated[0].Value("pDV"); got != "200" {
		t.Fatalf("pDV = %s", got)
	}

	snapshot := book.Snapshot()
	if len(snapshot) != 1 {
		t.Fatalf("snapshot size = %d", len(snapshot))
	}
	if got := snapshot[0].Value("pDV"); got != "200" {
		t.Fatalf("snapshot pDV = %s", got)
	}
}

func TestFDQuotes(t *testing.T) {
	fd := FD{
		Rows: []FDRow{
			{Row: 2, Fields: model.Attributes{"pDPP": "999"}},
		},
	}
	quotes := fd.Quotes(map[int]string{2: "6501"})
	if len(quotes) != 1 {
		t.Fatalf("quotes size = %d", len(quotes))
	}
	if quotes[0].Symbol != "6501" {
		t.Fatalf("symbol = %s", quotes[0].Symbol)
	}
}
