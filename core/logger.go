package core

import (
	"fmt"
	"time"
)

type LogType int

func logRequest(request string) {
	fmt.Printf("Timestamp: %s, Client Request: %s\r\n", time.Now().Format(time.RFC3339), request)
}
