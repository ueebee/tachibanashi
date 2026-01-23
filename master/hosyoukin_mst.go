package master

import "github.com/ueebee/tachibanashi/model"

const (
	HosyoukinMstFieldSystemKouzaKubun = "sSystemKouzaKubun"
	HosyoukinMstFieldIssueCode        = "sIssueCode"
	HosyoukinMstFieldZyouzyouSizyou   = "sZyouzyouSizyou"
	HosyoukinMstFieldHenkouDay        = "sHenkouDay"
	HosyoukinMstFieldDaiyoRitu        = "sDaiyoHosyokinRitu"
	HosyoukinMstFieldGenkinRitu       = "sGenkinHosyokinRitu"
	HosyoukinMstFieldCreateDate       = "sCreateDate"
	HosyoukinMstFieldUpdateNumber     = "sUpdateNumber"
	HosyoukinMstFieldUpdateDate       = "sUpdateDate"
)

type HosyoukinMst struct {
	Fields model.Attributes
}

func (h *HosyoukinMst) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &h.Fields)
}

func (h HosyoukinMst) IssueCode() string {
	return h.Fields.Value(HosyoukinMstFieldIssueCode)
}

func (h HosyoukinMst) MarketCode() string {
	return h.Fields.Value(HosyoukinMstFieldZyouzyouSizyou)
}
