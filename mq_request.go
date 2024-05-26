package xipc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joe-at-startupmedia/xipc/protos"
	"google.golang.org/protobuf/proto"
	"time"
)

type IMqRequester interface {
	Read() ([]byte, error)
	ReadTimed(time.Duration) ([]byte, error)
	Write([]byte) error
	Request([]byte) error
	RequestUsingMqRequest(req *MqRequest) error
	RequestUsingProto(req *proto.Message) error
	WaitForResponse() ([]byte, error)
	WaitForResponseTimed(duration time.Duration) ([]byte, error)
	WaitForProto(pbm proto.Message) (*proto.Message, error)
	WaitForProtoTimed(pbm proto.Message, duration time.Duration) (*proto.Message, error)
	WaitForMqResponse() (*MqResponse, error)
	WaitForMqResponseTimed(duration time.Duration) (*MqResponse, error)
	HasErrors() bool
	Error() error
	CloseRequester() error
}

type MqRequest protos.Request

func (mqr *MqRequest) HasId() bool {
	return len(mqr.Id) > 0
}

func (mqr *MqRequest) SetId() {
	mqr.Id = uuid.NewString()
}

// AsProtobuf used to convert the local type equivalent (MqRequest)
// back to its protobuf instance
func (mqr *MqRequest) AsProtobuf() *protos.Request {
	return (*protos.Request)(mqr)
}

// ProtoRequestToMqRequest used to convert the protobuf to the local
// type equivalent (MqRequest) for leveraging instance methods
func ProtoRequestToMqRequest(mqr *protos.Request) *MqRequest {
	return (*MqRequest)(mqr)
}

func WaitForMqResponse(mqs IMqRequester) (*MqResponse, error) {
	mqResp := &protos.Response{}
	_, err := mqs.WaitForProto(mqResp)
	if err != nil {
		return nil, err
	}
	return ProtoResponseToMqResponse(mqResp), err
}

func WaitForMqResponseTimed(mqs IMqRequester, duration time.Duration) (*MqResponse, error) {
	mqResp := &protos.Response{}
	_, err := mqs.WaitForProtoTimed(mqResp, duration)
	if err != nil {
		return nil, err
	}
	return ProtoResponseToMqResponse(mqResp), err
}

func RequestUsingMqRequest(mqs IMqRequester, req *MqRequest) error {
	if !req.HasId() {
		req.SetId()
	}
	pbm := proto.Message(req.AsProtobuf())
	return mqs.RequestUsingProto(&pbm)
}

func RequestUsingProto(mqs IMqRequester, req *proto.Message) error {
	data, err := proto.Marshal(*req)
	if err != nil {
		return fmt.Errorf("marshaling error: %w", err)
	}
	return mqs.Request(data)
}

func WaitForProto(mqs IMqRequester, pbm proto.Message) (*proto.Message, error) {
	data, err := mqs.Read()
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(data, pbm)
	return &pbm, err
}

func WaitForProtoTimed(mqs IMqRequester, pbm proto.Message, duration time.Duration) (*proto.Message, error) {
	data, err := mqs.ReadTimed(duration)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(data, pbm)
	return &pbm, err
}
