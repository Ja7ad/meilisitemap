package generator

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/Ja7ad/meilibridge/pkg/logger"
	"github.com/Ja7ad/meilisitemap/config"
	"github.com/Ja7ad/meilisitemap/internal/sched"
	"github.com/Ja7ad/meilisitemap/internal/server"
	"github.com/Ja7ad/meilisitemap/internal/sitemap"
	"github.com/meilisearch/meilisearch-go"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	_standardXmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

	_defaultWaitInterval   = 5 * time.Second
	_defaultHitSizePerPage = 100
	_dateLayout            = "2006-01-02"
)

type Sitemap struct {
	baseUrl          string
	storePath        string
	indexsitemapPath string
	fileName         string
	prefix           string
	stylesheet       config.Stylesheet
	pprof            *config.PprofConfig
	meili            *meilisearch.Client
	sitemaps         map[string]*config.SitemapConfig
	logger           logger.Logger
	wg               sync.WaitGroup
	sched            *sched.Sched
	ctx              context.Context
	server           *server.Server
	sm               *sitemap.Sitemap
}

func New(
	ctx context.Context,
	storePath string,
	general *config.GeneralConfig,
	logger logger.Logger,
	sitemaps map[string]*config.SitemapConfig,
) (*Sitemap, error) {
	s := new(Sitemap)
	s.baseUrl = general.BaseURL
	s.storePath = storePath
	s.indexsitemapPath = general.IndexSitemapPath
	s.fileName = general.FileName
	s.prefix = general.Prefix
	s.stylesheet = general.Stylesheet
	s.logger = logger
	s.ctx = ctx
	s.sched = sched.New(ctx, s.logger)
	s.sm = sitemap.New(s.baseUrl, s.stylesheet, sitemaps, s.logger)

	if _, err := os.Stat(filepath.Join(s.storePath, s.indexsitemapPath)); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(s.storePath, s.indexsitemapPath), 0777); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if general.Serve != nil && general.Serve.Enable {
		s.server = server.New(general.Serve, storePath)
	}

reconnect:
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   general.MeiliSearch.Host,
		APIKey: general.MeiliSearch.APIKey,
	})

	if !client.IsHealthy() {
		s.logger.Warn("failed connecting to Meilisearch, try connect...")
		t := time.NewTicker(_defaultWaitInterval)
		defer t.Stop()

		for {
			<-t.C
			goto reconnect
		}
	}

	s.logger.Info("successfully connected to Meilisearch")

	s.meili = client
	s.sitemaps = sitemaps

	return s, nil
}

func (s *Sitemap) Start() error {
	doneCh := make(chan struct{})
	sets := make([]string, 0)
	setCh := make(chan string)

	for idx, sm := range s.sitemaps {
		if err := s.existsIndex(idx); err != nil {
			return err
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			doFunc := func() {
				s.logger.Info("started fetching documents", "index", idx)
				results, err := s.fetchIndexDocuments(idx, sm.Filter)
				if err != nil {
					s.logger.Error("failed to fetch documents index", "index", idx, "err", err.Error())
					return
				}

				b, err := s.sm.CreateSitemap(idx, results)
				if err != nil {
					s.logger.Error("failed to create sitemap for index",
						"index", idx, "err", err.Error())
					return
				}

				sAddr, err := s.saveSitemap(b, idx, sm)
				if err != nil {
					s.logger.Fatal("failed to save sitemap", "index", idx, "err", err.Error())
				}

				setCh <- sAddr

				s.logger.Info("created sitemap for index", "index", idx)
			}

			if sm.LiveUpdate != nil && sm.LiveUpdate.Enabled {
				s.sched.AddJob(doFunc, time.Duration(sm.LiveUpdate.Interval)*time.Second)
			} else {
				doFunc()
			}

		}()
	}

	go func(log logger.Logger) {
		for set := range setCh {
			if !existsItem(sets, set) {
				sets = append(sets, set)

				if err := s.createSitemapIndex(sets); err != nil {
					log.Fatal("failed to create sitemap.xml", "err", err.Error())
				}
			}
		}
	}(s.logger)

	if s.server != nil {
		go func() {
			s.logger.Info("sitemaps served", "addr", "http://"+s.server.Addr())
			s.server.Start()
			for {
				select {
				case <-s.ctx.Done():
					_ = s.server.Shutdown(s.ctx)
					return
				case err := <-s.server.Notify():
					s.logger.Fatal(err.Error())
				}
			}
		}()
	}

	go func() {
		<-s.ctx.Done()
		doneCh <- struct{}{}
	}()

	s.wg.Wait()

	if s.sched.Len() != 0 {
		go s.sched.Start()
	}

	<-doneCh
	close(setCh)

	return nil
}

func (s *Sitemap) existsIndex(idx string) error {
	_, err := s.meili.GetIndex(idx)
	return err
}

func (s *Sitemap) createSitemapIndex(setsFilename []string) error {
	sitemapIdx := new(sitemap.SitemapIndex)
	sitemapIdx.Xmlns = _standardXmlns
	sitemapIdx.Sitemaps = make([]*sitemap.SMLoc, 0, len(setsFilename))

	now := time.Now().Format(_dateLayout)

	for _, fn := range setsFilename {
		smLoc := new(sitemap.SMLoc)

		loc, err := url.JoinPath(s.baseUrl, s.indexsitemapPath, fn)
		if err != nil {
			return err
		}

		if s.server != nil {
			loc, err = url.JoinPath("http://"+s.server.Addr(), s.indexsitemapPath, fn)
		}

		smLoc.Loc = loc
		smLoc.LastMod = now
		sitemapIdx.Sitemaps = append(sitemapIdx.Sitemaps, smLoc)
	}

	xmlData, err := xml.MarshalIndent(sitemapIdx, "", "  ")
	if err != nil {
		return err
	}

	fileName := "sitemap"

	if s.fileName != "" {
		fileName = s.fileName
	}

	filePath := filepath.Join(s.storePath, fileName+".xml")

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err := file.Write(xmlData); err != nil {
		return fmt.Errorf("error writing to file %s: %v", filePath, err)
	}

	return nil
}

func (s *Sitemap) fetchIndexDocuments(index, filter string) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)

	resp := new(meilisearch.DocumentsResult)
	if err := s.meili.Index(index).GetDocuments(&meilisearch.DocumentsQuery{
		Limit:  _defaultHitSizePerPage,
		Filter: filter,
	}, resp); err != nil {
		return nil, err
	}

	results = append(results, resp.Results...)

	if resp.Total > _defaultHitSizePerPage {
		totalOffset := (resp.Total + _defaultHitSizePerPage - 1) / _defaultHitSizePerPage
		for i := int64(1); i < totalOffset; i++ {
			var nextResp meilisearch.DocumentsResult
			if err := s.meili.Index(index).GetDocuments(&meilisearch.DocumentsQuery{
				Offset: i,
				Limit:  _defaultHitSizePerPage,
				Filter: filter,
			}, &nextResp); err != nil {
				return nil, err
			}

			if len(nextResp.Results) == 0 {
				break
			}

			results = append(results, nextResp.Results...)
		}
	}

	return results, nil
}

func (s *Sitemap) saveSitemap(data []byte, indexName string, cfg *config.SitemapConfig) (string, error) {
	fileName := indexName
	if cfg.SitemapFileName != "" {
		fileName = cfg.SitemapFileName
	}

	if s.prefix != "" {
		fileName = s.prefix + fileName
	}

	if cfg.Compress {
		fileName += ".xml.gz"
	} else {
		fileName += ".xml"
	}

	filePath := filepath.Join(s.storePath, s.indexsitemapPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err := file.Write(data); err != nil {
		return "", fmt.Errorf("error writing to file %s: %v", filePath, err)
	}

	return fileName, nil
}

func existsItem(items []string, item string) bool {
	isExists := false

	for i := range items {
		if items[i] == item {
			isExists = true
		}
	}

	return isExists
}
