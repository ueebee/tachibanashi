package event

import (
	"strconv"
	"strings"

	"github.com/ueebee/tachibanashi/model"
)

type EC struct {
	Frame
	Fields            model.Attributes
	Provider          string
	EventNo           string
	Alert             bool
	NoticeType        string
	OrderNumber       string
	BusinessDay       string
	ParentOrderNumber string
	OrderType         string
	SecurityType      string
	Symbol            string
	MarketCode        string
	Side              string
	TradeType         string
	OrderPrice        string
	OrderQuantity     string
	ExecutedPrice     string
	ExecutedQuantity  string
	ExecutedTime      string
	OrderStatus       string
}

func (e EC) Order() model.Order {
	order := model.Order{
		ID:     e.OrderNumber,
		Symbol: e.Symbol,
		Side:   e.Side,
		Status: e.OrderStatus,
		Raw:    cloneAttributes(e.Fields),
	}
	if price, ok := parsePrice(e.OrderPrice); ok {
		order.Price = model.Price(price)
	}
	if qty, ok := parseInt64Value(e.OrderQuantity); ok {
		order.Quantity = model.Quantity(qty)
	}
	return order
}

func (e EC) Execution() (model.Execution, bool) {
	exec := model.Execution{
		OrderID: e.OrderNumber,
		Symbol:  e.Symbol,
		Time:    e.ExecutedTime,
		Raw:     cloneAttributes(e.Fields),
	}
	if price, ok := parsePrice(e.ExecutedPrice); ok {
		exec.Price = model.Price(price)
	}
	if qty, ok := parseInt64Value(e.ExecutedQuantity); ok {
		exec.Quantity = model.Quantity(qty)
	}
	if exec.Price == 0 && exec.Quantity == 0 && exec.Time == "" {
		return model.Execution{}, false
	}
	return exec, true
}

func parseEC(frame Frame) (EC, error) {
	fields := frameAttributes(frame.Fields)

	ec := EC{
		Frame:             frame,
		Fields:            fields,
		Provider:          fields.Value("p_PV"),
		EventNo:           fields.Value("p_ENO"),
		Alert:             parseBool(fields.Value("p_ALT")),
		NoticeType:        fields.Value("p_NT"),
		OrderNumber:       fields.Value("p_ON"),
		BusinessDay:       fields.Value("p_ED"),
		ParentOrderNumber: fields.Value("p_OON"),
		OrderType:         fields.Value("p_OT"),
		SecurityType:      fields.Value("p_ST"),
		Symbol:            fields.Value("p_IC"),
		MarketCode:        fields.Value("p_MC"),
		Side:              fields.Value("p_BBKB"),
		TradeType:         fields.Value("p_CRSJ"),
		OrderPrice:        fields.Value("p_CRPR"),
		OrderQuantity:     fields.Value("p_CRSR"),
		ExecutedPrice:     fields.Value("p_EXPR"),
		ExecutedQuantity:  fields.Value("p_EXSR"),
		ExecutedTime:      fields.Value("p_EXDT"),
		OrderStatus:       fields.Value("p_ODST"),
	}

	return ec, nil
}

func frameAttributes(fields map[string][]string) model.Attributes {
	if fields == nil {
		return nil
	}
	out := make(model.Attributes, len(fields))
	for key, values := range fields {
		if len(values) == 0 {
			out[key] = ""
			continue
		}
		out[key] = values[0]
	}
	return out
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
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func parseInt64Value(value string) (int64, bool) {
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

func parseBool(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if value == "1" {
		return true
	}
	return strings.EqualFold(value, "true")
}
