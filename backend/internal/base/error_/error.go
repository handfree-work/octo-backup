package error_

import (
	"handfree-work/web-restic/internal/base/error_/code_"
	"handfree-work/web-restic/internal/base/log_"
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

func NewTextError(message string) *CodedError {
	err := &CodedError{Code: code_.ServerInternalError.Code, Message: message}
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

func (e *CodedError) Data() *CodedErrorResponse {
	return &CodedErrorResponse{
		Code:    e.Code,
		Message: e.Message,
	}
}
