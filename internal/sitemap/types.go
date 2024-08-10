package sitemap

import (
	"encoding/xml"
	"github.com/Ja7ad/meilisitemap/config"
)

const (
	xmlHeader        = `<?xml version="1.0" encoding="UTF-8"?>`
	stylesheetLayout = `<?xml-stylesheet type="text/xsl" href="%s"?>`
)

type SitemapIndex struct {
	XMLName  xml.Name `xml:"sitemapindex"`
	Xmlns    string   `xml:"xmlns,attr"`
	Sitemaps []*SMLoc `xml:"sitemap"`
}

type SMLoc struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

type URLSet struct {
	XMLName    xml.Name `xml:"urlset"`
	Xmlns      string   `xml:"xmlns,attr"`
	NewsXmlns  string   `xml:"xmlns:news,attr,omitempty"`
	VideoXmlns string   `xml:"xmlns:video,attr,omitempty"`
	ImageXmlns string   `xml:"xmlns:image,attr,omitempty"`
	URLs       []*URL   `xml:"url"`
}

type URL struct {
	Loc        string            `xml:"loc"`
	LastMod    string            `xml:"lastmod,omitempty"`
	ChangeFreq config.ChangeFreq `xml:"changefreq,omitempty"`
	Priority   string            `xml:"priority,omitempty"`
	Video      *Video            `xml:"video:video,omitempty"`
	Image      *Image            `xml:"image:image,omitempty"`
	News       *News             `xml:"news:news,omitempty"`
}

type Video struct {
	ThumbnailLoc            string `xml:"video:thumbnail_loc"`
	Title                   string `xml:"video:title"`
	Description             string `xml:"video:description"`
	ContentLoc              string `xml:"video:content_loc"`
	PlayerLoc               string `xml:"video:player_loc,omitempty"`
	PlayerLocAutoplay       string `xml:"autoplay,attr,omitempty"`
	Duration                string `xml:"video:duration,omitempty"`
	ExpirationDate          string `xml:"video:expiration_date,omitempty"`
	Rating                  string `xml:"video:rating,omitempty"`
	ViewCount               string `xml:"video:view_count,omitempty"`
	PublicationDate         string `xml:"video:publication_date,omitempty"`
	FamilyFriendly          string `xml:"video:family_friendly,omitempty"`
	Restriction             string `xml:"video:restriction,omitempty"`
	RestrictionRelationship string `xml:"relationship,attr,omitempty"`
	RequiresSubscription    string `xml:"video:requires_subscription,omitempty"`
	Live                    string `xml:"video:live,omitempty"`
}

type Image struct {
	Loc         string `xml:"image:loc"`
	Caption     string `xml:"image:caption,omitempty"`
	Title       string `xml:"image:title,omitempty"`
	License     string `xml:"image:license,omitempty"`
	GeoLocation string `xml:"image:geo_location,omitempty"`
}

type News struct {
	Publication *NewsPublication `xml:"news:publication"`
	PubDate     string           `xml:"news:publication_date"`
	Title       string           `xml:"news:title"`
	Keywords    string           `xml:"news:keywords,omitempty"`
	Description string           `xml:"news:description,omitempty"`
}

type NewsPublication struct {
	Name     string `xml:"news:name"`
	Language string `xml:"news:language"`
}

type RSS struct {
	XMLName xml.Name    `xml:"rss"`
	Xmlns   string      `xml:"xmlns:media,attr"`
	Version string      `xml:"version,attr"`
	Channel *RssChannel `xml:"channel"`
}

type RssChannel struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	Language    string     `xml:"language"`
	Image       *RssImage  `xml:"image"`
	Items       []*RssItem `xml:"item"`
}

type RssImage struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type RssItem struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description string        `xml:"description"`
	Enclosure   *RssEnclosure `xml:"enclosure"`
	Category    string        `xml:"category"`
}

type RssEnclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}
