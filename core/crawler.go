package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ScrapeGopherResponse(serverResourceLink, response string) ([]string, map[string][]string, []string, []string) {
	if response == "" {
		return nil, nil, nil, nil
	}

	reader := bufio.NewReader(strings.NewReader(response))
	var directories, externalServers, invalidReferences []string
	files := make(map[string][]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSuffix(line, "\r\n")

		// ignore malformed resources
		if len(line) < 2 {
			continue
		}

		itemType := string(line[0])

		if itemType == "3" {
			invalidReferences = append(invalidReferences, serverResourceLink)
			continue
		}

		parts := strings.Split(line[1:], "\t")

		// ignore malformed resources
		if len(parts) < 4 {
			continue
		}

		resource := parts[1]
		host := parts[2]
		port := parts[3]

		if itemType == "1" {
			if host != GopherServerHost || port != GopherServerPort {
				externalServers = append(externalServers, host+":"+port)
			} else if resource != "" {
				directories = append(directories, resource)
			}
		}

		switch itemType {
		case "0":
			files["text"] = append(files["text"], resource)
		case "4", "5", "6", "9":
			files["binary"] = append(files["binary"], resource)
		}
	}

	return directories, files, externalServers, invalidReferences
}

func ScanFilesAndExtServersInDir(server string, files map[string][]string, externalServers []string, serverResources map[string][]string, extServerStatuses map[string]bool) {
	for fileType, fileList := range files {
		for _, file := range fileList {
			size, err := FetchFileAttrs(server, file)
			if err != nil {
				fmt.Fprintln(os.Stderr, ""+err.Error())
				continue
			} else {
				serverResources[fileType] = append(serverResources[fileType], file)
			}

			fmt.Printf("File: %s, Type: %s, Size: %d B\n", file, fileType, size)
		}
	}

	for _, extServer := range externalServers {
		extServerErr := ConnectToExternalServer(extServer)

		if extServerErr != nil {
			fmt.Fprintln(os.Stderr, ""+extServerErr.Error())
			extServerStatuses[extServer] = false
			continue
		}

		fmt.Printf("External server: %s is up!\n", extServer)
		extServerStatuses[extServer] = true
	}
}

func StoreInvalidReferences(invalidReferences []string, serverResources map[string][]string) {
	for _, ref := range invalidReferences {
		fmt.Fprintln(os.Stderr, ""+FetchErrorResponse(InvalidReference, errors.New(""), ref).Error())
		serverResources["invalid"] = append(serverResources["invalid"], ref)
	}
}

func CrawlGopherServer(server string, initialServerResponse string) {
	// Scrape the initial response to get the directories, files, and external servers
	directories, files, externalServers, invalidReferences := ScrapeGopherResponse(server+"/", initialServerResponse)

	visited := make(map[string]bool)
	serverResources := make(map[string][]string)
	extServerStatuses := make(map[string]bool)

	// visited the root dir
	visited["/"] = true

	// Recursively scanning each directory within root directory
	for _, dir := range directories {
		ScanDirectories(server, dir, visited, serverResources, extServerStatuses)
	}

	// Scan files and references to external servers in root dir
	ScanFilesAndExtServersInDir(server, files, externalServers, serverResources, extServerStatuses)

	// Handle invalid references
	StoreInvalidReferences(invalidReferences, serverResources)

	fmt.Println("\n================= Gopher Server stats =======================")
	fmt.Println("No of Gopher directories on server: ", len(visited))
	fmt.Println("List of directories: ")
	for dir := range visited {
		fmt.Println(server + dir)
	}

	fmt.Println("\nTotal no of simple text files: ", len(serverResources["text"]))
	fmt.Println("List of simple text files(full path): ")
	for _, file := range serverResources["text"] {
		fmt.Println(server + file)
	}

	fmt.Println("\nTotal no of binary files: ", len(serverResources["binary"]))
	fmt.Println("List of binary files(full path): ")
	for _, file := range serverResources["binary"] {
		fmt.Println(server + file)
	}

	fmt.Println("\nTotal no of unique invalid references: ", len(serverResources["invalid"]))
	fmt.Println("List of invalid references(full path): ")
	for _, ref := range serverResources["invalid"] {
		fmt.Println(server + ref)
	}

	fmt.Println("\nList of external servers: ")
	for extServer, status := range extServerStatuses {
		upOrDown := "down"
		if status {
			upOrDown = "up"
		}
		fmt.Printf("External server: %s is %s\n", extServer, upOrDown)
	}
}

func ScanDirectories(server, directory string, visited map[string]bool, serverResources map[string][]string, extServerStatuses map[string]bool) {
	if visited[directory] {
		return
	}

	visited[directory] = true

	// Fetch the resources from the directory
	response, err := FetchResourcesFromDirectory(server, directory)
	if err != nil {
		fmt.Fprintln(os.Stderr, ""+err.Error()+" for directory: "+directory)
	}

	directories, files, externalServers, invalidReferences := ScrapeGopherResponse(server+directory, response)
	// Recursively scanning current dir
	for _, dir := range directories {
		ScanDirectories(server, dir, visited, serverResources, extServerStatuses)
	}

	// Scan files and external servers in current dir
	ScanFilesAndExtServersInDir(server, files, externalServers, serverResources, extServerStatuses)

	StoreInvalidReferences(invalidReferences, serverResources)
}
