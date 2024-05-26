package xipc

import (
	"fmt"
	"github.com/joe-at-startupmedia/xipc/protos"
	"google.golang.org/protobuf/proto"
	"time"
)

type IMqResponder interface {
	Read() ([]byte, error)
	ReadTimed(time.Duration) ([]byte, error)
	Write([]byte) error
	HandleRequest(ResponderCallback) error
	HandleRequestWithLag(ResponderCallback, int) error
	HandleMqRequest(ResponderMqRequestCallback) error
	HandleRequestFromProto(proto.Message, ResponderFromProtoMessageCallback) error
	HasErrors() bool
	Error() error
	CloseResponder() error
}

type MqResponse protos.Response

type ResponderCallback func(msq []byte) (processed []byte, err error)

type ResponderMqRequestCallback func(mqs *MqRequest) (mqr *MqResponse, err error)

type ResponderFromProtoMessageCallback func() (processed []byte, err error)

// AsProtobuf used to convert the local type equivalent (MqResponse)
// back to its protobuf instance
func (mqr *MqResponse) AsProtobuf() *protos.Response {
	return (*protos.Response)(mqr)
}

func (mqr *MqResponse) PrepareFromRequest(mqs *MqRequest) *MqResponse {
	mqr.RequestId = mqs.Id
	return mqr
}

// ProtoResponseToMqResponse used to convert the protobuf to the local
// type equivalent (MqResponse) for leveraging instance methods
func ProtoResponseToMqResponse(mqr *protos.Response) *MqResponse {
	return (*MqResponse)(mqr)
}

// HandleMqRequest provides a concrete implementation of HandleRequestFromProto using the local MqRequest type
func HandleMqRequest(mqr IMqResponder, requestProcessor ResponderMqRequestCallback) error {

	mqReq := &protos.Request{}

	return mqr.HandleRequestFromProto(mqReq, func() (processed []byte, err error) {

		mqResp, err := requestProcessor(ProtoRequestToMqRequest(mqReq))
		if err != nil {
			return nil, err
		}

		data, err := proto.Marshal(mqResp.AsProtobuf())

		if err != nil {
			return nil, fmt.Errorf("marshaling error: %w", err)
		}

		return data, nil
	})
}

// HandleRequestWithLag used for testing purposes to simulate lagging responder
func HandleRequestWithLag(mqr IMqResponder, msgHandler ResponderCallback, lag int) error {

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
func HandleRequestFromProto(mqr IMqResponder, protocMsg proto.Message, msgHandler ResponderFromProtoMessageCallback) error {

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
