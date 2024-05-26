package net

import (
	"fmt"
	"github.com/joe-at-startupmedia/gipc"
	"github.com/joe-at-startupmedia/xipc"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/protobuf/proto"
)

type MqResponder struct {
	MqResp  *gipc.Server
	ErrResp error
	Logger  *logrus.Logger
}

func NewResponder(config *QueueConfig) xipc.IMqResponder {

	logger := xipc.InitLogging(config.LogLevel)

	responder, errResp := gipc.StartServer(&gipc.ServerConfig{
		Name:              config.Name,
		UnmaskPermissions: config.ServerUnmaskPermissions,
		LogLevel:          config.LogLevel,
		Encryption:        false,
	})

	var mqr xipc.IMqResponder = &MqResponder{
		responder,
		errResp,
		logger,
	}

	return mqr
}

func (mqr *MqResponder) Read() ([]byte, error) {
	msg, err := mqr.MqResp.Read()
	if msg.MsgType < 1 {
		mqr.Logger.Debugln("using recursion")
		time.Sleep(2 * time.Microsecond)
		return mqr.Read()
	} else {
		return msg.Data, err
	}
}

func (mqr *MqResponder) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, err := mqr.MqResp.ReadTimed(duration)
	if msg.MsgType < 1 {
		mqr.Logger.Debugln("using recursion")
		time.Sleep(2 * time.Microsecond)
		return mqr.ReadTimed(duration)
	} else {
		return msg.Data, err
	}
}

func (mqr *MqResponder) Write(data []byte) error {
	err := mqr.MqResp.Write(DEFAULT_MSG_TYPE, data)
	if err != nil && err.Error() == "Connecting" {
		mqr.Logger.Infoln("Connecting error, reattempting")
		time.Sleep(xipc.REQUEST_RECURSION_WAITTIME * time.Second)
		return mqr.Write(data)
	}
	return err
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
	return mqr.ErrResp != nil
}

func (mqr *MqResponder) Error() error {
	return fmt.Errorf("requester: %w", mqr.ErrResp)
}

func (mqr *MqResponder) CloseResponder() error {
	mqr.MqResp.Close()
	return nil
}
