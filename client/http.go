package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	terrors "github.com/ueebee/tachibanashi/errors"
	"github.com/ueebee/tachibanashi/event"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func (c *Client) DoJSON(ctx context.Context, method, path string, req, resp any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if method == "" {
		method = http.MethodGet
	}
	if path == "" {
		return errors.New("tachibanashi: path is empty")
	}

	fullURL, err := c.resolveURL(path)
	if err != nil {
		return err
	}

	payload, err := c.preparePayload(req)
	if err != nil {
		return err
	}

	var body io.Reader
	if method == http.MethodGet || method == http.MethodHead {
		if len(payload) > 0 {
			fullURL, err = appendJSONQuery(fullURL, payload)
			if err != nil {
				return err
			}
		}
	} else if len(payload) > 0 {
		body = bytes.NewReader(payload)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return err
	}

	if c.cfg.UserAgent != "" {
		httpReq.Header.Set("User-Agent", c.cfg.UserAgent)
	}
	httpReq.Header.Set("Accept", "application/json")
	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	respBody = decodeResponseBody(httpResp, respBody)

	if httpResp.StatusCode >= http.StatusBadRequest {
		return &terrors.HTTPError{Status: httpResp.StatusCode, Body: respBody}
	}

	if apiErr := parseAPIError(respBody); apiErr != nil {
		return apiErr
	}

	if resp == nil || len(respBody) == 0 {
		return nil
	}
	if raw, ok := resp.(*[]byte); ok {
		*raw = respBody
		return nil
	}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return err
	}
	return nil
}

func (c *Client) DialEvent(ctx context.Context) (event.Conn, error) {
	return nil, terrors.ErrNotImplemented
}

func (c *Client) resolveURL(path string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}
	if c.cfg.BaseURL == "" {
		return "", errors.New("tachibanashi: base URL is empty")
	}
	base, err := url.Parse(c.cfg.BaseURL)
	if err != nil {
		return "", err
	}
	ref, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(ref).String(), nil
}

func (c *Client) preparePayload(req any) ([]byte, error) {
	if req == nil {
		return nil, nil
	}

	req = c.applyCommonParams(req)

	switch v := req.(type) {
	case []byte:
		return v, nil
	case json.RawMessage:
		return []byte(v), nil
	case string:
		return []byte(v), nil
	default:
		return json.Marshal(v)
	}
}

func (c *Client) applyCommonParams(req any) any {
	now := time.Now()

	switch v := req.(type) {
	case CommonParamsCarrier:
		params := v.Params()
		if params == nil {
			return req
		}
		if params.PNo == "" {
			params.PNo = strconv.FormatInt(c.token.Next(), 10)
		}
		if params.PSDDate == "" {
			params.PSDDate = formatTimestamp(now)
		}
		if params.JsonOfmt == "" {
			params.JsonOfmt = "5"
		}
		return req
	case map[string]any:
		if _, ok := v["p_no"]; !ok {
			v["p_no"] = strconv.FormatInt(c.token.Next(), 10)
		}
		if _, ok := v["p_sd_date"]; !ok {
			v["p_sd_date"] = formatTimestamp(now)
		}
		if _, ok := v["sJsonOfmt"]; !ok {
			v["sJsonOfmt"] = "5"
		}
		return v
	default:
		return req
	}
}

func appendJSONQuery(rawURL string, payload []byte) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	u.RawQuery = url.QueryEscape(string(payload))
	return u.String(), nil
}

func formatTimestamp(t time.Time) string {
	return t.Format("2006.01.02-15:04:05.000")
}

func parseAPIError(body []byte) *terrors.APIError {
	if len(body) == 0 {
		return nil
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}

	pErrNo := jsonString(raw["p_errno"])
	pErr := jsonString(raw["p_err"])
	resultCode := jsonString(raw["sResultCode"])
	resultText := jsonString(raw["sResultText"])

	if pErrNo != "" && pErrNo != "0" {
		return &terrors.APIError{
			Code:    pErrNo,
			Message: pErr,
			Detail:  resultText,
			Raw:     body,
		}
	}
	if resultCode != "" && resultCode != "0" {
		return &terrors.APIError{
			Code:    resultCode,
			Message: resultText,
			Detail:  pErr,
			Raw:     body,
		}
	}
	return nil
}

func decodeResponseBody(resp *http.Response, body []byte) []byte {
	if len(body) == 0 {
		return body
	}
	charset := responseCharset(resp)
	if charset == "" {
		if utf8.Valid(body) {
			return body
		}
		if decoded, err := decodeShiftJIS(body); err == nil {
			return decoded
		}
		return body
	}
	if isUTF8Charset(charset) {
		return body
	}
	if isShiftJISCharset(charset) {
		if decoded, err := decodeShiftJIS(body); err == nil {
			return decoded
		}
	}
	return body
}

func responseCharset(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	value := resp.Header.Get("Content-Type")
	if value == "" {
		return ""
	}
	_, params, err := mime.ParseMediaType(value)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(params["charset"]))
}

func isUTF8Charset(charset string) bool {
	switch charset {
	case "utf-8", "utf8":
		return true
	default:
		return false
	}
}

func isShiftJISCharset(charset string) bool {
	switch charset {
	case "shift_jis", "shift-jis", "sjis", "cp932", "windows-31j":
		return true
	default:
		return false
	}
}

func decodeShiftJIS(body []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(body), japanese.ShiftJIS.NewDecoder())
	return io.ReadAll(reader)
}

func jsonString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
