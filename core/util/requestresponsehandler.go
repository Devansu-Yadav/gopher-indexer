package util

import (
	"bufio"
	"net"
)

func CreateTCPConnection(server string) (net.Conn, error) {
	return net.Dial("tcp", server)
}

func HandleGopherServerResponse(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	var response string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			HandleServerResponseError(err)
		}

		response += line

		if line == ".\r\n" {
			break
		}
	}
	return response
}
