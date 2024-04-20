package core

import (
	"fmt"
)

type ErrorType int

const (
	ConnectionError ErrorType = iota
	ResponseError
)

type CustomError struct {
	Type ErrorType
	Err  error
}

func (e *CustomError) Error() string {
	switch e.Type {
	case ConnectionError:
		return fmt.Sprintf("Failed to connect to server: %v", e.Err)
	case ResponseError:
		return fmt.Sprintf("Failed to read response: %v", e.Err)
	default:
		return fmt.Sprintf("Unknown error: %v", e.Err)
	}
}

func FetchErrorResponse(errorType ErrorType, err error) error {
	return CustomError{Type: errorType, Err: err}.Err
}
