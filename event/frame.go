package event

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"strings"

	terrors "github.com/ueebee/tachibanashi/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	delimiterItem  = '\x01'
	delimiterKey   = '\x02'
	delimiterValue = '\x03'
)

type Frame struct {
	Raw     string
	Fields  map[string][]string
	No      int64
	Date    string
	Command Command
}

func (f Frame) Kind() string {
	return string(f.Command)
}

func (f Frame) Value(key string) string {
	values := f.Fields[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (f Frame) Values(key string) []string {
	values := f.Fields[key]
	if len(values) == 0 {
		return nil
	}
	out := make([]string, len(values))
	copy(out, values)
	return out
}

func DecodeFrame(data []byte) (Frame, error) {
	raw := strings.TrimRight(string(data), "\r\n")
	fields, err := parseFields(raw)
	if err != nil {
		return Frame{}, err
	}

	frame := Frame{
		Raw:    raw,
		Fields: fields,
	}

	cmd := strings.TrimSpace(frame.Value("p_cmd"))
	if cmd == "" {
		return Frame{}, &terrors.ValidationError{Field: "p_cmd", Reason: "required"}
	}
	frame.Command = Command(strings.ToUpper(cmd))
	frame.No = parseInt64(frame.Value("p_no"))
	frame.Date = frame.Value("p_date")

	return frame, nil
}

func DecodeEvent(data []byte) (Event, error) {
	frame, err := DecodeFrame(data)
	if err != nil {
		return nil, err
	}
	if _, ok := allowedCommands[frame.Command]; !ok {
		return nil, &terrors.ValidationError{Field: "p_cmd", Reason: "unsupported"}
	}
	if err := decodeBase64Fields(&frame); err != nil {
		return nil, err
	}

	switch frame.Command {
	case CommandST:
		return parseST(frame)
	case CommandKP:
		return parseKP(frame), nil
	case CommandFD:
		return parseFD(frame)
	case CommandEC:
		return parseEC(frame)
	default:
		return frame, nil
	}
}

func parseFields(raw string) (map[string][]string, error) {
	if raw == "" {
		return map[string][]string{}, nil
	}
	items := strings.Split(raw, string(delimiterItem))
	fields := make(map[string][]string, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		key, value, ok := strings.Cut(item, string(delimiterKey))
		if !ok {
			return nil, errors.New("tachibanashi: event frame missing delimiter")
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		values := strings.Split(value, string(delimiterValue))
		fields[key] = append(fields[key], values...)
	}
	return fields, nil
}

var base64Fields = map[Command][]string{
	CommandEC: {"p_IN"},
	CommandNS: {"p_HDL", "p_TX"},
}

func decodeBase64Fields(frame *Frame) error {
	keys := base64Fields[frame.Command]
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		values, ok := frame.Fields[key]
		if !ok {
			continue
		}
		for i, value := range values {
			decoded, err := decodeBase64ShiftJIS(value)
			if err != nil {
				return err
			}
			values[i] = decoded
		}
		frame.Fields[key] = values
	}
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
