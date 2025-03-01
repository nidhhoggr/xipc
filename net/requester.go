package net

import (
	"github.com/nidhhoggr/gipc"
	"github.com/nidhhoggr/xipc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"time"
)

type Requester struct {
	Rqst    *gipc.Client
	ErrRqst error
	Logger  *logrus.Logger
}

func NewRequester(config *QueueConfig) xipc.IRequester {

	logger := xipc.InitLogging(config.LogLevel)

	requester, errRqst := gipc.StartClient(&gipc.ClientConfig{
		Name:       config.Name,
		LogLevel:   config.LogLevel,
		RetryTimer: config.ClientRetryTimer,
		Timeout:    config.ClientTimeout,
		Encryption: false,
	})

	var mqs xipc.IRequester = &Requester{
		requester,
		errRqst,
		logger,
	}

	return mqs
}

func (mqs *Requester) Read() ([]byte, error) {
	msg, err := mqs.Rqst.Read()
	if msg.MsgType < 1 {
		time.Sleep(xipc.REQUEST_RECURSION_WAITTIME * time.Millisecond)
		return mqs.Read()
	} else {
		return msg.Data, err
	}
}

func (mqs *Requester) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, err := mqs.Rqst.ReadTimed(duration)
	if err != nil {
		return nil, err
	}
	if msg.MsgType < 1 {
		time.Sleep(xipc.REQUEST_RECURSION_WAITTIME * time.Millisecond)
		return mqs.ReadTimed(duration)
	} else if msg == gipc.TimeoutMessage {
		return nil, gipc.TimeoutMessage.Err
	} else {
		return msg.Data, err
	}
}

func (mqs *Requester) Write(data []byte) error {
	return mqs.Rqst.Write(DEFAULT_MSG_TYPE, data)
}

func (mqs *Requester) Request(data []byte) error {
	return mqs.Write(data)
}

func (mqs *Requester) RequestUsingRequest(req *xipc.Request) error {
	return xipc.RequestUsingRequest(mqs, req)
}

func (mqs *Requester) RequestUsingProto(req *proto.Message) error {
	return xipc.RequestUsingProto(mqs, req)
}

func (mqs *Requester) WaitForProto(pbm proto.Message) (*proto.Message, error) {
	return xipc.WaitForProto(mqs, pbm)
}

func (mqs *Requester) WaitForProtoTimed(pbm proto.Message, duration time.Duration) (*proto.Message, error) {
	return xipc.WaitForProtoTimed(mqs, pbm, duration)
}

func (mqs *Requester) WaitForResponseProto() (*xipc.Response, error) {
	return xipc.WaitForResponseProto(mqs)
}

func (mqs *Requester) WaitForResponseProtoTimed(duration time.Duration) (*xipc.Response, error) {
	return xipc.WaitForResponseProtoTimed(mqs, duration)
}

func (mqs *Requester) HasErrors() bool {
	return mqs.ErrRqst != nil
}

func (mqs *Requester) Error() error {
	return mqs.ErrRqst
}

func (mqs *Requester) CloseRequester() error {
	mqs.Rqst.Close()
	return nil
}
