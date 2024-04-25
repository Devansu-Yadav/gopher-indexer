package core

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func CreateTCPConnection(server string) (net.Conn, error) {
	conn, err := net.Dial("tcp", server)
	return conn, err
}

func CreateTCPConnectionWithTimeOut(server string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", server, MaxResponseTimeOut)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func FetchServerResponse(conn net.Conn, resource string) (string, error) {
	if resource != "" {
		fmt.Fprintf(conn, "%s\r\n", resource)
	} else {
		fmt.Fprintf(conn, "\r\n")
	}

	reader := bufio.NewReader(conn)
	var response string

	for {
		conn.SetReadDeadline(time.Now().Add(MaxResponseTimeOut))

		line, err := reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return "", FetchErrorResponse(ServerResponseTimeOut, err, resource)
			} else {
				return "", FetchErrorResponse(ResponseError, err, resource)
			}
		}

		response += line

		if line == ".\r\n" {
			break
		}
	}
	return response, nil
}

func FetchResourcesFromDirectory(server, resource string) (string, error) {
	conn, err := CreateTCPConnection(server)
	logRequest(server + resource)

	if err != nil {
		return "", FetchErrorResponse(ConnectionError, err, server+resource)
	}
	defer conn.Close()

	response, responseErr := FetchServerResponse(conn, resource)

	if responseErr != nil {
		return "", responseErr
	}
	return response, nil
}
