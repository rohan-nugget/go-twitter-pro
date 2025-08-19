package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

/**
	In order to run, the user will need to provide the bearer token, participant ID, and message text.
	This example sends a DM directly to a participant by their ID without needing an existing conversation.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	participantID := flag.String("participant", "", "participant user ID")
	text := flag.String("text", "", "message text")
	mediaID := flag.String("media", "", "media ID (optional)")
	flag.Parse()

	if *participantID == "" || (*text == "" && *mediaID == "") {
		log.Fatal("participant ID and either text or media ID are required")
	}

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	req := twitter.SendDMByParticipantRequest{
		Text:    *text,
		MediaID: *mediaID,
	}

	fmt.Println("Callout to send DM by participant ID")

	dmResponse, err := client.SendDMByParticipantID(context.Background(), *participantID, req)
	if err != nil {
		log.Panicf("Send DM by participant error: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(dmResponse); err != nil {
		log.Panic(err)
	}
}