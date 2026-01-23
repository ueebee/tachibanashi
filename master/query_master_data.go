package master

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ueebee/tachibanashi/model"
)

const clmGetMasterData = "CLMMfdsGetMasterData"

type MasterDataRequest struct {
	model.CommonParams
	CLMID        string `json:"sCLMID"`
	TargetCLMID  string `json:"sTargetCLMID,omitempty"`
	TargetColumn string `json:"sTargetColumn,omitempty"`
}

func (r *MasterDataRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type MasterDataResponse struct {
	model.CommonResponse
	Data map[MasterType][]model.Attributes
}

func (r *MasterDataResponse) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make(map[string]string, len(raw))
	dataMap := make(map[MasterType][]model.Attributes)
	for key, value := range raw {
		switch key {
		case "p_no", "p_sd_date", "p_rv_date", "p_errno", "p_err", "sCLMID", "sResultCode", "sResultText":
			values[key] = jsonString(value)
			continue
		default:
			if key == "" {
				continue
			}
			var entries []model.Attributes
			if err := decodeAttributesList(value, &entries); err != nil {
				return err
			}
			dataMap[MasterType(key)] = entries
		}
	}
	applyCommonResponse(&r.CommonResponse, values)
	r.Data = dataMap
	return nil
}

// MasterData queries master records by CLMID and optional columns.
func (s *Service) MasterData(ctx context.Context, clmids []MasterType, columns []string) (*MasterDataResponse, error) {
	url, err := s.masterURL()
	if err != nil {
		return nil, err
	}

	clmidValues := make([]string, 0, len(clmids))
	for _, clmid := range clmids {
		clmidValues = append(clmidValues, string(clmid))
	}
	clmidValues = normalizeList(clmidValues)
	columnValues := normalizeList(columns)

	req := MasterDataRequest{
		CommonParams: model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:        clmGetMasterData,
	}
	if len(clmidValues) > 0 {
		req.TargetCLMID = strings.Join(clmidValues, ",")
	}
	if len(columnValues) > 0 {
		req.TargetColumn = strings.Join(columnValues, ",")
	}

	var resp MasterDataResponse
	if err := s.client.DoJSON(ctx, http.MethodGet, url, &req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
