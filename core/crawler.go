package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ScrapeGopherResponse(serverResourceLink, response string) ([]string, map[string][]string, []string, []string, []string) {
	if response == "" {
		return nil, nil, nil, nil, nil
	}

	reader := bufio.NewReader(strings.NewReader(response))
	var directories, externalServers, invalidReferences, malformedReferences []string
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
			malformedReferences = append(malformedReferences, parts[1])
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
			files[TextFile] = append(files[TextFile], resource)
		case "4", "5", "6", "9":
			files[BinaryFile] = append(files[BinaryFile], resource)
		}
	}

	return directories, files, externalServers, invalidReferences, malformedReferences
}

func ScanFilesAndExtServersInDir(server string, files map[string][]string, externalServers []string, serverResources map[string][]string, extServerStatuses map[string]bool, fileSizes map[string][]int, smallestTextFileContents *string) {
	for fileType, fileList := range files {
		for _, file := range fileList {
			contents, size, err, isMalformed := FetchFileAttrs(server, file, fileType)
			if err != nil {
				fmt.Fprintln(os.Stderr, ""+err.Error())
				serverResources[ErrorFile] = append(serverResources[ErrorFile], file)
				continue
			}

			if isMalformed {
				serverResources[ErrorFile] = append(serverResources[ErrorFile], file)
			} else {
				serverResources[fileType] = append(serverResources[fileType], file)
			}

			fmt.Printf("File: %s, Type: %s, Size: %d B\n", file, fileType, size)

			// Update the size of the smallest text/binary file found so far
			if size < fileSizes[fileType][0] {
				fileSizes[fileType][0] = size
				if fileType == TextFile {
					*smallestTextFileContents = string(contents)
				}
			}

			// Update the size of the largest text/binary file found so far
			if size > fileSizes[fileType][1] {
				fileSizes[fileType][1] = size
			}
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
		serverResources[InvalidRef] = append(serverResources[InvalidRef], ref)
	}
}

func StoreMalformedReferences(malformedReferences []string, serverResources map[string][]string) {
	serverResources[ErrorFile] = append(serverResources[ErrorFile], malformedReferences...)
}

func CrawlGopherServer(server string, initialServerResponse string) {
	// Scrape the initial response to get the directories, files, and external servers
	directories, files, externalServers, invalidReferences, malformedReferences := ScrapeGopherResponse(server+"/", initialServerResponse)

	visited := make(map[string]bool)
	serverResources := make(map[string][]string)

	extServerStatuses := make(map[string]bool)
	fileSizes := map[string][]int{
		TextFile:   {MaxFileSize, 0},
		BinaryFile: {MaxFileSize, 0},
	}
	smallestTextFileContents := ""

	// visited the root dir
	visited["/"] = true

	// Recursively scanning each directory within root directory
	for _, dir := range directories {
		ScanDirectories(server, dir, visited, serverResources, extServerStatuses, fileSizes, &smallestTextFileContents)
	}

	// Scan files and references to external servers in root dir
	ScanFilesAndExtServersInDir(server, files, externalServers, serverResources, extServerStatuses, fileSizes, &smallestTextFileContents)

	// Handle invalid references
	StoreInvalidReferences(invalidReferences, serverResources)

	// Handle malformed references
	StoreMalformedReferences(malformedReferences, serverResources)

	fmt.Println("\n================= Gopher Server stats =======================")
	fmt.Println("a. No of Gopher directories on server: ", len(visited))
	fmt.Println("List of directories: ")
	for dir := range visited {
		fmt.Println(server + dir)
	}

	fmt.Println("\nb. Total no of simple text files: ", len(serverResources[TextFile]))
	fmt.Println("List of simple text files(full path): ")
	for _, file := range serverResources[TextFile] {
		fmt.Println(server + file)
	}

	fmt.Println("\nc. Total no of binary files: ", len(serverResources[BinaryFile]))
	fmt.Println("List of binary files(full path): ")
	for _, file := range serverResources[BinaryFile] {
		fmt.Println(server + file)
	}

	fmt.Println("\nd. The contents of the smallest text file:")
	fmt.Println(smallestTextFileContents)

	fmt.Printf("\ne. Size of the largest text file: %d B", fileSizes[TextFile][1])
	fmt.Printf("\nf. Size of the smallest binary file: %d B", fileSizes[BinaryFile][0])
	fmt.Printf("\nf. Size of the largest binary file: %d B\n", fileSizes[BinaryFile][1])

	fmt.Println("\ng. Total no of unique invalid references: ", len(serverResources[InvalidRef]))
	fmt.Println("List of invalid references(full path): ")
	for _, ref := range serverResources[InvalidRef] {
		fmt.Println(ref)
	}

	fmt.Println("\nh. List of external servers: ")
	for extServer, status := range extServerStatuses {
		upOrDown := "down"
		if status {
			upOrDown = "up"
		}
		fmt.Printf("External server: %s is %s\n", extServer, upOrDown)
	}

	fmt.Println("\ni. List of references with issues/errors(timeout, malformed or large file size): ")
	for _, errorRef := range serverResources[ErrorFile] {
		fmt.Println(server + errorRef)
	}
}

func ScanDirectories(server, directory string, visited map[string]bool, serverResources map[string][]string, extServerStatuses map[string]bool, fileSizes map[string][]int, smallestTextFileContents *string) {
	if visited[directory] {
		return
	}

	visited[directory] = true

	// Fetch the resources from the directory
	response, err := FetchResourcesFromDirectory(server, directory)
	if err != nil {
		fmt.Fprintln(os.Stderr, ""+err.Error()+" for directory: "+directory)
	}

	directories, files, externalServers, invalidReferences, malformedReferences := ScrapeGopherResponse(server+directory, response)
	// Recursively scanning current dir
	for _, dir := range directories {
		ScanDirectories(server, dir, visited, serverResources, extServerStatuses, fileSizes, smallestTextFileContents)
	}

	// Scan files and external servers in current dir
	ScanFilesAndExtServersInDir(server, files, externalServers, serverResources, extServerStatuses, fileSizes, smallestTextFileContents)

	// Handle invalid references
	StoreInvalidReferences(invalidReferences, serverResources)

	// Handle malformed references
	StoreMalformedReferences(malformedReferences, serverResources)
}
