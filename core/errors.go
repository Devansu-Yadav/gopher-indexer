package core

import (
	"fmt"
)

type ErrorType int

const (
	ConnectionError ErrorType = iota
	ResponseError
	FileSizeExceeded
	ReadFileTimeOut
	ServerResponseTimeOut
	ExternalServerConnection
	InvalidReference
)

type CustomError struct {
	Type     ErrorType
	Resource string
	Err      error
}

func (e *CustomError) Error() string {
	switch e.Type {
	case ConnectionError:
		return fmt.Sprintf("Error: Failed to connect to server - %s", e.Resource)
	case ResponseError:
		return fmt.Sprintf("Error: Failed to read response - %s", e.Resource)
	case FileSizeExceeded:
		return fmt.Sprintf("Error: File size exceeds 5MB limit - %s", e.Resource)
	case ReadFileTimeOut:
		return fmt.Sprintf("Error: File read timed out - %s", e.Resource)
	case ServerResponseTimeOut:
		return fmt.Sprintf("Error: Server response timed out - %s", e.Resource)
	case ExternalServerConnection:
		return fmt.Sprintf("Error: External server connection timed out - %s", e.Resource)
	case InvalidReference:
		return fmt.Sprintf("Error: Invalid Reference to resource - %s", e.Resource)
	default:
		return fmt.Sprintf("%v for resource - %s", e.Err, e.Resource)
	}
}

func FetchErrorResponse(errorType ErrorType, err error, resource string) error {
	customErr := CustomError{Type: errorType, Err: err, Resource: resource}
	customErrMsg := customErr.Error()
	return fmt.Errorf(customErrMsg)
}
