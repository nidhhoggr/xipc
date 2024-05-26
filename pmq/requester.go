package pmq

import (
	"fmt"
	"github.com/joe-at-startupmedia/posix_mq"
	"github.com/joe-at-startupmedia/xipc"
	"github.com/joe-at-startupmedia/xipc/protos"
	"google.golang.org/protobuf/proto"
	"time"
)

type MqRequester BidirectionalQueue

func NewRequester(config *QueueConfig, owner *Ownership) *MqRequester {
	requester, errRqst := openQueueForRequester(config, owner, "rqst")

	responder, errResp := openQueueForRequester(config, owner, "resp")

	mqs := MqRequester{
		requester,
		errRqst,
		responder,
		errResp,
	}

	return &mqs
}

func openQueueForRequester(config *QueueConfig, owner *Ownership, postfix string) (*posix_mq.MessageQueue, error) {
	if config.Flags == 0 {
		config.Flags = O_RDWR
	}
	return NewMessageQueueWithOwnership(*config, owner, postfix)
}

func (mqs *MqRequester) Request(data []byte) error {
	return mqs.MqRqst.Send(data, 0)
}

func (mqs *MqRequester) RequestUsingMqRequest(req *xipc.MqRequest) error {
	if !req.HasId() {
		req.SetId()
	}
	pbm := proto.Message(req.AsProtobuf())
	return mqs.RequestUsingProto(&pbm)
}

func (mqs *MqRequester) RequestUsingProto(req *proto.Message) error {
	data, err := proto.Marshal(*req)
	if err != nil {
		return fmt.Errorf("marshaling error: %w", err)
	}
	return mqs.Request(data)
}

func (mqs *MqRequester) WaitForResponse() ([]byte, error) {
	msg, _, err := mqs.MqResp.Receive()
	return msg, err
}

func (mqs *MqRequester) WaitForResponseTimed(duration time.Duration) ([]byte, error) {
	msg, _, err := mqs.MqResp.TimedReceive(duration)
	return msg, err
}

func (mqs *MqRequester) WaitForProto(pbm proto.Message) (*proto.Message, error) {
	data, _, err := mqs.MqResp.Receive()
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(data, pbm)
	return &pbm, err
}

func (mqs *MqRequester) WaitForProtoTimed(pbm proto.Message, duration time.Duration) (*proto.Message, error) {
	data, _, err := mqs.MqResp.TimedReceive(duration)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(data, pbm)
	return &pbm, err
}

func (mqs *MqRequester) WaitForMqResponse() (*xipc.MqResponse, error) {
	mqResp := &protos.Response{}
	_, err := mqs.WaitForProto(mqResp)
	if err != nil {
		return nil, err
	}
	return xipc.ProtoResponseToMqResponse(mqResp), err
}

func (mqs *MqRequester) WaitForMqResponseTimed(duration time.Duration) (*xipc.MqResponse, error) {
	mqResp := &protos.Response{}
	_, err := mqs.WaitForProtoTimed(mqResp, duration)
	if err != nil {
		return nil, err
	}
	return xipc.ProtoResponseToMqResponse(mqResp), err
}

func (mqs *MqRequester) CloseRequester() error {
	return (*BidirectionalQueue)(mqs).Close()
}

func (mqs *MqRequester) UnlinkRequester() error {
	return (*BidirectionalQueue)(mqs).Unlink()
}

func (mqs *MqRequester) HasErrors() bool {
	return (*BidirectionalQueue)(mqs).HasErrors()
}

func (mqs *MqRequester) Error() error {
	return (*BidirectionalQueue)(mqs).Error()
}

func CloseRequester(mqr *MqRequester) error {
	if mqr != nil {
		return mqr.CloseRequester()
	}
	return fmt.Errorf("pointer reference is nil")
}

func UnlinkRequester(mqr *MqRequester) error {
	if mqr != nil {
		return mqr.UnlinkRequester()
	}
	return fmt.Errorf("pointer reference is nil")
}
