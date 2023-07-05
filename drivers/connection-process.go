package drivers

import (
	"github.com/slack-go/slack/socketmode"
	"log"
)

func MiddlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	log.Print("Connecting to Slack with Socket Mode...")
}

func MiddlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	log.Print("Connection failed. Retrying later...")
}

func MiddlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	log.Print("Connected to Slack with Socket Mode.")
}

func MiddlewareHello(evt *socketmode.Event, client *socketmode.Client) {
	log.Print("Hello ! Slack.")
}
