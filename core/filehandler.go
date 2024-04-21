package core

import (
	"fmt"
)

func FetchFileAttrs(server, resource string) (int, error) {
	conn, err := CreateTCPConnection(server)
	if err != nil {
		return 0, FetchErrorResponse(ConnectionError, err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\r\n", resource)

	response, responseError := ReadFileFromServerAsBytes(conn)

	if responseError != nil {
		return 0, responseError
	}

	// Fetch file size
	fileSize := len(response)

	return fileSize, nil
}
