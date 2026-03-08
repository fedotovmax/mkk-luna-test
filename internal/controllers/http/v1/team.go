package v1

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/httpcommon"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/middlewares"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
	"github.com/fedotovmax/mkk-luna-test/internal/usecases"
	"github.com/go-chi/chi/v5"
)

type teams struct {
	log         *slog.Logger
	createTeam  *usecases.CreateTeam
	inviteUc    *usecases.Invite
	query       queries.Teams
	checkAuth   middlewares.Middleware
	rateLimiter middlewares.Middleware
}

func NewTeams(
	log *slog.Logger,
	createTeam *usecases.CreateTeam,
	inviteUc *usecases.Invite,
	query queries.Teams,
	checkAuth middlewares.Middleware,
	rateLimiter middlewares.Middleware,
) *teams {
	return &teams{log: log, createTeam: createTeam, inviteUc: inviteUc, query: query, checkAuth: checkAuth, rateLimiter: rateLimiter}
}

// @Summary      Создать команду
// @Description  Создать команду
// @Router       /api/v1/teams [post]
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param dto body inputs.CreateTeam true "Объект для создания команды"
// @Success      201  {object}  domain.IDResponse
// @Failure      400  {object}  domain.ValidatationErrors
// @Failure      401  {object}  httpcommon.MessageResponse
// @Failure      409  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (c *teams) create(w http.ResponseWriter, r *http.Request) {

	const op = "controllers.http.team.v1.create"

	l := c.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)

	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	var in inputs.CreateTeam

	err = httpcommon.DecodeJSON(r.Body, &in)

	if err != nil {
		httpcommon.WriteJSON(w, http.StatusBadRequest, httpcommon.Message(invalidBodyFormat))
		return
	}

	err = in.Validate()

	if err != nil {
		handleValidationErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	teamID, err := c.createTeam.Execute(ctx, local.UserID, &in)

	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenExecuteCreateTeam, err)
		return
	}

	res := domain.IDResponse{ID: teamID}

	httpcommon.WriteJSON(w, http.StatusCreated, res)

}

// @Summary Получить команды пользователя
// @Description олучить команды пользователя
// @Router /api/v1/teams [get]
// @Tags teams
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success      200  {object}  domain.FindTeamsResponse
// @Failure      401  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (c *teams) getUserTeams(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.http.team.v1.getUserTeams"

	l := c.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)

	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	res, err := c.query.FindMany(ctx, 0, 0, local.UserID)

	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenGetUserTeams, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, res)

}

// @Summary      Пригласить пользователя в команду
// @Description  Пригласить пользователя в команду
// @Router       /api/v1/teams/{id}/invite [post]
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                 true  "Team ID" format(uuid)
// @Param        dto  body      inputs.InviteMember    true  "Invite user to team dto"
// @Success      201  {object}  domain.IDResponse
// @Failure      400  {object}  domain.ValidatationErrors
// @Failure      401  {object}  httpcommon.MessageResponse
// @Failure      403  {object}  httpcommon.MessageResponse
// @Failure      404  {object}  httpcommon.MessageResponse
// @Failure      409  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (c *teams) invite(w http.ResponseWriter, r *http.Request) {

	const op = "controllers.http.team.v1.invite"

	l := c.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)

	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	teamID := inputs.UUID{
		ID: r.PathValue("id"),
	}

	err = teamID.Validate("team_id")

	if err != nil {
		handleValidationErrors(w, err)
		return
	}

	var in inputs.InviteMember

	err = httpcommon.DecodeJSON(r.Body, &in)

	if err != nil {
		httpcommon.WriteJSON(w, http.StatusBadRequest, httpcommon.Message(invalidBodyFormat))
		return
	}

	err = in.Validate()

	if err != nil {
		handleValidationErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	memberID, err := c.inviteUc.Execute(ctx, local.UserID, teamID.ID, in.UserID)

	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenInviteUser, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusCreated, domain.IDResponse{
		ID: memberID,
	})
}

func (c *teams) RegisterRoutes(r chi.Router) {

	r.Route("/teams", func(teamRouter chi.Router) {
		teamRouter.Post("/", c.create)
		teamRouter.Get("/", c.getUserTeams)

		teamRouter.Route("/{id}", func(oneTeamRouter chi.Router) {
			oneTeamRouter.Post("/invite", c.invite)
		})
	})

}
