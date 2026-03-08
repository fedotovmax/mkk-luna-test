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

// @Summary      Создать новый аккаунт пользователя
// @Description  Создать новый аккаунт пользователя
// @Router       /api/v1/register [post]
// @Tags         users
// @Accept       json
// @Produce      json
// @Param dto body inputs.CreateUser true "Объект для создания аккаунта пользователя"
// @Success      201  {object}  domain.IDResponse
// @Failure      400  {object}  domain.ValidatationErrors
// @Failure      409  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (c *users) register(w http.ResponseWriter, r *http.Request) {

	const op = "controllers.http.v1.user.register"

	l := c.log.With(slog.String("op", op))

	var in inputs.CreateUser

	err := httpcommon.DecodeJSON(r.Body, &in)

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

	id, err := c.registerUc.Execute(ctx, &in)

	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenExecuteRegister, err)
		return
	}

	res := domain.IDResponse{ID: id}

	httpcommon.WriteJSON(w, http.StatusCreated, res)

}

// @Summary      Войти в аккаунт
// @Description  Войти в аккаунт пользователя
// @Router       /api/v1/login [post]
// @Tags         users
// @Accept       json
// @Produce      json
// @Param dto body inputs.Login true "Войти в аккаунт с помощью электронной почты и пароля"
// @Success      201  {object}  domain.LoginResponse
// @Failure      400  {object}  domain.ValidatationErrors
// @Failure      404  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (c *users) login(w http.ResponseWriter, r *http.Request) {

	const op = "controllers.http.v1.user.login"

	l := c.log.With(slog.String("op", op))

	var in inputs.Login

	err := httpcommon.DecodeJSON(r.Body, &in)

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

	res, err := c.loginUc.Execute(ctx, &in)

	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenExecuteLogin, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusCreated, res)

}

func (c *users) RegisterRoutes(r chi.Router) {

	r.Post(registerRoute, c.register)
	r.Post(loginRoute, c.login)

}
