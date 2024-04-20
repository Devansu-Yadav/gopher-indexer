package main

import (
	"fmt"
	"log"

	"github.com/Devansu-Yadav/gopher-indexer/core"
)

func main() {
	// 1. Create initial connection
	conn, response, err := core.ConnectToServer("comp3310.ddns.net:70")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(response)

	conn.Close()
}
