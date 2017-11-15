package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"log"
)

func main() {
	//  Socket to talk to server
	fmt.Println("Connecting to eventsource server...")
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	err := requester.Connect("tcp://localhost:5555")
	if err != nil {
		log.Println("Failed to connect..", err)
	}

	// send hello
	msg := fmt.Sprintf("CardScanned:%s", "000001")
	fmt.Println("Sending ", msg)
	_, err = requester.Send(msg, 0)
	if err != nil {
		log.Println("Failed to send..", err)
	}

	// Wait for reply:
	reply, _ := requester.Recv(0)
	fmt.Println("Received ", reply)
}
