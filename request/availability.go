package request

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ueebee/tachibanashi/model"
)

const (
	clmZanShinkiKanoIjiritu = "CLMZanShinkiKanoIjiritu"
	clmZanUriKanousuu       = "CLMZanUriKanousuu"
)

type ZanShinkiKanoIjirituRequest struct {
	model.CommonParams
	CLMID     string `json:"sCLMID"`
	IssueCode string `json:"sIssueCode"`
	SizyouC   string `json:"sSizyouC"`
}

func (r *ZanShinkiKanoIjirituRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type ZanShinkiKanoIjirituResponse struct {
	model.CommonResponse
	WarningCode            string
	WarningText            string
	IssueCode              string
	SizyouC                string
	SummaryUpdate          string
	SummarySinyouSinkidate string
	Itakuhosyoukin         string
	OisyouKakuteiFlg       string
}

func (r *ZanShinkiKanoIjirituResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	for key, value := range raw {
		values[key] = jsonString(value)
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	r.IssueCode = values["sIssueCode"]
	r.SizyouC = values["sSizyouC"]
	r.SummaryUpdate = values["sSummaryUpdate"]
	r.SummarySinyouSinkidate = values["sSummarySinyouSinkidate"]
	r.Itakuhosyoukin = values["sItakuhosyoukin"]
	r.OisyouKakuteiFlg = values["sOisyouKakuteiFlg"]
	return nil
}

type MarginBuyingPowerSnapshot struct {
	Balance          model.Balance
	MaintenanceRatio float64
	Raw              *ZanShinkiKanoIjirituResponse
}

func (s *Service) ZanShinkiKanoIjiritu(ctx context.Context) (*ZanShinkiKanoIjirituResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}

	req := ZanShinkiKanoIjirituRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmZanShinkiKanoIjiritu,
		IssueCode:    "",
		SizyouC:      "",
	}

	var resp ZanShinkiKanoIjirituResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) MarginBuyingPower(ctx context.Context) (*MarginBuyingPowerSnapshot, error) {
	raw, err := s.ZanShinkiKanoIjiritu(ctx)
	if err != nil {
		return nil, err
	}
	attrs := model.Attributes{
		"sSummaryUpdate":          raw.SummaryUpdate,
		"sSummarySinyouSinkidate": raw.SummarySinyouSinkidate,
		"sItakuhosyoukin":         raw.Itakuhosyoukin,
		"sOisyouKakuteiFlg":       raw.OisyouKakuteiFlg,
	}
	balance := model.Balance{Raw: attrs}
	if value, ok := parseInt64(raw.SummarySinyouSinkidate); ok {
		balance.BuyingPower = value
	}
	ratio, _ := parseFloat64(raw.Itakuhosyoukin)
	return &MarginBuyingPowerSnapshot{Balance: balance, MaintenanceRatio: ratio, Raw: raw}, nil
}

type ZanUriKanousuuRequest struct {
	model.CommonParams
	CLMID     string `json:"sCLMID"`
	IssueCode string `json:"sIssueCode"`
}

func (r *ZanUriKanousuuRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type ZanUriKanousuuResponse struct {
	model.CommonResponse
	WarningCode                   string
	WarningText                   string
	IssueCode                     string
	SummaryUpdate                 string
	ZanKabuSuryouUriKanouIppan    string
	ZanKabuSuryouUriKanouTokutei  string
	ZanKabuSuryouUriKanouNisa     string
	ZanKabuSuryouUriKanouNseityou string
}

func (r *ZanUriKanousuuResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	for key, value := range raw {
		values[key] = jsonString(value)
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	r.IssueCode = values["sIssueCode"]
	r.SummaryUpdate = values["sSummaryUpdate"]
	r.ZanKabuSuryouUriKanouIppan = values["sZanKabuSuryouUriKanouIppan"]
	r.ZanKabuSuryouUriKanouTokutei = values["sZanKabuSuryouUriKanouTokutei"]
	r.ZanKabuSuryouUriKanouNisa = values["sZanKabuSuryouUriKanouNisa"]
	r.ZanKabuSuryouUriKanouNseityou = values["sZanKabuSuryouUriKanouNseityou"]
	return nil
}

type SellableQuantity struct {
	Ippan    model.Quantity
	Tokutei  model.Quantity
	Nisa     model.Quantity
	Nseityou model.Quantity
}

type SellableQuantitySnapshot struct {
	IssueCode string
	UpdateAt  string
	Quantity  SellableQuantity
	Raw       *ZanUriKanousuuResponse
}

func (s *Service) ZanUriKanousuu(ctx context.Context, issueCode string) (*ZanUriKanousuuResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}

	req := ZanUriKanousuuRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmZanUriKanousuu,
		IssueCode:    issueCode,
	}

	var resp ZanUriKanousuuResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) SellableQuantity(ctx context.Context, issueCode string) (*SellableQuantitySnapshot, error) {
	raw, err := s.ZanUriKanousuu(ctx, issueCode)
	if err != nil {
		return nil, err
	}
	snapshot := &SellableQuantitySnapshot{
		IssueCode: raw.IssueCode,
		UpdateAt:  raw.SummaryUpdate,
		Raw:       raw,
	}
	if value, ok := parseInt64(raw.ZanKabuSuryouUriKanouIppan); ok {
		snapshot.Quantity.Ippan = model.Quantity(value)
	}
	if value, ok := parseInt64(raw.ZanKabuSuryouUriKanouTokutei); ok {
		snapshot.Quantity.Tokutei = model.Quantity(value)
	}
	if value, ok := parseInt64(raw.ZanKabuSuryouUriKanouNisa); ok {
		snapshot.Quantity.Nisa = model.Quantity(value)
	}
	if value, ok := parseInt64(raw.ZanKabuSuryouUriKanouNseityou); ok {
		snapshot.Quantity.Nseityou = model.Quantity(value)
	}
	return snapshot, nil
}
