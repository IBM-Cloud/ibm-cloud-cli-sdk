package authentication

import (
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
)

type InvalidTokenError struct {
	Description string
}

func NewInvalidTokenError(description string) *InvalidTokenError {
	return &InvalidTokenError{Description: description}
}

func (e *InvalidTokenError) Error() string {
	return T("Invalid token: ") + e.Description
}

type ServerError struct {
	StatusCode  int
	ErrorCode   string
	Description string
}

func (s *ServerError) Error() string {
	return T("Remote server error. Status code: {{.StatusCode}}, error code: {{.ErrorCode}}, message: {{.Message}}",
		map[string]interface{}{"StatusCode": s.StatusCode, "ErrorCode": s.ErrorCode, "Message": s.Description})
}

func NewServerError(statusCode int, errorCode string, description string) *ServerError {
	return &ServerError{
		StatusCode:  statusCode,
		ErrorCode:   errorCode,
		Description: description,
	}
}
