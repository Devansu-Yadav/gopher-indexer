package util

import (
	"fmt"
	"net"
)

func HandleServerConnectionError(err error) (net.Conn, string, error) {
	return nil, "", fmt.Errorf("failed to connect: %v", err)
}

func HandleServerResponseError(err error) (net.Conn, string, error) {
	return nil, "", fmt.Errorf("failed to read response: %v", err)
}
