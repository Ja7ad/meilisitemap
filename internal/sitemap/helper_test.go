package sitemap

import (
	"github.com/Ja7ad/meilisitemap/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUniqueToSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Anatomy of a Fall", "anatomy-of-a-fall"},
		{"  Leading and trailing spaces  ", "leading-and-trailing-spaces"},
		{"NoSpaces", "nospaces"},
		{"--Leading and trailing--", "leading-and-trailing"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := uniqueToSlug(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetFileLoc(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		doc       map[string]interface{}
		expected  string
		expectErr bool
	}{
		{
			name:     "Valid image ID with base URL",
			key:      "image_id|https://cdn.example.com/images",
			doc:      map[string]interface{}{"image_id": "12345"},
			expected: "https://cdn.example.com/images/12345",
		},
		{
			name:     "Valid image file name with base URL",
			key:      "image_file_name|https://cdn.example.com/images",
			doc:      map[string]interface{}{"image_file_name": "foobar.jpg"},
			expected: "https://cdn.example.com/images/foobar.jpg",
		},
		{
			name:     "Valid video ID with base URL",
			key:      "video_id|https://cdn.example.com/videos",
			doc:      map[string]interface{}{"video_id": "67890"},
			expected: "https://cdn.example.com/videos/67890",
		},
		{
			name:     "Valid video file name with base URL",
			key:      "video_file_name|https://cdn.example.com/videos",
			doc:      map[string]interface{}{"video_file_name": "foobar.mp4"},
			expected: "https://cdn.example.com/videos/foobar.mp4",
		},
		{
			name:     "Valid image URL",
			key:      "image_url",
			doc:      map[string]interface{}{"image_url": "https://cdn.example.com/images/pic.jpg"},
			expected: "https://cdn.example.com/images/pic.jpg",
		},
		{
			name:     "Valid video URL",
			key:      "video_url",
			doc:      map[string]interface{}{"video_url": "https://cdn.example.com/videos/clip.mp4"},
			expected: "https://cdn.example.com/videos/clip.mp4",
		},
		{
			name:      "Invalid image ID key",
			key:       "image_id|https://cdn.example.com/images",
			doc:       map[string]interface{}{},
			expectErr: true,
		},
		{
			name:      "Invalid URL format in map",
			key:       "image_url",
			doc:       map[string]interface{}{"image_url": 12345},
			expectErr: true,
		},
		{
			name:      "Invalid key without base URL",
			key:       "non_existent_key",
			doc:       map[string]interface{}{"image_url": "https://cdn.example.com/images/pic.jpg"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getFileLoc(tt.key, tt.doc)
			if (err != nil) != tt.expectErr {
				t.Errorf("getFileLoc() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if result != tt.expected {
				t.Errorf("getFileLoc() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestImageFieldMapToSitemapImage(t *testing.T) {
	tests := []struct {
		name      string
		imgCfg    *config.ImageConfig
		doc       map[string]interface{}
		expected  *Image
		expectErr bool
	}{
		{
			name: "Valid image config",
			imgCfg: &config.ImageConfig{
				Loc:         "image_loc",
				Title:       "title_key",
				Caption:     "caption_key",
				GeoLocation: "geo_key",
				License:     "license_key",
			},
			doc: map[string]interface{}{
				"image_loc":   "https://example.com/image.jpg",
				"title_key":   "Sample Title",
				"caption_key": "Sample Caption",
				"geo_key":     "Sample GeoLocation",
				"license_key": "Sample License",
			},
			expected: &Image{
				Loc:         "https://example.com/image.jpg",
				Title:       "Sample Title",
				Caption:     "Sample Caption",
				GeoLocation: "Sample GeoLocation",
				License:     "Sample License",
			},
			expectErr: false,
		},
		{
			name: "Missing title key",
			imgCfg: &config.ImageConfig{
				Loc:   "image_loc",
				Title: "missing_key",
			},
			doc: map[string]interface{}{
				"image_loc": "https://example.com/image.jpg",
			},
			expectErr: true,
		},
		{
			name: "Empty doc",
			imgCfg: &config.ImageConfig{
				Loc: "image_loc",
			},
			doc:       map[string]interface{}{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := imageFieldMapToSitemapImage(tt.imgCfg, tt.doc)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
