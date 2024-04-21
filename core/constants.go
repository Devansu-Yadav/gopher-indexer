package core

import "time"

const (
	GopherServerConnectionString string        = "comp3310.ddns.net:70"
	GopherServerHost             string        = "comp3310.ddns.net"
	GopherServerPort             string        = "70"
	MaxFileSize                  int           = 1024 * 1024 * 5 // 5 MB
	MaxResponseTimeOut           time.Duration = 5 * time.Second
)
