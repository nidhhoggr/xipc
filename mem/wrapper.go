package mem

import (
	"errors"
	"fmt"
	"github.com/joe-at-startupmedia/shmemipc"
	"github.com/joe-at-startupmedia/xipc"
)

type QueueConfig struct {
	Name       string
	BasePath   string
	LogLevel   string
	MaxMsgSize uint64
	Mode       int
}

func (qc *QueueConfig) GetFile() string {
	return qc.BasePath + qc.Name
}

func NewResponderWithOwnership(config *QueueConfig, owner *xipc.Ownership, postfix string) (*shmemipc.IpcResponder, error) {

	if len(postfix) > 0 {
		config.Name = fmt.Sprintf("%s_%s", config.Name, postfix)
	}

	if config.Mode == 0 {
		if owner != nil && owner.IsValid() {
			config.Mode = 0660

		} else {
			config.Mode = 0666
		}
	}

	responder := shmemipc.NewResponder(config.GetFile(), config.MaxMsgSize)
	if err := responder.GetError(); err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create message queue %s: %-v", config.GetFile(), err))
	}

	err := owner.ApplyPermissions(config.GetFile(), config.Mode)
	if err != nil {
		return responder, errors.New(fmt.Sprintf("Could not apply permissions %s: %-v", config.GetFile(), err))
	}

	return responder, nil
}
