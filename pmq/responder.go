package pmq

import (
	"errors"
	"github.com/nidhhoggr/posix_mq"
	"github.com/nidhhoggr/xipc"
	"google.golang.org/protobuf/proto"
	"syscall"
	"time"
)

type Responder BidirectionalQueue

func NewResponder(config *QueueConfig, owner *xipc.Ownership) xipc.IResponder {

	requester, errRqst := openQueueForResponder(config, owner, "rqst")
	responder, errResp := openQueueForResponder(config, owner, "resp")

	var mqr xipc.IResponder = &Responder{
		requester,
		errRqst,
		responder,
		errResp,
	}

	return mqr
}

func openQueueForResponder(config *QueueConfig, owner *xipc.Ownership, postfix string) (*posix_mq.MessageQueue, error) {

	if config.Flags == 0 {
		config.Flags = O_RDWR | O_CREAT | O_NONBLOCK
	}
	return NewMessageQueueWithOwnership(*config, owner, postfix)
}

func (mqr *Responder) Read() ([]byte, error) {
	msg, _, err := mqr.Rqst.Receive()
	if err != nil {
		//EAGAIN simply means the queue is empty when O_NONBLOCK is set
		mqrAttr, _ := mqr.Rqst.GetAttr()
		if mqrAttr != nil && (mqrAttr.Flags&O_NONBLOCK == O_NONBLOCK) && errors.Is(err, syscall.EAGAIN) {
			return msg, nil
		}
	}
	return msg, err
}

func (mqr *Responder) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, _, err := mqr.Rqst.TimedReceive(duration)
	return msg, err
}

func (mqr *Responder) Write(data []byte) error {
	return mqr.Resp.Send(data, 0)
}

func (mqr *Responder) HandleRequest(msgHandler xipc.ResponderCallback) error {
	return xipc.HandleRequestWithLag(mqr, msgHandler, 0)
}

// HandleRequestWithLag used for testing purposes to simulate lagging responder
func (mqr *Responder) HandleRequestWithLag(msgHandler xipc.ResponderCallback, lag int) error {
	return xipc.HandleRequestWithLag(mqr, msgHandler, lag)
}

// HandleRequestProto provides a concrete implementation of HandleRequestFromProto using the local Request type
func (mqr *Responder) HandleRequestProto(requestProcessor xipc.ResponderRequestProtoCallback) error {
	return xipc.HandleRequestProto(mqr, requestProcessor)
}

// HandleRequestFromProto used to process arbitrary protobuf messages using a callback
func (mqr *Responder) HandleRequestFromProto(protocMsg proto.Message, msgHandler xipc.ResponderFromProtoMessageCallback) error {
	return xipc.HandleRequestFromProto(mqr, protocMsg, msgHandler)
}

func (mqr *Responder) HasErrors() bool {
	return (*BidirectionalQueue)(mqr).HasErrors()
}

func (mqr *Responder) Error() error {
	return (*BidirectionalQueue)(mqr).Error()
}

func (mqr *Responder) CloseResponder() error {
	return (*BidirectionalQueue)(mqr).Unlink()
}
