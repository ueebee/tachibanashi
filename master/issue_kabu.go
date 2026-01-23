package master

import "github.com/ueebee/tachibanashi/model"

const (
	IssueKabuFieldCode                  = "sIssueCode"
	IssueKabuFieldName                  = "sIssueName"
	IssueKabuFieldNameRyaku             = "sIssueNameRyaku"
	IssueKabuFieldNameKana              = "sIssueNameKana"
	IssueKabuFieldNameEizi              = "sIssueNameEizi"
	IssueKabuFieldTokuteiF              = "sTokuteiF"
	IssueKabuFieldHikazeiC              = "sHikazeiC"
	IssueKabuFieldZyouzyouHakkouKabusu  = "sZyouzyouHakkouKabusu"
	IssueKabuFieldKenriotiFlag          = "sKenriotiFlag"
	IssueKabuFieldKenritukiSaisyuDay    = "sKenritukiSaisyuDay"
	IssueKabuFieldZyouzyouNyusatuC      = "sZyouzyouNyusatuC"
	IssueKabuFieldNyusatuKaizyoDay      = "sNyusatuKaizyoDay"
	IssueKabuFieldNyusatuDay            = "sNyusatuDay"
	IssueKabuFieldBaibaiTani            = "sBaibaiTani"
	IssueKabuFieldBaibaiTaniYoku        = "sBaibaiTaniYoku"
	IssueKabuFieldBaibaiTeisiC          = "sBaibaiTeisiC"
	IssueKabuFieldHakkouKaisiDay        = "sHakkouKaisiDay"
	IssueKabuFieldHakkouSaisyuDay       = "sHakkouSaisyuDay"
	IssueKabuFieldKessanC               = "sKessanC"
	IssueKabuFieldKessanDay             = "sKessanDay"
	IssueKabuFieldZyouzyouOutouDay      = "sZyouzyouOutouDay"
	IssueKabuFieldNiruiKizituC          = "sNiruiKizituC"
	IssueKabuFieldOogutiKabusu          = "sOogutiKabusu"
	IssueKabuFieldOogutiKingaku         = "sOogutiKingaku"
	IssueKabuFieldBadenpyouOutputYNC    = "sBadenpyouOutputYNC"
	IssueKabuFieldHosyoukinDaiyouKakeme = "sHosyoukinDaiyouKakeme"
	IssueKabuFieldDaiyouHyoukaTanka     = "sDaiyouHyoukaTanka"
	IssueKabuFieldKikoSankaC            = "sKikoSankaC"
	IssueKabuFieldKarikessaiC           = "sKarikessaiC"
	IssueKabuFieldYusenSizyou           = "sYusenSizyou"
	IssueKabuFieldMukigenC              = "sMukigenC"
	IssueKabuFieldGyousyuCode           = "sGyousyuCode"
	IssueKabuFieldGyousyuName           = "sGyousyuName"
	IssueKabuFieldSorC                  = "sSorC"
	IssueKabuFieldCreateDate            = "sCreateDate"
	IssueKabuFieldUpdateDate            = "sUpdateDate"
	IssueKabuFieldUpdateNumber          = "sUpdateNumber"
)

type IssueMstKabu struct {
	Fields model.Attributes
}

func (i *IssueMstKabu) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &i.Fields)
}

func (i IssueMstKabu) IssueCode() string {
	return i.Fields.Value(IssueKabuFieldCode)
}

func (i IssueMstKabu) IssueName() string {
	return i.Fields.Value(IssueKabuFieldName)
}
