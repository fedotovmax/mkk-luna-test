package v1

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/httpcommon"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/middlewares"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
	"github.com/fedotovmax/mkk-luna-test/internal/usecases"
	"github.com/go-chi/chi/v5"
)

type tasks struct {
	log             *slog.Logger
	checkAuth       middlewares.Middleware
	rateLimiter     middlewares.Middleware
	createUc        *usecases.CreateTask
	updateUc        *usecases.UpdateTask
	historyUc       *usecases.GetTaskHistory
	commentsUc      *usecases.GetTaskComments
	createCommentUc *usecases.CreateComment // Добавлено
	getTaskUc       *usecases.GetTasks
	query           queries.Tasks
}

func NewTasks(
	log *slog.Logger,
	createUc *usecases.CreateTask,
	updateUc *usecases.UpdateTask,
	historyUc *usecases.GetTaskHistory,
	commentsUc *usecases.GetTaskComments,
	createCommentUc *usecases.CreateComment,
	getTaskUc *usecases.GetTasks,
	checkAuth middlewares.Middleware,
	rateLimiter middlewares.Middleware,
	query queries.Tasks,
) *tasks {
	return &tasks{
		log:             log,
		createUc:        createUc,
		updateUc:        updateUc,
		historyUc:       historyUc,
		commentsUc:      commentsUc,
		createCommentUc: createCommentUc,
		getTaskUc:       getTaskUc,
		checkAuth:       checkAuth,
		rateLimiter:     rateLimiter,
		query:           query,
	}
}

// @Summary      Создать задачу
// @Description  Создать задачу в команде, пользователь должен быть членом команды
// @Router       /api/v1/tasks [post]
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        dto body      inputs.CreateTask true "Create task dto"
// @Success      201 {object}  domain.IDResponse
// @Failure      400 {object}  domain.ValidatationErrors
// @Failure      401 {object}  httpcommon.MessageResponse
// @Failure      403 {object}  httpcommon.MessageResponse
// @Failure      404 {object}  httpcommon.MessageResponse
// @Failure      500 {object}  httpcommon.MessageResponse
func (t *tasks) create(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.http.tasks.v1.create"
	l := t.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)
	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	var in inputs.CreateTask
	if err := httpcommon.DecodeJSON(r.Body, &in); err != nil {
		l.Error("не удалось распарсить тело запроса", slog.Any("err", err))
		httpcommon.WriteJSON(w, http.StatusBadRequest, httpcommon.Message(invalidBodyFormat))
		return
	}

	if err := in.Validate(); err != nil {
		handleValidationErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	taskID, err := t.createUc.Execute(ctx, local.UserID, &in)
	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenExecuteCreateTask, err)
		return
	}
	res := domain.IDResponse{ID: taskID}

	httpcommon.WriteJSON(w, http.StatusCreated, res)
}

// @Summary      Обновить задачу
// @Description  Обновить задачу. Пользователь должен иметь права (owner/admin/assignee)
// @Router       /api/v1/tasks/{id} [put]
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path string           true "Task ID" format(uuid)
// @Param        dto body inputs.UpdateTask true "Update task dto"
// @Success      200 {object} httpcommon.MessageResponse
// @Failure      400 {object} domain.ValidatationErrors
// @Failure      401 {object} httpcommon.MessageResponse
// @Failure      403 {object} httpcommon.MessageResponse
// @Failure      404 {object} httpcommon.MessageResponse
// @Failure      500 {object} httpcommon.MessageResponse
func (t *tasks) update(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.http.tasks.v1.update"
	l := t.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)
	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	taskID := r.PathValue("id")
	uuidInput := inputs.UUID{ID: taskID}
	if err := uuidInput.Validate("id"); err != nil {
		handleValidationErrors(w, err)
		return
	}

	var in inputs.UpdateTask
	if err := httpcommon.DecodeJSON(r.Body, &in); err != nil {
		l.Error("не удалось распарсить тело запроса", slog.Any("err", err))
		httpcommon.WriteJSON(w, http.StatusBadRequest, httpcommon.Message(invalidBodyFormat))
		return
	}

	if err := in.Validate(); err != nil {
		handleValidationErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := t.updateUc.Execute(ctx, local.UserID, uuidInput.ID, &in); err != nil {
		handleErrors(w, l, unexpectedErrorWhenExecuteUpdateTask, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, httpcommon.Message(ok))
}

// @Summary      Получить список задач
// @Description  Получить список задач команды с фильтрацией и пагинацией. Доступно любому участнику команды
// @Router       /api/v1/tasks [get]
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id      query     string  true   "Team ID" format(uuid)
// @Param        status       query     string  false  "Статус задачи" Enums(todo,in_progress,done)
// @Param        assignee_id  query     string  false  "ID исполнителя" format(uuid)
// @Param        page         query     int     false  "Номер страницы" example(1)
// @Param        page_size    query     int     false  "Размер страницы" example(10)
// @Success      200  {object}  domain.FindTasksResponse
// @Failure      400  {object}  domain.ValidatationErrors
// @Failure      401  {object}  httpcommon.MessageResponse
// @Failure      403  {object}  httpcommon.MessageResponse
// @Failure      404  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (t *tasks) get(w http.ResponseWriter, r *http.Request) {

	const op = "controller.tasks.get"
	l := t.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)
	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	in := &inputs.FindManyTasks{
		TeamID:     query.Get("team_id"),
		Status:     domain.Status(query.Get("status")),
		AssigneeID: query.Get("assignee_id"),
	}

	if err := in.Validate(); err != nil {
		handleValidationErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	res, err := t.getTaskUc.Execute(ctx, local.UserID, limit, offset, in)
	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenGetTasks, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, res)
}

// @Summary      Получить историю изменений задачи
// @Description  Получить историю изменений задачи. Доступно любому участнику команды
// @Router       /api/v1/tasks/{id}/history [get]
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Task ID" format(uuid)
// @Success      200  {array}   domain.History
// @Failure      400  {object}  domain.ValidatationErrors
// @Failure      401  {object}  httpcommon.MessageResponse
// @Failure      403  {object}  httpcommon.MessageResponse
// @Failure      404  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (t *tasks) history(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.http.tasks.v1.history"

	l := t.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)
	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	taskID := r.PathValue("id")

	uuidInput := inputs.UUID{ID: taskID}
	if err := uuidInput.Validate("id"); err != nil {
		handleValidationErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	history, err := t.historyUc.Execute(ctx, local.UserID, uuidInput.ID)
	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenGetTaskHistory, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, history)
}

// @Summary      Получить комментарии задачи
// @Description  Получить список комментариев задачи. Доступно любому участнику команды
// @Router       /api/v1/tasks/{id}/comments [get]
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Task ID" format(uuid)
// @Success      200  {array}   domain.Comment
// @Failure      400  {object}  domain.ValidatationErrors
// @Failure      401  {object}  httpcommon.MessageResponse
// @Failure      403  {object}  httpcommon.MessageResponse
// @Failure      404  {object}  httpcommon.MessageResponse
// @Failure      500  {object}  httpcommon.MessageResponse
func (t *tasks) comments(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.http.tasks.v1.comments"
	l := t.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)
	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	taskID := r.PathValue("id")

	uuidInput := inputs.UUID{ID: taskID}
	if err := uuidInput.Validate("id"); err != nil {
		handleValidationErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	comments, err := t.commentsUc.Execute(ctx, local.UserID, uuidInput.ID)
	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenGetTaskComments, err)
		return
	}

	httpcommon.WriteJSON(w, http.StatusOK, comments)
}

// @Summary      Создать комментарий к задаче
// @Description  Создать комментарий к задаче. Пользователь должен быть членом команды задачи
// @Router       /api/v1/tasks/{id}/comments [post]
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string                  true "Task ID" format(uuid)
// @Param        dto body      inputs.CreateComment    true "Create comment dto"
// @Success      201 {object}  domain.IDResponse
// @Failure      400 {object}  domain.ValidatationErrors
// @Failure      401 {object}  httpcommon.MessageResponse
// @Failure      403 {object}  httpcommon.MessageResponse
// @Failure      404 {object}  httpcommon.MessageResponse
// @Failure      500 {object}  httpcommon.MessageResponse
func (t *tasks) createComment(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.http.tasks.v1.createComment"

	l := t.log.With(slog.String("op", op))

	local, err := httpcommon.GetLocalSession(r)

	if err != nil {
		httpcommon.WriteJSON(w, http.StatusUnauthorized, httpcommon.Message(unauthorized))
		return
	}

	uuidInput := inputs.UUID{ID: r.PathValue("id")}

	if err := uuidInput.Validate("id"); err != nil {
		handleValidationErrors(w, err)
		return
	}

	var in inputs.CreateComment
	if err := httpcommon.DecodeJSON(r.Body, &in); err != nil {
		httpcommon.WriteJSON(w, http.StatusBadRequest, httpcommon.Message(invalidBodyFormat))
		return
	}

	if err := in.Validate(); err != nil {
		handleValidationErrors(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	commentID, err := t.createCommentUc.Execute(ctx, local.UserID, uuidInput.ID, in.Text)
	if err != nil {
		handleErrors(w, l, unexpectedErrorWhenCreateComment, err)
		return
	}

	res := domain.IDResponse{ID: commentID}
	httpcommon.WriteJSON(w, http.StatusCreated, res)
}

func (t *tasks) RegisterRoutes(r chi.Router) {

	r.Route("/tasks", func(taskRouter chi.Router) {

		taskRouter.Post("/", t.create)
		taskRouter.Get("/", t.get)

		taskRouter.Route("/{id}", func(oneTaskRouter chi.Router) {

			oneTaskRouter.Put("/", t.update)
			oneTaskRouter.Get("/history", t.history)
			oneTaskRouter.Get("/comments", t.comments)
			oneTaskRouter.Post("/comments", t.createComment)
		})
	})

}
