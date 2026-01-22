package request

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/model"
)

const (
	clmKabuNewOrder       = "CLMKabuNewOrder"
	clmKabuCorrectOrder   = "CLMKabuCorrectOrder"
	clmKabuCancelOrder    = "CLMKabuCancelOrder"
	clmKabuCancelOrderAll = "CLMKabuCancelOrderAll"
	clmOrderList          = "CLMOrderList"
	clmOrderListDetail    = "CLMOrderListDetail"
)

type OrderParams map[string]any

type orderRequest struct {
	model.CommonParams
	CLMID  string
	Payload OrderParams
}

func (r *orderRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

func (r orderRequest) MarshalJSON() ([]byte, error) {
	payload := make(map[string]any, len(r.Payload)+4)
	if r.PNo != "" {
		payload["p_no"] = r.PNo
	}
	if r.PSDDate != "" {
		payload["p_sd_date"] = r.PSDDate
	}
	if r.JsonOfmt != "" {
		payload["sJsonOfmt"] = r.JsonOfmt
	}
	if r.CLMID != "" {
		payload["sCLMID"] = r.CLMID
	}
	for key, value := range r.Payload {
		if key == "" {
			continue
		}
		if _, exists := payload[key]; exists {
			continue
		}
		payload[key] = value
	}
	return json.Marshal(payload)
}

type OrderResponse struct {
	model.CommonResponse
	WarningCode string
	WarningText string
	OrderNumber string
	EigyouDay   string
	Fields      model.Attributes
}

func (r *OrderResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		values[key] = str
		r.Fields[key] = str
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	r.OrderNumber = values["sOrderNumber"]
	r.EigyouDay = values["sEigyouDay"]
	return nil
}

type OrderEntry struct {
	OrderID string
	Symbol  string
	Fields  model.Attributes
}

func (e *OrderEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		e.Fields[key] = jsonString(value)
	}
	if id := e.Fields.Value("sOrderOrderNumber"); id != "" {
		e.OrderID = id
	} else if id := e.Fields.Value("sOrderNumber"); id != "" {
		e.OrderID = id
	}
	if symbol := e.Fields.Value("sOrderIssueCode"); symbol != "" {
		e.Symbol = symbol
	} else if symbol := e.Fields.Value("sIssueCode"); symbol != "" {
		e.Symbol = symbol
	}
	return nil
}

func (e OrderEntry) Order() model.Order {
	order := model.Order{
		ID:     e.OrderID,
		Symbol: e.Symbol,
		Raw:    cloneAttributes(e.Fields),
	}
	if side := e.Fields.Value("sOrderBaibaiKubun"); side != "" {
		order.Side = side
	}
	if qty, ok := parseInt64(e.Fields.Value("sOrderOrderSuryou")); ok {
		order.Quantity = model.Quantity(qty)
	}
	if price, ok := parsePrice(e.Fields.Value("sOrderOrderPrice")); ok {
		order.Price = model.Price(price)
	}
	if status := e.Fields.Value("sOrderStatus"); status != "" {
		order.Status = status
	}
	return order
}

type OrderListResponse struct {
	model.CommonResponse
	WarningCode        string
	WarningText        string
	IssueCode          string
	OrderSyoukaiStatus string
	SikkouDay          string
	Entries            []OrderEntry
	Fields             model.Attributes
}

func (r *OrderListResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "aOrderList" {
			if err := decodeList(value, &r.Entries); err != nil {
				return err
			}
			continue
		}
		str := jsonString(value)
		values[key] = str
		r.Fields[key] = str
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	r.IssueCode = values["sIssueCode"]
	r.OrderSyoukaiStatus = values["sOrderSyoukaiStatus"]
	r.SikkouDay = values["sSikkouDay"]
	return nil
}

type OrdersSnapshot struct {
	Orders []model.Order
	Raw    *OrderListResponse
}

func (s *Service) Orders(ctx context.Context, params OrderParams) (*OrdersSnapshot, error) {
	raw, err := s.OrderList(ctx, params)
	if err != nil {
		return nil, err
	}
	orders := make([]model.Order, 0, len(raw.Entries))
	for _, entry := range raw.Entries {
		orders = append(orders, entry.Order())
	}
	return &OrdersSnapshot{Orders: orders, Raw: raw}, nil
}

type OrderListDetailResponse struct {
	model.CommonResponse
	WarningCode string
	WarningText string
	OrderNumber string
	EigyouDay   string
	IssueCode   string
	Fields      model.Attributes
}

func (r *OrderListDetailResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		values[key] = str
		r.Fields[key] = str
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	r.OrderNumber = values["sOrderNumber"]
	r.EigyouDay = values["sEigyouDay"]
	r.IssueCode = values["sIssueCode"]
	return nil
}

func (s *Service) KabuNewOrder(ctx context.Context, params OrderParams) (*OrderResponse, error) {
	if err := requireParams(params,
		"sZyoutoekiKazeiC",
		"sIssueCode",
		"sSizyouC",
		"sBaibaiKubun",
		"sCondition",
		"sOrderPrice",
		"sOrderSuryou",
		"sGenkinShinyouKubun",
		"sOrderExpireDay",
		"sGyakusasiOrderType",
		"sGyakusasiZyouken",
		"sGyakusasiPrice",
		"sTatebiType",
		"sTategyokuZyoutoekiKazeiC",
		"sSecondPassword",
	); err != nil {
		return nil, err
	}
	return s.submitOrder(ctx, clmKabuNewOrder, params)
}

func (s *Service) KabuCorrectOrder(ctx context.Context, params OrderParams) (*OrderResponse, error) {
	if err := requireParams(params, "sOrderNumber", "sEigyouDay", "sSecondPassword"); err != nil {
		return nil, err
	}
	return s.submitOrder(ctx, clmKabuCorrectOrder, params)
}

func (s *Service) KabuCancelOrder(ctx context.Context, params OrderParams) (*OrderResponse, error) {
	if err := requireParams(params, "sOrderNumber", "sEigyouDay", "sSecondPassword"); err != nil {
		return nil, err
	}
	return s.submitOrder(ctx, clmKabuCancelOrder, params)
}

func (s *Service) KabuCancelOrderAll(ctx context.Context, params OrderParams) (*OrderResponse, error) {
	if err := requireParams(params, "sSecondPassword"); err != nil {
		return nil, err
	}
	return s.submitOrder(ctx, clmKabuCancelOrderAll, params)
}

func (s *Service) OrderList(ctx context.Context, params OrderParams) (*OrderListResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}
	req := orderRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmOrderList,
		Payload:      params,
	}
	var resp OrderListResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) OrderListDetail(ctx context.Context, orderNumber, eigyouDay string) (*OrderListDetailResponse, error) {
	if strings.TrimSpace(orderNumber) == "" {
		return nil, &terrors.ValidationError{Field: "sOrderNumber", Reason: "required"}
	}
	if strings.TrimSpace(eigyouDay) == "" {
		return nil, &terrors.ValidationError{Field: "sEigyouDay", Reason: "required"}
	}
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}
	req := orderRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmOrderListDetail,
		Payload: OrderParams{
			"sOrderNumber": orderNumber,
			"sEigyouDay":   eigyouDay,
		},
	}
	var resp OrderListDetailResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) submitOrder(ctx context.Context, clmid string, params OrderParams) (*OrderResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}
	req := orderRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmid,
		Payload:      params,
	}
	var resp OrderResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func requireParams(params OrderParams, fields ...string) error {
	for _, field := range fields {
		if params == nil {
			return &terrors.ValidationError{Field: field, Reason: "required"}
		}
		value, ok := params[field]
		if !ok || value == nil {
			return &terrors.ValidationError{Field: field, Reason: "required"}
		}
		if s, ok := value.(string); ok {
			if strings.TrimSpace(s) == "" {
				return &terrors.ValidationError{Field: field, Reason: "required"}
			}
		}
	}
	return nil
}
