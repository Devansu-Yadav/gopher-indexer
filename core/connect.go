package core

import (
	"fmt"
)

func ConnectToServer(server string) (string, error) {
	conn, err := CreateTCPConnection(server)

	if err != nil {
		return "", FetchErrorResponse(ConnectionError, err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "\r\n")

	response, responseErr := FetchServerResponse(conn)

	if responseErr != nil {
		return "", responseErr
	}
	return response, nil
}
