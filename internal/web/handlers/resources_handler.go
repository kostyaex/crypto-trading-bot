package handlers

import (
	"crypto-trading-bot/internal/utils"
	"crypto-trading-bot/internal/web/ui"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// для обращения к API postgre (postgrest)
type ResourcesHandler struct {
	logger    *utils.Logger
	resources map[string]ui.Resource
}

func NewResourcesHandler(logger *utils.Logger) *ResourcesHandler {
	resources := make(map[string]ui.Resource)

	resources["market_data_statuss"] = ui.Resource{
		Name:        "market_data_statuss",
		Title:       "Статус загрузки данных бирж",
		FieldsOrder: []string{"id", "exchange", "symbol", "time_frame", "active", "actual_time", "status"},
		Fields: map[string]*ui.ResourceField{
			"id":          &ui.ResourceField{Name: "id", Title: "ID"},
			"exchange":    &ui.ResourceField{Name: "exchange", Title: "Биржа"},
			"symbol":      &ui.ResourceField{Name: "symbol", Title: "Пара"},
			"time_frame":  &ui.ResourceField{Name: "time_frame", Title: "Интервал"},
			"active":      &ui.ResourceField{Name: "active", Title: "Активность"},
			"actual_time": &ui.ResourceField{Name: "actual_time", Title: "Время актуальности"},
			"status":      &ui.ResourceField{Name: "status", Title: "Статус"},
		},
	}
	return &ResourcesHandler{
		logger:    logger,
		resources: resources,
	}
}

func (h *ResourcesHandler) GetResourcesListPage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	resourceParam := vars["resource"]

	var resource ui.Resource
	var ok bool
	if resource, ok = h.resources[resourceParam]; !ok {
		http.Error(w, fmt.Sprintf("Не заданы метаданные ресурса: %s", resourceParam), http.StatusInternalServerError)
		return
	}

	if err := ui.ResourceListPage(resource).Render(r.Context(), w); err != nil {
		h.logger.Errorf("Ошибка формирования страницы бектестинга: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
