package price

import "github.com/ueebee/tachibanashi/model"

var quoteFieldNames = map[string]struct{}{
	"pAAV":   {},
	"pABV":   {},
	"pAV":    {},
	"pBV":    {},
	"xDCFS":  {},
	"pDHF":   {},
	"pDHP":   {},
	"tDHP:T": {},
	"pDJ":    {},
	"pDLF":   {},
	"pDLP":   {},
	"tDLP:T": {},
	"pDOP":   {},
	"tDOP:T": {},
	"pDPG":   {},
	"pDPP":   {},
	"tDPP:T": {},
	"pDV":    {},
	"xDVES":  {},
	"pDYRP":  {},
	"pDYWP":  {},
	"pGAV10": {},
	"pGAP10": {},
	"pGAV9":  {},
	"pGAP9":  {},
	"pGAV8":  {},
	"pGAP8":  {},
	"pGAV7":  {},
	"pGAP7":  {},
	"pGAV6":  {},
	"pGAP6":  {},
	"pGAV5":  {},
	"pGAP5":  {},
	"pGAV4":  {},
	"pGAP4":  {},
	"pGAV3":  {},
	"pGAP3":  {},
	"pGAV2":  {},
	"pGAP2":  {},
	"pGAV1":  {},
	"pGAP1":  {},
	"pGBV1":  {},
	"pGBP1":  {},
	"pGBV2":  {},
	"pGBP2":  {},
	"pGBV3":  {},
	"pGBP3":  {},
	"pGBV4":  {},
	"pGBP4":  {},
	"pGBV5":  {},
	"pGBP5":  {},
	"pGBV6":  {},
	"pGBP6":  {},
	"pGBV7":  {},
	"pGBP7":  {},
	"pGBV8":  {},
	"pGBP8":  {},
	"pGBV9":  {},
	"pGBP9":  {},
	"pGBV10": {},
	"pGBP10": {},
	"xLISS":  {},
	"pPRP":   {},
	"pQAP":   {},
	"pQAS":   {},
	"pQBP":   {},
	"pQBS":   {},
	"pQOV":   {},
	"pQUV":   {},
	"pVWAP":  {},
}

func IsValidQuoteField(name string) bool {
	_, ok := quoteFieldNames[name]
	return ok
}

func ValidQuoteFields() []string {
	out := make([]string, 0, len(quoteFieldNames))
	for key := range quoteFieldNames {
		out = append(out, key)
	}
	return out
}

var defaultQuoteFields = []string{
	model.FieldLastPrice,
	model.FieldLastTime,
	model.FieldPrevClose,
}

func DefaultQuoteFields() []string {
	return append([]string(nil), defaultQuoteFields...)
}
