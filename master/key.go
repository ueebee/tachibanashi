package master

import (
	"strings"

	"github.com/ueebee/tachibanashi/model"
)

const (
	masterFieldSystemStatusKey = "sSystemStatusKey"
	masterFieldYobineTani      = "sYobineTaniNumber"
	masterFieldTekiyouDay      = "sTekiyouDay"
	masterFieldSystemKouza     = "sSystemKouzaKubun"
	masterFieldUnyouCategory   = "sUnyouCategory"
	masterFieldUnyouUnit       = "sUnyouUnit"
	masterFieldEigyouDayC      = "sEigyouDayC"
	masterFieldUnyouStatus     = "sUnyouStatus"
	masterFieldTaisyouGyoumu   = "sTaisyouGyoumu"
	masterFieldZyouzyouSizyou  = "sZyouzyouSizyou"
	masterFieldGensisanCode    = "sGensisanCode"
	masterFieldSyouhinType     = "sSyouhinType"
	masterFieldIssueCode       = "sIssueCode"
	masterFieldHenkouDay       = "sHenkouDay"
	masterFieldErrReasonCode   = "sErrReasonCode"
)

func MasterKey(typ MasterType, fields model.Attributes) (string, bool) {
	switch typ {
	case MasterSystemStatus:
		return valueKey(fields, masterFieldSystemStatusKey)
	case MasterDateZyouhou:
		return valueKey(fields, DateInfoFieldDayKey)
	case MasterYobine:
		return compositeKey(fields, masterFieldYobineTani, masterFieldTekiyouDay)
	case MasterUnyouStatus:
		return compositeKey(fields,
			masterFieldSystemKouza,
			masterFieldUnyouCategory,
			masterFieldUnyouUnit,
			masterFieldEigyouDayC,
			masterFieldUnyouStatus,
			masterFieldTaisyouGyoumu,
		)
	case MasterUnyouStatusKabu:
		return compositeKey(fields,
			masterFieldSystemKouza,
			masterFieldZyouzyouSizyou,
			masterFieldUnyouCategory,
			masterFieldUnyouUnit,
			masterFieldEigyouDayC,
		)
	case MasterUnyouStatusHasei:
		return compositeKey(fields,
			masterFieldSystemKouza,
			masterFieldZyouzyouSizyou,
			masterFieldGensisanCode,
			masterFieldSyouhinType,
			masterFieldUnyouCategory,
			masterFieldUnyouUnit,
			masterFieldEigyouDayC,
		)
	case MasterIssueMstKabu:
		return valueKey(fields, masterFieldIssueCode)
	case MasterIssueSizyouMstKabu:
		return compositeKey(fields, masterFieldIssueCode, masterFieldZyouzyouSizyou)
	case MasterIssueSizyouKiseiKabu:
		return compositeKey(fields, masterFieldSystemKouza, masterFieldIssueCode, masterFieldZyouzyouSizyou)
	case MasterIssueMstSak:
		return valueKey(fields, masterFieldIssueCode)
	case MasterIssueMstOp:
		return valueKey(fields, masterFieldIssueCode)
	case MasterIssueSizyouKiseiHasei:
		return compositeKey(fields, masterFieldSystemKouza, masterFieldIssueCode, masterFieldZyouzyouSizyou)
	case MasterDaiyouKakeme:
		return compositeKey(fields, masterFieldSystemKouza, masterFieldIssueCode, masterFieldTekiyouDay)
	case MasterHosyoukinMst:
		return compositeKey(fields, masterFieldSystemKouza, masterFieldIssueCode, masterFieldZyouzyouSizyou, masterFieldHenkouDay)
	case MasterOrderErrReason:
		return valueKey(fields, masterFieldErrReasonCode)
	case MasterIssueMstOther, MasterIssueMstIndex, MasterIssueMstFx:
		return issueCodeKey(fields)
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

func compositeKey(fields model.Attributes, keys ...string) (string, bool) {
	if len(keys) == 0 {
		return "", false
	}
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value := strings.TrimSpace(fields.Value(key))
		if value == "" {
			return "", false
		}
		parts = append(parts, value)
	}
	return JoinIndex(parts...), true
}

func issueCodeKey(fields model.Attributes) (string, bool) {
	if key, ok := compositeKey(fields, masterFieldIssueCode, masterFieldZyouzyouSizyou); ok {
		return key, true
	}
	return valueKey(fields, masterFieldIssueCode)
}
