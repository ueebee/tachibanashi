package master

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"

	"github.com/ueebee/tachibanashi/model"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const maxIssueCodes = 120

func normalizeList(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			item := strings.TrimSpace(part)
			if item == "" {
				continue
			}
			out = append(out, item)
		}
	}
	return out
}

func applyCommonResponse(target *model.CommonResponse, values map[string]string) {
	if target == nil {
		return
	}
	target.PNo = values["p_no"]
	target.PSDDate = values["p_sd_date"]
	target.PRVDate = values["p_rv_date"]
	target.PErrNo = values["p_errno"]
	target.PErr = values["p_err"]
	target.CLMID = values["sCLMID"]
	target.ResultCode = values["sResultCode"]
	target.ResultText = values["sResultText"]
}

func decodeList(raw json.RawMessage, out any) error {
	if len(raw) == 0 {
		return nil
	}
	trimmed := bytes.TrimSpace(raw)
	if bytes.Equal(trimmed, []byte("null")) {
		return nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if strings.TrimSpace(s) == "" {
			return nil
		}
		return json.Unmarshal([]byte(s), out)
	}
	return json.Unmarshal(raw, out)
}

func decodeAttributesList(raw json.RawMessage, out *[]model.Attributes) error {
	var items []map[string]json.RawMessage
	if err := decodeList(raw, &items); err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	attrs := make([]model.Attributes, 0, len(items))
	for _, item := range items {
		entry := make(model.Attributes, len(item))
		for key, value := range item {
			entry[key] = jsonString(value)
		}
		attrs = append(attrs, entry)
	}
	*out = attrs
	return nil
}

func decodeBase64ShiftJIS(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return value, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(value)
		if err != nil {
			return "", err
		}
	}
	return decodeShiftJIS(decoded)
}

func decodeShiftJIS(data []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), japanese.ShiftJIS.NewDecoder())
	out, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
