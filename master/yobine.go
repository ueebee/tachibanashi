package master

import "github.com/ueebee/tachibanashi/model"

const (
	YobineFieldTaniNumber   = "sYobineTaniNumber"
	YobineFieldTekiyouDay   = "sTekiyouDay"
	YobineFieldKizunPrice1  = "sKizunPrice_1"
	YobineFieldKizunPrice2  = "sKizunPrice_2"
	YobineFieldKizunPrice3  = "sKizunPrice_3"
	YobineFieldKizunPrice4  = "sKizunPrice_4"
	YobineFieldKizunPrice5  = "sKizunPrice_5"
	YobineFieldKizunPrice6  = "sKizunPrice_6"
	YobineFieldKizunPrice7  = "sKizunPrice_7"
	YobineFieldKizunPrice8  = "sKizunPrice_8"
	YobineFieldKizunPrice9  = "sKizunPrice_9"
	YobineFieldKizunPrice10 = "sKizunPrice_10"
	YobineFieldKizunPrice11 = "sKizunPrice_11"
	YobineFieldKizunPrice12 = "sKizunPrice_12"
	YobineFieldKizunPrice13 = "sKizunPrice_13"
	YobineFieldKizunPrice14 = "sKizunPrice_14"
	YobineFieldKizunPrice15 = "sKizunPrice_15"
	YobineFieldKizunPrice16 = "sKizunPrice_16"
	YobineFieldKizunPrice17 = "sKizunPrice_17"
	YobineFieldKizunPrice18 = "sKizunPrice_18"
	YobineFieldKizunPrice19 = "sKizunPrice_19"
	YobineFieldKizunPrice20 = "sKizunPrice_20"
	YobineFieldTanka1       = "sYobineTanka_1"
	YobineFieldTanka2       = "sYobineTanka_2"
	YobineFieldTanka3       = "sYobineTanka_3"
	YobineFieldTanka4       = "sYobineTanka_4"
	YobineFieldTanka5       = "sYobineTanka_5"
	YobineFieldTanka6       = "sYobineTanka_6"
	YobineFieldTanka7       = "sYobineTanka_7"
	YobineFieldTanka8       = "sYobineTanka_8"
	YobineFieldTanka9       = "sYobineTanka_9"
	YobineFieldTanka10      = "sYobineTanka_10"
	YobineFieldTanka11      = "sYobineTanka_11"
	YobineFieldTanka12      = "sYobineTanka_12"
	YobineFieldTanka13      = "sYobineTanka_13"
	YobineFieldTanka14      = "sYobineTanka_14"
	YobineFieldTanka15      = "sYobineTanka_15"
	YobineFieldTanka16      = "sYobineTanka_16"
	YobineFieldTanka17      = "sYobineTanka_17"
	YobineFieldTanka18      = "sYobineTanka_18"
	YobineFieldTanka19      = "sYobineTanka_19"
	YobineFieldTanka20      = "sYobineTanka_20"
	YobineFieldDecimal1     = "sDecimal_1"
	YobineFieldDecimal2     = "sDecimal_2"
	YobineFieldDecimal3     = "sDecimal_3"
	YobineFieldDecimal4     = "sDecimal_4"
	YobineFieldDecimal5     = "sDecimal_5"
	YobineFieldDecimal6     = "sDecimal_6"
	YobineFieldDecimal7     = "sDecimal_7"
	YobineFieldDecimal8     = "sDecimal_8"
	YobineFieldDecimal9     = "sDecimal_9"
	YobineFieldDecimal10    = "sDecimal_10"
	YobineFieldDecimal11    = "sDecimal_11"
	YobineFieldDecimal12    = "sDecimal_12"
	YobineFieldDecimal13    = "sDecimal_13"
	YobineFieldDecimal14    = "sDecimal_14"
	YobineFieldDecimal15    = "sDecimal_15"
	YobineFieldDecimal16    = "sDecimal_16"
	YobineFieldDecimal17    = "sDecimal_17"
	YobineFieldDecimal18    = "sDecimal_18"
	YobineFieldDecimal19    = "sDecimal_19"
	YobineFieldDecimal20    = "sDecimal_20"
	YobineFieldCreateDate   = "sCreateDate"
	YobineFieldUpdateDate   = "sUpdateDate"
)

type Yobine struct {
	Fields model.Attributes
}

func (y *Yobine) UnmarshalJSON(data []byte) error {
	return unmarshalAttributes(data, &y.Fields)
}

func (y Yobine) TaniNumber() string {
	return y.Fields.Value(YobineFieldTaniNumber)
}

func (y Yobine) TekiyouDay() string {
	return y.Fields.Value(YobineFieldTekiyouDay)
}
