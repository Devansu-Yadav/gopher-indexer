package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func CreateTCPConnection(server string) (net.Conn, error) {
	return net.Dial("tcp", server)
}

func CreateTCPConnectionWithTimeOut(server string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", server, MaxResponseTimeOut)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func FetchServerResponse(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	var response string

	for {
		conn.SetReadDeadline(time.Now().Add(MaxResponseTimeOut))

		line, err := reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Fprintln(os.Stderr, "Server Response timed out!")
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

	fmt.Fprintf(conn, "%s\r\n", resource)

	response, responseErr := FetchServerResponse(conn)

	if responseErr != nil {
		return "", responseErr
	}
	return response, nil
}

func FetchResourcesFromExternalServer(server, resource string) (string, error) {
	conn, err := CreateTCPConnectionWithTimeOut(server)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "External server connection timed out!")
			return "", nil
		}
		return "", FetchErrorResponse(ConnectionError, err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "%s\r\n", resource)

	response, responseErr := FetchServerResponse(conn)

	if responseErr != nil {
		return "", responseErr
	}
	return response, nil
}
