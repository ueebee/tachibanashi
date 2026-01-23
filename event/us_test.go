package event

import "testing"

func TestDecodeEventUS(t *testing.T) {
	raw := "p_no\x022\x01p_date\x022018.12.03-11:34:51.557\x01p_cmd\x02US\x01p_PV\x02MSGSV\x01p_ENO\x025227\x01p_ALT\x020\x01" +
		"p_CT\x0220181203075545\x01p_MC\x0200\x01p_GSCD\x02\x01p_SHSB\x02\x01p_UC\x0201\x01p_UU\x020101\x01p_EDK\x020\x01p_US\x02100"

	event, err := DecodeEvent([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}
	us, ok := event.(US)
	if !ok {
		t.Fatalf("event type mismatch")
	}
	if us.Provider != "MSGSV" {
		t.Fatalf("p_PV = %s", us.Provider)
	}
	if us.EventNo != "5227" {
		t.Fatalf("p_ENO = %s", us.EventNo)
	}
	if us.ChangedAt != "20181203075545" {
		t.Fatalf("p_CT = %s", us.ChangedAt)
	}
	if us.MarketCode != "00" {
		t.Fatalf("p_MC = %s", us.MarketCode)
	}
	if us.UnderlyingCode != "" {
		t.Fatalf("p_GSCD = %s", us.UnderlyingCode)
	}
	if us.InstrumentKind != "" {
		t.Fatalf("p_SHSB = %s", us.InstrumentKind)
	}
	if us.OperationCode != "01" {
		t.Fatalf("p_UC = %s", us.OperationCode)
	}
	if us.OperationUnit != "0101" {
		t.Fatalf("p_UU = %s", us.OperationUnit)
	}
	if us.BusinessDayKind != "0" {
		t.Fatalf("p_EDK = %s", us.BusinessDayKind)
	}
	if us.OperationStatus != "100" {
		t.Fatalf("p_US = %s", us.OperationStatus)
	}
}
