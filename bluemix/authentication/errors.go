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

// RefreshTokenExpiryError is an error when provided refresh token expires. This error normally requires
// the client to re-login.
type RefreshTokenExpiryError struct {
	Description string
}

func (e *RefreshTokenExpiryError) Error() string {
	return e.Description
}

// NewRefreshTokenExpiryError creates a RefreshTokenExpiryError
func NewRefreshTokenExpiryError(description string) *RefreshTokenExpiryError {
	return &RefreshTokenExpiryError{Description: description}
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

type InvalidGrantTypeError struct {
	Description string
}

func NewInvalidGrantTypeError(description string) *InvalidGrantTypeError {
	return &InvalidGrantTypeError{Description: description}
}

func (e *InvalidGrantTypeError) Error() string {
	return T("Invalid grant type: ") + e.Description
}

type ExternalAuthenticationError struct {
	ErrorCode    string
	ErrorMessage string
}

func (e ExternalAuthenticationError) Error() string {
	return T("External authentication failed. Error code: {{.ErrorCode}}, message: {{.Message}}",
		map[string]interface{}{"ErrorCode": e.ErrorCode, "Message": e.ErrorMessage})
}

type SessionInactiveError struct {
	Description string
}

func NewSessionInactiveError(description string) *SessionInactiveError {
	return &SessionInactiveError{Description: description}
}

func (e *SessionInactiveError) Error() string {
	return T("Session inactive: ") + e.Description
}
