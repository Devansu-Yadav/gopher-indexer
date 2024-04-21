package core

import (
	"fmt"
	"os"
	"time"
)

type LogType int

const (
	FileSizeExceeded LogType = iota
	ReadFileTimeOut
	ServerResponseTimeOut
	ExternalServerConnection
)

func logRequest(request string) {
	fmt.Printf("Timestamp: %s, Client Request: %s\r\n", time.Now().Format(time.RFC3339), request)
}

func logError(logType LogType, resource string) {
	switch logType {
	case FileSizeExceeded:
		fmt.Fprintln(os.Stderr, "Error: File size exceeds 5MB limit - "+resource)
	case ReadFileTimeOut:
		fmt.Fprintln(os.Stderr, "Error: File read timed out - "+resource)
	case ServerResponseTimeOut:
		fmt.Fprintln(os.Stderr, "Error: Server response timed out - "+resource)
	case ExternalServerConnection:
		fmt.Fprintln(os.Stderr, "Error: External server connection timed out - "+resource)
	}
}
