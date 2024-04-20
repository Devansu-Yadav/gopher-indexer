package core

import (
	"bufio"
	"io"
	"net"
)

func ReadFileFromServerAsBytes(conn net.Conn) ([]byte, error) {
	reader := bufio.NewReader(conn)
	var response []byte
	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, FetchErrorResponse(ResponseError, err)
		}

		response = append(response, buf[:n]...)
	}
	return response, nil
}
