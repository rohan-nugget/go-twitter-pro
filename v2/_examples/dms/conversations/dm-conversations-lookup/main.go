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
	In order to run, the user will need to provide the bearer token.
	This example retrieves DM conversations for the authenticated user.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	opts := twitter.DMConversationOpts{
		DMConversationFields: []twitter.DMConversationField{
			twitter.DMConversationFieldID,
			twitter.DMConversationFieldParticipantIDs,
			twitter.DMConversationFieldCreatedAt,
		},
		UserFields: []twitter.UserField{twitter.UserFieldUserName, twitter.UserFieldName},
		MaxResults: 25,
	}

	fmt.Println("Callout to DM conversations lookup")

	dmResponse, err := client.DMConversations(context.Background(), opts)
	if err != nil {
		log.Panicf("DM conversations error: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(dmResponse); err != nil {
		log.Panic(err)
	}
}