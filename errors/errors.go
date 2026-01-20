package errors

import (
	"errors"
	"fmt"
)

var ErrNotImplemented = errors.New("tachibanashi: not implemented")

type APIError struct {
	Code    string
	Message string
	Detail  string
	Raw     []byte
}

func (e *APIError) Error() string {
	if e == nil {
		return "tachibanashi: api error"
	}
	if e.Code == "" {
		return "tachibanashi: api error"
	}
	return fmt.Sprintf("tachibanashi: api error code=%s message=%s", e.Code, e.Message)
}

type HTTPError struct {
	Status int
	Body   []byte
}

func (e *HTTPError) Error() string {
	if e == nil {
		return "tachibanashi: http error"
	}
	return fmt.Sprintf("tachibanashi: http status=%d", e.Status)
}

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "tachibanashi: validation error"
	}
	if e.Field == "" {
		return "tachibanashi: validation error"
	}
	return fmt.Sprintf("tachibanashi: validation error field=%s reason=%s", e.Field, e.Reason)
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Status >= 500 && httpErr.Status <= 599
	}
	return false
}
