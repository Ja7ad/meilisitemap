package config

import (
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
)

func New(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	cfg := new(Config)

	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.General == nil {
		return ErrMissingGeneralConfig
	}

	if c.General.BaseURL == "" {
		return ErrInvalidBaseURL
	}

	_, err := url.Parse(c.General.BaseURL)
	if err != nil {
		return ErrInvalidBaseURL
	}

	if c.General.Stylesheet != "" {
		switch c.General.Stylesheet {
		case Style1, Style2:
		default:
			c.General.Stylesheet = Style1
		}
	}

	if c.General.MeiliSearch == nil {
		return ErrMissingMeilisearchConfig
	}

	if c.General.MeiliSearch.Host == "" {
		return ErrMeilisearchHostRequire
	}

	for name, sitemap := range c.Sitemaps {
		if err := validateSitemapConfig(name, sitemap); err != nil {
			return err
		}
	}

	return nil
}

func validateSitemapConfig(name string, sitemap *SitemapConfig) error {
	if name == "" {
		return ErrIndexNameIsEmpty
	}

	if !sitemap.Sitemap {
		return ErrInvalidSitemapConfig
	}

	if sitemap.BasePath == "" {
		return ErrMissingBasePathSitemap
	}

	if sitemap.FieldMap == nil {
		return ErrInvalidFieldMap
	}

	if sitemap.FieldMap.UniqueField == "" {
		return ErrInvalidUniqueField
	}

	switch sitemap.FieldMap.ChangeFreq {
	case Always, Hourly, Daily, Weekly, Monthly, Yearly, Never:
	default:
		sitemap.FieldMap.ChangeFreq = Daily
	}

	switch sitemap.FieldMap.Priority {
	case Low, Medium, High, Highest:
	default:
		sitemap.FieldMap.Priority = High
	}

	return nil
}
