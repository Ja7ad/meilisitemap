package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	config, err := New("../config.example.yml")
	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		expectErr error
	}{
		{
			name: "valid config",
			config: &Config{
				General: &GeneralConfig{
					BaseIndexURL: "https://example.com",
					FileName:     "validfilename",
					MeiliSearch: &MeiliSearchConfig{
						Host: "http://localhost:7700",
					},
				},
				Sitemaps: map[string]*SitemapConfig{
					"movies": {
						Sitemap:         true,
						BaseAddress:     "https://example.com/movies/",
						SitemapFileName: "validsitemap",
						FieldMap: &FieldMapConfig{
							UniqueField: "title",
							ChangeFreq:  Daily,
							Priority:    High,
						},
					},
				},
			},
			expectErr: nil,
		},
		{
			name: "missing general config",
			config: &Config{
				General: nil,
			},
			expectErr: ErrMissingGeneralConfig,
		},
		{
			name: "invalid base URL",
			config: &Config{
				General: &GeneralConfig{
					BaseIndexURL: "",
				},
			},
			expectErr: ErrInvalidBaseIndexURL,
		},
		{
			name: "missing MeiliSearch config",
			config: &Config{
				General: &GeneralConfig{
					BaseIndexURL: "https://example.com",
				},
				Sitemaps: map[string]*SitemapConfig{},
			},
			expectErr: ErrMissingMeilisearchConfig,
		},
		{
			name: "missing MeiliSearch host",
			config: &Config{
				General: &GeneralConfig{
					BaseIndexURL: "https://example.com",
					MeiliSearch: &MeiliSearchConfig{
						Host: "",
					},
				},
			},
			expectErr: ErrMeilisearchHostRequire,
		},
		{
			name: "empty index name",
			config: &Config{
				General: &GeneralConfig{
					BaseIndexURL: "https://example.com",
					MeiliSearch: &MeiliSearchConfig{
						Host: "http://localhost:7700",
					},
				},
				Sitemaps: map[string]*SitemapConfig{
					"": {
						Sitemap:     true,
						BaseAddress: "https://example.com/movies/",
						FieldMap: &FieldMapConfig{
							UniqueField: "title",
							ChangeFreq:  Daily,
							Priority:    High,
						},
					},
				},
			},
			expectErr: ErrIndexNameIsEmpty,
		},
		{
			name: "missing base path in sitemap",
			config: &Config{
				General: &GeneralConfig{
					BaseIndexURL: "https://example.com",
					MeiliSearch: &MeiliSearchConfig{
						Host: "http://localhost:7700",
					},
				},
				Sitemaps: map[string]*SitemapConfig{
					"movies": {
						Sitemap:     true,
						BaseAddress: "",
						FieldMap: &FieldMapConfig{
							UniqueField: "title",
							ChangeFreq:  Daily,
							Priority:    High,
						},
					},
				},
			},
			expectErr: ErrMissingBaseAddressSitemap,
		},
		{
			name: "missing field map in sitemap",
			config: &Config{
				General: &GeneralConfig{
					BaseIndexURL: "https://example.com",
					MeiliSearch: &MeiliSearchConfig{
						Host: "http://localhost:7700",
					},
				},
				Sitemaps: map[string]*SitemapConfig{
					"movies": {
						Sitemap:     true,
						BaseAddress: "https://example.com/movies/",
					},
				},
			},
			expectErr: ErrInvalidFieldMap,
		},
		{
			name: "missing unique field in field map",
			config: &Config{
				General: &GeneralConfig{
					BaseIndexURL: "https://example.com",
					MeiliSearch: &MeiliSearchConfig{
						Host: "http://localhost:7700",
					},
				},
				Sitemaps: map[string]*SitemapConfig{
					"movies": {
						Sitemap:     true,
						BaseAddress: "/movies/",
						FieldMap: &FieldMapConfig{
							UniqueField: "",
							ChangeFreq:  Daily,
							Priority:    High,
						},
					},
				},
			},
			expectErr: ErrInvalidUniqueField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
