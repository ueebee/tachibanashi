package event

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	terrors "github.com/ueebee/tachibanashi/errors"
)

type Command string

const (
	CommandST Command = "ST"
	CommandKP Command = "KP"
	CommandFD Command = "FD"
	CommandEC Command = "EC"
	CommandNS Command = "NS"
	CommandSS Command = "SS"
	CommandUS Command = "US"
)

const (
	defaultBoardNo = 1000
	maxSymbols     = 120
	maxRows        = 120
	maxRowSmall    = 20
)

var defaultCommands = []Command{
	CommandST,
	CommandKP,
	CommandEC,
	CommandSS,
	CommandUS,
}

var allowedCommands = map[Command]struct{}{
	CommandST: {},
	CommandKP: {},
	CommandFD: {},
	CommandEC: {},
	CommandNS: {},
	CommandSS: {},
	CommandUS: {},
}

type Params struct {
	RID         int
	BoardNo     int
	Rows        []int
	IssueCodes  []string
	MarketCodes []string
	Eno         int64
	Cmds        []Command
}

func (p Params) Validate() error {
	normalized := normalizeParams(p)
	return validateParams(normalized)
}

func BuildWSURL(base string, params Params) (string, error) {
	if strings.TrimSpace(base) == "" {
		return "", errors.New("tachibanashi: event base URL is empty")
	}
	normalized := normalizeParams(params)
	if err := validateParams(normalized); err != nil {
		return "", err
	}

	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	// Build query manually to keep comma-separated values intact.
	commaKeys := map[string]struct{}{
		"p_evt_cmd":    {},
		"p_issue_code": {},
		"p_gyou_no":    {},
		"p_mkt_code":   {},
	}
	paramsList := make([]string, 0, 8)
	add := func(key, value string) {
		if strings.TrimSpace(value) == "" {
			return
		}
		encoded := url.QueryEscape(value)
		if _, ok := commaKeys[key]; ok {
			encoded = strings.ReplaceAll(encoded, "%2C", ",")
		}
		paramsList = append(paramsList, key+"="+encoded)
	}

	add("p_rid", strconv.Itoa(normalized.RID))
	add("p_board_no", strconv.Itoa(normalized.BoardNo))
	if len(normalized.Rows) > 0 {
		add("p_gyou_no", joinInts(normalized.Rows))
	}
	if len(normalized.IssueCodes) > 0 {
		add("p_issue_code", strings.Join(normalized.IssueCodes, ","))
	}
	if len(normalized.MarketCodes) > 0 {
		add("p_mkt_code", strings.Join(normalized.MarketCodes, ","))
	}
	add("p_eno", strconv.FormatInt(normalized.Eno, 10))
	add("p_evt_cmd", joinCommands(normalized.Cmds))

	u.RawQuery = strings.Join(paramsList, "&")
	return u.String(), nil
}

func normalizeParams(p Params) Params {
	out := p
	out.Cmds = normalizeCommands(out.Cmds)
	if len(out.Cmds) == 0 {
		out.Cmds = append([]Command(nil), defaultCommands...)
	}

	out.Rows = normalizeIntList(out.Rows)
	out.IssueCodes = normalizeCodeList(out.IssueCodes)
	out.MarketCodes = normalizeCodeList(out.MarketCodes)

	if out.RID == 0 && (len(out.Rows) > 0 || len(out.IssueCodes) > 0 || len(out.MarketCodes) > 0) {
		out.RID = 22
	}
	if out.RID == 0 && out.BoardNo == 0 {
		out.BoardNo = defaultBoardNo
	}
	return out
}

func normalizeCommands(cmds []Command) []Command {
	if len(cmds) == 0 {
		return nil
	}
	seen := make(map[Command]struct{}, len(cmds))
	out := make([]Command, 0, len(cmds))
	for _, cmd := range cmds {
		normalized := Command(strings.ToUpper(strings.TrimSpace(string(cmd))))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	return out
}

func normalizeCodeList(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		item := strings.TrimSpace(value)
		if item == "" {
			continue
		}
		out = append(out, item)
	}
	return out
}

func normalizeIntList(values []int) []int {
	if len(values) == 0 {
		return nil
	}
	out := make([]int, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		out = append(out, value)
	}
	return out
}

func validateParams(p Params) error {
	if p.Eno < 0 {
		return &terrors.ValidationError{Field: "eno", Reason: "must be >= 0"}
	}
	if err := validateCommands(p.Cmds); err != nil {
		return err
	}
	if err := validateBoardCombo(p); err != nil {
		return err
	}
	return nil
}

func validateCommands(cmds []Command) error {
	if len(cmds) == 0 {
		return &terrors.ValidationError{Field: "evt_cmd", Reason: "required"}
	}
	for _, cmd := range cmds {
		if _, ok := allowedCommands[cmd]; !ok {
			return &terrors.ValidationError{Field: "evt_cmd", Reason: "unsupported"}
		}
	}
	return nil
}

func validateBoardCombo(p Params) error {
	if hasCommand(p.Cmds, CommandFD) && p.RID == 0 {
		return &terrors.ValidationError{Field: "rid", Reason: "fd requires price board"}
	}

	switch p.RID {
	case 0:
		if p.BoardNo != defaultBoardNo {
			return &terrors.ValidationError{Field: "board_no", Reason: "must be 1000 for rid=0"}
		}
		if len(p.Rows) > 0 || len(p.IssueCodes) > 0 || len(p.MarketCodes) > 0 {
			return &terrors.ValidationError{Field: "board_no", Reason: "rows/codes not allowed for rid=0"}
		}
	case 10, 11:
		if p.BoardNo < 1 || p.BoardNo > 10 {
			return &terrors.ValidationError{Field: "board_no", Reason: "must be 1-10"}
		}
		if len(p.Rows) > 0 || len(p.IssueCodes) > 0 || len(p.MarketCodes) > 0 {
			return &terrors.ValidationError{Field: "gyou_no", Reason: "rows/codes not allowed"}
		}
	case 12:
		if p.BoardNo != 120 {
			return &terrors.ValidationError{Field: "board_no", Reason: "must be 120"}
		}
		if len(p.Rows) > 0 || len(p.IssueCodes) > 0 || len(p.MarketCodes) > 0 {
			return &terrors.ValidationError{Field: "gyou_no", Reason: "rows/codes not allowed"}
		}
	case 13:
		if p.BoardNo != 20 && p.BoardNo != 30 {
			return &terrors.ValidationError{Field: "board_no", Reason: "must be 20 or 30"}
		}
		if err := validateRowList(p.Rows, 1, maxRows); err != nil {
			return err
		}
		if err := validateSymbolLists(p.IssueCodes, p.MarketCodes, len(p.Rows)); err != nil {
			return err
		}
	case 20:
		if p.BoardNo < 1 || p.BoardNo > 10 {
			return &terrors.ValidationError{Field: "board_no", Reason: "must be 1-10"}
		}
		if err := validateRowList(p.Rows, 1, maxRowSmall); err != nil {
			return err
		}
		if len(p.IssueCodes) > 0 || len(p.MarketCodes) > 0 {
			return &terrors.ValidationError{Field: "issue_code", Reason: "not allowed for rid=20"}
		}
	case 21:
		if p.BoardNo != 0 {
			return &terrors.ValidationError{Field: "board_no", Reason: "must be 0"}
		}
		if len(p.Rows) != 1 {
			return &terrors.ValidationError{Field: "gyou_no", Reason: "must be a single row"}
		}
		if err := validateRowList(p.Rows, 1, maxRowSmall); err != nil {
			return err
		}
		if len(p.IssueCodes) != 1 || len(p.MarketCodes) != 1 {
			return &terrors.ValidationError{Field: "issue_code", Reason: "must be 1"}
		}
	case 22:
		if p.BoardNo != defaultBoardNo {
			return &terrors.ValidationError{Field: "board_no", Reason: "must be 1000 for rid=22"}
		}
		if err := validateRowList(p.Rows, 1, maxRows); err != nil {
			return err
		}
		if err := validateSymbolLists(p.IssueCodes, p.MarketCodes, len(p.Rows)); err != nil {
			return err
		}
	default:
		return &terrors.ValidationError{Field: "rid", Reason: "unsupported"}
	}

	return nil
}

func validateRowList(rows []int, min, max int) error {
	if len(rows) == 0 {
		return &terrors.ValidationError{Field: "gyou_no", Reason: "required"}
	}
	if len(rows) > max {
		return &terrors.ValidationError{Field: "gyou_no", Reason: "max 120"}
	}
	for _, row := range rows {
		if row < min || row > max {
			return &terrors.ValidationError{Field: "gyou_no", Reason: "out of range"}
		}
	}
	return nil
}

func validateSymbolLists(issues, markets []string, rows int) error {
	if len(issues) == 0 {
		return &terrors.ValidationError{Field: "issue_code", Reason: "required"}
	}
	if len(issues) > maxSymbols {
		return &terrors.ValidationError{Field: "issue_code", Reason: "max 120"}
	}
	if len(markets) == 0 {
		return &terrors.ValidationError{Field: "mkt_code", Reason: "required"}
	}
	if len(markets) > maxSymbols {
		return &terrors.ValidationError{Field: "mkt_code", Reason: "max 120"}
	}
	if len(issues) != len(markets) {
		return &terrors.ValidationError{Field: "mkt_code", Reason: "count mismatch"}
	}
	if rows > 0 && len(issues) != rows {
		return &terrors.ValidationError{Field: "issue_code", Reason: "count mismatch"}
	}
	return nil
}

func joinCommands(cmds []Command) string {
	if len(cmds) == 0 {
		return ""
	}
	out := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		out = append(out, string(cmd))
	}
	return strings.Join(out, ",")
}

func joinInts(values []int) string {
	if len(values) == 0 {
		return ""
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, strconv.Itoa(value))
	}
	return strings.Join(out, ",")
}

func hasCommand(cmds []Command, target Command) bool {
	for _, cmd := range cmds {
		if cmd == target {
			return true
		}
	}
	return false
}
