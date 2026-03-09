package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/validation"
	"github.com/joho/godotenv"
)

type HTTPServer struct {
	Port uint16
}

type Database struct {
	RetryWait  time.Duration // by default is 500ms
	DSN        string
	MaxRetries uint8 // by default is 1
}

func (c *Database) SetMaxRetries(value uint8) {
	c.MaxRetries = value
}

func (c *Database) SetRetryWait(value time.Duration) {
	c.RetryWait = value
}

type Redis struct {
	RetryWait  time.Duration // by default is 200ms
	Addr       string
	Password   string
	DB         int
	MaxRetries uint8 // by default is 1
}

func (c *Redis) SetMaxRetries(value uint8) {
	c.MaxRetries = value
}

func (c *Redis) SetRetryWait(value time.Duration) {
	c.RetryWait = value
}

type Tokens struct {
	AccessExpDuration  time.Duration
	RefreshExpDuration time.Duration
	AccessSecret       string
	Issuer             string
}

type App struct {
	HTTPServer *HTTPServer
	Database   *Database
	Redis      *Redis
	Tokens     *Tokens
	Env        AppEnv
}

// Load config from file, when required APP_ENV variable provided and equal to development,
// And required flag "config_path" for *.env file with variables: -c or -config_path
// Else get env variables provided by operation system (defined by user/container/environment)
func New(path string) (*App, error) {

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

	accessTokenSecret, err := getEnv("ACCESS_TOKEN_SECRET")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	accessTokenExpDuration, err := getEnvAs[time.Duration]("ACCESS_TOKEN_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	refreshTokenExpDuration, err := getEnvAs[time.Duration]("REFRESH_TOKEN_DURATION")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tokenIssuer, err := getEnv("TOKEN_ISSUER")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	cfg := &App{
		HTTPServer: &HTTPServer{
			Port: httpServerPort,
		},
		Database: &Database{
			DSN:        mysqlDsn,
			MaxRetries: 1,
			RetryWait:  time.Millisecond * 500,
		},
		Redis: &Redis{
			Addr:       redisAddr,
			Password:   redisPassword,
			MaxRetries: 1,
			RetryWait:  time.Millisecond * 200,
		},
		Tokens: &Tokens{
			AccessExpDuration:  accessTokenExpDuration,
			RefreshExpDuration: refreshTokenExpDuration,
			AccessSecret:       accessTokenSecret,
			Issuer:             tokenIssuer,
		},
		Env: env,
	}

	err = cfg.validate()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cfg, nil

}

func (c *App) validate() error {

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

	err = validation.Min(c.Tokens.AccessExpDuration, 1)

	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("%s: %w", "Tokens.AccessExpDuration", err))
	}

	err = validation.Min(c.Tokens.RefreshExpDuration, 1)

	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("%s: %w", "Tokens.RefreshExpDuration", err))
	}

	err = validation.MinLength(c.Tokens.Issuer, 1)

	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("%s: %w", "Tokens.Issuer", err))
	}

	err = validation.MinLength(c.Tokens.AccessSecret, 1)

	if err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("%s: %w", "Tokens.AccessSecret", err))
	}

	return errors.Join(validationErrors...)
}
