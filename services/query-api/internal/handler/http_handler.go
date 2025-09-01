package handler

import (
	"context"
	"net/http"

	"github.com/DeadlyParkour777/logging-system/pkg/models"
	"github.com/DeadlyParkour777/logging-system/services/query-api/internal/util"
)

type QueryService interface {
	SearchLogs(ctx context.Context, params map[string]string) ([]models.LogEntry, error)
}

type HttpHandler struct {
	service QueryService
}

func NewHttpHandler(service QueryService) *HttpHandler {
	return &HttpHandler{service: service}
}

// SearchLogs godoc
// @Summary Поиск и фильтрация логов
// @Description Возвращает массив записей логов на основе query-параметров.
// @Tags Logs
// @Accept  json
// @Produce  json
// @Param service_name query string false "Фильтр по имени сервиса"
// @Param level query string false "Фильтр по уровню лога (INFO, ERROR)"
// @Param search query string false "Полнотекстовый поиск по сообщению"
// @Success 200 {array} models.LogEntry "Успешный ответ с массивом логов"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /logs [get]
func (h *HttpHandler) SearchLogs(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{
		"service_name": r.URL.Query().Get("service_name"),
		"level":        r.URL.Query().Get("level"),
		"search":       r.URL.Query().Get("search"),
	}

	logs, err := h.service.SearchLogs(r.Context(), params)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, logs)
}
