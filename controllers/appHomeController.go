package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go_slack/views"
	"log"
	"reflect"
)

type AppHomeController struct {
	EventHandler *socketmode.SocketmodeHandler
}

func AppHomeTest() string {
	return "Hello Test"
}

func NewAppHomeController(eventhandler *socketmode.SocketmodeHandler) AppHomeController {
	c := AppHomeController{
		EventHandler: eventhandler,
	}

	log.Println("----> AppHomeOpened Handler")
	c.EventHandler.HandleEvents(
		slackevents.AppHomeOpened,
		c.publishHomeTabView,
	)
	log.Println("----> return")
	return c
}

// HomeEvent is a helper struct for recovery from err bad message
type HomeEvent struct {
	Envelop string          `json:"envelop_id"`
	Payload json.RawMessage `json:"payload"`
}

func (c *AppHomeController) recoverAppHomeOpened(evt *socketmode.Event, clt *socketmode.Client) {
	log.Printf("Attempt to recover Bad Message %v", evt)

	var e *socketmode.ErrorBadMessage
	var ok bool
	var err error

	if e, ok = evt.Data.(*socketmode.ErrorBadMessage); !ok {
		log.Printf("Bad Message Not Cast: %+v, evt")
		return
	}
	var rawBytes []byte
	if rawBytes, err = e.Message.MarshalJSON(); err != nil {
		log.Printf("Bad Message Not Marshalled. Err: %+v\n Event: %+v", err, evt)
		return
	}

	rawMessage := bytes.Replace(rawBytes, []byte{34, 115, 116, 97, 116, 101, 34, 58, 123, 34, 118, 97, 108, 117, 101, 115, 34, 58, 91, 93, 125}, []byte{34, 115, 116, 97, 116, 101, 34, 58, 123, 34, 118, 97, 108, 117, 101, 115, 34, 58, 123, 125, 125}, 1)
	var hE HomeEvent
	if err := json.Unmarshal(rawMessage, &hE); err != nil {
		log.Printf("Raw Message Not Marshalled: %s", err)
		return
	}

	var newEvent slackevents.EventsAPIEvent
	if newEvent, err = slackevents.ParseEvent(hE.Payload, slackevents.OptionNoVerifyToken()); err != nil {
		log.Printf("Bad Message Not Parsed. Err: %+v\n Inner JSON: %+v", err, rawMessage)
		return
	}

	fabEvent := socketmode.Event{
		Type: socketmode.EventTypeEventsAPI,
		Data: newEvent,
		Request: &socketmode.Request{
			Type:       "events_api",
			EnvelopeID: hE.Envelop,
		},
	}

	c.EventHandler.Client.Events <- fabEvent
}

func (c *AppHomeController) publishHomeTabView(evt *socketmode.Event, clt *socketmode.Client) {
	log.Println("----> publishHomeTabView")
	evtApi, ok := evt.Data.(slackevents.EventsAPIEvent)

	if ok != true {
		log.Printf("ERROR converting event to slackevents.EventsAPIEvent")
	}

	evtAppHomeOpened, ok := evtApi.InnerEvent.Data.(slackevents.AppHomeOpenedEvent)

	var user string

	if ok != true {
		log.Printf("ERROR converting inner event to slackevents.AppHomeOpenedEvent")
		user = reflect.ValueOf(evtApi.InnerEvent.Data).Elem().FieldByName("User").Interface().(string)
	} else {
		user = evtAppHomeOpened.User
	}

	log.Printf("ERROR AppHomeOpenedEvent: %v", evtAppHomeOpened)

	log.Println("----> AppHomeTabView")
	view := views.AppHomeTabView()

	_, err := clt.PublishView(user, view, "")

	if err != nil {
		log.Printf("Error publishHomeTabView: %v", err)
	}
}
