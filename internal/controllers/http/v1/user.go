package v1

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/httpcommon"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/usecases"
	"github.com/fedotovmax/mkk-luna-test/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type users struct {
	log        *slog.Logger
	registerUc *usecases.Register
	loginUc    *usecases.Login
}

func NewUsers(register *usecases.Register, login *usecases.Login, log *slog.Logger) *users {
	return &users{registerUc: register, loginUc: login, log: log}
}

// @Summary      Create user account
// @Description  Create new user account
// @Router       /api/v1/register [post]
// @Tags         users
// @Accept       json
// @Produce      json
// @Param dto body inputs.CreateUser true "Create user account with body dto"
// @Success      201  {object}  domain.IDResponse
// @Failure      400  {object}  domain.ValidatationErrors
// @Failure      409  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (c *users) register(w http.ResponseWriter, r *http.Request) {

	const op = "controllers.http.user.register"

	l := c.log.With(slog.String("op", op))

	var in inputs.CreateUser

	err := httpcommon.DecodeJSON(r.Body, &in)

	if err != nil {
		l.Error("error when parse request body", logger.Err(err))
		httpcommon.WriteJSON(w, http.StatusBadRequest, httpcommon.Message(invalidBodyFormat))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	id, err := c.registerUc.Execute(ctx, &in)

	if err != nil {
		handleErrors(l, unexpectedErrorWhenExecuteRegister, err)
		return
	}

	res := domain.IDResponse{ID: id}

	httpcommon.WriteJSON(w, http.StatusCreated, res)

}

func (c *users) RegisterRoutes(r chi.Router) {

	r.Post(registerRoute, c.register)
	r.Post(loginRoute, c.register)

}
