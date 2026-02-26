package cmd

import (
	"flag"
	"log/slog"
	"os"

	"github.com/fedotovmax/mkk-luna-test/internal/config"
	"github.com/fedotovmax/mkk-luna-test/pkg/logger"
)

func setupLooger(env config.AppEnv) *slog.Logger {
	if env == config.Development {
		return logger.NewHandler(slog.LevelDebug)
	}
	return logger.NewHandler(slog.LevelWarn)
}

func loadConfigPathFlags() string {

	const op = "config.loadConfigPath"

	var configPath string

	flag.StringVar(&configPath, "config_path", "", "path to config file")
	flag.StringVar(&configPath, "c", "", "path to config file (shorthand)")

	flag.Parse()

	return configPath
}

func main() {

	configPath := loadConfigPathFlags()

	appConfig, err := config.New(configPath)

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	log := setupLooger(appConfig.Env)

	log.Info("Logger setup")

}
