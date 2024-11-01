package server

import (
	"context"
	"link-base/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

// NewServer initializes and returns a new HTTP server instance
// configured with the provided configuration and request handler.
//
// Parameters:
//   - cfg: An HTTPConfig struct containing the server configuration such as port,
//     maximum header bytes, read timeout, and write timeout.
//   - handler: An http.Handler instance that handles incoming HTTP requests.
//
// Returns:
//   - *Server: A pointer to the initialized HTTP server instance.
func NewServer(cfg config.HTTPConfig, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        handler,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
		},
	}
}

// Run starts the HTTP server and begins listening for incoming requests.
//
// The method returns an error if the server fails to start or if there is a
// problem with the underlying listener.
//
// Note: This method will block until the server is stopped by calling Stop.
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Stop gracefully stops the HTTP server and stops listening for incoming requests.
//
// The method returns an error if the server fails to stop gracefully.
//
// Note: This method will block until all active connections are closed or the context is canceled.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
