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

func (qc *QueueConfig) GetFile(configName string) string {
	if len(configName) > 0 {
		return qc.BasePath + configName
	} else {
		return qc.BasePath + qc.Name
	}
}

func NewResponderWithOwnership(config *QueueConfig, owner *xipc.Ownership, postfix string) (*shmemipc.IpcResponder, error) {

	if config.Mode == 0 {
		if owner != nil && owner.IsValid() {
			config.Mode = 0660

		} else {
			config.Mode = 0666
		}
	}

	responder := shmemipc.NewResponder(config.GetFile(""), config.MaxMsgSize)
	if err := responder.GetError(); err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create message queue %s: %-v", config.GetFile(""), err))
	}

	err := ApplyOwnership(config, owner, "_rqst")
	if err != nil {
		return responder, err
	}

	err = ApplyOwnership(config, owner, "_resp")
	if err != nil {
		return responder, err
	}

	return responder, nil
}

func ApplyOwnership(config *QueueConfig, owner *xipc.Ownership, postfix string) error {

	var configName string //make a copy so we don't modify the reference

	if len(postfix) > 0 {
		configName = fmt.Sprintf("%s_%s", config.Name, postfix)
	}

	err := owner.ApplyPermissions(config.GetFile(configName), config.Mode)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not apply permissions %s: %-v", config.GetFile(configName), err))
	}

	return nil
}
