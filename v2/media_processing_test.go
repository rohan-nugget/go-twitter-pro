package twitter

import (
	"testing"
)

// Test data structures based on real Twitter API responses (PII redacted)
func getTestMediaObjects() []*MediaObj {
	return []*MediaObj{
		// Test case 1: Photo media (should use direct URL)
		{
			Key:              "3_1234567890123456789",
			Type:             "photo",
			URL:              "https://pbs.twimg.com/media/ExamplePhotoKey.jpg",
			PreviewImageURL:  "https://pbs.twimg.com/media/ExamplePhotoKey?format=jpg&name=small",
			Width:            1200,
			Height:           800,
			AltText:          "Sample photo description",
			PublicMetrics: &MediaMetricsObj{
				Views: 15420,
			},
		},
		// Test case 2: Video with multiple variants (should select best MP4)
		{
			Key:              "13_9876543210987654321",
			Type:             "video",
			URL:              "",  // Often empty for videos
			PreviewImageURL:  "https://pbs.twimg.com/ext_tw_video_thumb/9876543210987654321/pu/img/ExampleThumb.jpg",
			DurationMS:       45000,
			Width:            1280,
			Height:           720,
			PublicMetrics: &MediaMetricsObj{
				Views: 89300,
			},
			Variants: []*MediaVariantObj{
				{
					BitRate:     256000,
					ContentType: "video/mp4",
					URL:         "https://video.twimg.com/ext_tw_video/9876543210987654321/pu/vid/480x270/ExampleVideo480.mp4",
				},
				{
					BitRate:     832000,
					ContentType: "video/mp4", 
					URL:         "https://video.twimg.com/ext_tw_video/9876543210987654321/pu/vid/720x404/ExampleVideo720.mp4",
				},
				{
					BitRate:     2176000,
					ContentType: "video/mp4",
					URL:         "https://video.twimg.com/ext_tw_video/9876543210987654321/pu/vid/1280x720/ExampleVideo1280.mp4",
				},
				{
					BitRate:     0,
					ContentType: "application/x-mpegURL", // HLS stream - should be skipped
					URL:         "https://video.twimg.com/ext_tw_video/9876543210987654321/pu/pl/ExamplePlaylist.m3u8",
				},
			},
		},
		// Test case 3: Animated GIF
		{
			Key:              "4_1122334455667788990",
			Type:             "animated_gif",
			URL:              "",
			PreviewImageURL:  "https://pbs.twimg.com/tweet_video_thumb/ExampleGif.jpg",
			Width:            498,
			Height:           280,
			Variants: []*MediaVariantObj{
				{
					BitRate:     0,
					ContentType: "video/mp4",
					URL:         "https://video.twimg.com/tweet_video/ExampleGif.mp4",
				},
			},
		},
		// Test case 4: Video with no variants (edge case)
		{
			Key:              "13_5544332211009988776",
			Type:             "video", 
			URL:              "https://video.twimg.com/ext_tw_video/5544332211009988776/fallback.mp4",
			PreviewImageURL:  "https://pbs.twimg.com/ext_tw_video_thumb/5544332211009988776/pu/img/FallbackThumb.jpg",
			Width:            640,
			Height:           360,
			DurationMS:       12000,
			Variants:         []*MediaVariantObj{}, // Empty variants
		},
		// Test case 5: Photo with no URL (edge case)
		{
			Key:              "3_9988776655443322110",
			Type:             "photo",
			URL:              "", // Empty URL
			PreviewImageURL:  "https://pbs.twimg.com/media/FallbackPreview.jpg",
			Width:            800,
			Height:           600,
		},
	}
}

func TestGetBestMediaURL(t *testing.T) {
	testMedia := getTestMediaObjects()
	
	testCases := []struct {
		name           string
		mediaIndex     int
		expectedURL    string
		expectedType   string
		description    string
	}{
		{
			name:         "Photo with direct URL",
			mediaIndex:   0,
			expectedURL:  "https://pbs.twimg.com/media/ExamplePhotoKey.jpg",
			expectedType: "image",
			description:  "Should return direct URL for photos",
		},
		{
			name:         "Video with variants - highest bitrate MP4",
			mediaIndex:   1,
			expectedURL:  "https://video.twimg.com/ext_tw_video/9876543210987654321/pu/vid/1280x720/ExampleVideo1280.mp4",
			expectedType: "video",
			description:  "Should select highest bitrate MP4 variant (2176000)",
		},
		{
			name:         "Animated GIF",
			mediaIndex:   2,
			expectedURL:  "https://video.twimg.com/tweet_video/ExampleGif.mp4",
			expectedType: "animated_gif",
			description:  "Should use variant URL for animated GIF",
		},
		{
			name:         "Video with no variants - fallback to direct URL",
			mediaIndex:   3,
			expectedURL:  "https://video.twimg.com/ext_tw_video/5544332211009988776/fallback.mp4",
			expectedType: "video",
			description:  "Should fallback to direct URL when no variants available",
		},
		{
			name:         "Photo with no URL - fallback to preview",
			mediaIndex:   4,
			expectedURL:  "https://pbs.twimg.com/media/FallbackPreview.jpg",
			expectedType: "image",
			description:  "Should fallback to preview URL when direct URL is empty",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			media := testMedia[tc.mediaIndex]
			
			// Test URL selection using new helper function
			actualURL := GetBestMediaURL(media)
			if actualURL != tc.expectedURL {
				t.Errorf("URL mismatch for %s:\nExpected: %s\nActual: %s\nDescription: %s", 
					tc.name, tc.expectedURL, actualURL, tc.description)
			}
			
			// Test media type conversion
			actualType := ConvertMediaType(media.Type)
			if actualType != tc.expectedType {
				t.Errorf("Type mismatch for %s:\nExpected: %s\nActual: %s", 
					tc.name, tc.expectedType, actualType)
			}
			
			t.Logf("✅ %s: URL=%s, Type=%s", tc.name, actualURL, actualType)
		})
	}
}

func TestSelectBestVariant(t *testing.T) {
	
	testCases := []struct {
		name            string
		variants        []*MediaVariantObj
		expectedURL     string
		expectedBitRate int
		description     string
	}{
		{
			name: "Multiple MP4 variants - select highest bitrate",
			variants: []*MediaVariantObj{
				{BitRate: 256000, ContentType: "video/mp4", URL: "https://example.com/low.mp4"},
				{BitRate: 832000, ContentType: "video/mp4", URL: "https://example.com/medium.mp4"},
				{BitRate: 2176000, ContentType: "video/mp4", URL: "https://example.com/high.mp4"},
			},
			expectedURL:     "https://example.com/high.mp4",
			expectedBitRate: 2176000,
			description:     "Should select MP4 with highest bitrate",
		},
		{
			name: "Mixed variants with HLS - skip HLS",
			variants: []*MediaVariantObj{
				{BitRate: 0, ContentType: "application/x-mpegURL", URL: "https://example.com/playlist.m3u8"},
				{BitRate: 832000, ContentType: "video/mp4", URL: "https://example.com/video.mp4"},
			},
			expectedURL:     "https://example.com/video.mp4",
			expectedBitRate: 832000,
			description:     "Should skip HLS streams and select MP4",
		},
		{
			name: "No MP4 variants - use first available",
			variants: []*MediaVariantObj{
				{BitRate: 0, ContentType: "application/x-mpegURL", URL: "https://example.com/playlist.m3u8"},
				{BitRate: 0, ContentType: "video/webm", URL: "https://example.com/video.webm"},
			},
			expectedURL:     "https://example.com/playlist.m3u8",
			expectedBitRate: 0,
			description:     "Should use first variant when no MP4 available",
		},
		{
			name:            "Empty variants",
			variants:        []*MediaVariantObj{},
			expectedURL:     "",
			expectedBitRate: 0,
			description:     "Should return nil for empty variants",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SelectBestVariant(tc.variants)
			
			if len(tc.variants) == 0 {
				if result != nil {
					t.Errorf("Expected nil for empty variants, got %+v", result)
				}
				return
			}
			
			if result == nil {
				t.Errorf("Expected variant, got nil")
				return
			}
			
			if result.URL != tc.expectedURL {
				t.Errorf("URL mismatch for %s:\nExpected: %s\nActual: %s\nDescription: %s",
					tc.name, tc.expectedURL, result.URL, tc.description)
			}
			
			if result.BitRate != tc.expectedBitRate {
				t.Errorf("BitRate mismatch for %s:\nExpected: %d\nActual: %d",
					tc.name, tc.expectedBitRate, result.BitRate)
			}
			
			t.Logf("✅ %s: URL=%s, BitRate=%d, ContentType=%s", 
				tc.name, result.URL, result.BitRate, result.ContentType)
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	testMedia := getTestMediaObjects()
	
	t.Run("Media dimension helpers", func(t *testing.T) {
		media := testMedia[0] // Photo with dimensions
		width, height, hasSize := GetMediaDimensions(media)
		
		if !hasSize {
			t.Error("Expected hasSize to be true for media with dimensions")
		}
		
		if width != 1200 || height != 800 {
			t.Errorf("Expected dimensions 1200x800, got %dx%d", width, height)
		}
		
		t.Logf("✅ Media dimensions: %dx%d", width, height)
	})
	
	t.Run("Video duration helper", func(t *testing.T) {
		media := testMedia[1] // Video with duration
		duration := GetVideoDuration(media)
		
		expectedDuration := 45 // 45000ms = 45s
		if duration != expectedDuration {
			t.Errorf("Expected duration %ds, got %ds", expectedDuration, duration)
		}
		
		t.Logf("✅ Video duration: %ds", duration)
	})
	
	t.Run("Media type helpers", func(t *testing.T) {
		photoMedia := testMedia[0]
		videoMedia := testMedia[1]
		gifMedia := testMedia[2]
		
		if !IsImageMedia(photoMedia.Type) {
			t.Error("Expected photo to be identified as image media")
		}
		
		if !IsVideoMedia(videoMedia.Type) {
			t.Error("Expected video to be identified as video media")
		}
		
		if !IsVideoMedia(gifMedia.Type) {
			t.Error("Expected animated_gif to be identified as video media")
		}
		
		t.Log("✅ Media type helpers working correctly")
	})
	
	t.Run("High quality variants detection", func(t *testing.T) {
		videoWithVariants := testMedia[1]
		videoWithoutVariants := testMedia[3]
		
		if !HasHighQualityVariants(videoWithVariants) {
			t.Error("Expected video with multiple MP4 variants to have high quality variants")
		}
		
		if HasHighQualityVariants(videoWithoutVariants) {
			t.Error("Expected video without variants to not have high quality variants")
		}
		
		t.Log("✅ High quality variants detection working")
	})
	
	t.Run("Video quality URLs extraction", func(t *testing.T) {
		videoMedia := testMedia[1] // Video with multiple variants
		qualities := GetAllVideoQualityURLs(videoMedia)
		
		if len(qualities) != 3 { // Should have 3 MP4 variants (excluding HLS)
			t.Errorf("Expected 3 MP4 qualities, got %d", len(qualities))
		}
		
		// Should be sorted by bitrate descending
		for i := 0; i < len(qualities)-1; i++ {
			if qualities[i].BitRate < qualities[i+1].BitRate {
				t.Error("Expected qualities to be sorted by bitrate descending")
			}
		}
		
		// Highest quality should be first
		if qualities[0].BitRate != 2176000 {
			t.Errorf("Expected highest bitrate 2176000, got %d", qualities[0].BitRate)
		}
		
		if qualities[0].Quality != "1080p" {
			t.Errorf("Expected quality label '1080p', got '%s'", qualities[0].Quality)
		}
		
		t.Logf("✅ Video qualities: %d variants, highest: %s (%d bps)", 
			len(qualities), qualities[0].Quality, qualities[0].BitRate)
	})
}

// Benchmark tests for performance validation
func BenchmarkGetBestMediaURL(b *testing.B) {
	media := getTestMediaObjects()[1] // Video with multiple variants
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetBestMediaURL(media)
	}
}

func BenchmarkSelectBestVariant(b *testing.B) {
	variants := getTestMediaObjects()[1].Variants // Video variants
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SelectBestVariant(variants)
	}
}