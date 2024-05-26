package pmq

import (
	"github.com/joe-at-startupmedia/posix_mq"
	"github.com/joe-at-startupmedia/xipc"
	"google.golang.org/protobuf/proto"
	"time"
)

type MqRequester BidirectionalQueue

func NewRequester(config *QueueConfig, owner *Ownership) xipc.IMqRequester {

	requester, errRqst := openQueueForRequester(config, owner, "rqst")
	responder, errResp := openQueueForRequester(config, owner, "resp")

	var mqs xipc.IMqRequester = &MqRequester{
		requester,
		errRqst,
		responder,
		errResp,
	}

	return mqs
}

func openQueueForRequester(config *QueueConfig, owner *Ownership, postfix string) (*posix_mq.MessageQueue, error) {
	if config.Flags == 0 {
		config.Flags = O_RDWR
	}
	return NewMessageQueueWithOwnership(*config, owner, postfix)
}

func (mqs *MqRequester) Read() ([]byte, error) {
	data, _, err := mqs.MqResp.Receive()
	return data, err
}

func (mqs *MqRequester) ReadTimed(duration time.Duration) ([]byte, error) {
	msg, _, err := mqs.MqResp.TimedReceive(duration)
	return msg, err
}

func (mqs *MqRequester) Write(data []byte) error {
	return mqs.MqRqst.Send(data, 0)
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
	return (*BidirectionalQueue)(mqs).HasErrors()
}

func (mqs *MqRequester) Error() error {
	return (*BidirectionalQueue)(mqs).Error()
}

func (mqs *MqRequester) CloseRequester() error {
	return (*BidirectionalQueue)(mqs).Close()
}
