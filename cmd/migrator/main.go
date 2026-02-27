package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/fedotovmax/mkk-luna-test/internal/config"
	"github.com/fedotovmax/mkk-luna-test/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migratorFlags struct {
	Cmd        *string
	Version    *int
	Steps      *int
	ConfigPath string
}

func loadMigratorFlags() (*migratorFlags, error) {

	const op = "config.loadMigratorFlags"

	migrationCommand := flag.String("m", "up", "migration command: up, down, force, version")
	version := flag.Int("version", 0, "version for force migration")
	steps := flag.Int("steps", 0, "number of steps for up/down migration")

	var configPath string

	flag.StringVar(&configPath, "config_path", "", "path to config file")
	flag.StringVar(&configPath, "c", "", "path to config file (shorthand)")

	flag.Parse()

	return &migratorFlags{
		Cmd:        migrationCommand,
		Version:    version,
		Steps:      steps,
		ConfigPath: configPath,
	}, nil
}

func main() {

	log := logger.NewHandler(slog.LevelDebug)

	flags, err := loadMigratorFlags()

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	cfg, err := config.NewMigrator(flags.ConfigPath)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	migrationsPath := "file://" + cfg.MigrationsPath

	url := fmt.Sprintf("mysql://%s&multiStatements=true&x-migrations-table=%s", cfg.DSN, cfg.MigrationsTable)

	m, err := migrate.New(
		migrationsPath,
		url,
	)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer m.Close()

	switch *flags.Cmd {
	case "up":
		if *flags.Steps > 0 {
			err = m.Steps(*flags.Steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *flags.Steps > 0 {
			err = m.Steps(-*flags.Steps)
		} else {
			err = m.Down()
		}
	case "force":
		err = m.Force(*flags.Version)
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		log.Info("current migration version", "version", version, "dirty", dirty)
		return
	default:
		log.Error(fmt.Sprintf("unknown migration command, command: %s", *flags.Cmd))
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Error(err.Error())
		os.Exit(1)
	}
	log.Info("migration completed successfully", "command", *flags.Cmd)
}
