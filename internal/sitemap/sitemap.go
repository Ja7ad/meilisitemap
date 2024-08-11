package sitemap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ja7ad/meilisitemap/config"
	"github.com/Ja7ad/meilisitemap/internal/logger"
	"github.com/Ja7ad/meilisitemap/utils"
	"github.com/klauspost/compress/gzip"
	"github.com/tdewolff/minify/v2"
	minXml "github.com/tdewolff/minify/v2/xml"
)

const (
	_standardXmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"
	_videoXmlns    = "http://www.google.com/schemas/sitemap-video/1.1"
	_imageXmlns    = "http://www.google.com/schemas/sitemap-image/1.1"
	_newsXmlns     = "http://www.google.com/schemas/sitemap-news/0.9"

	_datetimeLayout = "2006-01-02T15:04:05-07:00"
)

type Sitemap struct {
	indexes    map[string]*config.SitemapConfig
	stylesheet config.Stylesheet
	log        logger.Logger
}

func New(stylesheet config.Stylesheet,
	sitemaps map[string]*config.SitemapConfig, log logger.Logger,
) *Sitemap {
	return &Sitemap{
		indexes:    sitemaps,
		stylesheet: stylesheet,
		log:        log,
	}
}

func (s *Sitemap) CreateSitemap(index string, docs []map[string]any) ([]byte, error) {
	idxCfg := s.indexes[index]

	sitemap := new(URLSet)
	sitemap.Xmlns = _standardXmlns

	if idxCfg.FieldMap.Video != nil {
		sitemap.VideoXmlns = _videoXmlns
	}

	if idxCfg.FieldMap.Image != nil {
		sitemap.ImageXmlns = _imageXmlns
	}

	if idxCfg.FieldMap.News != nil {
		sitemap.NewsXmlns = _newsXmlns
	}

	sitemap.URLs = make([]*URL, 0)

	for _, doc := range docs {
		u, err := s.urlMaker(doc, idxCfg)
		if err != nil {
			return nil, err
		}
		if !isExistsLoc(u, sitemap.URLs) {
			sitemap.URLs = append(sitemap.URLs, u)
		}
	}

	xmlData, err := marshal(sitemap)
	if err != nil {
		return nil, err
	}

	header := []byte(xmlHeader + "\n")

	if s.stylesheet != "" {
		header = []byte(xmlHeader + fmt.Sprintf(stylesheetLayout, s.stylesheet.Link()) + "\n")
	}

	fullXmlData := append(header, xmlData...)

	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), minXml.Minify)
	b, err := m.Bytes("text/xml", fullXmlData)

	if idxCfg.Compress {
		return compress(b)
	}

	return b, nil
}

func (s *Sitemap) urlMaker(doc map[string]any, cfg *config.SitemapConfig) (*URL, error) {
	u := new(URL)

	unique := utils.PickByNestedKey(doc, cfg.FieldMap.UniqueField)
	if unique == nil {
		return nil, fmt.Errorf("failed to get value unique field %s", cfg.FieldMap.UniqueField)
	}

	var (
		loc string
		err error
	)

	switch unique.(type) {
	case string:
		slug := uniqueToSlug(unique.(string))
		if slug == "" {
			return nil, fmt.Errorf("failed to get value slug field %s", cfg.FieldMap.UniqueField)
		}
		if !strings.HasSuffix(cfg.BaseAddress, "=") {
			loc, err = url.JoinPath(cfg.BaseAddress, slug)
		} else {
			loc = fmt.Sprintf("%s%s", cfg.BaseAddress, slug)
		}
	case int:
		if !strings.HasSuffix(cfg.BaseAddress, "=") {
			loc, err = url.JoinPath(cfg.BaseAddress, strconv.Itoa(unique.(int)))
		} else {
			loc = fmt.Sprintf("%s%s", cfg.BaseAddress, strconv.Itoa(unique.(int)))
		}
	default:
		return nil, fmt.Errorf("not supported unique field %s, type is %T", cfg.FieldMap.UniqueField, unique)
	}

	if err != nil {
		return nil, err
	}

	u.Loc = loc
	u.Priority = fmt.Sprintf("%g", cfg.FieldMap.Priority.Rate())
	u.ChangeFreq = cfg.FieldMap.ChangeFreq

	if datetime, ok := doc[cfg.FieldMap.LastMod]; !ok {
		u.LastMod = time.Now().Format(_datetimeLayout)
	} else {
		lastMod, err := getDateTimeFromDoc(datetime)
		if err != nil {
			return nil, err
		}
		u.LastMod = lastMod
	}

	if cfg.FieldMap.Image != nil {
		u.Image, err = imageFieldMapToSitemapImage(cfg.FieldMap.Image, doc)
		if err != nil {
			s.log.Warn("failed to create image sitemap", "unique", unique, "err", err)
		}
	}

	if cfg.FieldMap.Video != nil {
		u.Video, err = videoFieldMapToSitemapVideo(cfg.FieldMap.Video, doc)
		if err != nil {
			s.log.Warn("failed to create video sitemap", "unique", unique, "err", err)
		}
	}

	if cfg.FieldMap.News != nil {
		u.News, err = newsFieldMapToSitemapNews(cfg.FieldMap.News, doc)
		if err != nil {
			s.log.Warn("failed to create news sitemap", "unique", unique, "err", err)
		}
	}

	return u, nil
}

func compress(b []byte) ([]byte, error) {
	var buf bytes.Buffer

	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write(b)
	if err != nil {
		return nil, fmt.Errorf("failed to write to gzip writer: %w", err)
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

func marshal(u *URLSet) ([]byte, error) {
	b, err := xml.MarshalIndent(u, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling XML: %v", err)
	}
	return b, nil
}

func isExistsLoc(item *URL, items []*URL) bool {
	for _, u := range items {
		if item.Loc == u.Loc {
			return true
		}
	}
	return false
}
