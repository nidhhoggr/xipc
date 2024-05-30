package net

import (
	"time"
)

// QueueConfig is used to configure an instance of the message queue.
type QueueConfig struct {
	Name                    string
	LogLevel                string
	ClientRetryTimer        time.Duration
	ClientTimeout           time.Duration
	ServerUnmaskPermissions bool
}

const (
	DEFAULT_MSG_TYPE = 1
)
