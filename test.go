package main

func test() {
	// 2. Test file size
	// size, err := core.FetchFileAttrs(core.GopherServerConnectionString, "/misc/godot")
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }

	// fmt.Println(size)

	// Test timeouts when fetching resources
	// response, err := core.FetchResourcesFromExternalServer("comp3310.ddns.net:71", "")
	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }

	// fmt.Println(response)

	// 3. Test some scraping
	// response, err := core.FetchDirsAndExternalServerResponses(core.GopherServerConnectionString, "/acme/products")

	// if err != nil {
	// 	log.Fatalf("Error: %v", err)
	// }

	// directories, files, externalServers := core.ScrapeGopherResponse(response)

	// if scrapeErr != nil {
	// 	log.Fatalf("Error: %v", scrapeErr)
	// }

	// // fmt.Println(response)
	// fmt.Println("Directories:", directories)
	// fmt.Println("Files:", files)
	// fmt.Println("External Servers:", externalServers)
}
