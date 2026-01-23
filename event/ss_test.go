package event

import "testing"

func TestDecodeEventSS(t *testing.T) {
	raw := "p_no\x021\x01p_date\x022020.06.18-07:30:34.810\x01p_cmd\x02SS\x01p_PV\x02MSGSV\x01p_ENO\x023\x01p_ALT\x020\x01" +
		"p_CT\x0220200618052959\x01p_LK\x021\x01p_SS\x021"

	event, err := DecodeEvent([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}
	ss, ok := event.(SS)
	if !ok {
		t.Fatalf("event type mismatch")
	}
	if ss.Provider != "MSGSV" {
		t.Fatalf("p_PV = %s", ss.Provider)
	}
	if ss.EventNo != "3" {
		t.Fatalf("p_ENO = %s", ss.EventNo)
	}
	if ss.ChangedAt != "20200618052959" {
		t.Fatalf("p_CT = %s", ss.ChangedAt)
	}
	if ss.LoginKind != "1" {
		t.Fatalf("p_LK = %s", ss.LoginKind)
	}
	if ss.SystemStatus != "1" {
		t.Fatalf("p_SS = %s", ss.SystemStatus)
	}
}
