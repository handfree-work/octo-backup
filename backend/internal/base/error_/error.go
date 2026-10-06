package error_

import (
	"fmt"
	"handfree-work/octo-backup/internal/base/error_/code_"
	"handfree-work/octo-backup/internal/base/log_"
)

type CodedError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CodedErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewApiError(err error) *CodedError {
	if _, ok := err.(*CodedError); ok {
		return err.(*CodedError)
	}
	return NewCodeTextError(code_.ServerInternalError, err.Error())
}

func NewWrapError(tag string, err error) *CodedError {
	if _, ok := err.(*CodedError); ok {
		e := err.(*CodedError)
		e.Message = tag + ": " + e.Message
		return e
	}
	return NewCodeTextError(code_.ServerInternalError, tag+": "+err.Error())
}

func NewCodeTextError(code code_.ErrorCode, message string) *CodedError {
	err := &CodedError{Code: code.Code, Message: message}
	log_.Sugar.Error(err)
	return err
}

func NewFormatError(code code_.ErrorCode, message string) *CodedError {
	err := &CodedError{Code: code.Code, Message: message}
	log_.Sugar.Error(err)
	return err
}

func NewTextError(format string, args ...any) *CodedError {
	err := &CodedError{Code: code_.ServerInternalError.Code, Message: fmt.Sprintf(format, args...)}
	log_.Sugar.Error(err)
	return err
}

func NewRequestError(message string) *CodedError {
	err := &CodedError{Code: code_.RequestError.Code, Message: message}
	log_.Sugar.Error(err)
	return err
}

func LogError(tag string, err error) {
	log_.Sugar.Error(tag, err.Error())
}

// -----------error 接口实现------------------
func (e *CodedError) Error() string {
	return e.Message
}

func (*CodedError) LoggedError() {}

func (e *CodedError) Data() *CodedErrorResponse {
	return &CodedErrorResponse{
		Code:    e.Code,
		Message: e.Message,
	}
}
