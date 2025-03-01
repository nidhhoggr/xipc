package mem

import (
	"github.com/nidhhoggr/shmemipc"
	"github.com/nidhhoggr/xipc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"time"
)

type Requester struct {
	Rqst    *shmemipc.IpcRequester
	ErrRqst error
	logger  *logrus.Logger
}

func NewRequester(config *QueueConfig) xipc.IRequester {

	logger := xipc.InitLogging(config.LogLevel)

	requester := shmemipc.NewRequester(config.BasePath + config.Name)

	var mqs xipc.IRequester = &Requester{
		Rqst:    requester,
		ErrRqst: requester.GetError(),
		logger:  logger,
	}

	return mqs
}

func (mqs *Requester) Read() ([]byte, error) {
	return mqs.Rqst.Read()
}

func (mqs *Requester) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, err := mqs.Rqst.ReadTimed(duration)
	if err != nil {
		if err.Error() == "timed_out" {
			mqs.logger.Info("requester timed out")
			time.Sleep(xipc.REQUEST_RECURSION_WAITTIME * time.Millisecond)
			return mqs.ReadTimed(duration)
		}
		return nil, err
	}
	return msg, nil
}

func (mqs *Requester) Write(data []byte) error {
	return mqs.Rqst.Write(data)
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
	return mqs.Rqst.Close()
}
