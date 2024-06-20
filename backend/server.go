package zebra

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/cors"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) RunTLS(port, certificate, key string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 5 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServeTLS(certificate, key)
}

func (s *Server) Run(port string, handler http.Handler) error {
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
		AllowedOrigins: []string{"*"},
	})
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        corsWrapper.Handler(handler),
		MaxHeaderBytes: 5 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
