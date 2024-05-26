package net

import (
	"github.com/joe-at-startupmedia/gipc"
	"github.com/joe-at-startupmedia/xipc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"time"
)

type MqRequester struct {
	MqRqst  *gipc.Client
	ErrRqst error
	Logger  *logrus.Logger
}

func NewRequester(config *QueueConfig) xipc.IMqRequester {

	logger := xipc.InitLogging(config.LogLevel)

	requester, errRqst := gipc.StartClient(&gipc.ClientConfig{
		Name:       config.Name,
		LogLevel:   config.LogLevel,
		RetryTimer: config.ClientRetryTimer,
		Timeout:    config.ClientTimeout,
		Encryption: false,
	})

	var mqs xipc.IMqRequester = &MqRequester{
		requester,
		errRqst,
		logger,
	}

	return mqs
}

func (mqs *MqRequester) Read() ([]byte, error) {
	msg, err := mqs.MqRqst.Read()
	if msg.MsgType < 1 {
		time.Sleep(xipc.REQUEST_RECURSION_WAITTIME * time.Millisecond)
		return mqs.WaitForResponse()
	} else {
		return msg.Data, err
	}
}

func (mqs *MqRequester) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, err := mqs.MqRqst.ReadTimed(duration)
	if err != nil {
		return nil, err
	}
	if msg.MsgType < 1 {
		time.Sleep(xipc.REQUEST_RECURSION_WAITTIME * time.Millisecond)
		return mqs.WaitForResponseTimed(duration)
	} else if msg == gipc.TimeoutMessage {
		return nil, gipc.TimeoutMessage.Err
	} else {
		return msg.Data, err
	}
}

func (mqs *MqRequester) Write(data []byte) error {
	return mqs.MqRqst.Write(DEFAULT_MSG_TYPE, data)
}

func (mqs *MqRequester) Request(data []byte) error {
	return mqs.Write(data)
}

func (mqs *MqRequester) RequestUsingMqRequest(req *xipc.MqRequest) error {
	return xipc.RequestUsingMqRequest(mqs, req)
}

func (mqs *MqRequester) RequestUsingProto(req *proto.Message) error {
	return xipc.RequestUsingProto(mqs, req)
}

func (mqs *MqRequester) WaitForResponse() ([]byte, error) {
	return mqs.Read()
}

func (mqs *MqRequester) WaitForResponseTimed(duration time.Duration) ([]byte, error) {
	return mqs.ReadTimed(duration)
}

func (mqs *MqRequester) WaitForProto(pbm proto.Message) (*proto.Message, error) {
	return xipc.WaitForProto(mqs, pbm)
}

func (mqs *MqRequester) WaitForProtoTimed(pbm proto.Message, duration time.Duration) (*proto.Message, error) {
	return xipc.WaitForProtoTimed(mqs, pbm, duration)
}

func (mqs *MqRequester) WaitForMqResponse() (*xipc.MqResponse, error) {
	return xipc.WaitForMqResponse(mqs)
}

func (mqs *MqRequester) WaitForMqResponseTimed(duration time.Duration) (*xipc.MqResponse, error) {
	return xipc.WaitForMqResponseTimed(mqs, duration)
}

func (mqs *MqRequester) HasErrors() bool {
	return mqs.ErrRqst != nil
}

func (mqs *MqRequester) Error() error {
	return mqs.ErrRqst
}

func (mqs *MqRequester) CloseRequester() error {
	mqs.MqRqst.Close()
	return nil
}
