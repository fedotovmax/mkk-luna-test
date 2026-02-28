package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/validation"
	"github.com/joho/godotenv"
)

type HTTPServerConfig struct {
	Port uint16
}

type DatabaseCofnig struct {
	RetryWait  time.Duration // by default is 500ms
	DSN        string
	MaxRetries uint8 // by default is 1
}

func (c *DatabaseCofnig) SetMaxRetries(value uint8) {
	c.MaxRetries = value
}

func (c *DatabaseCofnig) SetRetryWait(value time.Duration) {
	c.RetryWait = value
}

type RedisConfig struct {
	RetryWait  time.Duration // by default is 200ms
	Addr       string
	Password   string
	DB         int
	MaxRetries uint8 // by default is 1
}

func (c *RedisConfig) SetMaxRetries(value uint8) {
	c.MaxRetries = value
}

func (c *RedisConfig) SetRetryWait(value time.Duration) {
	c.RetryWait = value
}

type Jwt struct {
	AccessTokenExpDuration time.Duration
	AccessTokenSecret      string
}

type AppConfig struct {
	HTTPServer *HTTPServerConfig
	Database   *DatabaseCofnig
	Redis      *RedisConfig
	Jwt        *Jwt
	Env        AppEnv
}

// Load config from file, when required APP_ENV variable provided and equal to local or development,
// And required flag "config_path" for *.env file with variables: -c or -config_path
// Else get env variables provided by operation system (defined by user/container/environment)
func New(path string) (*AppConfig, error) {

	const op = "config.New"

	envString, err := getEnv("APP_ENV")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	env, err := parseEnvVariable(envString)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if env == Development {

		err := checkConfigPathExists(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		err = godotenv.Load(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

	}

	httpServerPort, err := getEnvAs[uint16]("HTTP_SERVER_PORT")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mysqlDsn, err := getEnv("MYSQL_DSN")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	redisAddr, err := getEnv("REDIS_ADDR")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	redisPassword, err := getEnv("REDIS_PASSWORD")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	cfg := &AppConfig{
		HTTPServer: &HTTPServerConfig{
			Port: httpServerPort,
		},
		Database: &DatabaseCofnig{
			DSN:        mysqlDsn,
			MaxRetries: 1,
			RetryWait:  time.Millisecond * 500,
		},
		Redis: &RedisConfig{
			Addr:       redisAddr,
			Password:   redisPassword,
			MaxRetries: 1,
			RetryWait:  time.Millisecond * 200,
		},
		Env: env,
	}

	err = cfg.validate()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cfg, nil

}

func (c *AppConfig) validate() error {

	var validationErrors []error

	err := validation.Range(c.HTTPServer.Port, 1024, 65535)

	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("%s: %w", "HTTPServer.Port", err))
	}

	err = validation.EmptyString(c.Redis.Password)

	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("%s: %w", "Redis.Password", err))
	}

	_, err = validation.IsURI(c.Redis.Addr)

	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("%s: %w", "Redis.Addr", err))
	}

	return errors.Join(validationErrors...)
}
