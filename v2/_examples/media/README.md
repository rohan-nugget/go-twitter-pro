# Media Upload Examples

This directory contains examples for using the Twitter API v2 Media Upload endpoints.

## Examples

### Upload
- **upload**: Upload media files to Twitter for use in tweets or DMs

### Integration
- **tweet-with-media**: Upload media and create a tweet with the uploaded media
- **dm-with-media**: Upload media and send a DM with the uploaded media

## Supported Media Types

### Images
- **JPEG** (`image/jpeg`) - Most common format
- **PNG** (`image/png`) - Supports transparency  
- **WebP** (`image/webp`) - Modern web format
- **BMP** (`image/bmp`) - Bitmap format
- **TIFF** (`image/tiff`) - High quality format

### Subtitles
- **SRT** (`text/srt`) - SubRip subtitle format
- **VTT** (`text/vtt`) - WebVTT subtitle format

## Media Categories

- **tweet_image** - For images in tweets
- **dm_image** - For images in direct messages  
- **subtitles** - For subtitle files

## OAuth2 Requirements

All media upload examples require OAuth2 authentication with appropriate scopes:
- `tweet.write` - For uploading media for tweets
- `dm.write` - For uploading media for DMs

## Usage

Each example can be run with:
```bash
go run main.go -token="your_access_token" -file="path/to/media/file"
```

Make sure to set up your OAuth2 credentials and tokens before running the examples.