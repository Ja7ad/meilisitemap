package config

import "errors"

var (
	ErrMissingMeilisearchConfig = errors.New("meilisearch configuration is missing")
	ErrMeilisearchHostRequire   = errors.New("meilisearch host is required")
	ErrInvalidBaseURL           = errors.New("invalid or missing base_url")
	ErrInvalidSitemapConfig     = errors.New("sitemap is required")
	ErrMissingBasePathSitemap   = errors.New("base_path sitemap is required")
	ErrInvalidFieldMap          = errors.New("invalid or missing field_map in sitemap config")
	ErrInvalidUniqueField       = errors.New("invalid or missing unique_field in field_map")
	ErrIndexNameIsEmpty         = errors.New("index name is empty")
	ErrInvalidFileName          = errors.New("file name contains invalid characters")
	ErrMissingGeneralConfig     = errors.New("general config is missing")
)
