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

func TestClient_UploadMedia(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		req MediaUploadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *MediaUploadResponse
		wantErr bool
	}{
		{
			name: "Success - Upload Image",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), string(mediaUploadEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), mediaUploadEndpoint)
					}
					body := `{
						"data": {
							"id": "1146654567674912769",
							"media_key": "3_1146654567674912769",
							"expires_after_secs": 86400,
							"size": 11065
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
				req: MediaUploadRequest{
					Media:         strings.NewReader("fake image data"),
					MediaCategory: MediaCategoryTweetImage,
					MediaType:     MediaTypeImageJPEG,
				},
			},
			want: &MediaUploadResponse{
				Data: &MediaUploadData{
					ID:               "1146654567674912769",
					MediaKey:         "3_1146654567674912769",
					ExpiresAfterSecs: 86400,
					Size:             11065,
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Upload DM Image with Processing Info",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					body := `{
						"data": {
							"id": "1146654567674912770",
							"media_key": "3_1146654567674912770",
							"expires_after_secs": 86400,
							"processing_info": {
								"state": "succeeded",
								"progress_percent": 100
							},
							"size": 25000
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
				req: MediaUploadRequest{
					Media:         strings.NewReader("fake image data for dm"),
					MediaCategory: MediaCategoryDMImage,
					MediaType:     MediaTypeImagePNG,
					Shared:        true,
				},
			},
			want: &MediaUploadResponse{
				Data: &MediaUploadData{
					ID:               "1146654567674912770",
					MediaKey:         "3_1146654567674912770",
					ExpiresAfterSecs: 86400,
					ProcessingInfo: &MediaUploadProcessingInfo{
						State:           ProcessingStateSucceeded,
						ProgressPercent: 100,
					},
					Size: 25000,
				},
			},
			wantErr: false,
		},
		{
			name: "Error - Missing Media",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client:     mockHTTPClient(func(req *http.Request) *http.Response { return nil }),
			},
			args: args{
				req: MediaUploadRequest{
					Media:         nil, // Missing media
					MediaCategory: MediaCategoryTweetImage,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error - Missing MediaCategory",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client:     mockHTTPClient(func(req *http.Request) *http.Response { return nil }),
			},
			args: args{
				req: MediaUploadRequest{
					Media:         strings.NewReader("fake image data"),
					MediaCategory: "", // Missing category
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success - Upload with Additional Owners",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					body := `{
						"data": {
							"id": "1146654567674912771",
							"media_key": "3_1146654567674912771",
							"expires_after_secs": 86400,
							"size": 15000
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
				req: MediaUploadRequest{
					Media:            strings.NewReader("fake image data"),
					MediaCategory:    MediaCategoryTweetImage,
					MediaType:        MediaTypeImageWebP,
					AdditionalOwners: []string{"123456789", "987654321"},
				},
			},
			want: &MediaUploadResponse{
				Data: &MediaUploadData{
					ID:               "1146654567674912771",
					MediaKey:         "3_1146654567674912771",
					ExpiresAfterSecs: 86400,
					Size:             15000,
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
			got, err := c.UploadMedia(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UploadMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("Client.UploadMedia() = %v, want %v", got.Data, tt.want.Data)
			}
		})
	}
}

func TestMediaUploadRequest_validate(t *testing.T) {
	tests := []struct {
		name    string
		r       MediaUploadRequest
		wantErr bool
	}{
		{
			name: "Valid - Tweet Image",
			r: MediaUploadRequest{
				Media:         strings.NewReader("fake image data"),
				MediaCategory: MediaCategoryTweetImage,
				MediaType:     MediaTypeImageJPEG,
			},
			wantErr: false,
		},
		{
			name: "Valid - DM Image",
			r: MediaUploadRequest{
				Media:         strings.NewReader("fake image data"),
				MediaCategory: MediaCategoryDMImage,
				MediaType:     MediaTypeImagePNG,
			},
			wantErr: false,
		},
		{
			name: "Valid - Subtitles",
			r: MediaUploadRequest{
				Media:         strings.NewReader("fake subtitle data"),
				MediaCategory: MediaCategorySubtitles,
				MediaType:     MediaTypeTextSRT,
			},
			wantErr: false,
		},
		{
			name: "Valid - With Additional Owners",
			r: MediaUploadRequest{
				Media:            strings.NewReader("fake image data"),
				MediaCategory:    MediaCategoryTweetImage,
				MediaType:        MediaTypeImageWebP,
				AdditionalOwners: []string{"123456789"},
				Shared:           true,
			},
			wantErr: false,
		},
		{
			name: "Invalid - No Media",
			r: MediaUploadRequest{
				Media:         nil,
				MediaCategory: MediaCategoryTweetImage,
			},
			wantErr: true,
		},
		{
			name: "Invalid - No MediaCategory",
			r: MediaUploadRequest{
				Media:         strings.NewReader("fake image data"),
				MediaCategory: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.validate(); (err != nil) != tt.wantErr {
				t.Errorf("MediaUploadRequest.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMediaCategory_String(t *testing.T) {
	tests := []struct {
		name string
		m    MediaCategory
		want string
	}{
		{
			name: "Tweet Image",
			m:    MediaCategoryTweetImage,
			want: "tweet_image",
		},
		{
			name: "DM Image",
			m:    MediaCategoryDMImage,
			want: "dm_image",
		},
		{
			name: "Subtitles",
			m:    MediaCategorySubtitles,
			want: "subtitles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("MediaCategory.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMediaType_String(t *testing.T) {
	tests := []struct {
		name string
		m    MediaType
		want string
	}{
		{
			name: "JPEG",
			m:    MediaTypeImageJPEG,
			want: "image/jpeg",
		},
		{
			name: "PNG",
			m:    MediaTypeImagePNG,
			want: "image/png",
		},
		{
			name: "WebP",
			m:    MediaTypeImageWebP,
			want: "image/webp",
		},
		{
			name: "SRT",
			m:    MediaTypeTextSRT,
			want: "text/srt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("MediaType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessingState_String(t *testing.T) {
	tests := []struct {
		name string
		p    ProcessingState
		want string
	}{
		{
			name: "Succeeded",
			p:    ProcessingStateSucceeded,
			want: "succeeded",
		},
		{
			name: "In Progress",
			p:    ProcessingStateInProgress,
			want: "in_progress",
		},
		{
			name: "Pending",
			p:    ProcessingStatePending,
			want: "pending",
		},
		{
			name: "Failed",
			p:    ProcessingStateFailed,
			want: "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("ProcessingState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}