package errors

import "fmt"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrNotFound     = &AppError{Code: "NOT_FOUND", Message: "Resource not found"}
	ErrInvalidInput = &AppError{Code: "INVALID_INPUT", Message: "Invalid input data"}
	ErrInternal     = &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error"}
)
