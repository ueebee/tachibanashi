package master

import "github.com/ueebee/tachibanashi/model"

const (
	SystemStatusFieldKey          = "sSystemStatusKey"
	SystemStatusFieldLoginKyoka   = "sLoginKyokaKubun"
	SystemStatusFieldStatus       = "sSystemStatus"
	SystemStatusFieldCreateTime   = "sCreateTime"
	SystemStatusFieldUpdateTime   = "sUpdateTime"
	SystemStatusFieldUpdateNumber = "sUpdateNumber"
	SystemStatusFieldDeleteFlag   = "sDeleteFlag"
	SystemStatusFieldDeleteTime   = "sDeleteTime"
)

type SystemStatus struct {
	Fields model.Attributes
}

func (s *SystemStatus) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &s.Fields)
}

func (s SystemStatus) Key() string {
	return s.Fields.Value(SystemStatusFieldKey)
}
