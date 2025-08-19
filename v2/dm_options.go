package twitter

import (
	"net/url"
	"strconv"
	"strings"
)

// DMEventField represents the fields that can be requested for DM events
type DMEventField string

const (
	DMEventFieldID               DMEventField = "id"
	DMEventFieldText             DMEventField = "text"
	DMEventFieldEventType        DMEventField = "event_type"
	DMEventFieldCreatedAt        DMEventField = "created_at"
	DMEventFieldSenderID         DMEventField = "sender_id"
	DMEventFieldDMConversationID DMEventField = "dm_conversation_id"
	DMEventFieldReferencedTweet  DMEventField = "referenced_tweets"
	DMEventFieldMediaKeys        DMEventField = "media_keys"
	DMEventFieldAttachments      DMEventField = "attachments"
	DMEventFieldEntities         DMEventField = "entities"
	DMEventFieldParticipantIDs   DMEventField = "participant_ids"
)

// DMConversationField represents the fields that can be requested for DM conversations
type DMConversationField string

const (
	DMConversationFieldID             DMConversationField = "id"
	DMConversationFieldParticipantIDs DMConversationField = "participant_ids"
	DMConversationFieldCreatedAt      DMConversationField = "created_at"
)

// DMEventOpts represents the options for fetching DM events
type DMEventOpts struct {
	DMEventFields        []DMEventField        `json:"dm_event.fields,omitempty"`
	DMConversationFields []DMConversationField `json:"dm_conversation.fields,omitempty"`
	UserFields           []UserField           `json:"user.fields,omitempty"`
	TweetFields          []TweetField          `json:"tweet.fields,omitempty"`
	MediaFields          []MediaField          `json:"media.fields,omitempty"`
	Expansions           []Expansion           `json:"expansions,omitempty"`
	MaxResults           int                   `json:"max_results,omitempty"`
	NextToken            string                `json:"next_token,omitempty"`
	PreviousToken        string                `json:"previous_token,omitempty"`
	SinceID              string                `json:"since_id,omitempty"`
	UntilID              string                `json:"until_id,omitempty"`
}

// DMConversationOpts represents the options for fetching DM conversations
type DMConversationOpts struct {
	DMConversationFields []DMConversationField `json:"dm_conversation.fields,omitempty"`
	UserFields           []UserField           `json:"user.fields,omitempty"`
	Expansions           []Expansion           `json:"expansions,omitempty"`
	MaxResults           int                   `json:"max_results,omitempty"`
	NextToken            string                `json:"next_token,omitempty"`
	PreviousToken        string                `json:"previous_token,omitempty"`
}

func (d DMEventOpts) addQuery(req *url.URL) {
	q := req.Query()
	if len(d.DMEventFields) > 0 {
		q.Add("dm_event.fields", dmEventFieldsToString(d.DMEventFields))
	}
	if len(d.DMConversationFields) > 0 {
		q.Add("dm_conversation.fields", dmConversationFieldsToString(d.DMConversationFields))
	}
	if len(d.UserFields) > 0 {
		var fields []string
		for _, field := range d.UserFields {
			fields = append(fields, string(field))
		}
		q.Add("user.fields", strings.Join(fields, ","))
	}
	if len(d.TweetFields) > 0 {
		var fields []string
		for _, field := range d.TweetFields {
			fields = append(fields, string(field))
		}
		q.Add("tweet.fields", strings.Join(fields, ","))
	}
	if len(d.MediaFields) > 0 {
		var fields []string
		for _, field := range d.MediaFields {
			fields = append(fields, string(field))
		}
		q.Add("media.fields", strings.Join(fields, ","))
	}
	if len(d.Expansions) > 0 {
		var fields []string
		for _, field := range d.Expansions {
			fields = append(fields, string(field))
		}
		q.Add("expansions", strings.Join(fields, ","))
	}
	if d.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(d.MaxResults))
	}
	if len(d.NextToken) > 0 {
		q.Add("next_token", d.NextToken)
	}
	if len(d.PreviousToken) > 0 {
		q.Add("previous_token", d.PreviousToken)
	}
	if len(d.SinceID) > 0 {
		q.Add("since_id", d.SinceID)
	}
	if len(d.UntilID) > 0 {
		q.Add("until_id", d.UntilID)
	}
	req.RawQuery = q.Encode()
}

func (d DMConversationOpts) addQuery(req *url.URL) {
	q := req.Query()
	if len(d.DMConversationFields) > 0 {
		q.Add("dm_conversation.fields", dmConversationFieldsToString(d.DMConversationFields))
	}
	if len(d.UserFields) > 0 {
		var fields []string
		for _, field := range d.UserFields {
			fields = append(fields, string(field))
		}
		q.Add("user.fields", strings.Join(fields, ","))
	}
	if len(d.Expansions) > 0 {
		var fields []string
		for _, field := range d.Expansions {
			fields = append(fields, string(field))
		}
		q.Add("expansions", strings.Join(fields, ","))
	}
	if d.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(d.MaxResults))
	}
	if len(d.NextToken) > 0 {
		q.Add("next_token", d.NextToken)
	}
	if len(d.PreviousToken) > 0 {
		q.Add("previous_token", d.PreviousToken)
	}
	req.RawQuery = q.Encode()
}

// dmEventFieldsToString converts a slice of DMEventField to a comma-separated string
func dmEventFieldsToString(fields []DMEventField) string {
	if len(fields) == 0 {
		return ""
	}
	var result []string
	for _, field := range fields {
		result = append(result, string(field))
	}
	return strings.Join(result, ",")
}

// dmConversationFieldsToString converts a slice of DMConversationField to a comma-separated string
func dmConversationFieldsToString(fields []DMConversationField) string {
	if len(fields) == 0 {
		return ""
	}
	var result []string
	for _, field := range fields {
		result = append(result, string(field))
	}
	return strings.Join(result, ",")
}

func (d DMEventField) String() string {
	return string(d)
}

func (d DMConversationField) String() string {
	return string(d)
}
