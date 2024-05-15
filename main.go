package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Devansu-Yadav/gopher-indexer/core"
)

func main() {
	flag.StringVar(&core.GopherServerHost, "host", "defaultHost", "Gopher server host")
	flag.StringVar(&core.GopherServerPort, "port", "defaultPort", "Gopher server port no")
	flag.Parse()

	if core.GopherServerHost == "" {
		fmt.Println("Error: You must provide a valid host")
		os.Exit(1)
	}

	if core.GopherServerPort == "" {
		fmt.Println("Error: You must provide a valid port no")
		os.Exit(1)
	}

	GopherServerConnectionString := fmt.Sprintf("%s:%s", core.GopherServerHost, core.GopherServerPort)

	// 1. Create initial connection
	response, err := core.ConnectToServer(GopherServerConnectionString)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// fmt.Println(response)

	// 2. Scan each directory on the server
	core.CrawlGopherServer(GopherServerConnectionString, response)
}
