package client

import "github.com/ueebee/tachibanashi/model"

type CommonParamsCarrier interface {
	Params() *model.CommonParams
}
