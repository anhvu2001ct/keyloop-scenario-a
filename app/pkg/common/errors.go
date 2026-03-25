package common

import (
	"fmt"
)

type SafeError struct {
	// short message describe the error
	Short string

	// user-facing message
	UserMsg string

	Cause error

	// additional context
	Metadata map[string]any
}

func (e *SafeError) Error() string {
	return e.UserMsg
}

func (e *SafeError) Unwrap() error {
	return e.Cause
}

func (e *SafeError) LogString() string {
	return fmt.Sprintf("Short: %s | Msg: %s | Cause: %v | Meta: %+v",
		e.Short, e.UserMsg, e.Cause, e.Metadata)
}

type appErrorType int

const (
	AppErrorUnknown appErrorType = iota
	AppErrorRecordNotFound
	AppErrorValidationFailed
	AppErrorRecordExisted
)

type AppError struct {
	Message string
	Type    appErrorType
	Cause   error
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return ""
}

func (e *AppError) Unwrap() error {
	return e.Cause
}
