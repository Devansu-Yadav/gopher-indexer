package core

func ConnectToServer(server string) (string, error) {
	conn, err := CreateTCPConnection(server)

	if err != nil {
		return "", FetchErrorResponse(ConnectionError, err)
	}
	defer conn.Close()

	response, responseErr := FetchServerResponse(conn, "")

	if responseErr != nil {
		return "", responseErr
	}
	return response, nil
}
