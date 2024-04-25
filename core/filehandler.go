package core

func FetchFileAttrs(server, resource string) (int, error) {
	conn, err := CreateTCPConnection(server)
	logRequest(server + resource)

	if err != nil {
		return 0, FetchErrorResponse(ConnectionError, err, server+resource)
	}
	defer conn.Close()

	response, responseError := ReadFileFromServerAsBytes(conn, resource)

	if responseError != nil {
		return 0, responseError
	}

	// Fetch file size
	fileSize := len(response)

	return fileSize, nil
}
