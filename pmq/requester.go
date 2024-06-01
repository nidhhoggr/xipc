package pmq

import (
	"github.com/joe-at-startupmedia/posix_mq"
	"github.com/joe-at-startupmedia/xipc"
	"google.golang.org/protobuf/proto"
	"time"
)

type Requester BidirectionalQueue

func NewRequester(config *QueueConfig, owner *xipc.Ownership) xipc.IRequester {

	requester, errRqst := openQueueForRequester(config, owner, "rqst")
	responder, errResp := openQueueForRequester(config, owner, "resp")

	var mqs xipc.IRequester = &Requester{
		requester,
		errRqst,
		responder,
		errResp,
	}

	return mqs
}

func openQueueForRequester(config *QueueConfig, owner *xipc.Ownership, postfix string) (*posix_mq.MessageQueue, error) {
	if config.Flags == 0 {
		config.Flags = O_RDWR
	}
	return NewMessageQueueWithOwnership(*config, owner, postfix)
}

func (mqs *Requester) Read() ([]byte, error) {
	data, _, err := mqs.Resp.Receive()
	return data, err
}

func (mqs *Requester) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, _, err := mqs.Resp.TimedReceive(duration)
	return msg, err
}

func (mqs *Requester) Write(data []byte) error {
	return mqs.Rqst.Send(data, 0)
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
	return (*BidirectionalQueue)(mqs).HasErrors()
}

func (mqs *Requester) Error() error {
	return (*BidirectionalQueue)(mqs).Error()
}

func (mqs *Requester) CloseRequester() error {
	return (*BidirectionalQueue)(mqs).Close()
}
