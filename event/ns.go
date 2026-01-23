package event

import "github.com/ueebee/tachibanashi/model"

type NS struct {
	Frame
	Fields        model.Attributes
	Provider      string
	EventNo       string
	Alert         bool
	NewsID        string
	NewsDate      string
	NewsTime      string
	CategoryCount int64
	Categories    []string
	GenreCount    int64
	Genres        []string
	IssueCount    int64
	Issues        []string
	SkipFlag      string
	UpdateFlag    string
	Headline      string
	Body          string
}

func parseNS(frame Frame) (NS, error) {
	fields := frameAttributes(frame.Fields)

	ns := NS{
		Frame:         frame,
		Fields:        fields,
		Provider:      fields.Value("p_PV"),
		EventNo:       fields.Value("p_ENO"),
		Alert:         parseBool(fields.Value("p_ALT")),
		NewsID:        fields.Value("p_ID"),
		NewsDate:      fields.Value("p_DT"),
		NewsTime:      fields.Value("p_TM"),
		CategoryCount: parseInt64(fields.Value("p_CGN")),
		Categories:    frame.Values("p_CGL"),
		GenreCount:    parseInt64(fields.Value("p_GRN")),
		Genres:        frame.Values("p_GRL"),
		IssueCount:    parseInt64(fields.Value("p_ISN")),
		Issues:        frame.Values("p_ISL"),
		SkipFlag:      fields.Value("p_SKF"),
		UpdateFlag:    fields.Value("p_UPD"),
		Headline:      fields.Value("p_HDL"),
		Body:          fields.Value("p_TX"),
	}

	return ns, nil
}
