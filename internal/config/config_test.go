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
		t.Setenv("REDIS_PASSWORD", "123456789")
		t.Setenv("ACCESS_TOKEN_SECRET", "sdfsdfdsfsdfsdf")
		t.Setenv("ACCESS_TOKEN_DURATION", "1m")
		t.Setenv("TOKEN_ISSUER", "app")
		t.Setenv("REFRESH_TOKEN_DURATION", "2m")

		cfg, err := New("")
		require.NoError(t, err)

		assert.Equal(t, Release, cfg.Env)
		assert.Equal(t, uint16(9000), cfg.HTTPServer.Port)
	})

	var envFile = `
	HTTP_SERVER_PORT=5000
	MYSQL_DSN=file_dsn
	REDIS_ADDR=localhost:6381
	REDIS_PASSWORD=123456789
	ACCESS_TOKEN_SECRET=dfgdfgdfgdfgdfg
	ACCESS_TOKEN_DURATION=2m
	TOKEN_ISSUER=app
	REFRESH_TOKEN_DURATION=3m
	`

	t.Run("development mode with file", func(t *testing.T) {
		tmpDir := t.TempDir()
		envPath := filepath.Join(tmpDir, ".env")
		content := []byte(envFile)
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
	baseConfig := func() *App {
		return &App{
			HTTPServer: &HTTPServer{
				Port: 8080,
			},
			Database: &Database{
				DSN: "user:pass@tcp(localhost:3306)/db",
			},
			Redis: &Redis{
				Addr:     "redis://localhost:6379",
				Password: "secret",
			},
			Tokens: &Tokens{
				AccessExpDuration:  time.Minute * 1,
				RefreshExpDuration: time.Minute * 2,
				AccessSecret:       "secret",
				Issuer:             "issuer",
			},
		}
	}

	tests := []struct {
		name        string
		modify      func(cfg *App)
		wantErr     bool
		errorFields []string
	}{
		{
			name:    "valid config",
			modify:  func(cfg *App) {},
			wantErr: false,
		},
		{
			name: "invalid port",
			modify: func(cfg *App) {
				cfg.HTTPServer.Port = 80
			},
			wantErr:     true,
			errorFields: []string{"HTTPServer.Port"},
		},
		{
			name: "empty redis password",
			modify: func(cfg *App) {
				cfg.Redis.Password = ""
			},
			wantErr:     true,
			errorFields: []string{"Redis.Password"},
		},
		{
			name: "multiple errors",
			modify: func(cfg *App) {
				cfg.HTTPServer.Port = 80
				cfg.Redis.Password = ""
			},
			wantErr:     true,
			errorFields: []string{"HTTPServer.Port", "Redis.Password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := baseConfig()
			tt.modify(cfg)

			err := cfg.validate()

			if tt.wantErr {
				require.Error(t, err)

				for _, field := range tt.errorFields {
					assert.Contains(t, err.Error(), field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
