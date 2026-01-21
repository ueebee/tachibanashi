package request

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ueebee/tachibanashi/model"
)

const (
	clmZanKaiKanougakuSuii          = "CLMZanKaiKanougakuSuii"
	clmZanKaiGenbutuKaitukeSyousai  = "CLMZanKaiGenbutuKaitukeSyousai"
	clmZanKaiSinyouSinkidateSyousai = "CLMZanKaiSinyouSinkidateSyousai"
	clmZanRealHosyoukinRitu         = "CLMZanRealHosyoukinRitu"
)

type ZanKaiKanougakuSuiiRequest struct {
	model.CommonParams
	CLMID string `json:"sCLMID"`
}

func (r *ZanKaiKanougakuSuiiRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type KanougakuSuiiEntry struct {
	Fields model.Attributes
}

func (e *KanougakuSuiiEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		e.Fields[key] = jsonString(value)
	}
	return nil
}

type ZanKaiKanougakuSuiiResponse struct {
	model.CommonResponse
	WarningCode string
	WarningText string
	IssueCode   string
	UpdateDate  string
	NearaiKubun string
	Entries     []KanougakuSuiiEntry
	Fields      model.Attributes
}

func (r *ZanKaiKanougakuSuiiResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "aKanougakuSuiiList" {
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
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	r.IssueCode = values["sIssueCode"]
	r.UpdateDate = values["sUpdateDate"]
	r.NearaiKubun = values["sNearaiKubun"]
	return nil
}

func (s *Service) ZanKaiKanougakuSuii(ctx context.Context) (*ZanKaiKanougakuSuiiResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}
	req := ZanKaiKanougakuSuiiRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmZanKaiKanougakuSuii,
	}
	var resp ZanKaiKanougakuSuiiResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZanKaiGenbutuKaitukeSyousaiRequest struct {
	model.CommonParams
	CLMID       string `json:"sCLMID"`
	HitukeIndex string `json:"sHitukeIndex"`
}

func (r *ZanKaiGenbutuKaitukeSyousaiRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type ZanKaiGenbutuKaitukeSyousaiResponse struct {
	model.CommonResponse
	WarningCode string
	WarningText string
	HitukeIndex string
	Hituke      string
	Fields      model.Attributes
}

func (r *ZanKaiGenbutuKaitukeSyousaiResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		values[key] = str
		r.Fields[key] = str
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	r.HitukeIndex = values["sHitukeIndex"]
	r.Hituke = values["sHituke"]
	return nil
}

func (s *Service) ZanKaiGenbutuKaitukeSyousai(ctx context.Context, hitukeIndex string) (*ZanKaiGenbutuKaitukeSyousaiResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}
	req := ZanKaiGenbutuKaitukeSyousaiRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmZanKaiGenbutuKaitukeSyousai,
		HitukeIndex:  hitukeIndex,
	}
	var resp ZanKaiGenbutuKaitukeSyousaiResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZanKaiSinyouSinkidateSyousaiRequest struct {
	model.CommonParams
	CLMID       string `json:"sCLMID"`
	HitukeIndex string `json:"sHitukeIndex"`
}

func (r *ZanKaiSinyouSinkidateSyousaiRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type ZanKaiSinyouSinkidateSyousaiResponse struct {
	model.CommonResponse
	WarningCode string
	WarningText string
	HitukeIndex string
	Hituke      string
	Fields      model.Attributes
}

func (r *ZanKaiSinyouSinkidateSyousaiResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		values[key] = str
		r.Fields[key] = str
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	r.HitukeIndex = values["sHitukeIndex"]
	r.Hituke = values["sHituke"]
	return nil
}

func (s *Service) ZanKaiSinyouSinkidateSyousai(ctx context.Context, hitukeIndex string) (*ZanKaiSinyouSinkidateSyousaiResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}
	req := ZanKaiSinyouSinkidateSyousaiRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmZanKaiSinyouSinkidateSyousai,
		HitukeIndex:  hitukeIndex,
	}
	var resp ZanKaiSinyouSinkidateSyousaiResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZanRealHosyoukinRituRequest struct {
	model.CommonParams
	CLMID string `json:"sCLMID"`
}

func (r *ZanRealHosyoukinRituRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type ZanRealHosyoukinRituResponse struct {
	model.CommonResponse
	WarningCode string
	WarningText string
	Fields      model.Attributes
}

func (r *ZanRealHosyoukinRituResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		values[key] = str
		r.Fields[key] = str
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.WarningCode = values["sWarningCode"]
	r.WarningText = values["sWarningText"]
	return nil
}

func (s *Service) ZanRealHosyoukinRitu(ctx context.Context) (*ZanRealHosyoukinRituResponse, error) {
	path, err := s.requestURL()
	if err != nil {
		return nil, err
	}
	req := ZanRealHosyoukinRituRequest{
		CommonParams: model.CommonParams{JsonOfmt: "5"},
		CLMID:        clmZanRealHosyoukinRitu,
	}
	var resp ZanRealHosyoukinRituResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
