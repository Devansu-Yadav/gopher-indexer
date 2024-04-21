package core

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func CreateTCPConnection(server string) (net.Conn, error) {
	conn, err := net.Dial("tcp", server)
	logRequest(server)

	return conn, err
}

func CreateTCPConnectionWithTimeOut(server string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", server, MaxResponseTimeOut)
	logRequest(server)

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
				logError(ServerResponseTimeOut, resource)
				break
			} else {
				return "", FetchErrorResponse(ResponseError, err)
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
	if err != nil {
		return "", FetchErrorResponse(ConnectionError, err)
	}
	defer conn.Close()

	response, responseErr := FetchServerResponse(conn, resource)

	if responseErr != nil {
		return "", responseErr
	}
	return response, nil
}

func FetchResourcesFromExternalServer(server, resource string) (string, error) {
	conn, err := CreateTCPConnectionWithTimeOut(server)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			logError(ExternalServerConnection, resource)
			return "", nil
		}
		return "", FetchErrorResponse(ConnectionError, err)
	}
	defer conn.Close()

	response, responseErr := FetchServerResponse(conn, resource)

	if responseErr != nil {
		return "", responseErr
	}
	return response, nil
}
