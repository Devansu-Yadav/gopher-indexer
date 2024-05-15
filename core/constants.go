package core

import "time"

const (
	MaxFileSize        int           = 1024 * 1024 * 5 // 5 MB
	MaxResponseTimeOut time.Duration = 5 * time.Second
	TextFile           string        = "text"
	BinaryFile         string        = "binary"
	ErrorFile          string        = "error"
	InvalidRef         string        = "invalid"
)
