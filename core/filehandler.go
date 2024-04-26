package core

// Connects to the server, fetches a file, and returns its attributes.
// Takes the server address, resource string, and file type as inputs.
// Returns the file content as a byte slice, the file size, any error that occurred,
// and a boolean indicating if the file is malformed. If an error occurs while reading the
// file, it returns the error and the malformed status.
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
