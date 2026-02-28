package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/fedotovmax/mkk-luna-test/docs"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/auth/jwt"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/cache/redis"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/mysql"
	mysqlTx "github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction/mysql"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/server/http"
	"github.com/fedotovmax/mkk-luna-test/internal/config"
	"github.com/fedotovmax/mkk-luna-test/pkg/logger"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type App struct {
	cfg        *config.App
	redispool  *redis.RedisDb
	httpserver *http.Server
	log        *slog.Logger
	dbpool     db.StdSQLDriver
}

func New(cfg *config.App, log *slog.Logger) (*App, error) {

	const op = "app.New"

	redisCtx, cancelRedisCtx := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelRedisCtx()

	redisConn, err := redis.New(redisCtx, cfg.Redis, log)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mysqlCtx, cancelMysqlCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelMysqlCtx()

	mysqlConn, err := mysql.New(mysqlCtx, log, cfg.Database)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	transactionManager, err := mysqlTx.Init(mysqlConn, log.With(slog.String("op", "transaction.manager")))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	txExtractor := transactionManager.GetExtractor()

	_ = txExtractor

	tokenManager := jwt.New(cfg.Tokens.AccessSecret, cfg.Tokens.AccessExpDuration)

	_ = tokenManager

	r := chi.NewRouter()

	r.Handle("/swagger/*", httpSwagger.WrapHandler)

	//TODO: init routes

	httpServer := http.New(cfg.HTTPServer, r)

	app := &App{
		redispool:  redisConn,
		dbpool:     mysqlConn,
		httpserver: httpServer,
		log:        log,
		cfg:        cfg,
	}

	return app, nil
}

func (a *App) Start() <-chan error {
	const op = "app.Start"

	log := a.log.With(slog.String("op", op))

	errChan := make(chan error, 1)

	go func() {
		log.Info(
			"Starting HTTP server...",
			slog.String("addr", fmt.Sprintf("http://localhost:%d", a.cfg.HTTPServer.Port)),
		)
		if err := a.httpserver.Start(); err != nil {
			errChan <- fmt.Errorf("%s: %w", op, err)
		}
	}()

	return errChan
}

func (a *App) Stop(ctx context.Context) {
	const op = "app.Start"

	log := a.log.With(slog.String("op", op))

	if err := a.httpserver.Stop(ctx); err != nil {
		log.Error("Error when shutdown HTTP server", logger.Err(err))
	} else {
		log.Info("HTTP server stopped successfully!")
	}

	if err := a.redispool.Stop(ctx); err != nil {
		log.Error("Error when stop redis", logger.Err(err))
	} else {
		log.Info("Redis stopped successfully!")
	}

	if err := a.dbpool.Stop(ctx); err != nil {
		log.Error("Error when stop DB pool", logger.Err(err))
	} else {
		log.Info("DB pool stopped successfully!")
	}
}
