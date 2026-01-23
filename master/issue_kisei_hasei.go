package master

import "github.com/ueebee/tachibanashi/model"

const (
	IssueKiseiHaseiFieldSystemKouzaKubun = "sSystemKouzaKubun"
	IssueKiseiHaseiFieldIssueCode        = "sIssueCode"
	IssueKiseiHaseiFieldZyouzyouSizyou   = "sZyouzyouSizyou"
	IssueKiseiHaseiFieldTeisiKubun       = "sTeisiKubun"
	IssueKiseiHaseiFieldKaitate          = "sKaitate"
	IssueKiseiHaseiFieldKaitateYoku      = "sKaitateYoku"
	IssueKiseiHaseiFieldUritate          = "sUritate"
	IssueKiseiHaseiFieldUritateYoku      = "sUritateYoku"
	IssueKiseiHaseiFieldKaiHensai        = "sKaiHensai"
	IssueKiseiHaseiFieldKaiHensaiYoku    = "sKaiHensaiYoku"
	IssueKiseiHaseiFieldUriHensai        = "sUriHensai"
	IssueKiseiHaseiFieldUriHensaiYoku    = "sUriHensaiYoku"
	IssueKiseiHaseiFieldCreateDate       = "sCreateDate"
	IssueKiseiHaseiFieldUpdateDate       = "sUpdateDate"
	IssueKiseiHaseiFieldUpdateNumber     = "sUpdateNumber"
)

type IssueSizyouKiseiHasei struct {
	Fields model.Attributes
}

func (i *IssueSizyouKiseiHasei) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &i.Fields)
}

func (i IssueSizyouKiseiHasei) IssueCode() string {
	return i.Fields.Value(IssueKiseiHaseiFieldIssueCode)
}

func (i IssueSizyouKiseiHasei) MarketCode() string {
	return i.Fields.Value(IssueKiseiHaseiFieldZyouzyouSizyou)
}
