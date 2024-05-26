package pmq

import (
	"errors"
	"github.com/joe-at-startupmedia/posix_mq"
	"github.com/joe-at-startupmedia/xipc"
	"google.golang.org/protobuf/proto"
	"syscall"
	"time"
)

type MqResponder BidirectionalQueue

func NewResponder(config *QueueConfig, owner *Ownership) xipc.IMqResponder {

	requester, errRqst := openQueueForResponder(config, owner, "rqst")
	responder, errResp := openQueueForResponder(config, owner, "resp")

	var mqr xipc.IMqResponder = &MqResponder{
		requester,
		errRqst,
		responder,
		errResp,
	}

	return mqr
}

func openQueueForResponder(config *QueueConfig, owner *Ownership, postfix string) (*posix_mq.MessageQueue, error) {

	if config.Flags == 0 {
		config.Flags = O_RDWR | O_CREAT | O_NONBLOCK
	}
	return NewMessageQueueWithOwnership(*config, owner, postfix)
}

func (mqr *MqResponder) Read() ([]byte, error) {
	msg, _, err := mqr.MqRqst.Receive()
	if err != nil {
		//EAGAIN simply means the queue is empty when O_NONBLOCK is set
		mqrAttr, _ := mqr.MqRqst.GetAttr()
		if mqrAttr != nil && (mqrAttr.Flags&O_NONBLOCK == O_NONBLOCK) && errors.Is(err, syscall.EAGAIN) {
			return msg, nil
		}
	}
	return msg, err
}

func (mqr *MqResponder) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, _, err := mqr.MqRqst.TimedReceive(duration)
	return msg, err
}

func (mqr *MqResponder) Write(data []byte) error {
	return mqr.MqResp.Send(data, 0)
}

func (mqr *MqResponder) HandleRequest(msgHandler xipc.ResponderCallback) error {
	return xipc.HandleRequestWithLag(mqr, msgHandler, 0)
}

// HandleRequestWithLag used for testing purposes to simulate lagging responder
func (mqr *MqResponder) HandleRequestWithLag(msgHandler xipc.ResponderCallback, lag int) error {
	return xipc.HandleRequestWithLag(mqr, msgHandler, lag)
}

// HandleMqRequest provides a concrete implementation of HandleRequestFromProto using the local MqRequest type
func (mqr *MqResponder) HandleMqRequest(requestProcessor xipc.ResponderMqRequestCallback) error {
	return xipc.HandleMqRequest(mqr, requestProcessor)
}

// HandleRequestFromProto used to process arbitrary protobuf messages using a callback
func (mqr *MqResponder) HandleRequestFromProto(protocMsg proto.Message, msgHandler xipc.ResponderFromProtoMessageCallback) error {
	return xipc.HandleRequestFromProto(mqr, protocMsg, msgHandler)
}

func (mqr *MqResponder) HasErrors() bool {
	return (*BidirectionalQueue)(mqr).HasErrors()
}

func (mqr *MqResponder) Error() error {
	return (*BidirectionalQueue)(mqr).Error()
}

func (mqr *MqResponder) CloseResponder() error {
	return (*BidirectionalQueue)(mqr).Unlink()
}
