package core

import "net"

func ConnectToServer(server string) (string, error) {
	conn, err := CreateTCPConnection(server)

	if err != nil {
		return "", FetchErrorResponse(ConnectionError, err, server)
	}
	defer conn.Close()

	response, responseErr := FetchServerResponse(conn, "")

	if responseErr != nil {
		return "", responseErr
	}
	return response, nil
}

func ConnectToExternalServer(server string) error {
	conn, err := CreateTCPConnectionWithTimeOut(server)
	logRequest(server)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return FetchErrorResponse(ExternalServerConnection, err, server)
		}
		return FetchErrorResponse(ConnectionError, err, server)
	}
	defer conn.Close()

	return nil
}
