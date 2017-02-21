package authentication

import "fmt"

type InvalidTokenError struct {
	Description string
}

func NewInvalidTokenError(description string) *InvalidTokenError {
	return &InvalidTokenError{Description: description}
}

func (e *InvalidTokenError) Error() string {
	return "Invalid token: " + e.Description
}

type ServerError struct {
	StatusCode  int
	ErrorCode   string
	Description string
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("Server error, Status code: %d; error code: %s, message: %s", s.StatusCode, s.ErrorCode, s.Description)
}

func NewServerError(statusCode int, errorCode string, description string) *ServerError {
	return &ServerError{
		StatusCode:  statusCode,
		ErrorCode:   errorCode,
		Description: description,
	}
}
