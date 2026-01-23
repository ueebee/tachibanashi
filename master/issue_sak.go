package master

import "github.com/ueebee/tachibanashi/model"

const (
	IssueSakFieldIssueCode           = "sIssueCode"
	IssueSakFieldIssueName           = "sIssueName"
	IssueSakFieldIssueNameEizi       = "sIssueNameEizi"
	IssueSakFieldSakOpSyouhin        = "sSakOpSyouhin"
	IssueSakFieldGensisanKubun       = "sGensisanKubun"
	IssueSakFieldGensisanCode        = "sGensisanCode"
	IssueSakFieldGengetu             = "sGengetu"
	IssueSakFieldZyouzyouSizyou      = "sZyouzyouSizyou"
	IssueSakFieldTorihikiStartDay    = "sTorihikiStartDay"
	IssueSakFieldLastBaibaiDay       = "sLastBaibaiDay"
	IssueSakFieldTaniSuryou          = "sTaniSuryou"
	IssueSakFieldYobineTaniNumber    = "sYobineTaniNumber"
	IssueSakFieldZyouhouSource       = "sZyouhouSource"
	IssueSakFieldZyouhouCode         = "sZyouhouCode"
	IssueSakFieldNehabaMin           = "sNehabaMin"
	IssueSakFieldNehabaMax           = "sNehabaMax"
	IssueSakFieldIssueKisei1C        = "sIssueKisei1C"
	IssueSakFieldBaibaiTeisiC        = "sBaibaiTeisiC"
	IssueSakFieldZenzituOwarine      = "sZenzituOwarine"
	IssueSakFieldBaDenpyouOutputUmuC = "sBaDenpyouOutputUmuC"
	IssueSakFieldCreateDate          = "sCreateDate"
	IssueSakFieldUpdateDate          = "sUpdateDate"
	IssueSakFieldUpdateNumber        = "sUpdateNumber"
)

type IssueMstSak struct {
	Fields model.Attributes
}

func (i *IssueMstSak) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &i.Fields)
}

func (i IssueMstSak) IssueCode() string {
	return i.Fields.Value(IssueSakFieldIssueCode)
}
