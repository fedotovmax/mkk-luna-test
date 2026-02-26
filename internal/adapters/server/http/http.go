package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fedotovmax/mkk-luna-test/internal/config"
)

type Server struct {
	instance *http.Server
}

func New(httpServerConfig *config.HTTPServerConfig, handler http.Handler) *Server {
	inst := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpServerConfig.Port),
		Handler: handler,
	}

	return &Server{
		instance: inst,
	}
}

func (srv *Server) Start() error {

	const op = "adapters.server.http.Start"

	if err := srv.instance.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (srv *Server) Stop(ctx context.Context) error {
	const op = "adapters.server.http.Stop"

	if err := srv.instance.Shutdown(ctx); err != nil {

		if closeErr := srv.instance.Close(); closeErr != nil {
			return fmt.Errorf(
				"%s: shutdown error: %v, force close error: %w",
				op, err, closeErr,
			)
		}

		return fmt.Errorf("%s: graceful shutdown failed: %w", op, err)
	}

	return nil
}
