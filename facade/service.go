package facade

import (
	"context"
	"errors"

	"github.com/ueebee/tachibanashi/model"
	"github.com/ueebee/tachibanashi/request"
)

type Requester interface {
	Orders(ctx context.Context, params request.OrderParams) (*request.OrdersSnapshot, error)
	CashPositions(ctx context.Context, issueCode string) (*request.CashPositionsSnapshot, error)
	MarginPositions(ctx context.Context, issueCode string) (*request.MarginPositionsSnapshot, error)
	BuyingPower(ctx context.Context) (*request.BuyingPowerSnapshot, error)
	MarginBuyingPower(ctx context.Context) (*request.MarginBuyingPowerSnapshot, error)
	ZanKaiSummary(ctx context.Context) (*request.ZanKaiSummaryResponse, error)
}

type Service struct {
	req Requester
}

func New(req Requester) *Service {
	return &Service{req: req}
}

type OrdersSnapshot struct {
	Orders []model.Order
	Raw    *request.OrderListResponse
}

type PositionsSnapshot struct {
	Cash   []model.Position
	Margin []model.Position
	All    []model.Position
	Raw    *PositionsRaw
}

type PositionsRaw struct {
	Cash   *request.GenbutuKabuListResponse
	Margin *request.ShinyouTategyokuListResponse
}

type AccountSnapshot struct {
	BuyingPower       model.Balance
	MarginBuyingPower model.Balance
	MaintenanceRatio  float64
	Summary           model.Attributes
	Raw               *AccountRaw
}

type AccountRaw struct {
	BuyingPower       *request.ZanKaiKanougakuResponse
	MarginBuyingPower *request.ZanShinkiKanoIjirituResponse
	Summary           *request.ZanKaiSummaryResponse
}

func (s *Service) Orders(ctx context.Context, params request.OrderParams) (*OrdersSnapshot, error) {
	if s.req == nil {
		return nil, errors.New("tachibanashi: facade request service is nil")
	}
	raw, err := s.req.Orders(ctx, params)
	if err != nil {
		return nil, err
	}
	return &OrdersSnapshot{Orders: raw.Orders, Raw: raw.Raw}, nil
}

func (s *Service) Positions(ctx context.Context, issueCode string) (*PositionsSnapshot, error) {
	if s.req == nil {
		return nil, errors.New("tachibanashi: facade request service is nil")
	}
	cash, err := s.req.CashPositions(ctx, issueCode)
	if err != nil {
		return nil, err
	}
	margin, err := s.req.MarginPositions(ctx, issueCode)
	if err != nil {
		return nil, err
	}

	all := make([]model.Position, 0, len(cash.Positions)+len(margin.Positions))
	all = append(all, cash.Positions...)
	all = append(all, margin.Positions...)

	return &PositionsSnapshot{
		Cash:   cash.Positions,
		Margin: margin.Positions,
		All:    all,
		Raw: &PositionsRaw{
			Cash:   cash.Raw,
			Margin: margin.Raw,
		},
	}, nil
}

func (s *Service) Account(ctx context.Context) (*AccountSnapshot, error) {
	if s.req == nil {
		return nil, errors.New("tachibanashi: facade request service is nil")
	}
	buyingPower, err := s.req.BuyingPower(ctx)
	if err != nil {
		return nil, err
	}
	marginPower, err := s.req.MarginBuyingPower(ctx)
	if err != nil {
		return nil, err
	}
	summary, err := s.req.ZanKaiSummary(ctx)
	if err != nil {
		return nil, err
	}

	snapshot := &AccountSnapshot{
		BuyingPower:       buyingPower.Balance,
		MarginBuyingPower: marginPower.Balance,
		MaintenanceRatio:  marginPower.MaintenanceRatio,
		Summary:           summary.Fields,
		Raw: &AccountRaw{
			BuyingPower:       buyingPower.Raw,
			MarginBuyingPower: marginPower.Raw,
			Summary:           summary,
		},
	}
	return snapshot, nil
}
