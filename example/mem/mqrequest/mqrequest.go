package main

import (
	"fmt"
	"github.com/nidhhoggr/xipc"
	"github.com/nidhhoggr/xipc/mem"
	"log"
	"time"
)

const maxRequestTickNum = 10

const queue_name = "mem_example_Request"

var config = &mem.QueueConfig{
	Name:       queue_name,
	BasePath:   "/tmp/",
	MaxMsgSize: 1024,
}

func main() {
	resp_c := make(chan int)
	go responder(resp_c)
	//wait for the responder to create the posix_mq files
	time.Sleep(1 * time.Second)
	request_c := make(chan int)
	go requester(request_c)
	<-resp_c
	<-request_c
	//gives time for deferred functions to complete
	time.Sleep(2 * time.Second)
}

func responder(c chan int) {
	mqr := mem.NewResponder(config, nil)
	defer func() {
		mqr.CloseResponder()
		fmt.Println("Responder: finished and unlinked")
		c <- 0
	}()
	if mqr.HasErrors() {
		log.Printf("Responder: could not initialize: %s", mqr.Error())
		c <- 1
		return
	}

	count := 0
	for {
		count++
		if err := mqr.HandleRequestProto(requestProcessor); err != nil {
			fmt.Printf("Responder: error handling request: %s\n", err)
			continue
		}

		fmt.Println("Responder: Sent a response")

		if count >= maxRequestTickNum {
			break
		}
	}
}

func requester(c chan int) {
	mqs := mem.NewRequester(config)
	defer func() {
		mqs.CloseRequester()
		fmt.Println("Requester: finished and closed")
		c <- 0
	}()
	if mqs.HasErrors() {
		log.Printf("Requester: could not initialize: %s", mqs.Error())
		c <- 1
		return
	}

	count := 0
	for {
		count++
		request := fmt.Sprintf("Hello, World : %d\n", count)
		if err := mqs.RequestUsingRequest(&xipc.Request{
			Arg1: request,
		}); err != nil {
			fmt.Printf("Requester: error requesting request: %s\n", err)
			continue
		}

		fmt.Printf("Requester: sent a new request: %s", request)

		msg, err := mqs.WaitForResponseProto()

		if err != nil {
			fmt.Printf("Requester: error getting response: %s\n", err)
			continue
		}

		fmt.Printf("Requester: got a response: %s\n", msg.ValueStr)
		//fmt.Printf("Requester: got a response: %-v\n", msg)

		if count >= maxRequestTickNum {
			break
		}
	}
}

func requestProcessor(request *xipc.Request) (*xipc.Response, error) {
	response := xipc.Response{}
	//assigns the response.request_id
	response.PrepareFromRequest(request)
	response.ValueStr = fmt.Sprintf("I recieved request: %s\n", request.Arg1)
	return &response, nil
}
