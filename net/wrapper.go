package net

import (
	"time"
)

// QueueConfig is used to configure an instance of the message queue.
type QueueConfig struct {
	Name                    string
	ServerUnmaskPermissions bool
	ClientRetryTimer        time.Duration
	ClientTimeout           time.Duration
	LogLevel                string
}

const (
	DEFAULT_MSG_TYPE = 1
)
