package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// MediaCategory represents the media use-case category
type MediaCategory string

const (
	// MediaCategoryTweetImage for images used in tweets
	MediaCategoryTweetImage MediaCategory = "tweet_image"
	// MediaCategoryDMImage for images used in direct messages
	MediaCategoryDMImage MediaCategory = "dm_image"
	// MediaCategorySubtitles for subtitle files
	MediaCategorySubtitles MediaCategory = "subtitles"
)

// MediaType represents supported media types
type MediaType string

const (
	// MediaTypeImageJPEG for JPEG images
	MediaTypeImageJPEG MediaType = "image/jpeg"
	// MediaTypeImagePNG for PNG images
	MediaTypeImagePNG MediaType = "image/png"
	// MediaTypeImageWebP for WebP images
	MediaTypeImageWebP MediaType = "image/webp"
	// MediaTypeImageBMP for BMP images
	MediaTypeImageBMP MediaType = "image/bmp"
	// MediaTypeImagePJPEG for Progressive JPEG images
	MediaTypeImagePJPEG MediaType = "image/pjpeg"
	// MediaTypeImageTIFF for TIFF images
	MediaTypeImageTIFF MediaType = "image/tiff"
	// MediaTypeTextSRT for SRT subtitle files
	MediaTypeTextSRT MediaType = "text/srt"
	// MediaTypeTextVTT for VTT subtitle files
	MediaTypeTextVTT MediaType = "text/vtt"
)

// ProcessingState represents the media processing state
type ProcessingState string

const (
	// ProcessingStateSucceeded indicates processing completed successfully
	ProcessingStateSucceeded ProcessingState = "succeeded"
	// ProcessingStateInProgress indicates processing is ongoing
	ProcessingStateInProgress ProcessingState = "in_progress"
	// ProcessingStatePending indicates processing is queued
	ProcessingStatePending ProcessingState = "pending"
	// ProcessingStateFailed indicates processing failed
	ProcessingStateFailed ProcessingState = "failed"
)

// MediaUploadRequest represents a request to upload media
type MediaUploadRequest struct {
	Media            io.Reader       `json:"-"`
	MediaCategory    MediaCategory   `json:"media_category"`
	AdditionalOwners []string        `json:"additional_owners,omitempty"`
	MediaType        MediaType       `json:"media_type,omitempty"`
	Shared           bool            `json:"shared,omitempty"`
}

// MediaUploadProcessingInfo represents processing information for uploaded media
type MediaUploadProcessingInfo struct {
	CheckAfterSecs   int             `json:"check_after_secs,omitempty"`
	ProgressPercent  int             `json:"progress_percent,omitempty"`
	State            ProcessingState `json:"state,omitempty"`
}

// MediaUploadData represents the data returned from media upload
type MediaUploadData struct {
	ID               string                     `json:"id"`
	MediaKey         string                     `json:"media_key"`
	ExpiresAfterSecs int                        `json:"expires_after_secs,omitempty"`
	ProcessingInfo   *MediaUploadProcessingInfo `json:"processing_info,omitempty"`
	Size             int                        `json:"size,omitempty"`
}

// MediaUploadResponse represents the response from media upload
type MediaUploadResponse struct {
	Data      *MediaUploadData `json:"data,omitempty"`
	Errors    []ErrorObj       `json:"errors,omitempty"`
	RateLimit *RateLimit       `json:"-"`
}

// UploadMedia uploads media to Twitter API v2
func (c *Client) UploadMedia(ctx context.Context, req MediaUploadRequest) (*MediaUploadResponse, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}

	// Create multipart form
	body, contentType, err := createMediaUploadForm(req)
	if err != nil {
		return nil, fmt.Errorf("media upload form creation error: %w", err)
	}

	ep := mediaUploadEndpoint.url(c.Host)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, body)
	if err != nil {
		return nil, fmt.Errorf("media upload request: %w", err)
	}

	httpReq.Header.Set("Content-Type", contentType)
	httpReq.Header.Set("Accept", "application/json")
	c.Authorizer.Add(httpReq)

	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("media upload response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &MediaUploadResponse{}
	if err := decoder.Decode(respBody); err != nil {
		return nil, fmt.Errorf("media upload decode: %w", err)
	}
	respBody.RateLimit = rl

	return respBody, nil
}

// validate validates the MediaUploadRequest
func (r MediaUploadRequest) validate() error {
	if r.Media == nil {
		return fmt.Errorf("media upload: media is required: %w", ErrParameter)
	}
	if r.MediaCategory == "" {
		return fmt.Errorf("media upload: media_category is required: %w", ErrParameter)
	}
	return nil
}

// createMediaUploadForm creates a multipart form for media upload
func createMediaUploadForm(req MediaUploadRequest) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add media file
	mediaWriter, err := writer.CreateFormFile("media", "media")
	if err != nil {
		return nil, "", fmt.Errorf("create media form file: %w", err)
	}

	if _, err := io.Copy(mediaWriter, req.Media); err != nil {
		return nil, "", fmt.Errorf("copy media to form: %w", err)
	}

	// Add media_category
	if err := writer.WriteField("media_category", string(req.MediaCategory)); err != nil {
		return nil, "", fmt.Errorf("write media_category field: %w", err)
	}

	// Add optional fields
	if req.MediaType != "" {
		if err := writer.WriteField("media_type", string(req.MediaType)); err != nil {
			return nil, "", fmt.Errorf("write media_type field: %w", err)
		}
	}

	if req.Shared {
		if err := writer.WriteField("shared", "true"); err != nil {
			return nil, "", fmt.Errorf("write shared field: %w", err)
		}
	}

	// Add additional_owners if present
	for i, owner := range req.AdditionalOwners {
		fieldName := fmt.Sprintf("additional_owners[%d]", i)
		if err := writer.WriteField(fieldName, owner); err != nil {
			return nil, "", fmt.Errorf("write additional_owners field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close multipart writer: %w", err)
	}

	return &buf, writer.FormDataContentType(), nil
}

// String returns the string representation of MediaCategory
func (m MediaCategory) String() string {
	return string(m)
}

// String returns the string representation of MediaType
func (m MediaType) String() string {
	return string(m)
}

// String returns the string representation of ProcessingState
func (p ProcessingState) String() string {
	return string(p)
}