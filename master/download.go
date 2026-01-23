package master

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/model"
)

const (
	clmEventDownload  = "CLMEventDownload"
	defaultJSONFormat = "5"

	updateTimeKey   = "sUpdateTime"
	updateNumberKey = "sUpdateNumber"
	deleteFlagKey   = "sDeleteFlag"
	deleteTimeKey   = "sDeleteTime"
)

type EventDownloadRequest struct {
	model.CommonParams
	CLMID string `json:"sCLMID"`
}

func (r *EventDownloadRequest) Params() *model.CommonParams {
	return &r.CommonParams
}

type DownloadHeader struct {
	PNo        string
	PSDDate    string
	PRVDate    string
	ErrNo      string
	Err        string
	ResultCode string
	ResultText string
}

type DownloadMessage struct {
	Type   MasterType
	Key    string
	Fields model.Attributes
	Meta   UpdateMeta
	Header DownloadHeader
}

type DownloadHandler func(message DownloadMessage) error

// Download streams master data until CLMEventDownloadComplete is received.
// Either store or handler must be provided.
func (s *Service) Download(ctx context.Context, store MasterStore, handler DownloadHandler) error {
	if store == nil && handler == nil {
		return &terrors.ValidationError{Field: "handler", Reason: "store or handler required"}
	}
	url, err := s.masterURL()
	if err != nil {
		return err
	}

	req := EventDownloadRequest{
		CommonParams: model.CommonParams{JsonOfmt: defaultJSONFormat},
		CLMID:        clmEventDownload,
	}

	resp, reader, err := s.client.DoStream(ctx, http.MethodGet, url, &req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(reader)
	for {
		var raw map[string]json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		message, err := parseDownloadMessage(raw)
		if err != nil {
			return err
		}

		if message.Type == MasterEventDownloadComplete {
			if handler != nil {
				if err := handler(message); err != nil {
					return err
				}
			}
			return nil
		}

		if store != nil && message.Key != "" {
			store.Upsert(message.Type, message.Key, message.Fields, message.Meta)
		}
		if handler != nil {
			if err := handler(message); err != nil {
				return err
			}
		}
	}
}

var reservedKeys = map[string]struct{}{
	"p_no":          {},
	"p_sd_date":     {},
	"p_rv_date":     {},
	"p_errno":       {},
	"p_err":         {},
	"sCLMID":        {},
	"sResultCode":   {},
	"sResultText":   {},
	updateTimeKey:   {},
	updateNumberKey: {},
	deleteFlagKey:   {},
	deleteTimeKey:   {},
}

func parseDownloadMessage(raw map[string]json.RawMessage) (DownloadMessage, error) {
	typ := MasterType(jsonString(raw["sCLMID"]))
	if typ == "" {
		return DownloadMessage{}, errors.New("tachibanashi: master download missing sCLMID")
	}

	header := DownloadHeader{
		PNo:        jsonString(raw["p_no"]),
		PSDDate:    jsonString(raw["p_sd_date"]),
		PRVDate:    jsonString(raw["p_rv_date"]),
		ErrNo:      jsonString(raw["p_errno"]),
		Err:        jsonString(raw["p_err"]),
		ResultCode: jsonString(raw["sResultCode"]),
		ResultText: jsonString(raw["sResultText"]),
	}

	if header.ResultCode != "" && header.ResultCode != "0" {
		return DownloadMessage{}, &terrors.APIError{
			Code:    header.ResultCode,
			Message: header.ResultText,
			Detail:  header.Err,
		}
	}

	meta := UpdateMeta{
		Serial:    parseInt64(jsonString(raw[updateNumberKey])),
		UpdatedAt: normalizeTimestamp(jsonString(raw[updateTimeKey])),
		Deleted:   parseDeleteFlag(jsonString(raw[deleteFlagKey])),
	}

	fields := make(model.Attributes, len(raw))
	for key, value := range raw {
		if _, reserved := reservedKeys[key]; reserved {
			continue
		}
		fields[key] = jsonString(value)
	}

	key, _ := MasterKey(typ, fields)

	return DownloadMessage{
		Type:   typ,
		Key:    key,
		Fields: fields,
		Meta:   meta,
		Header: header,
	}, nil
}

func parseInt64(value string) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func normalizeTimestamp(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	for i := 0; i < len(value); i++ {
		if value[i] != '0' {
			return value
		}
	}
	return ""
}

func parseDeleteFlag(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if value == "1" {
		return true
	}
	return strings.EqualFold(value, "true")
}
