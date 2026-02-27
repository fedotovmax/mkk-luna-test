package config

import (
	"errors"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/validation"
	"github.com/joho/godotenv"
)

type MigratorConfig struct {
	DSN             string
	MigrationsPath  string
	MigrationsTable string
	Env             AppEnv
}

func NewMigrator(path string) (*MigratorConfig, error) {
	const op = "config.NewMigrator"

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

	mysqlDsn, err := getEnv("MYSQL_DSN")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	migrationsPath, err := getEnv("MIGRATIONS_PATH")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	migrationsTable, err := getEnv("MIGRATIONS_TABLE")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	c := &MigratorConfig{
		DSN:             mysqlDsn,
		MigrationsPath:  migrationsPath,
		MigrationsTable: migrationsTable,
	}

	err = c.validate()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil

}

func (mc *MigratorConfig) validate() error {
	var verrs []error

	err := validation.IsFilePath(mc.MigrationsPath)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "MigrationsPath", err))
	}

	err = validation.MinLength(mc.MigrationsTable, 1)

	if err != nil {
		verrs = append(verrs, fmt.Errorf("%s: %w", "MigrationsTable", err))
	}

	return errors.Join(verrs...)

}
