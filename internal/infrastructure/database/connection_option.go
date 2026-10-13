package db

import (
	"time"

	"github.com/ABA-Developer/go-backend-template/constant"
)

type connectionOption struct {
	maxIdle           int
	maxOpen           int
	maxLifetime       time.Duration
	maxConnIdleTime   time.Duration
	keepAliveInterval time.Duration
}

func defaultConnectionOption() *connectionOption {
	return &connectionOption{
		maxIdle:           constant.DefaultDBMaxIdleConns,
		maxOpen:           constant.DefaultDBMaxOpenConns,
		maxLifetime:       constant.DefaultDBMaxConnLifetime,
		maxConnIdleTime:   constant.DefaultDBMaxConnIdleTime,
		keepAliveInterval: constant.DefaultDBKeepAliveInterval,
	}
}
