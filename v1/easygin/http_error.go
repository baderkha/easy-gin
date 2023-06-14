package easygin

import (
	"fmt"
)

type HTTPError struct {
	Err        error
	StatusCode int
}

func (r *HTTPError) Error() string {
	return r.Err.Error()
}

func NewHTTPError(sCode int, err error) *HTTPError {
	return &HTTPError{
		Err:        err,
		StatusCode: sCode,
	}
}

func NewComposedHTTPError(sCode int, tmpl string) func(args ...any) *HTTPError {
	return func(args ...any) *HTTPError {
		return NewHTTPError(sCode, fmt.Errorf(tmpl, args...))
	}
}
