package master

import "github.com/ueebee/tachibanashi/model"

const (
	IssueOpFieldIssueCode           = "sIssueCode"
	IssueOpFieldIssueName           = "sIssueName"
	IssueOpFieldIssueNameEizi       = "sIssueNameEizi"
	IssueOpFieldSakOpSyouhin        = "sSakOpSyouhin"
	IssueOpFieldGensisanKubun       = "sGensisanKubun"
	IssueOpFieldGensisanCode        = "sGensisanCode"
	IssueOpFieldGengetu             = "sGengetu"
	IssueOpFieldZyouzyouSizyou      = "sZyouzyouSizyou"
	IssueOpFieldKousiPrice          = "sKousiPrice"
	IssueOpFieldPutCall             = "sPutCall"
	IssueOpFieldTorihikiStartDay    = "sTorihikiStartDay"
	IssueOpFieldLastBaibaiDay       = "sLastBaibaiDay"
	IssueOpFieldKenrikousiLastDay   = "sKenrikousiLastDay"
	IssueOpFieldTaniSuryou          = "sTaniSuryou"
	IssueOpFieldYobineTaniNumber    = "sYobineTaniNumber"
	IssueOpFieldZyouhouSource       = "sZyouhouSource"
	IssueOpFieldZyouhouCode         = "sZyouhouCode"
	IssueOpFieldNehabaMin           = "sNehabaMin"
	IssueOpFieldNehabaMax           = "sNehabaMax"
	IssueOpFieldIssueKisei1C        = "sIssueKisei1C"
	IssueOpFieldZenzituOwarine      = "sZenzituOwarine"
	IssueOpFieldZenzituRironPrice   = "sZenzituRironPrice"
	IssueOpFieldBaDenpyouOutputUmuC = "sBaDenpyouOutputUmuC"
	IssueOpFieldCreateDate          = "sCreateDate"
	IssueOpFieldUpdateDate          = "sUpdateDate"
	IssueOpFieldUpdateNumber        = "sUpdateNumber"
	IssueOpFieldATMFlag             = "sATMFlag"
)

type IssueMstOp struct {
	Fields model.Attributes
}

func (i *IssueMstOp) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &i.Fields)
}

func (i IssueMstOp) IssueCode() string {
	return i.Fields.Value(IssueOpFieldIssueCode)
}
