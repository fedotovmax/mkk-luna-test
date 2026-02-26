package http

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestServer_StartAndStop(t *testing.T) {
	// используем случайный порт)
	cfg := &config.HTTPServerConfig{Port: 0}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := New(cfg, handler)

	startErr := make(chan error, 1)
	go func() {
		err := srv.Start()
		if err != nil {
			startErr <- err
		}
	}()

	time.Sleep(100 * time.Millisecond)

	t.Run("Graceful Shutdown", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := srv.Stop(ctx)
		assert.NoError(t, err, "http server should stop gracefully")
	})

	select {
	case err := <-startErr:
		t.Fatalf("http server start error: %v", err)
	default:
	}
}
