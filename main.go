package main

import (
	"log"

	"github.com/Devansu-Yadav/gopher-indexer/core"
)

func main() {
	// 1. Create initial connection
	response, err := core.ConnectToServer(core.GopherServerConnectionString)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// fmt.Println(response)

	// 2. Scan each directory on the server
	core.CrawlGopherServer(core.GopherServerConnectionString, response)
}
