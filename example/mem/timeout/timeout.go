package main

import (
	"fmt"
	"github.com/joe-at-startupmedia/xipc"
	"github.com/joe-at-startupmedia/xipc/mem"
	"log"
	"time"
)

const maxRequestTickNum = 10

const queue_name = "mem_example_timeout"

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
	requester(resp_c)
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
		log.Fatalf("Responder: could not initialize: %s", mqr.Error())
	}

	count := 0
	for {
		count++
		var err error
		if count > 5 {
			err = mqr.HandleRequestWithLag(handleMessage, count-4)
		} else {
			err = mqr.HandleRequest(handleMessage)
		}

		if err != nil {
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
	}()
	if mqs.HasErrors() {
		log.Fatalf("Requester: could not initialize: %s", mqs.Error())
	}

	count := 0
	ch := make(chan pResponse)
	for {
		count++
		request := fmt.Sprintf("Hello, World : %d\n", count)
		go requestResponse(mqs, request, ch)

		if count >= maxRequestTickNum {
			break
		}
	}

	result := make([]pResponse, maxRequestTickNum)
	for i := range result {
		result[i] = <-ch
		if result[i].status {
			fmt.Println(result[i].response)
		} else {
			fmt.Printf("Requester: Got error: %s \n", result[i].response)
		}
	}
	<-c
}

func requestResponse(mqs xipc.IRequester, msg string, c chan pResponse) {
	if err := mqs.Request([]byte(msg)); err != nil {
		c <- pResponse{fmt.Sprintf("%s", err), false}
		return
	}
	fmt.Printf("Requester: sent a new request: %s", msg)

	resp, err := mqs.ReadTimed(time.Second * 10)

	if err != nil {
		c <- pResponse{fmt.Sprintf("%s", err), false}
		return
	}

	c <- pResponse{fmt.Sprintf("Requester: got a response: %s\n", resp), true}
}

type pResponse struct {
	response string
	status   bool
}

func handleMessage(request []byte) (processed []byte, err error) {
	return []byte(fmt.Sprintf("I recieved request: %s\n", request)), nil
}
