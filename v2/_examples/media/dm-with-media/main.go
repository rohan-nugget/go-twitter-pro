package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	In order to run, the user will need to provide the bearer token, file path, and participant IDs.
	This example uploads media and creates a DM conversation with the uploaded media.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	filePath := flag.String("file", "", "path to media file")
	participants := flag.String("participants", "", "comma-separated participant user IDs")
	text := flag.String("text", "", "optional message text to accompany the media")
	flag.Parse()

	if *filePath == "" || *participants == "" {
		log.Fatal("file path and participant IDs are required")
	}

	// Open the file
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	// Step 1: Upload media
	mediaType := detectMediaType(*filePath)
	
	mediaReq := twitter.MediaUploadRequest{
		Media:         file,
		MediaCategory: twitter.MediaCategoryDMImage,
		MediaType:     mediaType,
	}

	fmt.Printf("Step 1: Uploading media file: %s\n", *filePath)

	mediaResponse, err := client.UploadMedia(context.Background(), mediaReq)
	if err != nil {
		log.Panicf("Media upload error: %v", err)
	}

	fmt.Printf("Media uploaded successfully! ID: %s\n", mediaResponse.Data.ID)

	// Step 2: Create DM conversation with media
	dmReq := twitter.CreateDMConversationRequest{
		ConversationType: "Group",
		ParticipantIDs:   strings.Split(*participants, ","),
		MediaID:          mediaResponse.Data.ID,
	}

	// Add text if provided
	if *text != "" {
		dmReq.Text = *text
	}

	fmt.Printf("Step 2: Creating DM conversation with media\n")

	dmResponse, err := client.CreateDMConversation(context.Background(), dmReq)
	if err != nil {
		log.Panicf("Create DM conversation error: %v", err)
	}

	fmt.Printf("DM conversation created successfully!\n")

	// Output the result
	result := map[string]interface{}{
		"media_upload":    mediaResponse,
		"dm_conversation": dmResponse,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(result); err != nil {
		log.Panic(err)
	}
}

// detectMediaType detects media type from file extension
func detectMediaType(filePath string) twitter.MediaType {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return twitter.MediaTypeImageJPEG
	case ".png":
		return twitter.MediaTypeImagePNG
	case ".webp":
		return twitter.MediaTypeImageWebP
	case ".bmp":
		return twitter.MediaTypeImageBMP
	case ".tiff", ".tif":
		return twitter.MediaTypeImageTIFF
	default:
		return twitter.MediaTypeImageJPEG // Default fallback
	}
}