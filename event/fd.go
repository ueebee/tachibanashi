package event

import (
	"encoding/hex"
	"sort"
	"strconv"
	"strings"

	"github.com/ueebee/tachibanashi/model"
)

type FD struct {
	Frame
	Rows []FDRow
}

type FDRow struct {
	Row    int
	Fields model.Attributes
}

type QuoteBook struct {
	rows map[int]model.Attributes
}

func NewQuoteBook() *QuoteBook {
	return &QuoteBook{rows: make(map[int]model.Attributes)}
}

func (b *QuoteBook) Apply(event FD) []model.Quote {
	if b.rows == nil {
		b.rows = make(map[int]model.Attributes)
	}

	updated := make([]model.Quote, 0, len(event.Rows))
	for _, row := range event.Rows {
		if row.Fields == nil {
			continue
		}
		merged := mergeAttributes(b.rows[row.Row], row.Fields)
		b.rows[row.Row] = merged
		updated = append(updated, model.Quote{Fields: cloneAttributes(merged)})
	}
	return updated
}

func (b *QuoteBook) Snapshot() []model.Quote {
	if len(b.rows) == 0 {
		return nil
	}
	keys := make([]int, 0, len(b.rows))
	for key := range b.rows {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	out := make([]model.Quote, 0, len(keys))
	for _, key := range keys {
		out = append(out, model.Quote{Fields: cloneAttributes(b.rows[key])})
	}
	return out
}

func (f FD) Quotes(symbols map[int]string) []model.Quote {
	if len(f.Rows) == 0 {
		return nil
	}
	out := make([]model.Quote, 0, len(f.Rows))
	for _, row := range f.Rows {
		symbol := ""
		if symbols != nil {
			symbol = symbols[row.Row]
		}
		out = append(out, model.Quote{
			Symbol: symbol,
			Fields: cloneAttributes(row.Fields),
		})
	}
	return out
}

func parseFD(frame Frame) (FD, error) {
	rows := make(map[int]model.Attributes)
	for key, values := range frame.Fields {
		row, field, ok := parseFDKey(key)
		if !ok {
			continue
		}
		value := ""
		if len(values) > 0 {
			value = values[0]
		}
		if strings.HasPrefix(field, "x") {
			value = decodeHexShiftJIS(value)
		}
		attrs := rows[row]
		if attrs == nil {
			attrs = make(model.Attributes)
			rows[row] = attrs
		}
		attrs[field] = value
	}

	rowKeys := make([]int, 0, len(rows))
	for row := range rows {
		rowKeys = append(rowKeys, row)
	}
	sort.Ints(rowKeys)

	out := make([]FDRow, 0, len(rowKeys))
	for _, row := range rowKeys {
		out = append(out, FDRow{
			Row:    row,
			Fields: cloneAttributes(rows[row]),
		})
	}
	return FD{Frame: frame, Rows: out}, nil
}

func parseFDKey(key string) (int, string, bool) {
	if len(key) < 4 {
		return 0, "", false
	}
	if key[1] != '_' {
		return 0, "", false
	}
	prefix := key[:1]
	rest := key[2:]
	index := strings.IndexByte(rest, '_')
	if index <= 0 {
		return 0, "", false
	}
	row, err := strconv.Atoi(rest[:index])
	if err != nil {
		return 0, "", false
	}
	field := prefix + rest[index+1:]
	return row, field, true
}

func decodeHexShiftJIS(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		return value
	}
	out, err := decodeShiftJIS(decoded)
	if err != nil {
		return value
	}
	return out
}

func mergeAttributes(base, update model.Attributes) model.Attributes {
	if base == nil && update == nil {
		return nil
	}
	out := make(model.Attributes, len(base)+len(update))
	for key, value := range base {
		out[key] = value
	}
	for key, value := range update {
		out[key] = value
	}
	return out
}

func cloneAttributes(attrs model.Attributes) model.Attributes {
	if attrs == nil {
		return nil
	}
	out := make(model.Attributes, len(attrs))
	for key, value := range attrs {
		out[key] = value
	}
	return out
}
