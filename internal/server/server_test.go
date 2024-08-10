package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Ja7ad/meilisitemap/config"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	serveConfig := &config.ServeConfig{
		Listen: "127.0.0.1:8080",
		PPROF:  false,
	}
	storePath := "./testdata"

	server := New(serveConfig, storePath)

	assert.NotNil(t, server)
	assert.Equal(t, serveConfig.Listen, server.server.Addr)
	assert.NotNil(t, server.server.Handler)
	assert.NotNil(t, server.notify)
}

func TestServerStartAndNotify(t *testing.T) {
	serveConfig := &config.ServeConfig{
		Listen: "127.0.0.1:8081",
		PPROF:  false,
	}
	storePath := "./testdata"

	server := New(serveConfig, storePath)
	assert.NotNil(t, server)

	server.Start()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:8081/testdata/sitemap.xml")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	select {
	case err := <-server.Notify():
		assert.NoError(t, err)
	default:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServerWithPPROF(t *testing.T) {
	serveConfig := &config.ServeConfig{
		Listen: "127.0.0.1:8082",
		PPROF:  true,
	}
	storePath := "./testdata"

	server := New(serveConfig, storePath)
	assert.NotNil(t, server)

	server.Start()

	time.Sleep(100 * time.Millisecond)

	endpoints := []string{
		"/debug/pprof/",
	}

	for _, endpoint := range endpoints {
		resp, err := http.Get("http://127.0.0.1:8082" + endpoint)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}
