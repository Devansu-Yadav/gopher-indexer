package core

import (
	"fmt"
	"net"

	"github.com/Devansu-Yadav/gopher-indexer/core/util"
)

func ConnectToServer(server string) (net.Conn, string, error) {
	conn, err := util.CreateTCPConnection(server)

	if err != nil {
		util.HandleServerConnectionError(err)
	}

	fmt.Fprintf(conn, "\r\n")

	response := util.HandleGopherServerResponse(conn)
	return conn, response, nil
}

func CloseConnection(conn net.Conn) {
	conn.Close()
}
