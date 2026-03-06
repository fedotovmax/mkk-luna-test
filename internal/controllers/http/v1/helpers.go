package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/httpcommon"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
)

func handleErrors(w http.ResponseWriter, log *slog.Logger, fallbackMessage string, err error) {
	if err == nil {
		return
	}

	log.Debug("request failed", slog.String("error", err.Error()))

	if errors.Is(err, errs.ErrSessionExpired) {
		httpcommon.WriteJSON(w, http.StatusUnauthorized,
			httpcommon.Message("Сессия истекла, выполните вход заново"))
		return
	}

	if errors.Is(err, errs.ErrBadCredentials) {
		httpcommon.WriteJSON(w, http.StatusNotFound,
			httpcommon.Message("Неверный email или пароль"))
		return
	}

	if errors.Is(err, errs.ErrNoRightsToDeleteTaskComment) {
		httpcommon.WriteJSON(w, http.StatusForbidden,

			httpcommon.Message("Недостаточно прав для удаления комментария"))
		return
	}

	if errors.Is(err, errs.ErrNoRightsToUpdateTask) {
		httpcommon.WriteJSON(w, http.StatusForbidden,
			httpcommon.Message("Недостаточно прав для обновления задачи"))
		return
	}

	if errors.Is(err, errs.ErrNoRightsToDeleteMember) {
		httpcommon.WriteJSON(w, http.StatusForbidden,
			httpcommon.Message("Недостаточно прав для удаления участника"))
		return
	}

	if errors.Is(err, errs.ErrNoRightsToInviteMember) {
		httpcommon.WriteJSON(w, http.StatusForbidden,
			httpcommon.Message("Недостаточно прав для приглашения участника"))
		return
	}

	if errors.Is(err, errs.ErrUserNotInTaskTeam) {
		httpcommon.WriteJSON(w, http.StatusForbidden,
			httpcommon.Message("Пользователь не состоит в команде этой задачи"))
		return
	}

	if errors.Is(err, errs.ErrSessionNotFound) {
		httpcommon.WriteJSON(w, http.StatusNotFound,
			httpcommon.Message("Сессия не найдена"))
		return
	}

	if errors.Is(err, errs.ErrTaskNotFound) {
		httpcommon.WriteJSON(w, http.StatusNotFound,
			httpcommon.Message("Задача не найдена"))
		return
	}

	if errors.Is(err, errs.ErrTeamNotFound) {
		httpcommon.WriteJSON(w, http.StatusNotFound,
			httpcommon.Message("Команда не найдена"))
		return
	}

	if errors.Is(err, errs.ErrTeamMemberNotFound) {
		httpcommon.WriteJSON(w, http.StatusNotFound,
			httpcommon.Message("Участник команды не найден"))
		return
	}

	if errors.Is(err, errs.ErrUserNotFound) {
		httpcommon.WriteJSON(w, http.StatusNotFound,
			httpcommon.Message("Пользователь не найден"))
		return
	}

	if errors.Is(err, errs.ErrTeamAlreadyExists) {
		httpcommon.WriteJSON(w, http.StatusConflict,
			httpcommon.Message("Команда с таким названием уже существует"))
		return
	}

	if errors.Is(err, errs.ErrUserAlreadyInTeam) {
		httpcommon.WriteJSON(w, http.StatusConflict,
			httpcommon.Message("Пользователь уже состоит в команде"))
		return
	}

	if errors.Is(err, errs.ErrUserAlreadyExists) {
		httpcommon.WriteJSON(w, http.StatusConflict,
			httpcommon.Message("Пользователь с таким email уже существует"))
		return
	}

	msg := fallbackMessage
	if msg == "" {
		msg = "Произошла ошибка непредвиденная ошибка"
	}

	httpcommon.WriteJSON(w, http.StatusInternalServerError, httpcommon.Message(msg))
}

func handleValidationErrors(w http.ResponseWriter, err error) {

	if ve, ok := errors.AsType[*domain.ValidatationErrors](err); ok {
		httpcommon.WriteJSON(w, http.StatusBadRequest, ve)
		return
	}

	httpcommon.WriteJSON(w, http.StatusBadRequest, httpcommon.Message(err.Error()))
}
