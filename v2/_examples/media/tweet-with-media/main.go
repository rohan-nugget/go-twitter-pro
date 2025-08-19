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
	In order to run, the user will need to provide the bearer token, file path, and tweet text.
	This example uploads media and creates a tweet with the uploaded media.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	filePath := flag.String("file", "", "path to media file")
	text := flag.String("text", "", "tweet text")
	flag.Parse()

	if *filePath == "" || *text == "" {
		log.Fatal("file path and tweet text are required")
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
		MediaCategory: twitter.MediaCategoryTweetImage,
		MediaType:     mediaType,
	}

	fmt.Printf("Step 1: Uploading media file: %s\n", *filePath)

	mediaResponse, err := client.UploadMedia(context.Background(), mediaReq)
	if err != nil {
		log.Panicf("Media upload error: %v", err)
	}

	fmt.Printf("Media uploaded successfully! ID: %s\n", mediaResponse.Data.ID)

	// Step 2: Create tweet with media
	tweetReq := twitter.CreateTweetRequest{
		Text: *text,
		Media: &twitter.CreateTweetMedia{
			IDs: []string{mediaResponse.Data.ID},
		},
	}

	fmt.Printf("Step 2: Creating tweet with media\n")

	tweetResponse, err := client.CreateTweet(context.Background(), tweetReq)
	if err != nil {
		log.Panicf("Create tweet error: %v", err)
	}

	fmt.Printf("Tweet created successfully!\n")

	// Output the result
	result := map[string]interface{}{
		"media_upload": mediaResponse,
		"tweet":        tweetResponse,
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