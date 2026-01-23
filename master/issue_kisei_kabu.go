package master

import "github.com/ueebee/tachibanashi/model"

const (
	IssueKiseiKabuFieldSystemKouzaKubun            = "sSystemKouzaKubun"
	IssueKiseiKabuFieldIssueCode                   = "sIssueCode"
	IssueKiseiKabuFieldZyouzyouSizyou              = "sZyouzyouSizyou"
	IssueKiseiKabuFieldTeisiKubun                  = "sTeisiKubun"
	IssueKiseiKabuFieldGenbutuKaituke              = "sGenbutuKaituke"
	IssueKiseiKabuFieldGenbutuKaitukeYoku          = "sGenbutuKaitukeYoku"
	IssueKiseiKabuFieldGenbutuUrituke              = "sGenbutuUrituke"
	IssueKiseiKabuFieldGenbutuUritukeYoku          = "sGenbutuUritukeYoku"
	IssueKiseiKabuFieldSeidoSinyouSinkiKaitate     = "sSeidoSinyouSinkiKaitate"
	IssueKiseiKabuFieldSeidoSinyouSinkiKaitateYoku = "sSeidoSinyouSinkiKaitateYoku"
	IssueKiseiKabuFieldSeidoSinyouSinkiUritate     = "sSeidoSinyouSinkiUritate"
	IssueKiseiKabuFieldSeidoSinyouSinkiUritateYoku = "sSeidoSinyouSinkiUritateYoku"
	IssueKiseiKabuFieldSeidoSinyouKaiHensai        = "sSeidoSinyouKaiHensai"
	IssueKiseiKabuFieldSeidoSinyouKaiHensaiYoku    = "sSeidoSinyouKaiHensaiYoku"
	IssueKiseiKabuFieldSeidoSinyouUriHensai        = "sSeidoSinyouUriHensai"
	IssueKiseiKabuFieldSeidoSinyouUriHensaiYoku    = "sSeidoSinyouUriHensaiYoku"
	IssueKiseiKabuFieldSeidoSinyouGenbiki          = "sSeidoSinyouGenbiki"
	IssueKiseiKabuFieldSeidoSinyouGenbikiYoku      = "sSeidoSinyouGenbikiYoku"
	IssueKiseiKabuFieldSeidoSinyouGenwatasi        = "sSeidoSinyouGenwatasi"
	IssueKiseiKabuFieldSeidoSinyouGenwatasiYoku    = "sSeidoSinyouGenwatasiYoku"
	IssueKiseiKabuFieldIppanSinyouSinkiKaitate     = "sIppanSinyouSinkiKaitate"
	IssueKiseiKabuFieldIppanSinyouSinkiKaitateYoku = "sIppanSinyouSinkiKaitateYoku"
	IssueKiseiKabuFieldIppanSinyouSinkiUritate     = "sIppanSinyouSinkiUritate"
	IssueKiseiKabuFieldIppanSinyouSinkiUritateYoku = "sIppanSinyouSinkiUritateYoku"
	IssueKiseiKabuFieldIppanSinyouKaiHensai        = "sIppanSinyouKaiHensai"
	IssueKiseiKabuFieldIppanSinyouKaiHensaiYoku    = "sIppanSinyouKaiHensaiYoku"
	IssueKiseiKabuFieldIppanSinyouUriHensai        = "sIppanSinyouUriHensai"
	IssueKiseiKabuFieldIppanSinyouUriHensaiYoku    = "sIppanSinyouUriHensaiYoku"
	IssueKiseiKabuFieldIppanSinyouGenbiki          = "sIppanSinyouGenbiki"
	IssueKiseiKabuFieldIppanSinyouGenbikiYoku      = "sIppanSinyouGenbikiYoku"
	IssueKiseiKabuFieldIppanSinyouGenwatasi        = "sIppanSinyouGenwatasi"
	IssueKiseiKabuFieldIppanSinyouGenwatasiYoku    = "sIppanSinyouGenwatasiYoku"
	IssueKiseiKabuFieldZizenCyouseiC               = "sZizenCyouseiC"
	IssueKiseiKabuFieldZizenCyouseiCYoku           = "sZizenCyouseiCYoku"
	IssueKiseiKabuFieldSokuzituNyukinC             = "sSokuzituNyukinC"
	IssueKiseiKabuFieldSokuzituNyukinCYoku         = "sSokuzituNyukinCYoku"
	IssueKiseiKabuFieldSokuzituNyukinKiseiDate     = "sSokuzituNyukinKiseiDate"
	IssueKiseiKabuFieldSinyouSyutyuKubun           = "sSinyouSyutyuKubun"
	IssueKiseiKabuFieldSinyouSyutyuKubunYoku       = "sSinyouSyutyuKubunYoku"
	IssueKiseiKabuFieldCreateDate                  = "sCreateDate"
	IssueKiseiKabuFieldUpdateDate                  = "sUpdateDate"
	IssueKiseiKabuFieldUpdateNumber                = "sUpdateNumber"
)

type IssueSizyouKiseiKabu struct {
	Fields model.Attributes
}

func (i *IssueSizyouKiseiKabu) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &i.Fields)
}

func (i IssueSizyouKiseiKabu) IssueCode() string {
	return i.Fields.Value(IssueKiseiKabuFieldIssueCode)
}

func (i IssueSizyouKiseiKabu) MarketCode() string {
	return i.Fields.Value(IssueKiseiKabuFieldZyouzyouSizyou)
}
