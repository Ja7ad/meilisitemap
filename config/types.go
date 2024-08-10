package config

import "time"

type Config struct {
	General  *GeneralConfig            `yaml:"general"`
	Sitemaps map[string]*SitemapConfig `yaml:"sitemaps"`
}

type GeneralConfig struct {
	BaseURL          string             `yaml:"base_url"`
	IndexSitemapPath string             `yaml:"indexsitemap_path"`
	FileName         string             `yaml:"file_name"`
	Prefix           string             `yaml:"prefix"`
	Stylesheet       Stylesheet         `yaml:"stylesheet"`
	Serve            *ServeConfig       `yaml:"serve"`
	MeiliSearch      *MeiliSearchConfig `yaml:"meilisearch"`
}

type ServeConfig struct {
	Enable bool   `yaml:"enable"`
	Listen string `yaml:"listen"`
	PPROF  bool   `yaml:"pprof"`
}

type MeiliSearchConfig struct {
	Host   string `yaml:"host"`
	APIKey string `yaml:"api_key"`
}

type PprofConfig struct {
	Enable bool   `yaml:"enable"`
	Listen string `yaml:"listen"`
}

type SitemapConfig struct {
	Sitemap         bool            `yaml:"sitemap"`
	HTMLSitemap     bool            `yaml:"html_sitemap"`
	RSS             bool            `yaml:"rss"`
	Filter          string          `yaml:"filter"`
	BasePath        string          `yaml:"base_path"`
	Compress        bool            `yaml:"compress"`
	SitemapFileName string          `yaml:"sitemap_file_name"`
	LiveUpdate      *LiveConfig     `yaml:"live_update"`
	FieldMap        *FieldMapConfig `yaml:"field_map"`
}

type LiveConfig struct {
	Enabled  bool  `yaml:"enabled"`
	Interval int64 `yaml:"interval"`
}

type FieldMapConfig struct {
	UniqueField string       `yaml:"unique_field"`
	LastMod     string       `yaml:"lastmod"`
	ChangeFreq  ChangeFreq   `yaml:"changefreq"`
	Priority    Priority     `yaml:"priority"`
	Video       *VideoConfig `yaml:"video,omitempty"`
	Image       *ImageConfig `yaml:"image,omitempty"`
	News        *NewsConfig  `yaml:"news,omitempty"`
}

type VideoConfig struct {
	ThumbnailLoc            string `yaml:"thumbnail_loc"`
	Title                   string `yaml:"title"`
	Description             string `yaml:"description"`
	ContentLoc              string `yaml:"content_loc"`
	PlayerLoc               string `yaml:"player_loc"`
	PlayerAutoPlay          string `yaml:"player_auto_play"`
	Duration                string `yaml:"duration"`
	ExpirationDate          string `yaml:"expiration_date"`
	Rating                  string `yaml:"rating"`
	ViewCount               string `yaml:"view_count"`
	PublicationDate         string `yaml:"publication_date"`
	FamilyFriendly          string `yaml:"family_friendly"`
	RestrictionRelationship string `yaml:"relationship"`
	Restriction             string `yaml:"restriction"`
	RequiresSubscription    string `yaml:"requires_subscription"`
	Live                    string `yaml:"live"`
}

type ImageConfig struct {
	Loc         string `yaml:"loc"`
	Caption     string `yaml:"caption"`
	Title       string `yaml:"title"`
	License     string `yaml:"license"`
	GeoLocation string `yaml:"geo_location"`
}

type NewsConfig struct {
	Publication *NewsPublicationConfig `yaml:"publication"`
	PubDate     string                 `yaml:"pub_date"`
	Title       string                 `yaml:"title"`
	Keywords    string                 `yaml:"keywords"`
	Description string                 `yaml:"description"`
}

type NewsPublicationConfig struct {
	Name     string `yaml:"name"`
	Language string `yaml:"language"`
}

type (
	ChangeFreq string
	Priority   string
	Stylesheet string
)

const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

const (
	Low     Priority = "low"
	Medium  Priority = "medium"
	High    Priority = "high"
	Highest Priority = "highest"
)

const (
	Style1 Stylesheet = "style1"
	Style2 Stylesheet = "style2"
)

func (c ChangeFreq) Interval() time.Duration {
	switch c {
	case Always:
		return 5 * time.Minute
	case Daily, Weekly, Monthly, Yearly, Never:
		return time.Hour * 24
	default:
		return time.Hour
	}
}

func (s Stylesheet) Link() string {
	switch s {
	case Style1:
		return "https://raw.githubusercontent.com/Ja7ad/meilisitemap/main/stylesheets/style1.xsl"
	case Style2:
		return "https://raw.githubusercontent.com/Ja7ad/meilisitemap/main/stylesheets/style2.xsl"
	default:
		return ""
	}
}

func (p Priority) Rate() float64 {
	switch p {
	case Low:
		return 0.3
	case Medium:
		return 0.5
	case High:
		return 0.8
	case Highest:
		return 1.0
	default:
		return 0.8
	}
}
