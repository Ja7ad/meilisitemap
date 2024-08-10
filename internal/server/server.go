package server

import (
	"context"
	"expvar"
	"github.com/Ja7ad/meilisitemap/config"
	"net/http"
	"net/http/pprof"
)

type Server struct {
	server *http.Server
	notify chan error
	listen string
}

func New(serve *config.ServeConfig, storePath string) *Server {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(storePath))
	mux.Handle("/", http.StripPrefix("/", fileServer))

	if serve.PPROF {
		debuggerHandler(mux)
	}

	return &Server{
		server: &http.Server{
			Addr:    serve.Listen,
			Handler: mux,
		},
		listen: serve.Listen,
		notify: make(chan error, 1),
	}
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.listen
}

func debuggerHandler(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())
	return mux
}
