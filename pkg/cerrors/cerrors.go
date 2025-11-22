package cerrors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ErrorType string

const (
	UserError         ErrorType = "validation_error"
	InternalError     ErrorType = "internal_error"
	BusinessError     ErrorType = "business_error"
	UnauthorizedError ErrorType = "unauthorized_error"
)

type Error struct {
	Type     ErrorType
	Err      error
	messages []string
	params   map[string][]any
}

func New(errorType ErrorType, err error) Error {
	return Error{
		Type: errorType,
		Err:  err,
		params: map[string][]any{
			"error_type": {errorType},
		},
	}
}

func As(err error) (Error, bool) {
	var target Error
	if !errors.As(err, &target) {
		return target, false
	}
	return target, true
}

func AsUserError(err error) (Error, bool) {
	return asTError(err, UserError)
}

func AsBusinessError(err error) (Error, bool) {
	return asTError(err, BusinessError)
}

func AsUnauthorizedErrorError(err error) (Error, bool) {
	return asTError(err, UnauthorizedError)
}

func AsInternalError(err error) (Error, bool) {
	return asTError(err, InternalError)
}

func asTError(err error, t ErrorType) (Error, bool) {
	var target Error
	if !errors.As(err, &target) {
		return target, false
	}

	if target.Type != t {
		return target, false
	}

	return target, true
}

func NewUserError(err error) Error {
	return New(UserError, err)
}

func NewBusinessError(err error) Error {
	return New(BusinessError, err)
}

func NewUnauthorizedError(err error) Error {
	return New(UnauthorizedError, err)
}

func NewInternalError(err error) Error {
	return New(InternalError, err)
}

func (v Error) WithParam(key string, value any) Error {
	v.params[key] = append(v.params[key], value)
	return v
}

func (v Error) WithParams(params map[string]any) Error {
	for key, value := range params {
		v.params[key] = append(v.params[key], value)
	}
	return v
}

func (v Error) WithMessage(msg string) Error {
	v.messages = append(v.messages, msg)
	return v
}

func (v Error) Message() string {
	return strings.Join(v.messages, "\n")
}

func (v Error) Params() map[string]any {
	params := make(map[string]any, len(v.params))
	for key, values := range v.params {
		if len(values) == 1 {
			params[key] = values[0]
			continue
		}

		for i, value := range values {
			params[key+"_"+strconv.Itoa(i)] = value
		}
	}
	return params
}

func (v Error) Error() string {
	if v.Err != nil {
		return fmt.Sprintf("[%s]: %s", v.Type, v.Err.Error())
	}
	return "nil"
}

func (v Error) Unwrap() error {
	return v.Err
}

func (v Error) GetType() ErrorType {
	return v.Type
}
