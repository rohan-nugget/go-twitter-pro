package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

/**
	In order to run, the user will need to provide the bearer token, participant IDs, and message text.
	This example creates a new DM conversation.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	participants := flag.String("participants", "", "comma-separated participant user IDs")
	text := flag.String("text", "", "message text")
	conversationType := flag.String("type", "Group", "conversation type (Group or OneToOne)")
	flag.Parse()

	if *participants == "" || *text == "" {
		log.Fatal("participants and text are required")
	}

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	req := twitter.CreateDMConversationRequest{
		ConversationType: *conversationType,
		ParticipantIDs:   strings.Split(*participants, ","),
		Text:             *text,
	}

	fmt.Println("Callout to create DM conversation")

	dmResponse, err := client.CreateDMConversation(context.Background(), req)
	if err != nil {
		log.Panicf("Create DM conversation error: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(dmResponse); err != nil {
		log.Panic(err)
	}
}