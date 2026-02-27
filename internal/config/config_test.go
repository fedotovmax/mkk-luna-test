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
	tests := []struct {
		name        string
		port        uint16
		redisAddr   string
		redisPass   string
		wantErr     bool
		errorFields []string // для проверки, какие поля упали
	}{
		{
			name:      "valid config",
			port:      8080,
			redisAddr: "redis://localhost:6379",
			redisPass: "secret",
			wantErr:   false,
		},
		{
			name:        "invalid port",
			port:        80,
			redisAddr:   "redis://localhost:6379",
			redisPass:   "secret",
			wantErr:     true,
			errorFields: []string{"HTTPServer.Port"},
		},
		{
			name:        "empty redis password",
			port:        8080,
			redisAddr:   "redis://localhost:6379",
			redisPass:   "",
			wantErr:     true,
			errorFields: []string{"Redis.Password"},
		},
		{
			name:        "multiple errors",
			port:        80,
			redisAddr:   "localhost:6379",
			redisPass:   "",
			wantErr:     true,
			errorFields: []string{"HTTPServer.Port", "Redis.Password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppConfig{
				HTTPServer: &HTTPServerConfig{Port: tt.port},
				Redis: &RedisConfig{
					Addr:     tt.redisAddr,
					Password: tt.redisPass,
				},
			}

			err := cfg.validate()
			if tt.wantErr {
				require.Error(t, err)

				// проверим, что все ожидаемые поля упали
				for _, field := range tt.errorFields {
					assert.Contains(t, err.Error(), field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
