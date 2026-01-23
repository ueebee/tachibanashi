package master

import (
	"encoding/json"

	"github.com/ueebee/tachibanashi/model"
)

const (
	UnyouStatusFieldSystemKouzaKubun = "sSystemKouzaKubun"
	UnyouStatusFieldUnyouCategory    = "sUnyouCategory"
	UnyouStatusFieldUnyouUnit        = "sUnyouUnit"
	UnyouStatusFieldEigyouDayC       = "sEigyouDayC"
	UnyouStatusFieldUnyouStatus      = "sUnyouStatus"
	UnyouStatusFieldTaisyouGyoumu    = "sTaisyouGyoumu"
	UnyouStatusFieldGyoumuZyoutai    = "sGyoumuZyoutai"
	UnyouStatusFieldEventName        = "sEventName"
	UnyouStatusFieldMeyasuTime       = "sMeyasuTime"
	UnyouStatusFieldCreateTime       = "sCreateTime"
	UnyouStatusFieldUpdateTime       = "sUpdateTime"
	UnyouStatusFieldUpdateNumber     = "sUpdateNumber"
	UnyouStatusFieldDeleteFlag       = "sDeleteFlag"
	UnyouStatusFieldDeleteTime       = "sDeleteTime"
	UnyouStatusFieldZyouzyouSizyou   = "sZyouzyouSizyou"
	UnyouStatusFieldGensisanCode     = "sGensisanCode"
	UnyouStatusFieldSyouhinType      = "sSyouhinType"
)

type UnyouStatus struct {
	Fields model.Attributes
}

type UnyouStatusKabu struct {
	Fields model.Attributes
}

type UnyouStatusHasei struct {
	Fields model.Attributes
}

func (u *UnyouStatus) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &u.Fields)
}

func (u *UnyouStatusKabu) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &u.Fields)
}

func (u *UnyouStatusHasei) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &u.Fields)
}

func unmarshalAttributes(data []byte, dest *model.Attributes) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	attrs := make(model.Attributes, len(raw))
	for key, value := range raw {
		attrs[key] = jsonString(value)
	}
	*dest = attrs
	return nil
}
