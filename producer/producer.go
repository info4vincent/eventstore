package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
)

func main() {
	//  Socket to talk to server
	fmt.Println("Connecting to eventsource server...")
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://localhost:5555")

	// send hello
	msg := fmt.Sprintf("CardScanned:%s", "000001")
	fmt.Println("Sending ", msg)
	requester.Send(msg, 0)

	// Wait for reply:
	reply, _ := requester.Recv(0)
	fmt.Println("Received ", reply)
}
