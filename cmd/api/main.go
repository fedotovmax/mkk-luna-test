package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/app"
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

// @title Swagger Documentation for mkk_luna_test_rest_api
// @version 1.0
// @description Swagger Documentation for mkk_luna_test_rest_api (тестовое задание для компании "МКК ЛУНА")
// @contact.name Fedotv Maxim (developer)
// @contact.email f3d0t0tvmax@yandex.ru
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите "Bearer [ваш JWT токен]"
func main() {

	configPath := loadConfigPathFlags()

	appConfig, err := config.New(configPath)

	if err != nil {
		logger.GetFallback().Error(err.Error())
		os.Exit(1)
	}

	log := setupLooger(appConfig.Env)

	app, err := app.New(appConfig, log)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	notifyCtx, cancelNotifyCtx := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancelNotifyCtx()

	startErrChan := app.Start()

	select {
	case err := <-startErrChan:
		log.Error("Error recivied when starting application", logger.Err(err))
		cancelNotifyCtx()
	case <-notifyCtx.Done():
		log.Info("OS signal recevied")
	}

	log.Info("Starting to shutdown all resources")

	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdownCtx()

	app.Stop(shutdownCtx)

}
