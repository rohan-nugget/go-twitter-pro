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
	In order to run, the user will need to provide the bearer token and file path.
	This example uploads media to Twitter for use in tweets or DMs.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	filePath := flag.String("file", "", "path to media file")
	category := flag.String("category", "tweet_image", "media category (tweet_image, dm_image, subtitles)")
	shared := flag.Bool("shared", false, "whether media is shared")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("file path is required")
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

	// Detect media type from file extension
	mediaType := detectMediaType(*filePath)
	
	req := twitter.MediaUploadRequest{
		Media:         file,
		MediaCategory: twitter.MediaCategory(*category),
		MediaType:     mediaType,
		Shared:        *shared,
	}

	fmt.Printf("Uploading media file: %s\n", *filePath)
	fmt.Printf("Media type: %s\n", mediaType)
	fmt.Printf("Category: %s\n", *category)

	mediaResponse, err := client.UploadMedia(context.Background(), req)
	if err != nil {
		log.Panicf("Media upload error: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(mediaResponse); err != nil {
		log.Panic(err)
	}

	fmt.Printf("\nMedia uploaded successfully!\n")
	fmt.Printf("Media ID: %s\n", mediaResponse.Data.ID)
	fmt.Printf("Media Key: %s\n", mediaResponse.Data.MediaKey)
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
	case ".srt":
		return twitter.MediaTypeTextSRT
	case ".vtt":
		return twitter.MediaTypeTextVTT
	default:
		return twitter.MediaTypeImageJPEG // Default fallback
	}
}