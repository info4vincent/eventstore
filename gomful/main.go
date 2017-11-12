package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	nxprd "github.com/hohlerde/go-nxprd"
	zmq "github.com/pebbe/zmq4"
)

func getTag(t nxprd.TagType) string {
	switch t {
	case nxprd.TagType1:
		return "1"
	case nxprd.TagType2:
		return "2"
	case nxprd.TagType3:
		return "3"
	case nxprd.TagType4A:
		return "4A"
	case nxprd.TagTypeP2P:
		return "P2P"
	case nxprd.TagTypeNFCDEP4A:
		return "NFC_DEP & 4A"
	default:
		return "Undefined"
	}
}

func getTech(t nxprd.TechType) string {
	switch t {
	case nxprd.TechA:
		return "A"
	case nxprd.TechB:
		return "B"
	case nxprd.TechF:
		return "F"
	case nxprd.TechV15693T5T:
		return "V / ISO 15693 / T5T"
	case nxprd.Tech18000p3m3EPCGen2:
		return "ISO 18000p3m3 / EPC Gen2"
	default:
		return "Undefined"
	}
}

func slice2Str(arr []byte) string {
	var buffer bytes.Buffer

	for i := 0; i < len(arr); i++ {
		buffer.WriteString(fmt.Sprintf("0x%02X ", arr[i]))
	}

	return strings.TrimSpace(buffer.String())
}

func slice2StrHex(arr []byte) string {
	var buffer bytes.Buffer

	for i := 0; i < len(arr); i++ {
		buffer.WriteString(fmt.Sprintf("%02X", arr[i]))
	}

	return strings.TrimSpace(buffer.String())
}

func printInfo(dev *nxprd.Device) {
	fmt.Printf("Card            : %s\n", dev.Params.DevType)
	fmt.Printf("Tag type        : %s\n", getTag(dev.Params.TagType))
	fmt.Printf("Technology type : %s\n", getTech(dev.Params.TechType))
	fmt.Printf("UID             : %s\n", slice2StrHex(dev.Params.UID))
	fmt.Printf("ATQ(A)          : %s\n", slice2Str(dev.Params.ATQ))
	fmt.Printf("SAK             : 0x%02X\n", dev.Params.SAK)
}

func publishCardScanned(uid string) {
	//  Socket to talk to server
	fmt.Println("Connecting to eventstore server...")
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://localhost:5555")

	// send hello
	msg := fmt.Sprintf("CardScanned:%s", uid)
	fmt.Println("Sending card scanned msg -> ", msg)
	requester.Send(msg, 0)

	// Wait for reply:
	reply, _ := requester.Recv(0)
	fmt.Println("Received ", reply)

	fmt.Println("Disconnecting from eventstore.")
}

func main() {

	fmt.Println("Mifare Ultralight Test")
	fmt.Println("\nPlace a card near the reader....")
	//fmt.Println("\nPlace a card near the reader and press <Enter>")
	//termReader := bufio.NewReader(os.Stdin)
	//termReader.ReadString('\n')

	// Initialize the library. We need to do that once.
	if err := nxprd.Init(); err != nil {
		fmt.Println("\nError: Initializing NXP library failed")
		fmt.Println(err)
		return
	}

	// In order to cleanup the C part of the wrapper DeInit need to be called.
	defer nxprd.DeInit()

	for {
		time.Sleep(1000 * time.Millisecond)
		// Try to detect/discover a card/tag for 1000ms. Discover will block.
		// 1000ms is the default timeout.
		dev, err := nxprd.Discover(40000)
		if err != nil {
			fmt.Println("\nCouldn't detect card")
			fmt.Println(err)
			//return
			continue
		}

		if dev.Params.DevType == nxprd.Unknown {
			// A card/tag could be detected, but the wrapper doesn't support it yet.
			// So we can't read or write blocks, but we can access some parameters.
			fmt.Println("\nFound an unknown card for data storage retrieval.")
			fmt.Println("")
			printInfo(dev)
			publishCardScanned(slice2StrHex(dev.Params.UID))
			//return
			continue
		}

		fmt.Println("")
		printInfo(dev)

		// The standard Mifare Ultralight card has a minimum of 15 blocks (pages),
		// each block consisting of 4 bytes.
		fmt.Println("\nTrying to read memory pages 0-15:")
		fmt.Println("")

		for i := 0; i < 16; i++ {
			fmt.Printf("Page %02d: ", i)
			// Read each block
			buffer, err := dev.Reader.ReadBlock(i)
			if err != nil {
				fmt.Println(err)
			} else {
				// and print the bytes in hex format.
				fmt.Println(slice2Str(buffer))
			}
		}
		publishCardScanned(slice2StrHex(dev.Params.UID))
	}
}
