package core

import (
	"bufio"
	"fmt"
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

func FetchDirsAndExternalServerResponses(server, resource string) (net.Conn, string, error) {
	conn, err := CreateTCPConnection(server)
	if err != nil {
		return nil, "", FetchErrorResponse(ConnectionError, err)
	}

	fmt.Fprintf(conn, "%s\r\n", resource)

	response, responseErr := FetchGopherServerResponse(conn)

	if responseErr != nil {
		return nil, "", responseErr
	}
	return conn, response, nil
}
