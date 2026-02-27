package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvAs_Integrated(t *testing.T) {

	t.Run("success parse int", func(t *testing.T) {
		t.Setenv("TEST_PORT", "8080")
		val, err := getEnvAs[int]("TEST_PORT")
		assert.NoError(t, err)
		assert.Equal(t, 8080, val)
	})

	t.Run("success parse duration", func(t *testing.T) {
		t.Setenv("TEST_TIMEOUT", "30s")
		val, err := getEnvAs[time.Duration]("TEST_TIMEOUT")
		assert.NoError(t, err)
		assert.Equal(t, 30*time.Second, val)
	})

	t.Run("error empty variable", func(t *testing.T) {
		t.Setenv("EMPTY_VAR", "")
		_, err := getEnvAs[string]("EMPTY_VAR")
		assert.ErrorIs(t, err, ErrVariableIsEmpty)
	})

	t.Run("error not provided", func(t *testing.T) {
		_, err := getEnvAs[string]("TOTALLY_MISSING")
		assert.ErrorIs(t, err, ErrVariableNotProvided)
	})
}

func TestNew_Environments(t *testing.T) {
	t.Run("release mode success", func(t *testing.T) {
		t.Setenv("APP_ENV", "release")
		t.Setenv("HTTP_SERVER_PORT", "9000")
		t.Setenv("MYSQL_DSN", "user:pass@tcp(localhost:3306)/db")
		t.Setenv("REDIS_ADDR", "localhost:6381")

		cfg, err := New("")
		require.NoError(t, err)

		assert.Equal(t, Release, cfg.Env)
		assert.Equal(t, uint16(9000), cfg.HTTPServer.Port)
	})

	t.Run("development mode with file", func(t *testing.T) {
		tmpDir := t.TempDir()
		envPath := filepath.Join(tmpDir, ".env")
		content := []byte("HTTP_SERVER_PORT=5000\nMYSQL_DSN=file_dsn\nREDIS_ADDR=localhost:6381")
		err := os.WriteFile(envPath, content, 0644)
		require.NoError(t, err)

		t.Setenv("APP_ENV", "development")

		cfg, err := New(envPath)
		require.NoError(t, err)
		assert.Equal(t, uint16(5000), cfg.HTTPServer.Port)
		assert.Equal(t, "file_dsn", cfg.Database.DSN)
	})
}

func TestAppConfig_Validate(t *testing.T) {

	tests := []struct {
		name    string
		port    uint16
		wantErr bool
	}{
		{"valid port", 8080, false},
		{"invalid port", 80, true},
		{"min port", 1024, false},
		{"max port", 65535, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppConfig{
				HTTPServer: &HTTPServerConfig{Port: tt.port},
			}
			err := cfg.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
