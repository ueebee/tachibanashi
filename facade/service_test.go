package facade

import (
	"context"
	"errors"
	"testing"

	"github.com/ueebee/tachibanashi/model"
	"github.com/ueebee/tachibanashi/request"
)

type stubRequest struct {
	orders       *request.OrdersSnapshot
	cash         *request.CashPositionsSnapshot
	margin       *request.MarginPositionsSnapshot
	buyingPower  *request.BuyingPowerSnapshot
	marginPower  *request.MarginBuyingPowerSnapshot
	summary      *request.ZanKaiSummaryResponse
	errOrders    error
	errPositions error
	errBalance   error
}

func (s *stubRequest) Orders(ctx context.Context, params request.OrderParams) (*request.OrdersSnapshot, error) {
	if s.errOrders != nil {
		return nil, s.errOrders
	}
	return s.orders, nil
}

func (s *stubRequest) CashPositions(ctx context.Context, issueCode string) (*request.CashPositionsSnapshot, error) {
	if s.errPositions != nil {
		return nil, s.errPositions
	}
	return s.cash, nil
}

func (s *stubRequest) MarginPositions(ctx context.Context, issueCode string) (*request.MarginPositionsSnapshot, error) {
	if s.errPositions != nil {
		return nil, s.errPositions
	}
	return s.margin, nil
}

func (s *stubRequest) BuyingPower(ctx context.Context) (*request.BuyingPowerSnapshot, error) {
	if s.errBalance != nil {
		return nil, s.errBalance
	}
	return s.buyingPower, nil
}

func (s *stubRequest) MarginBuyingPower(ctx context.Context) (*request.MarginBuyingPowerSnapshot, error) {
	if s.errBalance != nil {
		return nil, s.errBalance
	}
	return s.marginPower, nil
}

func (s *stubRequest) ZanKaiSummary(ctx context.Context) (*request.ZanKaiSummaryResponse, error) {
	if s.errBalance != nil {
		return nil, s.errBalance
	}
	return s.summary, nil
}

func TestOrders(t *testing.T) {
	stub := &stubRequest{
		orders: &request.OrdersSnapshot{
			Orders: []model.Order{{ID: "1", Symbol: "6501"}},
		},
	}
	svc := New(stub)

	snapshot, err := svc.Orders(context.Background(), request.OrderParams{"sIssueCode": "6501"})
	if err != nil {
		t.Fatalf("Orders() error = %v", err)
	}
	if len(snapshot.Orders) != 1 {
		t.Fatalf("orders length mismatch: %d", len(snapshot.Orders))
	}
	if snapshot.Orders[0].Symbol != "6501" {
		t.Fatalf("symbol mismatch: %s", snapshot.Orders[0].Symbol)
	}
}

func TestPositions(t *testing.T) {
	stub := &stubRequest{
		cash: &request.CashPositionsSnapshot{
			Positions: []model.Position{{Symbol: "6501"}},
		},
		margin: &request.MarginPositionsSnapshot{
			Positions: []model.Position{{Symbol: "6502"}},
		},
	}
	svc := New(stub)

	snapshot, err := svc.Positions(context.Background(), "")
	if err != nil {
		t.Fatalf("Positions() error = %v", err)
	}
	if len(snapshot.All) != 2 {
		t.Fatalf("all positions mismatch: %d", len(snapshot.All))
	}
	if snapshot.All[0].Symbol != "6501" || snapshot.All[1].Symbol != "6502" {
		t.Fatalf("positions mismatch: %#v", snapshot.All)
	}
}

func TestAccount(t *testing.T) {
	stub := &stubRequest{
		buyingPower: &request.BuyingPowerSnapshot{
			Balance: model.Balance{BuyingPower: 100},
		},
		marginPower: &request.MarginBuyingPowerSnapshot{
			Balance:          model.Balance{BuyingPower: 200},
			MaintenanceRatio: 1.5,
		},
		summary: &request.ZanKaiSummaryResponse{
			Fields: model.Attributes{"sUpdateDate": "20240101"},
		},
	}
	svc := New(stub)

	snapshot, err := svc.Account(context.Background())
	if err != nil {
		t.Fatalf("Account() error = %v", err)
	}
	if snapshot.BuyingPower.BuyingPower != 100 {
		t.Fatalf("buying power mismatch: %d", snapshot.BuyingPower.BuyingPower)
	}
	if snapshot.MarginBuyingPower.BuyingPower != 200 {
		t.Fatalf("margin buying power mismatch: %d", snapshot.MarginBuyingPower.BuyingPower)
	}
	if snapshot.MaintenanceRatio != 1.5 {
		t.Fatalf("maintenance ratio mismatch: %v", snapshot.MaintenanceRatio)
	}
	if snapshot.Summary.Value("sUpdateDate") != "20240101" {
		t.Fatalf("summary mismatch: %s", snapshot.Summary.Value("sUpdateDate"))
	}
}

func TestNilRequester(t *testing.T) {
	svc := New(nil)
	if _, err := svc.Orders(context.Background(), nil); err == nil {
		t.Fatalf("expected error for nil requester")
	}
	if _, err := svc.Positions(context.Background(), ""); err == nil {
		t.Fatalf("expected error for nil requester")
	}
	if _, err := svc.Account(context.Background()); err == nil {
		t.Fatalf("expected error for nil requester")
	}
}

func TestErrorPropagation(t *testing.T) {
	err := errors.New("boom")
	stub := &stubRequest{errOrders: err}
	svc := New(stub)
	if _, got := svc.Orders(context.Background(), nil); got == nil {
		t.Fatalf("expected error")
	}
}
