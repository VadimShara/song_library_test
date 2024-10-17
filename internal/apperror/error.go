package apperror

import (
	"encoding/json"
)

var(
	ErrNotFound = NewAppError(nil, "not found", "", "US-000003")
	DatabaseErr = NewAppError(nil, "database error", "method of bd return error", "US-000004")
	MarshalJSONErr = NewAppError(nil, "marshal json error", "", "US-000005")
	BadRequestErr = NewAppError(nil, "bad request error", "", "US-000006")
	DecodingErr = NewAppError(nil, "object decoding error", "", "US-000007")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code			 string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}
 
func NewAppError(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:				err,
		Message:			message,
		DeveloperMessage: 	developerMessage,
		Code:				code,
	}
}

func systemError(err error) *AppError {
	return NewAppError(err, "internal system error", err.Error(), "US-000000")
}