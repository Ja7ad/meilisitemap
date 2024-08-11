package config

import "errors"

var (
	ErrMissingMeilisearchConfig  = errors.New("meilisearch configuration is missing")
	ErrMeilisearchHostRequire    = errors.New("meilisearch host is required")
	ErrInvalidBaseIndexURL       = errors.New("invalid or missing base_index_url")
	ErrInvalidSitemapConfig      = errors.New("sitemap is required")
	ErrMissingBaseAddressSitemap = errors.New("base_address sitemap is required")
	ErrInvalidFieldMap           = errors.New("invalid or missing field_map in sitemap config")
	ErrInvalidUniqueField        = errors.New("invalid or missing unique_field in field_map")
	ErrIndexNameIsEmpty          = errors.New("index name is empty")
	ErrMissingGeneralConfig      = errors.New("general config is missing")
)
