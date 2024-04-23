package core

import (
	"bufio"
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

func ScanFilesAndExtServersInDir(server string, files map[string][]string, externalServers []string) {
	for fileType, fileList := range files {
		for _, file := range fileList {
			size, err := FetchFileAttrs(server, file)
			if err != nil {
				fmt.Fprintf(os.Stderr, ""+err.Error()+" for file: "+file)
			}

			fmt.Printf("File: %s, Type: %s, Size: %d\n", file, fileType, size)
		}
	}

	for _, extServer := range externalServers {
		extServerRes, extServerErr := FetchResourcesFromExternalServer(extServer)
		if extServerErr != nil {
			fmt.Fprintf(os.Stderr, ""+extServerErr.Error()+" for external server: "+extServer)
		}

		if extServerRes != "" {
			fmt.Printf("External server: %s is up!", extServer)
		}
	}
}

func CrawlGopherServer(server string, initialServerResponse string) {
	// Scrape the initial response to get the directories, files, and external servers
	directories, files, externalServers, invalidReferences := ScrapeGopherResponse(server+"/", initialServerResponse)

	visited := make(map[string]bool)

	// Recursively scanning each directory within root directory
	for _, dir := range directories {
		ScanDirectories(server, dir, visited)
	}

	// Scan files and references to external servers in root dir
	ScanFilesAndExtServersInDir(server, files, externalServers)

	for _, ref := range invalidReferences {
		fmt.Fprintln(os.Stderr, "Invalid reference found for resource - ", ref)
	}
}

func ScanDirectories(server, directory string, visited map[string]bool) {
	if visited[directory] {
		return
	}

	visited[directory] = true

	// Fetch the resources from the directory
	response, err := FetchResourcesFromDirectory(server, directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, ""+err.Error()+" for directory: "+directory)
	}

	directories, files, externalServers, invalidReferences := ScrapeGopherResponse(server+directory, response)
	// Recursively scanning current dir
	for _, dir := range directories {
		ScanDirectories(server, dir, visited)
	}

	// Scan files and external servers in current dir
	ScanFilesAndExtServersInDir(server, files, externalServers)

	for _, ref := range invalidReferences {
		logError(InvalidReference, ref)
	}
}
