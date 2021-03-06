package main

import (
	zmq "github.com/pebbe/zmq4"

	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/info4vincent/eventstore/commands"

	"github.com/boltdb/bolt"
)

const boltdbname string = "myevents.db"
const boltEventsourceBucketName string = "MyEventsBucket"

func initDb() {
	fmt.Println("Initializing bolt db.")
	db, err := bolt.Open(boltdbname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(boltEventsourceBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}
func storeInDb(event string) {
	db, err := bolt.Open(boltdbname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//db.Update(func(tx *bolt.Tx) error {
	//		b := tx.Bucket([]byte("MyEventsBucket"))
	//		err := b.Put([]byte(event), []byte("42"))
	//		return err
	//	})

	u := EventSource{}
	u.Event = event
	u.Timestamp = time.Now().UTC().String()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltEventsourceBucketName))

		// Generate ID for the event.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		u.ID = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(u.ID), buf)
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	//	b := make([]byte, 8)
	//	binary.BigEndian.PutUint64(b, uint64(v))
	//	return b
	return []byte(strconv.Itoa(v))
}

type EventSource struct {
	ID        int
	Event     string
	Timestamp string
}

func main() {
	initDb()
	fmt.Println("starting eventstore....")
	responder, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		log.Fatal(err)
	}

	defer responder.Close()
	err = responder.Bind("tcp://*:5555")

	if err != nil {
		log.Fatal(err)
	}

	cardHandler := commands.NewCardScannedCommand()
	msgHandler := commands.NewMessageofthedayCommand()

	fmt.Println("Bind *:5555 succesful")

	//  Prepare our publisher
	publisher, _ := zmq.NewSocket(zmq.PUB)
	defer publisher.Close()
	publisher.Bind("tcp://*:5556")
	publisher.Bind("ipc://weather.ipc")

	// Wait for messages
	for {
		pubRec, _ := publisher.Recv(0)
		msg, _ := responder.Recv(0)
		if msg == "" && pubRec == "" {
			continue
		}
		println("Received :", string(msg))
		storeInDb(msg)

		println("Received publish:", string(pubRec))

		actionToSend := cardHandler.HandleCommand(msg)

		if actionToSend == "" {
			actionToSend = msgHandler.HandleCommand(msg)
		}
		if actionToSend == "" {
			responder.Send("No known action to take, continuing", 0)
			continue
		}

		reply := "eventsource received and broadcasting action."
		responder.Send(reply, 0)
		publisher.Send(actionToSend, 0)
	}
}
