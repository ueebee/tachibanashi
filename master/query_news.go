package master

import (
	"context"
	"encoding/json"
	"net/http"

	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/model"
)

const (
	clmGetNewsHead = "CLMMfdsGetNewsHead"
	clmGetNewsBody = "CLMMfdsGetNewsBody"
)

type NewsHeadRequest struct {
	model.CommonParams
	CLMID        string `json:"sCLMID"`
	Category     string `json:"p_CG,omitempty"`
	IssueCode    string `json:"p_IS,omitempty"`
	DateFrom     string `json:"p_DT_FROM,omitempty"`
	DateTo       string `json:"p_DT_TO,omitempty"`
	RecordOffset string `json:"p_REC_OFST,omitempty"`
	RecordLimit  string `json:"p_REC_LIMT,omitempty"`
}

func (r *NewsHeadRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type NewsHeadResponse struct {
	model.CommonResponse
	RecordMax string
	Entries   []NewsHeadEntry
	Fields    model.Attributes
}

type NewsHeadEntry struct {
	ID           string
	Date         string
	Time         string
	CategoryList string
	GenreList    string
	IssueList    string
	Headline     string
	Fields       model.Attributes
}

func (r *NewsHeadResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "aCLMMfdsNewsHead" {
			if err := decodeList(value, &r.Entries); err != nil {
				return err
			}
			continue
		}
		str := jsonString(value)
		values[key] = str
		r.Fields[key] = str
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.RecordMax = values["p_REC_MAX"]
	return nil
}

func (e *NewsHeadEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		e.Fields[key] = jsonString(value)
	}
	e.ID = e.Fields.Value("p_ID")
	e.Date = e.Fields.Value("p_DT")
	e.Time = e.Fields.Value("p_TM")
	e.CategoryList = e.Fields.Value("p_CGL")
	e.GenreList = e.Fields.Value("p_GNL")
	e.IssueList = e.Fields.Value("p_ISL")
	if encoded := e.Fields.Value("p_HDL"); encoded != "" {
		decoded, err := decodeBase64ShiftJIS(encoded)
		if err != nil {
			return err
		}
		e.Headline = decoded
		e.Fields["p_HDL"] = decoded
	}
	return nil
}

func (s *Service) NewsHead(ctx context.Context, req NewsHeadRequest) (*NewsHeadResponse, error) {
	url, err := s.masterURL()
	if err != nil {
		return nil, err
	}
	if req.CLMID == "" {
		req.CLMID = clmGetNewsHead
	}
	if req.JsonOfmt == "" {
		req.JsonOfmt = defaultJSONFormat
	}

	var resp NewsHeadResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, url, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type NewsBodyRequest struct {
	model.CommonParams
	CLMID string `json:"sCLMID"`
	ID    string `json:"p_ID"`
}

func (r *NewsBodyRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type NewsBodyResponse struct {
	model.CommonResponse
	RequestID string
	Entries   []NewsBodyEntry
	Fields    model.Attributes
}

type NewsBodyEntry struct {
	ID           string
	Date         string
	Time         string
	CategoryList string
	GenreList    string
	IssueList    string
	Headline     string
	Body         string
	Fields       model.Attributes
}

func (r *NewsBodyResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "aCLMMfdsNewsBody" {
			if err := decodeList(value, &r.Entries); err != nil {
				return err
			}
			continue
		}
		str := jsonString(value)
		values[key] = str
		r.Fields[key] = str
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.RequestID = values["p_ID"]
	return nil
}

func (e *NewsBodyEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		e.Fields[key] = jsonString(value)
	}
	e.ID = e.Fields.Value("p_ID")
	e.Date = e.Fields.Value("p_DT")
	e.Time = e.Fields.Value("p_TM")
	e.CategoryList = e.Fields.Value("p_CGL")
	e.GenreList = e.Fields.Value("p_GNL")
	e.IssueList = e.Fields.Value("p_ISL")
	if encoded := e.Fields.Value("p_HDL"); encoded != "" {
		decoded, err := decodeBase64ShiftJIS(encoded)
		if err != nil {
			return err
		}
		e.Headline = decoded
		e.Fields["p_HDL"] = decoded
	}
	if encoded := e.Fields.Value("p_TX"); encoded != "" {
		decoded, err := decodeBase64ShiftJIS(encoded)
		if err != nil {
			return err
		}
		e.Body = decoded
		e.Fields["p_TX"] = decoded
	}
	return nil
}

func (s *Service) NewsBody(ctx context.Context, id string) (*NewsBodyResponse, error) {
	if id == "" {
		return nil, &terrors.ValidationError{Field: "id", Reason: "required"}
	}
	url, err := s.masterURL()
	if err != nil {
		return nil, err
	}

	req := NewsBodyRequest{
		CommonParams: model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:        clmGetNewsBody,
		ID:           id,
	}

	var resp NewsBodyResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, url, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
