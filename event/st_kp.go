package event

import (
	"strings"

	terrors "github.com/ueebee/tachibanashi/errors"
)

type ST struct {
	Frame
	ErrNo string
	Err   string
}

type KP struct {
	Frame
}

func parseST(frame Frame) (ST, error) {
	errNo := strings.TrimSpace(frame.Value("p_errno"))
	if errNo == "" {
		return ST{}, &terrors.ValidationError{Field: "p_errno", Reason: "required"}
	}
	return ST{
		Frame: frame,
		ErrNo: errNo,
		Err:   frame.Value("p_err"),
	}, nil
}

func parseKP(frame Frame) KP {
	return KP{Frame: frame}
}
