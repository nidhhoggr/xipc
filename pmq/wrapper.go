package pmq

import (
	"errors"
	"fmt"
	"github.com/nidhhoggr/posix_mq"
	"github.com/nidhhoggr/xipc"
)

type QueueConfig posix_mq.QueueConfig

const (
	O_RDONLY = posix_mq.O_RDONLY
	O_WRONLY = posix_mq.O_WRONLY
	O_RDWR   = posix_mq.O_RDWR

	O_CLOEXEC  = posix_mq.O_CLOEXEC
	O_CREAT    = posix_mq.O_CREAT
	O_EXCL     = posix_mq.O_EXCL
	O_NONBLOCK = posix_mq.O_NONBLOCK
)

// GetFile gets the file on the OS where the queues are stored
func (config *QueueConfig) GetFile() string {
	return (*posix_mq.QueueConfig)(config).GetFile()
}

// NewMessageQueue returns an instance of the message queue given a QueueConfig.
func NewMessageQueue(config *QueueConfig) (*posix_mq.MessageQueue, error) {
	return posix_mq.NewMessageQueue((*posix_mq.QueueConfig)(config))
}

func NewMessageQueueWithOwnership(config QueueConfig, owner *xipc.Ownership, postfix string) (*posix_mq.MessageQueue, error) {

	if len(postfix) > 0 {
		config.Name = fmt.Sprintf("%s_%s", config.Name, postfix)
	}

	var (
		messageQueue *posix_mq.MessageQueue
	)

	if config.Mode == 0 {
		if owner != nil && owner.IsValid() {
			config.Mode = 0660

		} else {
			config.Mode = 0666
		}
	}

	messageQueue, err := NewMessageQueue(&config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create message queue %s: %-v", config.GetFile(), err))
	}

	err = owner.ApplyPermissions(config.GetFile(), config.Mode)
	if err != nil {
		return messageQueue, errors.New(fmt.Sprintf("Could not apply permissions %s: %-v", config.GetFile(), err))
	}

	return messageQueue, nil
}

type BidirectionalQueue struct {
	Rqst    *posix_mq.MessageQueue
	ErrRqst error
	Resp    *posix_mq.MessageQueue
	ErrResp error
}

func (bdr *BidirectionalQueue) Close() error {
	if err := bdr.Rqst.Close(); err != nil {
		return err
	}
	return bdr.Resp.Close()
}

func (bdr *BidirectionalQueue) Unlink() error {
	if err := bdr.Rqst.Unlink(); err != nil {
		return err
	}
	return bdr.Resp.Unlink()
}

func (bdr *BidirectionalQueue) HasErrors() bool {
	return bdr.ErrResp != nil || bdr.ErrRqst != nil
}

func (bdr *BidirectionalQueue) Error() error {
	if bdr.ErrResp != nil {
		return bdr.ErrResp
	}
	return bdr.ErrRqst
}

func ForceRemoveQueue(queueFile string) error {
	return posix_mq.ForceRemoveQueue(queueFile)
}
