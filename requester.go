package xipc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joe-at-startupmedia/xipc/protos"
	"google.golang.org/protobuf/proto"
	"time"
)

type IRequester interface {
	Read() ([]byte, error)
	ReadTimed(time.Duration) ([]byte, error)
	Write([]byte) error
	Request([]byte) error
	RequestUsingRequest(req *Request) error
	RequestUsingProto(req *proto.Message) error
	WaitForProto(pbm proto.Message) (*proto.Message, error)
	WaitForProtoTimed(pbm proto.Message, duration time.Duration) (*proto.Message, error)
	WaitForResponseProto() (*Response, error)
	WaitForResponseProtoTimed(duration time.Duration) (*Response, error)
	HasErrors() bool
	Error() error
	CloseRequester() error
}

type Request protos.Request

func (mqr *Request) HasId() bool {
	return len(mqr.Id) > 0
}

func (mqr *Request) SetId() {
	mqr.Id = uuid.NewString()
}

// AsProtobuf used to convert the local type equivalent (Request)
// back to its protobuf instance
func (mqr *Request) AsProtobuf() *protos.Request {
	return (*protos.Request)(mqr)
}

// ProtoRequestToRequest used to convert the protobuf to the local
// type equivalent (Request) for leveraging instance methods
func ProtoRequestToRequest(mqr *protos.Request) *Request {
	return (*Request)(mqr)
}

func WaitForResponseProto(mqs IRequester) (*Response, error) {
	Resp := &protos.Response{}
	_, err := mqs.WaitForProto(Resp)
	if err != nil {
		return nil, err
	}
	return ProtoResponseToResponse(Resp), err
}

func WaitForResponseProtoTimed(mqs IRequester, duration time.Duration) (*Response, error) {
	Resp := &protos.Response{}
	_, err := mqs.WaitForProtoTimed(Resp, duration)
	if err != nil {
		return nil, err
	}
	return ProtoResponseToResponse(Resp), err
}

func RequestUsingRequest(mqs IRequester, req *Request) error {
	if !req.HasId() {
		req.SetId()
	}
	pbm := proto.Message(req.AsProtobuf())
	return mqs.RequestUsingProto(&pbm)
}

func RequestUsingProto(mqs IRequester, req *proto.Message) error {
	data, err := proto.Marshal(*req)
	if err != nil {
		return fmt.Errorf("marshaling error: %w", err)
	}
	return mqs.Request(data)
}

func WaitForProto(mqs IRequester, pbm proto.Message) (*proto.Message, error) {
	data, err := mqs.Read()
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(data, pbm)
	return &pbm, err
}

func WaitForProtoTimed(mqs IRequester, pbm proto.Message, duration time.Duration) (*proto.Message, error) {
	data, err := mqs.ReadTimed(duration)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(data, pbm)
	return &pbm, err
}
