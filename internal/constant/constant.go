package constant

import "time"

// runtime.
const (
	DefaultAppName = "GO Backend Template API"
	DefaultAppPort = 8000
	DefaultAppHost = "127.0.0.1"
)

// db connection.
const (
	DefaultDBSSLMode         = "disable"
	DefaultDBMaxOpenConns    = 30
	DefaultDBMaxIdleConns    = 30
	DefaultDBMaxIdleTime     = 15 * time.Second
	DefaultDBMaxLifetime     = 60 * time.Second
	DefaultDBMaxConnWaitTime = 1 * time.Second
	DefaultDBMaxConnLifetime = 3600 * time.Second
	DefaultDBMaxConnIdleTime = 60 * time.Second
)
