package twitter

// GetBestMediaURL determines the best URL for media based on type and available variants
// This is a public helper function for external use
func GetBestMediaURL(media *MediaObj) string {
	if media == nil {
		return ""
	}

	// For images (photos), use direct URL if available
	if media.Type == "photo" {
		if media.URL != "" {
			return media.URL
		}
		return media.PreviewImageURL
	}

	// For videos and animated GIFs, prefer variants over direct URL
	if len(media.Variants) > 0 {
		bestVariant := SelectBestVariant(media.Variants)
		if bestVariant != nil && bestVariant.URL != "" {
			return bestVariant.URL
		}
	}

	// Fallback to direct URL or preview URL
	if media.URL != "" {
		return media.URL
	}
	return media.PreviewImageURL
}

// SelectBestVariant chooses the best video variant based on quality and format preferences
// This is a public helper function for external use
func SelectBestVariant(variants []*MediaVariantObj) *MediaVariantObj {
	if len(variants) == 0 {
		return nil
	}

	var bestVariant *MediaVariantObj
	var bestBitRate int

	for _, variant := range variants {
		// Skip HLS streams, prefer MP4
		if variant.ContentType == "application/x-mpegURL" {
			continue
		}

		// Prefer MP4 videos with highest bitrate
		if variant.ContentType == "video/mp4" && (bestVariant == nil || variant.BitRate > bestBitRate) {
			bestVariant = variant
			bestBitRate = variant.BitRate
		}
	}

	// If no MP4 found, use the first available variant
	if bestVariant == nil {
		return variants[0]
	}

	return bestVariant
}

// ConvertMediaType converts Twitter's media type to a more standard format
// This is a public helper function for external use
func ConvertMediaType(twitterType string) string {
	switch twitterType {
	case "photo":
		return "image"
	default:
		return twitterType
	}
}

// IsVideoMedia checks if the media type represents a video (including animated GIFs)
func IsVideoMedia(mediaType string) bool {
	return mediaType == "video" || mediaType == "animated_gif"
}

// IsImageMedia checks if the media type represents an image/photo
func IsImageMedia(mediaType string) bool {
	return mediaType == "photo"
}

// GetMediaDimensions returns width and height if available
func GetMediaDimensions(media *MediaObj) (width int, height int, hasSize bool) {
	if media == nil {
		return 0, 0, false
	}
	if media.Width > 0 && media.Height > 0 {
		return media.Width, media.Height, true
	}
	return 0, 0, false
}

// GetVideoDuration returns video duration in seconds, or 0 if not available
func GetVideoDuration(media *MediaObj) int {
	if media == nil || media.DurationMS <= 0 {
		return 0
	}
	return media.DurationMS / 1000
}

// HasHighQualityVariants checks if the media has multiple quality variants
func HasHighQualityVariants(media *MediaObj) bool {
	if media == nil || len(media.Variants) <= 1 {
		return false
	}
	
	// Check if there are MP4 variants with different bitrates
	mp4Count := 0
	for _, variant := range media.Variants {
		if variant.ContentType == "video/mp4" {
			mp4Count++
		}
	}
	return mp4Count > 1
}

// GetAllVideoQualityURLs returns all MP4 variants sorted by bitrate (highest first)
func GetAllVideoQualityURLs(media *MediaObj) []VideoQuality {
	if media == nil {
		return nil
	}

	var qualities []VideoQuality
	for _, variant := range media.Variants {
		if variant.ContentType == "video/mp4" {
			qualities = append(qualities, VideoQuality{
				URL:     variant.URL,
				BitRate: variant.BitRate,
				Quality: getQualityLabel(variant.BitRate),
			})
		}
	}

	// Sort by bitrate descending
	for i := 0; i < len(qualities)-1; i++ {
		for j := i + 1; j < len(qualities); j++ {
			if qualities[i].BitRate < qualities[j].BitRate {
				qualities[i], qualities[j] = qualities[j], qualities[i]
			}
		}
	}

	return qualities
}

// VideoQuality represents a video quality option
type VideoQuality struct {
	URL     string `json:"url"`
	BitRate int    `json:"bit_rate"`
	Quality string `json:"quality"`
}

// getQualityLabel converts bitrate to human-readable quality label
func getQualityLabel(bitRate int) string {
	switch {
	case bitRate >= 2000000:
		return "1080p"
	case bitRate >= 800000:
		return "720p"
	case bitRate >= 500000:
		return "480p"
	case bitRate >= 200000:
		return "360p"
	default:
		return "240p"
	}
}