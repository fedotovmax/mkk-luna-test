package v1

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/httpcommon"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type statistic struct {
	log   *slog.Logger
	teams queries.Teams
}

func NewStatistic(log *slog.Logger, teams queries.Teams) *statistic {
	return &statistic{
		log:   log,
		teams: teams,
	}
}

func (s *statistic) RegisterRoutes(r chi.Router) {

	r.Route("/statistics", func(statRouter chi.Router) {

		statRouter.Get("/teams", s.teamStats)
		statRouter.Get("/top-users", s.topUsers)

	})
}

// @Summary      Статистика по командам
// @Description  Возвращает статистику по всем командам
// @Tags         statistics
// @Produce      json
// @Success      200 {array} domain.TeamStats
// @Failure      500 {object} httpcommon.MessageResponse
// @Router       /api/v1/statistics/teams [get]
func (s *statistic) teamStats(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.http.statistics.v1.teamStats"

	log := s.log.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	stats, err := s.teams.Stats(ctx)
	if err != nil {
		handleErrors(w, log, unexpectedErrorWhenGetTeamStats, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, stats)
}

// @Summary      Топ пользователей
// @Description  Возвращает пользователей с наибольшим количеством задач в командах
// @Tags         statistics
// @Produce      json
// @Success      200 {array} domain.TopUserInTeam
// @Failure      500 {object} httpcommon.MessageResponse
// @Router       /api/v1/statistics/top-users [get]
func (s *statistic) topUsers(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.http.statistics.v1.topUsers"

	log := s.log.With(slog.String("op", op))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users, err := s.teams.TopUsers(ctx)
	if err != nil {
		handleErrors(w, log, unexpectedErrorWhenGetTopUsers, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, users)
}
