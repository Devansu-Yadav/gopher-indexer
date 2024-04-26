package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func ReadFileFromServerAsBytes(conn net.Conn, resource string, fileType string) ([]byte, error, bool) {
	fmt.Fprintf(conn, "%s\r\n", resource)

	reader := bufio.NewReader(conn)
	var response []byte
	buf := make([]byte, 1024)
	var isMalformed bool = false

	for {
		conn.SetReadDeadline(time.Now().Add(MaxResponseTimeOut))

		n, err := reader.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return nil, FetchErrorResponse(ReadFileTimeOut, err, resource), isMalformed
			} else {
				return nil, FetchErrorResponse(ResponseError, err, resource), isMalformed
			}
		}

		response = append(response, buf[:n]...)

		if len(response) >= MaxFileSize {
			return nil, FetchErrorResponse(FileSizeExceeded, errors.New(""), resource), isMalformed
		}
	}

	if fileType == TextFile {
		responseStr := string(response)

		// Check if the file is properly terminated
		if !strings.HasSuffix(responseStr, ".\r\n") {
			isMalformed = true
		}

		responseStr = strings.TrimSuffix(responseStr, ".\r\n")
		response = []byte(responseStr)
	}
	return response, nil, isMalformed
}
