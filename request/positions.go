package request

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ueebee/tachibanashi/model"
)

const (
	clmGenbutuKabuList      = "CLMGenbutuKabuList"
	clmShinyouTategyokuList = "CLMShinyouTategyokuList"
)

type GenbutuKabuListRequest struct {
	model.CommonParams
	CLMID     string `json:"sCLMID"`
	IssueCode string `json:"sIssueCode"`
}

func (r *GenbutuKabuListRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type GenbutuKabuListResponse struct {
	model.CommonResponse
	WarningCode                      string             `json:"sWarningCode"`
	WarningText                      string             `json:"sWarningText"`
	IssueCode                        string             `json:"sIssueCode"`
	IppanGaisanHyoukagakuGoukei      string             `json:"sIppanGaisanHyoukagakuGoukei"`
	IppanGaisanHyoukaSonekiGoukei    string             `json:"sIppanGaisanHyoukaSonekiGoukei"`
	NisaGaisanHyoukagakuGoukei       string             `json:"sNisaGaisanHyoukagakuGoukei"`
	NisaGaisanHyoukaSonekiGoukei     string             `json:"sNisaGaisanHyoukaSonekiGoukei"`
	NseityouGaisanHyoukagakuGoukei   string             `json:"sNseityouGaisanHyoukagakuGoukei"`
	NseityouGaisanHyoukaSonekiGoukei string             `json:"sNseityouGaisanHyoukaSonekiGoukei"`
	TokuteiGaisanHyoukagakuGoukei    string             `json:"sTokuteiGaisanHyoukagakuGoukei"`
	TokuteiGaisanHyoukaSonekiGoukei  string             `json:"sTokuteiGaisanHyoukaSonekiGoukei"`
	TotalGaisanHyoukagakuGoukei      string             `json:"sTotalGaisanHyoukagakuGoukei"`
	TotalGaisanHyoukaSonekiGoukei    string             `json:"sTotalGaisanHyoukaSonekiGoukei"`
	Entries                          []GenbutuKabuEntry `json:"aGenbutuKabuList"`
}

func (r *GenbutuKabuListResponse) UnmarshalJSON(data []byte) error {
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
	r.IppanGaisanHyoukagakuGoukei = values["sIppanGaisanHyoukagakuGoukei"]
	r.IppanGaisanHyoukaSonekiGoukei = values["sIppanGaisanHyoukaSonekiGoukei"]
	r.NisaGaisanHyoukagakuGoukei = values["sNisaGaisanHyoukagakuGoukei"]
	r.NisaGaisanHyoukaSonekiGoukei = values["sNisaGaisanHyoukaSonekiGoukei"]
	r.NseityouGaisanHyoukagakuGoukei = values["sNseityouGaisanHyoukagakuGoukei"]
	r.NseityouGaisanHyoukaSonekiGoukei = values["sNseityouGaisanHyoukaSonekiGoukei"]
	r.TokuteiGaisanHyoukagakuGoukei = values["sTokuteiGaisanHyoukagakuGoukei"]
	r.TokuteiGaisanHyoukaSonekiGoukei = values["sTokuteiGaisanHyoukaSonekiGoukei"]
	r.TotalGaisanHyoukagakuGoukei = values["sTotalGaisanHyoukagakuGoukei"]
	r.TotalGaisanHyoukaSonekiGoukei = values["sTotalGaisanHyoukaSonekiGoukei"]

	if err := decodeList(raw["aGenbutuKabuList"], &r.Entries); err != nil {
		return err
	}
	return nil
}

type GenbutuKabuEntry struct {
	IssueCode string
	Fields    model.Attributes
}

func (e *GenbutuKabuEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "sUriOrderIssueCode" {
			e.IssueCode = value
			continue
		}
		e.Fields[key] = value
	}
	return nil
}

func (e GenbutuKabuEntry) Position() model.Position {
	pos := model.Position{
		Symbol: e.IssueCode,
		Raw:    cloneAttributes(e.Fields),
	}
	if qty, ok := parseInt64(e.Fields.Value("sUriOrderZanKabuSuryou")); ok {
		pos.Quantity = model.Quantity(qty)
	}
	if price, ok := parsePrice(e.Fields.Value("sUriOrderGaisanBokaTanka")); ok {
		pos.AvgPrice = model.Price(price)
	}
	if pnl, ok := parseInt64(e.Fields.Value("sUriOrderGaisanHyoukaSoneki")); ok {
		pos.UnrealPnL = pnl
	}
	return pos
}

type CashPositionsSnapshot struct {
	Positions []model.Position
	Raw       *GenbutuKabuListResponse
}

func (s *Service) GenbutuKabuList(ctx context.Context, issueCode string) (*GenbutuKabuListResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}

	req := GenbutuKabuListRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmGenbutuKabuList,
		IssueCode:    issueCode,
	}

	var resp GenbutuKabuListResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) CashPositions(ctx context.Context, issueCode string) (*CashPositionsSnapshot, error) {
	raw, err := s.GenbutuKabuList(ctx, issueCode)
	if err != nil {
		return nil, err
	}
	positions := make([]model.Position, 0, len(raw.Entries))
	for _, entry := range raw.Entries {
		positions = append(positions, entry.Position())
	}
	return &CashPositionsSnapshot{Positions: positions, Raw: raw}, nil
}

type ShinyouTategyokuListRequest struct {
	model.CommonParams
	CLMID     string `json:"sCLMID"`
	IssueCode string `json:"sIssueCode"`
}

func (r *ShinyouTategyokuListRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type ShinyouTategyokuListResponse struct {
	model.CommonResponse
	WarningCode               string                  `json:"sWarningCode"`
	WarningText               string                  `json:"sWarningText"`
	IssueCode                 string                  `json:"sIssueCode"`
	UritateDaikin             string                  `json:"sUritateDaikin"`
	KaitateDaikin             string                  `json:"sKaitateDaikin"`
	TotalDaikin               string                  `json:"sTotalDaikin"`
	HyoukaSonekiGoukeiUridate string                  `json:"sHyoukaSonekiGoukeiUridate"`
	HyoukaSonekiGoukeiKaidate string                  `json:"sHyoukaSonekiGoukeiKaidate"`
	TokuteiHyoukaSonekiGoukei string                  `json:"sTokuteiHyoukaSonekiGoukei"`
	TotalHyoukaSonekiGoukei   string                  `json:"sTotalHyoukaSonekiGoukei"`
	IppanHyoukaSonekiGoukei   string                  `json:"sIppanHyoukaSonekiGoukei"`
	Entries                   []ShinyouTategyokuEntry `json:"aShinyouTategyokuList"`
}

func (r *ShinyouTategyokuListResponse) UnmarshalJSON(data []byte) error {
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
	r.UritateDaikin = values["sUritateDaikin"]
	r.KaitateDaikin = values["sKaitateDaikin"]
	r.TotalDaikin = values["sTotalDaikin"]
	r.HyoukaSonekiGoukeiUridate = values["sHyoukaSonekiGoukeiUridate"]
	r.HyoukaSonekiGoukeiKaidate = values["sHyoukaSonekiGoukeiKaidate"]
	r.TokuteiHyoukaSonekiGoukei = values["sTokuteiHyoukaSonekiGoukei"]
	r.TotalHyoukaSonekiGoukei = values["sTotalHyoukaSonekiGoukei"]
	r.IppanHyoukaSonekiGoukei = values["sIppanHyoukaSonekiGoukei"]

	if err := decodeList(raw["aShinyouTategyokuList"], &r.Entries); err != nil {
		return err
	}
	return nil
}

type ShinyouTategyokuEntry struct {
	IssueCode string
	Fields    model.Attributes
}

func (e *ShinyouTategyokuEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "sOrderIssueCode" {
			e.IssueCode = value
			continue
		}
		e.Fields[key] = value
	}
	return nil
}

func (e ShinyouTategyokuEntry) Position() model.Position {
	pos := model.Position{
		Symbol: e.IssueCode,
		Raw:    cloneAttributes(e.Fields),
	}
	if qty, ok := parseInt64(e.Fields.Value("sOrderTategyokuSuryou")); ok {
		pos.Quantity = model.Quantity(qty)
	}
	if price, ok := parsePrice(e.Fields.Value("sOrderTategyokuTanka")); ok {
		pos.AvgPrice = model.Price(price)
	}
	if pnl, ok := parseInt64(e.Fields.Value("sOrderGaisanHyoukaSoneki")); ok {
		pos.UnrealPnL = pnl
	}
	return pos
}

type MarginPositionsSnapshot struct {
	Positions []model.Position
	Raw       *ShinyouTategyokuListResponse
}

func (s *Service) ShinyouTategyokuList(ctx context.Context, issueCode string) (*ShinyouTategyokuListResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}

	req := ShinyouTategyokuListRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmShinyouTategyokuList,
		IssueCode:    issueCode,
	}

	var resp ShinyouTategyokuListResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) MarginPositions(ctx context.Context, issueCode string) (*MarginPositionsSnapshot, error) {
	raw, err := s.ShinyouTategyokuList(ctx, issueCode)
	if err != nil {
		return nil, err
	}
	positions := make([]model.Position, 0, len(raw.Entries))
	for _, entry := range raw.Entries {
		positions = append(positions, entry.Position())
	}
	return &MarginPositionsSnapshot{Positions: positions, Raw: raw}, nil
}

func parseInt64(value string) (int64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func parsePrice(value string) (int64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	if strings.Contains(value, ".") {
		parts := strings.SplitN(value, ".", 2)
		if len(parts) != 2 {
			return 0, false
		}
		if strings.TrimRight(parts[1], "0") != "" {
			return 0, false
		}
		value = parts[0]
	}
	return parseInt64(value)
}

func parseFloat64(value string) (float64, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func cloneAttributes(fields model.Attributes) model.Attributes {
	if fields == nil {
		return nil
	}
	clone := make(model.Attributes, len(fields))
	for key, value := range fields {
		clone[key] = value
	}
	return clone
}

func applyCommonResponse(target *model.CommonResponse, values map[string]string) {
	if target == nil {
		return
	}
	target.PNo = values["p_no"]
	target.PSDDate = values["p_sd_date"]
	target.PRVDate = values["p_rv_date"]
	target.PErrNo = values["p_errno"]
	target.PErr = values["p_err"]
	target.CLMID = values["sCLMID"]
	target.ResultCode = values["sResultCode"]
	target.ResultText = values["sResultText"]
}

func jsonString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func decodeList(raw json.RawMessage, out any) error {
	if len(raw) == 0 {
		return nil
	}
	trimmed := bytes.TrimSpace(raw)
	if bytes.Equal(trimmed, []byte("null")) {
		return nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		return json.Unmarshal([]byte(s), out)
	}
	return json.Unmarshal(raw, out)
}
