package commands

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CardUID struct {
	_id       string
	Uid       string
	User      string
	Actionuri string
}

const DBName = "mybluemarvin"
const CollectionName = "CardUID"

var session *mgo.Session

func connectDB() *mgo.Session {
	uri := os.Getenv("MONGODB_URL")
	if uri == "" {
		fmt.Println("No connection string provided - set MONGODB_URL = mongodb://{user}:{password}@mongodb.documents.azure.com:{port}")
		os.Exit(1)
	}
	uri = strings.TrimSuffix(uri, "?ssl=true")

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	dialInfo, err := mgo.ParseURL(uri)

	if err != nil {
		fmt.Println("Failed to parse URI: ", err)
		os.Exit(1)
	}

	maxWait := time.Duration(5 * time.Second)
	dialInfo.Timeout = maxWait

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Failed to connect: ", err)
		os.Exit(1)
	}

	dbnames, err := session.DB("").CollectionNames()
	if err != nil {
		fmt.Println("Couldn't query for collections names: ", err)
		os.Exit(1)
	}

	fmt.Println(dbnames)

	return session
}

func insertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('-')
		}
	}
	return buffer.String()
}

func GetCommandForCard(session *mgo.Session, uid string) CardUID {
	c := session.DB(DBName).C(CollectionName)

	var results []CardUID
	err := c.Find(nil).All(&results)
	err = c.Find(bson.M{"uid": uid}).All(&results)

	log.Println("aantal results:", len(results))
	readAbleUid := insertNth(uid, 1)
	actionUri := fmt.Sprintf("Say:Hoi, deze kaart ken ik niet. Het nummer van de kaart is:%v", readAbleUid)
	if err == nil {
		if len(results) > 0 {
			log.Println("Card found returning details now.")
			return results[0]
		} else {
			log.Println("Could not find card")
			return CardUID{uid, uid, uid, actionUri}
		}
	} else if err != nil {
		log.Println("Could not find card")
		return CardUID{uid, uid, uid, actionUri}
	}

	return results[0]
}

type CardScannedCommand struct {
	cards map[string][]byte
}

func NewCardScannedCommand() *CardScannedCommand {
	return &CardScannedCommand{
		cards: map[string][]byte{},
	}
}

func (n *CardScannedCommand) HandleCommand(event string) string {

	if !strings.HasPrefix(event, "CardScanned:") {
		return ""
	}

	fmt.Println("CardScanned command executed with card:", event)

	session = connectDB()
	defer session.Close()

	dayPlan := GetCommandForCard(session, strings.TrimPrefix(event, "CardScanned:"))

	fmt.Println(dayPlan.Actionuri)
	return dayPlan.Actionuri
}

func (n *CardScannedCommand) Type() string {
	return "CardScannedCommand"
}
