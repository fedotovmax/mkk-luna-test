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
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/mysql/sessions"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/mysql/tasks"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/mysql/teams"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/mysql/users"
	mysqlTx "github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction/mysql"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/server/http"
	"github.com/fedotovmax/mkk-luna-test/internal/config"
	v1 "github.com/fedotovmax/mkk-luna-test/internal/controllers/http/v1"
	"github.com/fedotovmax/mkk-luna-test/internal/middlewares"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
	"github.com/fedotovmax/mkk-luna-test/internal/usecases"
	"github.com/fedotovmax/mkk-luna-test/pkg/logger"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"golang.org/x/time/rate"
)

type App struct {
	cfg         *config.App
	redispool   *redis.RedisDb
	httpserver  *http.Server
	rateLimiter *middlewares.UserRateLimiter
	log         *slog.Logger
	dbpool      db.StdSQLDriver
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

	teamsMysql := teams.New(txExtractor)
	tasksMysql := tasks.New(txExtractor)
	usersMysql := users.New(txExtractor)
	sessionsMysql := sessions.New(txExtractor)

	teamsQuery := queries.NewTeams(teamsMysql)
	tasksQuery := queries.NewTasks(tasksMysql)
	usersQuery := queries.NewUsers(usersMysql)

	tokenManager := jwt.New(cfg.Tokens.AccessSecret)

	loginUsecase := usecases.NewLogin(log, usersQuery, sessionsMysql, tokenManager, cfg.Tokens)
	registerUsecase := usecases.NewRegister(log, usersMysql, usersQuery)
	createTeamUsecase := usecases.NewCreateTeam(log, transactionManager, teamsMysql, teamsQuery)
	inviteUsecase := usecases.NewInvite(log, teamsMysql, teamsQuery)
	createTaskUsecase := usecases.NewCreateTask(log, transactionManager, tasksMysql, teamsQuery)
	updateTaskUsecase := usecases.NewUpdateTask(log, transactionManager, tasksMysql, tasksQuery, teamsQuery)
	getTaskHistoryUsecase := usecases.NewGetTaskHistory(log, tasksQuery, teamsQuery)
	getTaskCommentsUsecase := usecases.NewGetTaskComments(log, tasksQuery, teamsQuery)
	getTasksUsecase := usecases.NewGetTasks(log, tasksQuery, teamsQuery)
	createCommentUsecase := usecases.NewCreateComment(log, tasksMysql, tasksQuery, teamsQuery)

	authMiddleware := middlewares.NewAuthMiddleware(log, tokenManager, cfg.Tokens.Issuer)
	rateLimiter := middlewares.NewUserRateLimiter(rate.Every(time.Minute/100), 10)

	r := chi.NewRouter()

	r.Handle("/swagger/*", httpSwagger.WrapHandler)

	usersController := v1.NewUsers(registerUsecase, loginUsecase, log)
	teamController := v1.NewTeams(
		log,
		createTeamUsecase,
		inviteUsecase,
		teamsQuery,
		authMiddleware,
		rateLimiter,
	)
	taskController := v1.NewTasks(
		log,
		createTaskUsecase,
		updateTaskUsecase,
		getTaskHistoryUsecase,
		getTaskCommentsUsecase,
		createCommentUsecase,
		getTasksUsecase,
		authMiddleware,
		rateLimiter,
		tasksQuery,
	)
	statsController := v1.NewStatistic(log, teamsQuery)

	r.Route("/api/v1", func(r chi.Router) {
		usersController.RegisterRoutes(r)
		statsController.RegisterRoutes(r)
		r.Group(func(withMiddlewaresRouter chi.Router) {

			withMiddlewaresRouter.Use(authMiddleware.Middleware)
			withMiddlewaresRouter.Use(rateLimiter.Middleware)

			taskController.RegisterRoutes(withMiddlewaresRouter)
			teamController.RegisterRoutes(withMiddlewaresRouter)
		})
	})

	httpServer := http.New(cfg.HTTPServer, r)

	app := &App{
		redispool:   redisConn,
		dbpool:      mysqlConn,
		httpserver:  httpServer,
		rateLimiter: rateLimiter,
		log:         log,
		cfg:         cfg,
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

		log.Info("Swagger documentation is available at", slog.String("addr", fmt.Sprintf("http://localhost:%d/swagger/", a.cfg.HTTPServer.Port)))

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

	a.rateLimiter.Stop()
	log.Info("Rate limiter stopped")
}
