package core

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func ReadFileFromServerAsBytes(conn net.Conn) ([]byte, error) {
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
				fmt.Fprintln(os.Stderr, "Read file time out - server is too slow!")
				break
			} else {
				return nil, FetchErrorResponse(ResponseError, err)
			}
		}

		response = append(response, buf[:n]...)

		if len(response) >= MaxFileSize {
			fmt.Fprintln(os.Stderr, "File too long! This client can handle max file sizes of up to 5 MBs.")
			break
		}
	}
	return response, nil
}
