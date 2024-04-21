package core

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func ReadFileFromServerAsBytes(conn net.Conn, resource string) ([]byte, error) {
	fmt.Fprintf(conn, "%s\r\n", resource)

	reader := bufio.NewReader(conn)
	var response []byte
	buf := make([]byte, 1024)

	for {
		conn.SetReadDeadline(time.Now().Add(MaxResponseTimeOut))

		n, err := reader.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				logError(ReadFileTimeOut, resource)
				break
			} else {
				return nil, FetchErrorResponse(ResponseError, err)
			}
		}

		response = append(response, buf[:n]...)

		if len(response) >= MaxFileSize {
			logError(FileSizeExceeded, resource)
			break
		}
	}
	return response, nil
}
