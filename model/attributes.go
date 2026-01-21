package model

import "strconv"

type Attributes map[string]string

func (a Attributes) Value(key string) string {
	if a == nil {
		return ""
	}
	return a[key]
}

func (a Attributes) Int64(key string) (int64, bool) {
	value, ok := a.lookup(key)
	if !ok {
		return 0, false
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func (a Attributes) Float64(key string) (float64, bool) {
	value, ok := a.lookup(key)
	if !ok {
		return 0, false
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func (a Attributes) lookup(key string) (string, bool) {
	if a == nil {
		return "", false
	}
	value, ok := a[key]
	if !ok || value == "" {
		return "", false
	}
	return value, true
}
