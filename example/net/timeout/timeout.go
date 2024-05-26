package main

import (
	"errors"
	"fmt"
	"github.com/joe-at-startupmedia/gipc"
	"github.com/joe-at-startupmedia/xipc"
	"github.com/joe-at-startupmedia/xipc/net"
	"log"
	"time"
)

const maxRequestTickNum = 10

const queue_name = "goqr_example_timeout"

var mqr xipc.IMqResponder
var mqs xipc.IMqRequester
var config = net.QueueConfig{
	Name:             queue_name,
	ClientTimeout:    time.Second * 10,
	ClientRetryTimer: time.Second * 1,
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

	mqr.CloseResponder()
	mqs.CloseRequester()
	//gives time for deferred functions to complete
	xipc.Sleep()
}

func responder(c chan int) {

	mqr = net.NewResponder(&config)

	defer func() {
		c <- 0
		log.Println("Responder: finished")
	}()
	if mqr.HasErrors() {
		log.Printf("Responder: could not initialize: %s", mqr.Error())
		c <- 1
		return
	}

	xipc.Sleep()

	count := 0
	for {
		//time.Sleep(1 * time.Second)
		count++
		var err error
		if count > 5 {
			err = mqr.HandleRequestWithLag(handleMessage, count-4)
		} else {
			err = mqr.HandleRequest(handleMessage)
		}

		if err != nil {
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
		c <- 0
		log.Println("Requester: finished and closed")
	}()
	if mqs.HasErrors() {
		log.Printf("Requester: could not initialize: %s", mqs.Error())
		c <- 1
		return
	}
	xipc.Sleep()

	count := 0
	ch := make(chan goqResponse, 13)
	for {
		count++
		request := fmt.Sprintf("Hello, World : %d\n", count)
		go requestResponse(mqs, request, ch)

		if count >= maxRequestTickNum {
			break
		}
		time.Sleep(1 * time.Second)
	}

	result := make([]goqResponse, maxRequestTickNum)
	for i := range result {
		result[i] = <-ch
		if result[i].status {
			log.Println(result[i].response)
		} else {
			log.Printf("Requester: Got error: %s \n", result[i].response)
		}
	}

}

func requestResponse(mqs xipc.IMqRequester, msg string, c chan goqResponse) {

	if len(msg) > 0 {
		err := mqs.Request([]byte(msg))
		if err != nil {
			c <- goqResponse{fmt.Sprintf("%s", err), false}
			return
		}
		log.Printf("Requester: sent a new request: %s", msg)
	}

	resp, err := mqs.WaitForResponseTimed(time.Second * 5)

	if err != nil {

		if errors.Is(err, gipc.TimeoutMessage.Err) {
			log.Printf("Requester: requestResponse timedout, using recursion")
			go requestResponse(mqs, "", c)
			return
		}

		c <- goqResponse{fmt.Sprintf("%s", err), false}
		return
	}

	c <- goqResponse{fmt.Sprintf("Requester: got a response: %s\n", resp), true}
}

type goqResponse struct {
	response string
	status   bool
}

func handleMessage(request []byte) (processed []byte, err error) {
	return []byte(fmt.Sprintf("I recieved request: %s\n", request)), nil
}
