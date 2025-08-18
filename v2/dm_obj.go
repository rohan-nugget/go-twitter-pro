package twitter

import "time"

// DMEventType represents the type of DM event
type DMEventType string

const (
	MessageCreate         DMEventType = "MessageCreate"
	ParticipantsJoin      DMEventType = "ParticipantsJoin" 
	ParticipantsLeave     DMEventType = "ParticipantsLeave"
)

// DMEvent represents a direct message event
type DMEvent struct {
	ID                 string              `json:"id"`
	Text               string              `json:"text,omitempty"`
	EventType          string              `json:"event_type"`
	CreatedAt          time.Time           `json:"created_at,omitempty"`
	SenderID           string              `json:"sender_id,omitempty"`
	DMConversationID   string              `json:"dm_conversation_id"`
	ReferencedTweet    *DMReferencedTweet  `json:"referenced_tweet,omitempty"`
	MediaKeys          []string            `json:"media_keys,omitempty"`
	Attachments        *DMAttachments      `json:"attachments,omitempty"`
}

// DMReferencedTweet represents a referenced tweet in a DM
type DMReferencedTweet struct {
	ID string `json:"id"`
}

// DMAttachments represents attachments in a DM
type DMAttachments struct {
	MediaKeys []string `json:"media_keys,omitempty"`
}

// DMConversation represents a DM conversation
type DMConversation struct {
	ID              string    `json:"id"`
	ParticipantIDs  []string  `json:"participant_ids,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

// DMEventResponse represents the response when fetching DM events
type DMEventResponse struct {
	Data      []DMEvent       `json:"data,omitempty"`
	Includes  *DMIncludes     `json:"includes,omitempty"`
	Meta      *DMEventMeta    `json:"meta,omitempty"`
	RateLimit *RateLimit      `json:"-"`
}

// DMConversationResponse represents the response when fetching DM conversations
type DMConversationResponse struct {
	Data      []DMConversation `json:"data,omitempty"`
	Includes  *DMIncludes      `json:"includes,omitempty"`
	Meta      *DMConversationMeta `json:"meta,omitempty"`
	RateLimit *RateLimit       `json:"-"`
}

// DMIncludes represents the includes in DM responses
type DMIncludes struct {
	Users  []UserObj  `json:"users,omitempty"`
	Tweets []TweetObj `json:"tweets,omitempty"`
	Media  []MediaObj `json:"media,omitempty"`
}

// DMEventMeta represents metadata for DM event responses
type DMEventMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token,omitempty"`
	PreviousToken string `json:"previous_token,omitempty"`
}

// DMConversationMeta represents metadata for DM conversation responses
type DMConversationMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token,omitempty"`
	PreviousToken string `json:"previous_token,omitempty"`
}

// CreateDMConversationRequest represents a request to create a new DM conversation
type CreateDMConversationRequest struct {
	ConversationType string   `json:"conversation_type"`
	ParticipantIDs   []string `json:"participant_ids"`
	Text             string   `json:"text,omitempty"`
	MediaID          string   `json:"media_id,omitempty"`
}

// CreateDMConversationResponse represents the response when creating a DM conversation
type CreateDMConversationResponse struct {
	Data      *CreateDMConversationData `json:"data,omitempty"`
	RateLimit *RateLimit                `json:"-"`
}

// CreateDMConversationData represents the data returned when creating a DM conversation
type CreateDMConversationData struct {
	DMConversationID string `json:"dm_conversation_id"`
	DMEventID        string `json:"dm_event_id"`
}

// SendDMRequest represents a request to send a DM to an existing conversation
type SendDMRequest struct {
	Text    string `json:"text,omitempty"`
	MediaID string `json:"media_id,omitempty"`
}

// SendDMResponse represents the response when sending a DM
type SendDMResponse struct {
	Data      *SendDMData `json:"data,omitempty"`
	RateLimit *RateLimit  `json:"-"`
}

// SendDMData represents the data returned when sending a DM
type SendDMData struct {
	DMEventID        string `json:"dm_event_id"`
	DMConversationID string `json:"dm_conversation_id"`
}