package core

func FetchFileAttrs(server, resource string, fileType string) ([]byte, int, error, bool) {
	conn, err := CreateTCPConnection(server)
	logRequest(server + resource)

	if err != nil {
		return nil, 0, FetchErrorResponse(ConnectionError, err, server+resource), false
	}
	defer conn.Close()

	response, responseError, isMalformed := ReadFileFromServerAsBytes(conn, resource, fileType)

	if responseError != nil {
		return nil, 0, responseError, isMalformed
	}

	// Fetch file size
	fileSize := len(response)

	return response, fileSize, nil, isMalformed
}
