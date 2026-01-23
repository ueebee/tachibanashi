package master

import "github.com/ueebee/tachibanashi/model"

const (
	DaiyouKakemeFieldSystemKouzaKubun = "sSystemKouzaKubun"
	DaiyouKakemeFieldIssueCode        = "sIssueCode"
	DaiyouKakemeFieldTekiyouDay       = "sTekiyouDay"
	DaiyouKakemeFieldHosyokinKakeme   = "sHosyokinDaiyoKakeme"
	DaiyouKakemeFieldDeleteDay        = "sDeleteDay"
	DaiyouKakemeFieldCreateDate       = "sCreateDate"
	DaiyouKakemeFieldUpdateNumber     = "sUpdateNumber"
	DaiyouKakemeFieldUpdateDate       = "sUpdateDate"
)

type DaiyouKakeme struct {
	Fields model.Attributes
}

func (d *DaiyouKakeme) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &d.Fields)
}

func (d DaiyouKakeme) IssueCode() string {
	return d.Fields.Value(DaiyouKakemeFieldIssueCode)
}

func (d DaiyouKakeme) TekiyouDay() string {
	return d.Fields.Value(DaiyouKakemeFieldTekiyouDay)
}
