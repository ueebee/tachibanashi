package master

import (
	"strings"

	"github.com/ueebee/tachibanashi/model"
)

func MasterKey(typ MasterType, fields model.Attributes) (string, bool) {
	switch typ {
	case MasterDateZyouhou:
		return valueKey(fields, DateInfoFieldDayKey)
	default:
		return "", false
	}
}

func valueKey(fields model.Attributes, key string) (string, bool) {
	value := strings.TrimSpace(fields.Value(key))
	if value == "" {
		return "", false
	}
	return value, true
}
