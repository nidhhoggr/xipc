package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joe-at-startupmedia/xipc"
	protos2 "github.com/joe-at-startupmedia/xipc/example/protos"
	"github.com/joe-at-startupmedia/xipc/net"
	"google.golang.org/protobuf/proto"
	"log"
)

const maxRequestTickNum = 10

const queue_name = "goqr_example_protobuf"

var mqr xipc.IResponder
var mqs xipc.IRequester
var config = net.QueueConfig{
	Name:             queue_name,
	ClientTimeout:    0,
	ClientRetryTimer: 0,
}

func main() {
	resp_c := make(chan int)
	go responder(resp_c)
	//wait for the responder to create the posix_mq files
	xipc.Sleep()
	request_c := make(chan int)
	go requester(request_c)
	<-resp_c
	<-request_c

	mqs.CloseRequester()
	mqr.CloseResponder()
	//gives time for deferred functions to complete
	xipc.Sleep()
}

func responder(c chan int) {
	mqr = net.NewResponder(&config)
	defer func() {
		log.Println("Responder: finished")
		c <- 0
	}()

	if mqr.HasErrors() {
		log.Printf("Responder: could not initialize: %s", mqr.Error())
		c <- 1
		return
	}

	xipc.Sleep()

	count := 0
	for {
		count++
		if err := handleCmdRequest(mqr); err != nil {
			log.Printf("Responder: error handling request: %s\n", err)
			continue
		}

		log.Println("Responder: Sent a response")

		if count >= maxRequestTickNum {
			break
		}
	}
}

func requester(c chan int) {
	mqs = net.NewRequester(&config)
	defer func() {
		log.Println("Requester: finished and closed")
		c <- 0
	}()
	if mqs.HasErrors() {
		log.Printf("Requester: could not initialize: %s", mqs.Error())
		c <- 1
		return
	}
	xipc.Sleep()

	count := 0
	for {
		count++
		cmd := &protos2.Cmd{
			Name: "restart",
			Arg1: fmt.Sprintf("%d", count), //using count as the id of the process
			ExecFlags: &protos2.ExecFlags{
				User: "nonroot",
			},
		}
		if err := requestUsingCmd(mqs, cmd); err != nil {
			log.Printf("Requester: error requesting request: %s\n", err)
			continue
		}

		log.Printf("Requester: sent a new request: %s \n", cmd.String())

		cmdResp, err := waitForCmdResponse(mqs)

		if err != nil {
			log.Printf("Requester: error getting response: %s\n", err)
			continue
		}

		log.Printf("Requester: got a response: %s\n", cmdResp.ValueStr)

		if count >= maxRequestTickNum {
			break
		}
	}
}

func requestUsingCmd(mqs xipc.IRequester, req *protos2.Cmd) error {
	if len(req.Id) == 0 {
		req.Id = uuid.NewString()
	}
	pbm := proto.Message(req)
	return mqs.RequestUsingProto(&pbm)
}

func waitForCmdResponse(mqs xipc.IRequester) (*protos2.CmdResp, error) {
	Resp := &protos2.CmdResp{}
	_, err := mqs.WaitForProto(Resp)
	if err != nil {
		return nil, err
	}
	return Resp, err
}

// handleCmdRequest provides a concrete implementation of HandleRequestFromProto using the local Cmd protobuf type
func handleCmdRequest(mqr xipc.IResponder) error {

	cmd := &protos2.Cmd{}

	return mqr.HandleRequestFromProto(cmd, func() (processed []byte, err error) {

		cmdResp := protos2.CmdResp{}
		cmdResp.Id = cmd.Id
		cmdResp.ValueStr = fmt.Sprintf("I recieved request: %s(%s) - %s\n", cmd.Name, cmd.Id, cmd.Arg1)

		data, err := proto.Marshal(&cmdResp)
		if err != nil {
			return nil, fmt.Errorf("marshaling error: %w", err)
		}

		return data, nil
	})
}
