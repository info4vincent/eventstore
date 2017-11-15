package commands

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const msgCollectionName = "weekplan"

type msgData struct {
	_id     string
	Userid  string
	Eventid string
	Action  string
	Data    string
}

func GetMessageOfUserForEvent(session *mgo.Session, user string, eventId string) msgData {
	c := session.DB(DBName).C(msgCollectionName)

	var results []msgData
	err := c.Find(nil).All(&results)
	err = c.Find(bson.M{"userid": user, "eventid": eventId}).All(&results)

	log.Println("aantal results:", len(results))
	// actionUri := fmt.Sprintf("Say:Hoi, ik heb geen message gevonden voor user:%v en eventId", user, eventId)
	if err == nil {
		if len(results) > 0 {
			log.Println("Messagefound found returning details now.")
			return results[0]
		} else {
			log.Println("Could not find the user or event..")
			return msgData{user, user, eventId, "say", "user en of event niet gevonden."}
		}
	} else if err != nil {
		log.Println("Could not find user or event")
		return msgData{user, user, eventId, "say", "user en of event niet gevonden."}
	}

	return results[0]
}

type MessageofthedayCommand struct {
	cards map[string][]byte
}

func NewMessageofthedayCommand() *MessageofthedayCommand {
	return &MessageofthedayCommand{
		cards: map[string][]byte{},
	}
}

func (n *MessageofthedayCommand) HandleCommand(event string) string {

	if !strings.HasPrefix(event, "Say:") {
		return ""
	}

	// fmt.Println("Messageoftheday command executed with data:", event)

	// session = connectDB()
	// defer session.Close()

	// eventData := strings.TrimPrefix(event, "Messageoftheday:")
	// eventDataFields := strings.Fields(eventData)

	// msgoftheDayData := GetMessageOfUserForEvent(session, eventDataFields[0], eventDataFields[1])

	// fmt.Println(msgoftheDayData.Data)

	// fmt.Println(msgoftheDayData.Action)
	return event
}

func (n *MessageofthedayCommand) Type() string {
	return "MessageofthedayCommand"
}
