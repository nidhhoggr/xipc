package mem

import (
	"github.com/joe-at-startupmedia/shmemipc"
	"github.com/joe-at-startupmedia/xipc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"time"
)

type Responder struct {
	Resp    *shmemipc.IpcResponder
	ErrResp error
	logger  *logrus.Logger
}

func NewResponder(config *QueueConfig, owner *xipc.Ownership) xipc.IResponder {

	logger := xipc.InitLogging(config.LogLevel)

	responder, err := NewResponderWithOwnership(config, owner, "")
	if err != nil {
		logger.Fatal(err)
	}

	var mqr xipc.IResponder = &Responder{
		Resp:    responder,
		ErrResp: responder.GetError(),
		logger:  logger,
	}

	return mqr
}

func (mqr *Responder) Read() ([]byte, error) {
	return mqr.Resp.Read()
}

func (mqr *Responder) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, err := mqr.Resp.ReadTimed(duration)
	if err != nil {
		if err.Error() == "timed_out" {
			mqr.logger.Info("responder timed out")
			time.Sleep(xipc.REQUEST_RECURSION_WAITTIME * time.Millisecond)
			return mqr.ReadTimed(duration)
		}
		return nil, err
	}
	return msg, nil
}

func (mqr *Responder) Write(data []byte) error {
	return mqr.Resp.Write(data)
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
	return mqr.ErrResp != nil
}

func (mqr *Responder) Error() error {
	return mqr.ErrResp
}

func (mqr *Responder) CloseResponder() error {
	return mqr.Resp.Close()
}
