package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/fedotovmax/mkk-luna-test/internal/validation"
	"github.com/joho/godotenv"
)

var ErrInvalidAppEnv = errors.New("app env is invalid or not supported")
var ErrConfigPathNotExists = errors.New("config path for dev env is not exists")

type AppEnv string

const (
	Development AppEnv = "development"
	Release     AppEnv = "release"
)

func parseEnvVariable(env string) (AppEnv, error) {
	switch env {
	case string(Development):
		return Development, nil
	case string(Release):
		return Release, nil
	default:
		return "", ErrInvalidAppEnv
	}
}

func checkConfigPathExists(path string) error {

	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %v", ErrConfigPathNotExists, err)
		}
		return err
	}

	return nil
}

type HTTPServerConfig struct {
	Port uint16
}

type AppConfig struct {
	HTTPServer *HTTPServerConfig
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

	cfg := &AppConfig{
		HTTPServer: &HTTPServerConfig{
			Port: httpServerPort,
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

	return errors.Join(validationErrors...)
}
