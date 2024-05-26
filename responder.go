package xipc

import (
	"fmt"
	"github.com/joe-at-startupmedia/xipc/protos"
	"google.golang.org/protobuf/proto"
	"time"
)

type IResponder interface {
	Read() ([]byte, error)
	ReadTimed(time.Duration) ([]byte, error)
	Write([]byte) error
	HandleRequest(ResponderCallback) error
	HandleRequestWithLag(ResponderCallback, int) error
	HandleRequestProto(ResponderRequestProtoCallback) error
	HandleRequestFromProto(proto.Message, ResponderFromProtoMessageCallback) error
	HasErrors() bool
	Error() error
	CloseResponder() error
}

type Response protos.Response

type ResponderCallback func(msq []byte) (processed []byte, err error)

type ResponderRequestProtoCallback func(mqs *Request) (mqr *Response, err error)

type ResponderFromProtoMessageCallback func() (processed []byte, err error)

// AsProtobuf used to convert the local type equivalent (Response)
// back to its protobuf instance
func (mqr *Response) AsProtobuf() *protos.Response {
	return (*protos.Response)(mqr)
}

func (mqr *Response) PrepareFromRequest(mqs *Request) *Response {
	mqr.RequestId = mqs.Id
	return mqr
}

// ProtoResponseToResponse used to convert the protobuf to the local
// type equivalent (Response) for leveraging instance methods
func ProtoResponseToResponse(mqr *protos.Response) *Response {
	return (*Response)(mqr)
}

// HandleRequestProto provides a concrete implementation of HandleRequestFromProto using the local Request type
func HandleRequestProto(mqr IResponder, requestProcessor ResponderRequestProtoCallback) error {

	mqReq := &protos.Request{}

	return mqr.HandleRequestFromProto(mqReq, func() (processed []byte, err error) {

		Resp, err := requestProcessor(ProtoRequestToRequest(mqReq))
		if err != nil {
			return nil, err
		}

		data, err := proto.Marshal(Resp.AsProtobuf())

		if err != nil {
			return nil, fmt.Errorf("marshaling error: %w", err)
		}

		return data, nil
	})
}

// HandleRequestWithLag used for testing purposes to simulate lagging responder
func HandleRequestWithLag(mqr IResponder, msgHandler ResponderCallback, lag int) error {

	msg, err := mqr.Read()
	if err != nil {
		return err
	}

	processed, err := msgHandler(msg)
	if err != nil {
		return err
	}

	if lag > 0 {
		time.Sleep(time.Duration(lag) * time.Second)
	}

	err = mqr.Write(processed)
	return err
}

// HandleRequestFromProto used to process arbitrary protobuf messages using a callback
func HandleRequestFromProto(mqr IResponder, protocMsg proto.Message, msgHandler ResponderFromProtoMessageCallback) error {

	msg, err := mqr.Read()
	if err != nil {
		return err
	}

	err = proto.Unmarshal(msg, protocMsg)
	if err != nil {
		return fmt.Errorf("unmarshaling error: %w", err)
	}

	processed, err := msgHandler()
	if err != nil {
		return err
	}

	return mqr.Write(processed)
}
