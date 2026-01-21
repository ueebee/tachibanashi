package model

// Quote is a generic snapshot for a single symbol.
type Quote struct {
	Symbol string
	Fields map[string]string
}

func (q Quote) Value(key string) string {
	if q.Fields == nil {
		return ""
	}
	return q.Fields[key]
}
