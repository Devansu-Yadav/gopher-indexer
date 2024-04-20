package core

import (
	"fmt"
	"net"
)

func ConnectToServer(server string) (net.Conn, string, error) {
	conn, err := CreateTCPConnection(server)

	if err != nil {
		return nil, "", FetchErrorResponse(ConnectionError, err)
	}

	fmt.Fprintf(conn, "\r\n")

	response, responseErr := FetchGopherServerResponse(conn)

	if responseErr != nil {
		return nil, "", responseErr
	}
	return conn, response, nil
}

func CloseConnection(conn net.Conn) {
	conn.Close()
}
