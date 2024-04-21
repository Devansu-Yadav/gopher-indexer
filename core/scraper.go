package core

import (
	"bufio"
	"strings"
)

func ScrapeGopherResponse(response string) ([]string, map[string][]string, []string, error) {
	if response == "" {
		return nil, nil, nil, nil
	}

	reader := bufio.NewReader(strings.NewReader(response))
	var directories, externalServers []string
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

		// ignore resources without host name and port
		if host == "" || port == "" {
			continue
		}

		if itemType == "1" {
			if resource != "" && host == GopherServerHost && port == GopherServerPort {
				directories = append(directories, resource)
			} else {
				externalServers = append(externalServers, host+":"+port)
			}
		}

		switch itemType {
		case "0":
			files["text"] = append(files["text"], resource)
		case "4", "5", "6", "9":
			files["binary"] = append(files["binary"], resource)
		}
	}

	return directories, files, externalServers, nil
}
