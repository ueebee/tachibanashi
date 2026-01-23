package master

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/model"
)

const (
	clmGetIssueDetail = "CLMMfdsGetIssueDetail"
	clmGetSyoukinZan  = "CLMMfdsGetSyoukinZan"
	clmGetShinyouZan  = "CLMMfdsGetShinyouZan"
	clmGetHibuInfo    = "CLMMfdsGetHibuInfo"
)

type IssueDetailRequest struct {
	model.CommonParams
	CLMID           string `json:"sCLMID"`
	TargetIssueCode string `json:"sTargetIssueCode"`
}

func (r *IssueDetailRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type IssueDetailResponse struct {
	model.CommonResponse
	Entries []IssueDetailEntry
	Fields  model.Attributes
}

type IssueDetailEntry struct {
	IssueCode string
	Fields    model.Attributes
}

func (r *IssueDetailResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "aCLMMfdsIssueDetail" {
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
	return nil
}

func (e *IssueDetailEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		e.Fields[key] = str
		if key == "sIssueCode" {
			e.IssueCode = str
		}
	}
	return nil
}

func (s *Service) IssueDetail(ctx context.Context, issueCodes []string) (*IssueDetailResponse, error) {
	codes := normalizeList(issueCodes)
	if len(codes) == 0 {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "required"}
	}
	if len(codes) > maxIssueCodes {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "max 120"}
	}

	url, err := s.masterURL()
	if err != nil {
		return nil, err
	}

	req := IssueDetailRequest{
		CommonParams:    model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:           clmGetIssueDetail,
		TargetIssueCode: strings.Join(codes, ","),
	}

	var resp IssueDetailResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, url, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type SyoukinZanRequest struct {
	model.CommonParams
	CLMID           string `json:"sCLMID"`
	TargetIssueCode string `json:"sTargetIssueCode"`
}

func (r *SyoukinZanRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type SyoukinZanResponse struct {
	model.CommonResponse
	Entries []SyoukinZanEntry
	Fields  model.Attributes
}

type SyoukinZanEntry struct {
	IssueCode string
	Fields    model.Attributes
}

func (r *SyoukinZanResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "aCLMMfdsSyoukinZan" {
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
	return nil
}

func (e *SyoukinZanEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		e.Fields[key] = str
		if key == "sIssueCode" {
			e.IssueCode = str
		}
	}
	return nil
}

func (s *Service) SyoukinZan(ctx context.Context, issueCodes []string) (*SyoukinZanResponse, error) {
	codes := normalizeList(issueCodes)
	if len(codes) == 0 {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "required"}
	}
	if len(codes) > maxIssueCodes {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "max 120"}
	}

	url, err := s.masterURL()
	if err != nil {
		return nil, err
	}

	req := SyoukinZanRequest{
		CommonParams:    model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:           clmGetSyoukinZan,
		TargetIssueCode: strings.Join(codes, ","),
	}

	var resp SyoukinZanResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, url, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ShinyouZanRequest struct {
	model.CommonParams
	CLMID           string `json:"sCLMID"`
	TargetIssueCode string `json:"sTargetIssueCode"`
}

func (r *ShinyouZanRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type ShinyouZanResponse struct {
	model.CommonResponse
	Entries []ShinyouZanEntry
	Fields  model.Attributes
}

type ShinyouZanEntry struct {
	IssueCode string
	Fields    model.Attributes
}

func (r *ShinyouZanResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "aCLMMfdsShinyouZan" {
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
	return nil
}

func (e *ShinyouZanEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		e.Fields[key] = str
		if key == "sIssueCode" {
			e.IssueCode = str
		}
	}
	return nil
}

func (s *Service) ShinyouZan(ctx context.Context, issueCodes []string) (*ShinyouZanResponse, error) {
	codes := normalizeList(issueCodes)
	if len(codes) == 0 {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "required"}
	}
	if len(codes) > maxIssueCodes {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "max 120"}
	}

	url, err := s.masterURL()
	if err != nil {
		return nil, err
	}

	req := ShinyouZanRequest{
		CommonParams:    model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:           clmGetShinyouZan,
		TargetIssueCode: strings.Join(codes, ","),
	}

	var resp ShinyouZanResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, url, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type HibuInfoRequest struct {
	model.CommonParams
	CLMID           string `json:"sCLMID"`
	TargetIssueCode string `json:"sTargetIssueCode"`
}

func (r *HibuInfoRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type HibuInfoResponse struct {
	model.CommonResponse
	Entries []HibuInfoEntry
	Fields  model.Attributes
}

type HibuInfoEntry struct {
	IssueCode string
	Fields    model.Attributes
}

func (r *HibuInfoResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	r.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		if key == "aCLMMfdsHibuInfo" {
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
	return nil
}

func (e *HibuInfoEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Fields = make(model.Attributes, len(raw))
	for key, value := range raw {
		str := jsonString(value)
		e.Fields[key] = str
		if key == "sIssueCode" {
			e.IssueCode = str
		}
	}
	return nil
}

func (s *Service) HibuInfo(ctx context.Context, issueCodes []string) (*HibuInfoResponse, error) {
	codes := normalizeList(issueCodes)
	if len(codes) == 0 {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "required"}
	}
	if len(codes) > maxIssueCodes {
		return nil, &terrors.ValidationError{Field: "issue_codes", Reason: "max 120"}
	}

	url, err := s.masterURL()
	if err != nil {
		return nil, err
	}

	req := HibuInfoRequest{
		CommonParams:    model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:           clmGetHibuInfo,
		TargetIssueCode: strings.Join(codes, ","),
	}

	var resp HibuInfoResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, url, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
