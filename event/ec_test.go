package event

import (
	"fmt"
	"testing"
)

func TestDecodeEventEC(t *testing.T) {
	raw := fmt.Sprintf("p_no\x0270\x01p_date\x022018.12.03-13:22:43.921\x01p_cmd\x02EC\x01p_PV\x02MSGSV\x01p_ENO\x0210507\x01p_ALT\x021\x01p_NT\x02100\x01p_ON\x023000945\x01p_ED\x0220181203\x01p_OON\x020\x01p_OT\x021\x01p_ST\x021\x01p_IC\x022468\x01p_MC\x0200\x01p_BBKB\x021\x01p_CRSJ\x020\x01p_CRPR\x02850.000000\x01p_CRSR\x025300\x01p_ODST\x021\x01p_EXPR\x02851.000000\x01p_EXSR\x0210\x01p_EXDT\x0220181203132243")

	event, err := DecodeEvent([]byte(raw))
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}
	ec, ok := event.(EC)
	if !ok {
		t.Fatalf("event type mismatch")
	}
	if ec.OrderNumber != "3000945" {
		t.Fatalf("p_ON = %s", ec.OrderNumber)
	}
	if ec.ParentOrderNumber != "0" {
		t.Fatalf("p_OON = %s", ec.ParentOrderNumber)
	}
	if ec.OrderType != "1" {
		t.Fatalf("p_OT = %s", ec.OrderType)
	}

	order := ec.Order()
	if order.Price != 850 {
		t.Fatalf("order price = %d", order.Price)
	}
	if order.Quantity != 5300 {
		t.Fatalf("order quantity = %d", order.Quantity)
	}
	if order.Status != "1" {
		t.Fatalf("order status = %s", order.Status)
	}

	exec, ok := ec.Execution()
	if !ok {
		t.Fatalf("expected execution")
	}
	if exec.Price != 851 {
		t.Fatalf("execution price = %d", exec.Price)
	}
	if exec.Quantity != 10 {
		t.Fatalf("execution quantity = %d", exec.Quantity)
	}
	if exec.Time != "20181203132243" {
		t.Fatalf("execution time = %s", exec.Time)
	}
}
