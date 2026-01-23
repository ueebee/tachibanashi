package event

import "github.com/ueebee/tachibanashi/model"

type US struct {
	Frame
	Fields          model.Attributes
	Provider        string
	EventNo         string
	Alert           bool
	ChangedAt       string
	MarketCode      string
	UnderlyingCode  string
	InstrumentKind  string
	OperationCode   string
	OperationUnit   string
	BusinessDayKind string
	OperationStatus string
}

func parseUS(frame Frame) (US, error) {
	fields := frameAttributes(frame.Fields)

	us := US{
		Frame:           frame,
		Fields:          fields,
		Provider:        fields.Value("p_PV"),
		EventNo:         fields.Value("p_ENO"),
		Alert:           parseBool(fields.Value("p_ALT")),
		ChangedAt:       fields.Value("p_CT"),
		MarketCode:      fields.Value("p_MC"),
		UnderlyingCode:  fields.Value("p_GSCD"),
		InstrumentKind:  fields.Value("p_SHSB"),
		OperationCode:   fields.Value("p_UC"),
		OperationUnit:   fields.Value("p_UU"),
		BusinessDayKind: fields.Value("p_EDK"),
		OperationStatus: fields.Value("p_US"),
	}

	return us, nil
}
