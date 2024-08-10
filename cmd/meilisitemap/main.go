package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ja7ad/meilisitemap/config"
	"github.com/Ja7ad/meilisitemap/internal/generator"
	"github.com/Ja7ad/meilisitemap/internal/logger"
)

const _defaultStoreDir = "sitemap"

func main() {
	configPath := flag.String("config", "./config.json", "path to config file")
	storePath := flag.String("store", _defaultStoreDir, "path to store sitemap")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	log := logger.DefaultLogger

	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatal("failed to load config", "err", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal("invalid config", "err", err)
	}

	log.Info("configuration file loaded")

	if _, err := os.Stat(*storePath); os.IsNotExist(err) {
		if err := os.Mkdir(*storePath, 0o777); err != nil {
			log.Fatal("failed to create directory", "err", err)
		}
	} else if err != nil {
		log.Fatal("failed to check if directory exists", "err", err)
	}
	sm, err := generator.New(
		ctx,
		*storePath,
		cfg.General,
		log,
		cfg.Sitemaps,
	)
	if err != nil {
		log.Fatal("failed to initialize generator", "err", err)
	}

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		log.Warn("cancellation signal received.")
		cancel()
	}()

	if err := sm.Start(); err != nil {
		log.Fatal("failed to start sitemap", "err", err)
	}
}
