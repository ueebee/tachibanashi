package request

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/ueebee/tachibanashi/model"
)

const (
	clmZanKaiKanougaku = "CLMZanKaiKanougaku"
	clmZanKaiSummary   = "CLMZanKaiSummary"
)

type ZanKaiKanougakuRequest struct {
	model.CommonParams
	CLMID     string `json:"sCLMID"`
	IssueCode string `json:"sIssueCode"`
	SizyouC   string `json:"sSizyouC"`
}

func (r *ZanKaiKanougakuRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type ZanKaiKanougakuResponse struct {
	model.CommonResponse
	WarningCode                   string `json:"sWarningCode"`
	WarningText                   string `json:"sWarningText"`
	IssueCode                     string `json:"sIssueCode"`
	SizyouC                       string `json:"sSizyouC"`
	SummaryUpdate                 string `json:"sSummaryUpdate"`
	SummaryGenkabuKaituke         string `json:"sSummaryGenkabuKaituke"`
	SummaryNseityouTousiKanougaku string `json:"sSummaryNseityouTousiKanougaku"`
	HusokukinHasseiFlg            string `json:"sHusokukinHasseiFlg"`
}

type BuyingPowerSnapshot struct {
	Balance model.Balance
	Raw     *ZanKaiKanougakuResponse
}

func (s *Service) ZanKaiKanougaku(ctx context.Context) (*ZanKaiKanougakuResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}

	req := ZanKaiKanougakuRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmZanKaiKanougaku,
		IssueCode:    "",
		SizyouC:      "",
	}

	var resp ZanKaiKanougakuResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) BuyingPower(ctx context.Context) (*BuyingPowerSnapshot, error) {
	raw, err := s.ZanKaiKanougaku(ctx)
	if err != nil {
		return nil, err
	}
	attrs := model.Attributes{
		"sSummaryUpdate":                 raw.SummaryUpdate,
		"sSummaryGenkabuKaituke":         raw.SummaryGenkabuKaituke,
		"sSummaryNseityouTousiKanougaku": raw.SummaryNseityouTousiKanougaku,
		"sHusokukinHasseiFlg":            raw.HusokukinHasseiFlg,
	}
	balance := model.Balance{Raw: attrs}
	if value, ok := parseInt64(raw.SummaryGenkabuKaituke); ok {
		balance.BuyingPower = value
	}
	return &BuyingPowerSnapshot{Balance: balance, Raw: raw}, nil
}

type ZanKaiSummaryRequest struct {
	model.CommonParams
	CLMID string `json:"sCLMID"`
}

func (r *ZanKaiSummaryRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type HikazeiKouzaEntry struct {
	TekiyouYear           string `json:"sHikazeiTekiyouYear"`
	SeityouTousiKanougaku string `json:"sSeityouTousiKanougaku"`
}

type ZanKaiSummaryResponse struct {
	model.CommonResponse
	WarningCode                 string
	WarningText                 string
	Fields                      model.Attributes
	HikazeiKouzaList            []HikazeiKouzaEntry `json:"aHikazeiKouzaList"`
	OisyouHasseiZyoukyouList    json.RawMessage     `json:"aOisyouHasseiZyoukyouList"`
	HosyoukinSeikyuZyoukyouList json.RawMessage     `json:"aHosyoukinSeikyuZyoukyouList"`
}

func (r *ZanKaiSummaryResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Fields = make(model.Attributes, len(raw))

	for key, value := range raw {
		switch key {
		case "aHikazeiKouzaList":
			if len(value) == 0 || bytes.Equal(value, []byte(`""`)) {
				continue
			}
			var list []HikazeiKouzaEntry
			if err := json.Unmarshal(value, &list); err != nil {
				return err
			}
			r.HikazeiKouzaList = list
		case "aOisyouHasseiZyoukyouList":
			r.OisyouHasseiZyoukyouList = value
		case "aHosyoukinSeikyuZyoukyouList":
			r.HosyoukinSeikyuZyoukyouList = value
		default:
			var s string
			if err := json.Unmarshal(value, &s); err == nil {
				r.Fields[key] = s
			}
		}
	}

	r.PNo = r.Fields["p_no"]
	r.PSDDate = r.Fields["p_sd_date"]
	r.PRVDate = r.Fields["p_rv_date"]
	r.PErrNo = r.Fields["p_errno"]
	r.PErr = r.Fields["p_err"]
	r.CLMID = r.Fields["sCLMID"]
	r.ResultCode = r.Fields["sResultCode"]
	r.ResultText = r.Fields["sResultText"]
	r.WarningCode = r.Fields["sWarningCode"]
	r.WarningText = r.Fields["sWarningText"]

	return nil
}

func (s *Service) ZanKaiSummary(ctx context.Context) (*ZanKaiSummaryResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}

	req := ZanKaiSummaryRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmZanKaiSummary,
	}

	var resp ZanKaiSummaryResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
