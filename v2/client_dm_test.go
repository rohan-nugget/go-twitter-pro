package twitter

import (
	"context"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestClient_DMEvents(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		opts DMEventOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DMEventResponse
		wantErr bool
	}{
		{
			name: "Success - Get DM Events",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(dmEventsEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), dmEventsEndpoint)
					}
					body := `{
						"data": [
							{
								"id": "29515892301193216-1697637602605010945-29515892301193216",
								"text": "Hello World!",
								"event_type": "MessageCreate",
								"created_at": "2023-09-01T16:27:44.000Z",
								"sender_id": "29515892301193216",
								"dm_conversation_id": "29515892301193216-1697637602605010945"
							}
						],
						"includes": {
							"users": [
								{
									"id": "29515892301193216",
									"name": "Test User",
									"username": "testuser"
								}
							]
						},
						"meta": {
							"result_count": 1,
							"next_token": "7140dibdnow9c7btw3z403c0tkjkbdtlu8m8cg5r96gq2"
						}
					}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header: func() http.Header {
							h := http.Header{}
							h.Set("Content-Type", "application/json")
							return h
						}(),
					}
				}),
			},
			args: args{
				opts: DMEventOpts{
					MaxResults: 10,
				},
			},
			want: &DMEventResponse{
				Data: []DMEvent{
					{
						ID:               "29515892301193216-1697637602605010945-29515892301193216",
						Text:             "Hello World!",
						EventType:        "MessageCreate",
						SenderID:         "29515892301193216",
						DMConversationID: "29515892301193216-1697637602605010945",
					},
				},
				Includes: &DMIncludes{
					Users: []UserObj{
						{
							ID:       "29515892301193216",
							Name:     "Test User",
							UserName: "testuser",
						},
					},
				},
				Meta: &DMEventMeta{
					ResultCount: 1,
					NextToken:   "7140dibdnow9c7btw3z403c0tkjkbdtlu8m8cg5r96gq2",
				},
			},
			wantErr: false,
		},
		{
			name: "Error - Too Many Results",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client:     mockHTTPClient(func(req *http.Request) *http.Response { return nil }),
			},
			args: args{
				opts: DMEventOpts{
					MaxResults: 101, // Greater than dmEventMaxResults
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.DMEvents(context.Background(), tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DMEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			// Compare everything except timestamps which are correctly parsed
			if len(got.Data) != len(tt.want.Data) {
				t.Errorf("Client.DMEvents() Data length = %v, want %v", len(got.Data), len(tt.want.Data))
			} else if len(got.Data) > 0 {
				if got.Data[0].ID != tt.want.Data[0].ID ||
					got.Data[0].Text != tt.want.Data[0].Text ||
					got.Data[0].EventType != tt.want.Data[0].EventType ||
					got.Data[0].SenderID != tt.want.Data[0].SenderID ||
					got.Data[0].DMConversationID != tt.want.Data[0].DMConversationID {
					t.Errorf("Client.DMEvents() Data fields mismatch")
				}
			}
			if !reflect.DeepEqual(got.Meta, tt.want.Meta) {
				t.Errorf("Client.DMEvents() Meta = %v, want %v", got.Meta, tt.want.Meta)
			}
		})
	}
}

func TestClient_DMConversations(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		opts DMConversationOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DMConversationResponse
		wantErr bool
	}{
		{
			name: "Success - Get DM Conversations",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(dmConversationsEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), dmConversationsEndpoint)
					}
					body := `{
						"data": [
							{
								"id": "29515892301193216-1697637602605010945",
								"participant_ids": ["29515892301193216", "1697637602605010945"],
								"created_at": "2023-09-01T16:27:44.000Z"
							}
						],
						"includes": {
							"users": [
								{
									"id": "29515892301193216",
									"name": "Test User",
									"username": "testuser"
								}
							]
						},
						"meta": {
							"result_count": 1
						}
					}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header: func() http.Header {
							h := http.Header{}
							h.Set("Content-Type", "application/json")
							return h
						}(),
					}
				}),
			},
			args: args{
				opts: DMConversationOpts{
					MaxResults: 10,
				},
			},
			want: &DMConversationResponse{
				Data: []DMConversation{
					{
						ID:             "29515892301193216-1697637602605010945",
						ParticipantIDs: []string{"29515892301193216", "1697637602605010945"},
					},
				},
				Meta: &DMConversationMeta{
					ResultCount: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.DMConversations(context.Background(), tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DMConversations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			// Compare everything except timestamps which are correctly parsed
			if len(got.Data) != len(tt.want.Data) {
				t.Errorf("Client.DMConversations() Data length = %v, want %v", len(got.Data), len(tt.want.Data))
			} else if len(got.Data) > 0 {
				if got.Data[0].ID != tt.want.Data[0].ID ||
					!reflect.DeepEqual(got.Data[0].ParticipantIDs, tt.want.Data[0].ParticipantIDs) {
					t.Errorf("Client.DMConversations() Data fields mismatch")
				}
			}
			if !reflect.DeepEqual(got.Meta, tt.want.Meta) {
				t.Errorf("Client.DMConversations() Meta = %v, want %v", got.Meta, tt.want.Meta)
			}
		})
	}
}

func TestClient_CreateDMConversation(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		req CreateDMConversationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateDMConversationResponse
		wantErr bool
	}{
		{
			name: "Success - Create DM Conversation",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), string(dmConversationCreateEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), dmConversationCreateEndpoint)
					}
					body := `{
						"data": {
							"dm_conversation_id": "29515892301193216-1697637602605010945",
							"dm_event_id": "29515892301193216-1697637602605010945-29515892301193216"
						}
					}`
					return &http.Response{
						StatusCode: http.StatusCreated,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header: func() http.Header {
							h := http.Header{}
							h.Set("Content-Type", "application/json")
							return h
						}(),
					}
				}),
			},
			args: args{
				req: CreateDMConversationRequest{
					ConversationType: "Group",
					ParticipantIDs:   []string{"1697637602605010945"},
					Text:             "Hello World!",
				},
			},
			want: &CreateDMConversationResponse{
				Data: &CreateDMConversationData{
					DMConversationID: "29515892301193216-1697637602605010945",
					DMEventID:        "29515892301193216-1697637602605010945-29515892301193216",
				},
			},
			wantErr: false,
		},
		{
			name: "Error - Missing Participant IDs",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client:     mockHTTPClient(func(req *http.Request) *http.Response { return nil }),
			},
			args: args{
				req: CreateDMConversationRequest{
					ConversationType: "Group",
					ParticipantIDs:   []string{}, // Empty
					Text:             "Hello World!",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error - Missing Text and MediaID",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client:     mockHTTPClient(func(req *http.Request) *http.Response { return nil }),
			},
			args: args{
				req: CreateDMConversationRequest{
					ConversationType: "Group",
					ParticipantIDs:   []string{"1697637602605010945"},
					// Missing both Text and MediaID
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.CreateDMConversation(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateDMConversation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("Client.CreateDMConversation() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestClient_SendDM(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		conversationID string
		req            SendDMRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SendDMResponse
		wantErr bool
	}{
		{
			name: "Success - Send DM",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if !strings.Contains(req.URL.String(), "dm_conversations") {
						log.Panicf("the url is not correct %s", req.URL.String())
					}
					body := `{
						"data": {
							"dm_event_id": "29515892301193216-1697637602605010945-29515892301193216",
							"dm_conversation_id": "29515892301193216-1697637602605010945"
						}
					}`
					return &http.Response{
						StatusCode: http.StatusCreated,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header: func() http.Header {
							h := http.Header{}
							h.Set("Content-Type", "application/json")
							return h
						}(),
					}
				}),
			},
			args: args{
				conversationID: "29515892301193216-1697637602605010945",
				req: SendDMRequest{
					Text: "Reply message",
				},
			},
			want: &SendDMResponse{
				Data: &SendDMData{
					DMEventID:        "29515892301193216-1697637602605010945-29515892301193216",
					DMConversationID: "29515892301193216-1697637602605010945",
				},
			},
			wantErr: false,
		},
		{
			name: "Error - Missing Text and MediaID",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client:     mockHTTPClient(func(req *http.Request) *http.Response { return nil }),
			},
			args: args{
				conversationID: "29515892301193216-1697637602605010945",
				req:            SendDMRequest{}, // Empty request
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.SendDM(context.Background(), tt.args.conversationID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SendDM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("Client.SendDM() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestCreateDMConversationRequest_validate(t *testing.T) {
	tests := []struct {
		name    string
		r       CreateDMConversationRequest
		wantErr bool
	}{
		{
			name: "Valid - With Text",
			r: CreateDMConversationRequest{
				ParticipantIDs: []string{"123456"},
				Text:           "Hello World!",
			},
			wantErr: false,
		},
		{
			name: "Valid - With MediaID",
			r: CreateDMConversationRequest{
				ParticipantIDs: []string{"123456"},
				MediaID:        "media123",
			},
			wantErr: false,
		},
		{
			name: "Invalid - No Participants",
			r: CreateDMConversationRequest{
				ParticipantIDs: []string{},
				Text:           "Hello World!",
			},
			wantErr: true,
		},
		{
			name: "Invalid - No Text or MediaID",
			r: CreateDMConversationRequest{
				ParticipantIDs: []string{"123456"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateDMConversationRequest.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendDMRequest_validate(t *testing.T) {
	tests := []struct {
		name    string
		r       SendDMRequest
		wantErr bool
	}{
		{
			name: "Valid - With Text",
			r: SendDMRequest{
				Text: "Hello World!",
			},
			wantErr: false,
		},
		{
			name: "Valid - With MediaID",
			r: SendDMRequest{
				MediaID: "media123",
			},
			wantErr: false,
		},
		{
			name: "Invalid - No Text or MediaID",
			r:       SendDMRequest{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.validate(); (err != nil) != tt.wantErr {
				t.Errorf("SendDMRequest.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SendDMByParticipantID(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		participantID string
		req           SendDMByParticipantRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SendDMByParticipantResponse
		wantErr bool
	}{
		{
			name: "Success - Send DM by Participant ID",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if !strings.Contains(req.URL.String(), "dm_conversations/by/participant_id") {
						log.Panicf("the url is not correct %s", req.URL.String())
					}
					body := `{
						"data": {
							"dm_event_id": "29515892301193216-1697637602605010945-29515892301193216",
							"dm_conversation_id": "29515892301193216-1697637602605010945"
						}
					}`
					return &http.Response{
						StatusCode: http.StatusCreated,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header: func() http.Header {
							h := http.Header{}
							h.Set("Content-Type", "application/json")
							return h
						}(),
					}
				}),
			},
			args: args{
				participantID: "1697637602605010945",
				req: SendDMByParticipantRequest{
					Text: "Direct message to participant",
				},
			},
			want: &SendDMByParticipantResponse{
				Data: &SendDMByParticipantData{
					DMEventID:        "29515892301193216-1697637602605010945-29515892301193216",
					DMConversationID: "29515892301193216-1697637602605010945",
				},
			},
			wantErr: false,
		},
		{
			name: "Error - Missing Text and MediaID",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client:     mockHTTPClient(func(req *http.Request) *http.Response { return nil }),
			},
			args: args{
				participantID: "1697637602605010945",
				req:           SendDMByParticipantRequest{}, // Empty request
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.SendDMByParticipantID(context.Background(), tt.args.participantID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SendDMByParticipantID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("Client.SendDMByParticipantID() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestSendDMByParticipantRequest_validate(t *testing.T) {
	tests := []struct {
		name    string
		r       SendDMByParticipantRequest
		wantErr bool
	}{
		{
			name: "Valid - With Text",
			r: SendDMByParticipantRequest{
				Text: "Hello World!",
			},
			wantErr: false,
		},
		{
			name: "Valid - With MediaID",
			r: SendDMByParticipantRequest{
				MediaID: "media123",
			},
			wantErr: false,
		},
		{
			name: "Invalid - No Text or MediaID",
			r:       SendDMByParticipantRequest{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.validate(); (err != nil) != tt.wantErr {
				t.Errorf("SendDMByParticipantRequest.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}