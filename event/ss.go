package event

import "github.com/ueebee/tachibanashi/model"

type SS struct {
	Frame
	Fields       model.Attributes
	Provider     string
	EventNo      string
	Alert        bool
	ChangedAt    string
	LoginKind    string
	SystemStatus string
}

func parseSS(frame Frame) (SS, error) {
	fields := frameAttributes(frame.Fields)

	ss := SS{
		Frame:        frame,
		Fields:       fields,
		Provider:     fields.Value("p_PV"),
		EventNo:      fields.Value("p_ENO"),
		Alert:        parseBool(fields.Value("p_ALT")),
		ChangedAt:    fields.Value("p_CT"),
		LoginKind:    fields.Value("p_LK"),
		SystemStatus: fields.Value("p_SS"),
	}

	return ss, nil
}
