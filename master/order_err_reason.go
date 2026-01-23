package master

import "github.com/ueebee/tachibanashi/model"

const (
	OrderErrReasonFieldCode = "sErrReasonCode"
	OrderErrReasonFieldText = "sErrReasonText"
)

type OrderErrReason struct {
	Fields model.Attributes
}

func (o *OrderErrReason) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &o.Fields)
}

func (o OrderErrReason) Code() string {
	return o.Fields.Value(OrderErrReasonFieldCode)
}

func (o OrderErrReason) Text() string {
	return o.Fields.Value(OrderErrReasonFieldText)
}
