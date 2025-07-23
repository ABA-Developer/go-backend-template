package config

import "time"

type DBConfig struct {
	MigratorDriver    string
	Username          string
	Password          string
	Name              string
	Host              string
	Port              string
	SSLMode           string
	MaxOpenConn       int
	MaxIdleConn       int
	MaxIdleTime       int
	MaxLifetime       int
	MaxConnWaitTime   int
	MaxConnLifetime   int
	MaxConnIdleTime   int
	KeepAliveInterval time.Duration
}
