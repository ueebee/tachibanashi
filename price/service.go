package price

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ueebee/tachibanashi/auth"
	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/model"
)

const (
	clmMarketPrice        = "CLMMfdsGetMarketPrice"
	clmMarketPriceHistory = "CLMMfdsGetMarketPriceHistory"

	maxIssueCodes = 120
)

type Client interface {
	DoJSON(ctx context.Context, method, path string, req, resp any) error
	VirtualURLs() auth.VirtualURLs
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}

type MarketPriceRequest struct {
	model.CommonParams
	CLMID           string `json:"sCLMID"`
	TargetIssueCode string `json:"sTargetIssueCode"`
	TargetColumn    string `json:"sTargetColumn"`
}

func (r *MarketPriceRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type MarketPriceResponse struct {
	model.CommonResponse
	Prices []MarketPriceEntry `json:"aCLMMfdsMarketPrice"`
}

type MarketPriceEntry struct {
	IssueCode string
	Fields    model.Attributes
}

func (e *MarketPriceEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "sIssueCode" {
			e.IssueCode = value
			continue
		}
		e.Fields[key] = value
	}
	return nil
}

func (e MarketPriceEntry) Value(key string) string {
	return e.Fields.Value(key)
}

type QuoteSnapshot struct {
	Quotes []model.Quote
	Raw    *MarketPriceResponse
}

type MarketPriceHistoryRequest struct {
	model.CommonParams
	CLMID      string `json:"sCLMID"`
	IssueCode  string `json:"sIssueCode"`
	MarketCode string `json:"sSizyouC,omitempty"`
}

func (r *MarketPriceHistoryRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type MarketPriceHistoryResponse struct {
	model.CommonResponse
	IssueCode  string                    `json:"sIssueCode"`
	MarketCode string                    `json:"sSizyouC"`
	Entries    []MarketPriceHistoryEntry `json:"aCLMMfdsGetMarketPriceHistory"`
}

type MarketPriceHistoryEntry struct {
	Date   string
	Fields model.Attributes
}

func (e *MarketPriceHistoryEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "sDate" {
			e.Date = value
			continue
		}
		e.Fields[key] = value
	}
	return nil
}

// Snapshot fetches market prices for multiple codes in a single request (max 120).
func (s *Service) Snapshot(ctx context.Context, issueCodes []string, columns []string) (*MarketPriceResponse, error) {
	codes := normalizeList(issueCodes)
	if len(codes) == 0 {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "required"}
	}
	if len(codes) > maxIssueCodes {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "max 120"}
	}

	cols := normalizeList(columns)
	if len(cols) == 0 {
		return nil, &terrors.ValidationError{Field: "columns", Reason: "required"}
	}

	urls := s.client.VirtualURLs()
	if urls.Price == "" {
		return nil, errors.New("tachibanashi: virtual price URL not set")
	}

	req := MarketPriceRequest{
		CommonParams:    model.CommonParams{JsonOfmt: "5"},
		CLMID:           clmMarketPrice,
		TargetIssueCode: strings.Join(codes, ","),
		TargetColumn:    strings.Join(cols, ","),
	}

	var resp MarketPriceResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, urls.Price, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QuoteSnapshot fetches quotes using generic terms and keeps the raw response.
func (s *Service) QuoteSnapshot(ctx context.Context, symbols []string, fields []string) (*QuoteSnapshot, error) {
	raw, err := s.Snapshot(ctx, symbols, fields)
	if err != nil {
		return nil, err
	}

	quotes := make([]model.Quote, 0, len(raw.Prices))
	for _, entry := range raw.Prices {
		quotes = append(quotes, model.Quote{
			Symbol: entry.IssueCode,
			Fields: cloneAttributes(entry.Fields),
		})
	}

	return &QuoteSnapshot{Quotes: quotes, Raw: raw}, nil
}

// History fetches daily price history for a single issue code.
func (s *Service) History(ctx context.Context, issueCode, marketCode string) (*MarketPriceHistoryResponse, error) {
	codes := normalizeList([]string{issueCode})
	if len(codes) == 0 {
		return nil, &terrors.ValidationError{Field: "issue_code", Reason: "required"}
	}
	if len(codes) > 1 {
		return nil, &terrors.ValidationError{Field: "issue_code", Reason: "single issue only"}
	}

	marketCodes := normalizeList([]string{marketCode})
	if len(marketCodes) > 1 {
		return nil, &terrors.ValidationError{Field: "market_code", Reason: "single market only"}
	}

	urls := s.client.VirtualURLs()
	if urls.Price == "" {
		return nil, errors.New("tachibanashi: virtual price URL not set")
	}

	req := MarketPriceHistoryRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmMarketPriceHistory,
		IssueCode:    codes[0],
	}
	if len(marketCodes) == 1 {
		req.MarketCode = marketCodes[0]
	}

	var resp MarketPriceHistoryResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, urls.Price, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func normalizeList(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			item := strings.TrimSpace(part)
			if item == "" {
				continue
			}
			out = append(out, item)
		}
	}
	return out
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
