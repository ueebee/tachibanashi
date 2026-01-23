package master

import "github.com/ueebee/tachibanashi/model"

const (
	IssueSizyouKabuFieldIssueCode                = "sIssueCode"
	IssueSizyouKabuFieldZyouzyouSizyou           = "sZyouzyouSizyou"
	IssueSizyouKabuFieldSystemC                  = "sSystemC"
	IssueSizyouKabuFieldNehabaMin                = "sNehabaMin"
	IssueSizyouKabuFieldNehabaMax                = "sNehabaMax"
	IssueSizyouKabuFieldIssueKubunC              = "sIssueKubunC"
	IssueSizyouKabuFieldNehabaSizyouC            = "sNehabaSizyouC"
	IssueSizyouKabuFieldSinyouC                  = "sSinyouC"
	IssueSizyouKabuFieldSinkiZyouzyouDay         = "sSinkiZyouzyouDay"
	IssueSizyouKabuFieldNehabaKigenDay           = "sNehabaKigenDay"
	IssueSizyouKabuFieldNehabaKiseiC             = "sNehabaKiseiC"
	IssueSizyouKabuFieldNehabaKiseiTi            = "sNehabaKiseiTi"
	IssueSizyouKabuFieldNehabaCheckKahiC         = "sNehabaCheckKahiC"
	IssueSizyouKabuFieldIssueBubetuC             = "sIssueBubetuC"
	IssueSizyouKabuFieldZenzituOwarine           = "sZenzituOwarine"
	IssueSizyouKabuFieldNehabaSansyutuSizyouC    = "sNehabaSansyutuSizyouC"
	IssueSizyouKabuFieldIssueKisei1C             = "sIssueKisei1C"
	IssueSizyouKabuFieldIssueKisei2C             = "sIssueKisei2C"
	IssueSizyouKabuFieldZyouzyouKubun            = "sZyouzyouKubun"
	IssueSizyouKabuFieldZyouzyouHaisiDay         = "sZyouzyouHaisiDay"
	IssueSizyouKabuFieldSizyoubetuBaibaiTani     = "sSizyoubetuBaibaiTani"
	IssueSizyouKabuFieldSizyoubetuBaibaiTaniYoku = "sSizyoubetuBaibaiTaniYoku"
	IssueSizyouKabuFieldYobineTaniNumber         = "sYobineTaniNumber"
	IssueSizyouKabuFieldYobineTaniNumberYoku     = "sYobineTaniNumberYoku"
	IssueSizyouKabuFieldZyouhouSource            = "sZyouhouSource"
	IssueSizyouKabuFieldZyouhouCode              = "sZyouhouCode"
	IssueSizyouKabuFieldKouboPrice               = "sKouboPrice"
	IssueSizyouKabuFieldCreateDate               = "sCreateDate"
	IssueSizyouKabuFieldUpdateDate               = "sUpdateDate"
	IssueSizyouKabuFieldUpdateNumber             = "sUpdateNumber"
)

type IssueSizyouMstKabu struct {
	Fields model.Attributes
}

func (i *IssueSizyouMstKabu) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &i.Fields)
}

func (i IssueSizyouMstKabu) IssueCode() string {
	return i.Fields.Value(IssueSizyouKabuFieldIssueCode)
}

func (i IssueSizyouMstKabu) MarketCode() string {
	return i.Fields.Value(IssueSizyouKabuFieldZyouzyouSizyou)
}
