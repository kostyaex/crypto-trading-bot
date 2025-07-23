package handlers

import (
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/web/ui"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// для обращения к API postgre (postgrest)
type ResourcesHandler struct {
	logger    *logger.Logger
	resources map[string]ui.Resource
}

func NewResourcesHandler(logger *logger.Logger) *ResourcesHandler {
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

	resources["market_data"] = ui.Resource{
		Name:        "market_data",
		Title:       "Данные с бирж",
		FieldsOrder: []string{"id", "exchange", "symbol", "time_frame", "active", "actual_time", "status"},
		Fields: map[string]*ui.ResourceField{
			"timestamp":   &ui.ResourceField{Name: "timestamp", Title: "timestamp"},
			"exchange":    &ui.ResourceField{Name: "exchange", Title: "exchange"},
			"symbol":      &ui.ResourceField{Name: "symbol", Title: "symbol"},
			"time_frame":  &ui.ResourceField{Name: "time_frame", Title: "time_frame"},
			"open_price":  &ui.ResourceField{Name: "open_price", Title: "open_price"},
			"hight_price": &ui.ResourceField{Name: "hight_price", Title: "hight_price"},
			"low_price":   &ui.ResourceField{Name: "low_price", Title: "low_price"},
			"close_price": &ui.ResourceField{Name: "close_price", Title: "close_price"},
			"volume":      &ui.ResourceField{Name: "volume", Title: "volume"},
			"buy_volume":  &ui.ResourceField{Name: "buy_volume", Title: "buy_volume"},
			"sell_volume": &ui.ResourceField{Name: "sell_volume", Title: "sell_volume"},
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
