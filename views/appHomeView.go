package views

import (
	"embed"
	"encoding/json"
	"github.com/slack-go/slack"
	"log"
)

var appHomeAssets embed.FS

func AppHomeTabView() slack.HomeTabViewRequest {
	str, err := appHomeAssets.ReadFile("assets/AppHomeView.json")
	if err != nil {
		log.Printf("Unable to read view `AppHomeView`: %v", err)
	}
	view := slack.HomeTabViewRequest{}
	json.Unmarshal([]byte(str), &view)

	log.Printf("-----> view page\n%v", view)

	return view
}
