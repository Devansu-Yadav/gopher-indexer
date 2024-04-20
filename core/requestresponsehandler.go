package core

import (
	"bufio"
	"net"
)

func CreateTCPConnection(server string) (net.Conn, error) {
	return net.Dial("tcp", server)
}

func FetchGopherServerResponse(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	var response string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", FetchErrorResponse(ResponseError, err)
		}

		response += line

		if line == ".\r\n" {
			break
		}
	}
	return response, nil
}
