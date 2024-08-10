package sitemap

import (
	"github.com/Ja7ad/meilisitemap/config"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var docsForTest = []map[string]interface{}{
	{
		"id":        1,
		"createdAt": time.Now(),
		"video": map[string]interface{}{
			"thumbnail":             1234,
			"title":                 "foobar",
			"description":           "foo bar x y z",
			"file_id":               "123456789",
			"player":                "https://www.example.com/videoplayer.mp4?video=123",
			"autoplay":              true,
			"duration":              5000,
			"expire_at":             time.Now().Add(24 * time.Hour),
			"rating":                10,
			"view":                  5000,
			"published_at":          time.Now(),
			"family_friendly":       true,
			"restriction":           "GB USCA",
			"relationship":          "deny",
			"requires_subscription": false,
			"live":                  false,
		},
		"image": map[string]interface{}{
			"file_id": 123245353543,
			"title":   "foobar",
			"license": "https://creativecommons.org/licenses/by-sa/2.0/",
			"geo":     "Berlin, Germany",
			"caption": "Funny cat on the table is looking at photographer.",
		},
		"news": map[string]interface{}{
			"publication_date": time.Now(),
			"title":            "foobar news",
			"keywords":         []string{"foo", "bar", "baz"},
			"description":      "description",
			"publication": []map[string]interface{}{
				{
					"name":     "foo",
					"language": "english",
				},
				{
					"name":     "bar",
					"language": "spanish",
				},
			},
		},
	},
}

func TestSitemap_CreateSitemap(t *testing.T) {
	tests := []struct {
		name       string
		baseUrl    string
		stylesheet config.Stylesheet
		sitemaps   map[string]*config.SitemapConfig
	}{
		{
			name:       "full test normal",
			baseUrl:    "https://foobar.com",
			stylesheet: config.Style1,
			sitemaps: map[string]*config.SitemapConfig{
				"index1": {
					Sitemap:  true,
					BasePath: "test",
					FieldMap: &config.FieldMapConfig{
						UniqueField: "id",
						LastMod:     "created_at",
						ChangeFreq:  config.Daily,
						Priority:    config.High,
						Video: &config.VideoConfig{
							ThumbnailLoc:            "video.thumbnail|https://cdn.example.com/image|.jpg",
							Title:                   "video.title",
							Description:             "video.description",
							ContentLoc:              "video.file_id|https://cdn.example.com/video|.mp4",
							PlayerLoc:               "video.player",
							Duration:                "video.duration",
							ExpirationDate:          "video.expire_at",
							Rating:                  "video.rating",
							ViewCount:               "video.view",
							PublicationDate:         "video.published_at",
							FamilyFriendly:          "video.family_friendly",
							PlayerAutoPlay:          "video.autoplay",
							RestrictionRelationship: "video.relationship",
							Restriction:             "video.restriction",
							RequiresSubscription:    "video.requires_subscription",
							Live:                    "video.live",
						},
						Image: &config.ImageConfig{
							Loc:         "image.file_id|https://cdn.example.com/image|.jpg",
							Title:       "image.title",
							Caption:     "image.caption",
							License:     "image.license",
							GeoLocation: "image.geo",
						},
						News: &config.NewsConfig{
							Publication: &config.NewsPublicationConfig{
								Name:     "news.publication.0.name",
								Language: "news.publication.0.language",
							},
							PubDate:     "news.publication_date",
							Title:       "news.title",
							Keywords:    "news.keywords",
							Description: "news.description",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sm := New(test.baseUrl, test.stylesheet, test.sitemaps)
			res, err := sm.CreateSitemap("index1", docsForTest)
			require.NoError(t, err)
			require.NotNil(t, res)
		})
	}
}
