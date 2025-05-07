package handlers

import (
	"crypto-trading-bot/internal/utils"
	"crypto-trading-bot/internal/web/ui"
	"encoding/json"
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

	// Заполняем метаданные для различных таблиц

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

	resources["strategies"] = ui.Resource{
		Name:        "strategies",
		Title:       "Стратегии",
		FieldsOrder: []string{"id", "name", "description", "active", "config"},
		Fields: map[string]*ui.ResourceField{
			"id":          &ui.ResourceField{Name: "id", Title: "ID"},
			"name":        &ui.ResourceField{Name: "name", Title: "Наименование"},
			"description": &ui.ResourceField{Name: "description", Title: "Описание"},
			"active":      &ui.ResourceField{Name: "active", Title: "Активность"},
			"config":      &ui.ResourceField{Name: "config", Title: "Конфигурация", Component: "strategysettings"},
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

	// Формируем пустой JSON объект для нового элемента
	newItem := make(map[string]string)
	for _, field := range resource.Fields {
		if field.Name == "id" {
			continue
		}
		newItem[field.Name] = ""
	}
	newItemJSON, err := json.Marshal(newItem)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка формирования JSON: %s", err), http.StatusInternalServerError)
		return
	}

	if err := ui.ResourceListPage(resource, string(newItemJSON[:])).Render(r.Context(), w); err != nil {
		h.logger.Errorf("Ошибка формирования страницы бектестинга: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
