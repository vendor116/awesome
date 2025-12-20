package httpserver

import "time"

const (
	readTimeout       = 5 * time.Second
	readHeaderTimeout = 2 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 30 * time.Second
	shutdownTimeout   = 3 * time.Second
	maxHeaderBytes    = 1 << 20
)
