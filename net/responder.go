package net

import (
	"github.com/nidhhoggr/gipc"
	"github.com/nidhhoggr/xipc"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/protobuf/proto"
)

type Responder struct {
	Resp    *gipc.Server
	ErrResp error
	Logger  *logrus.Logger
}

func NewResponder(config *QueueConfig) xipc.IResponder {

	logger := xipc.InitLogging(config.LogLevel)

	responder, errResp := gipc.StartServer(&gipc.ServerConfig{
		Name:              config.Name,
		UnmaskPermissions: config.ServerUnmaskPermissions,
		LogLevel:          config.LogLevel,
		Encryption:        false,
	})

	var mqr xipc.IResponder = &Responder{
		responder,
		errResp,
		logger,
	}

	return mqr
}

func (mqr *Responder) Read() ([]byte, error) {
	msg, err := mqr.Resp.Read()
	if msg.MsgType < 1 {
		mqr.Logger.Debugln("using recursion")
		time.Sleep(2 * time.Microsecond)
		return mqr.Read()
	} else {
		return msg.Data, err
	}
}

func (mqr *Responder) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, err := mqr.Resp.ReadTimed(duration)
	if msg.MsgType < 1 {
		mqr.Logger.Debugln("using recursion")
		time.Sleep(2 * time.Microsecond)
		return mqr.ReadTimed(duration)
	} else {
		return msg.Data, err
	}
}

func (mqr *Responder) Write(data []byte) error {
	err := mqr.Resp.Write(DEFAULT_MSG_TYPE, data)
	if err != nil && err.Error() == "Connecting" {
		mqr.Logger.Infoln("Connecting error, reattempting")
		time.Sleep(xipc.REQUEST_RECURSION_WAITTIME * time.Second)
		return mqr.Write(data)
	}
	return err
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
	mqr.Resp.Close()
	return nil
}
