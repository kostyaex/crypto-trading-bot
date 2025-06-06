package handlers

import (
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/trading/strategy"
)

type StrategyHandler struct {
	strategyService strategy.StrategyService
	logger          *logger.Logger
}

func NewStrategyHandler(strategyService strategy.StrategyService, logger *logger.Logger) *StrategyHandler {
	return &StrategyHandler{strategyService: strategyService, logger: logger}
}

// func (h *StrategyHandler) GetStrategiesListPage(w http.ResponseWriter, r *http.Request) {
// 	strategies, err := h.strategyService.GetActiveStrategies()

// 	if err != nil {
// 		h.logger.Errorf("Failed to GetActiveStrategies: %v", err)
// 		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
// 	}

// 	if err = ui.StrategyListComponent(strategies).Render(r.Context(), w); err != nil {
// 		h.logger.Errorf("Failed to render index page: %v", err)
// 		http.Error(w, "Failed to render page", http.StatusInternalServerError)
// 	}
// }

// func (h *StrategyHandler) GetNewStrategyPage(w http.ResponseWriter, r *http.Request) {

// 	if err := ui.StrategyNewPage().Render(r.Context(), w); err != nil {
// 		h.logger.Errorf("Failed to render new strategy page: %v", err)
// 		http.Error(w, "Failed to render page", http.StatusInternalServerError)
// 	}
// }

// func (h *StrategyHandler) GetEditStrategyPage(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	idStr := vars["id"]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		h.logger.Errorf("Invalid strategy ID: %v", err)
// 		http.Error(w, "Invalid strategy ID", http.StatusBadRequest)
// 		return
// 	}

// 	strat, err := h.strategyService.GetStrategyByID(id)
// 	if err != nil {
// 		h.logger.Errorf("Failed to get strategy by ID %d: %v", id, err)
// 		http.Error(w, "Failed to get strategy", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := ui.StrategyEditPage(strat).Render(r.Context(), w); err != nil {
// 		h.logger.Errorf("Failed to render edit strategy page: %v", err)
// 		http.Error(w, "Failed to render page", http.StatusInternalServerError)
// 	}
// }

// func (h *StrategyHandler) PostCreateStrategy(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		h.logger.Errorf("Failed to parse form: %v", err)
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	name := r.Form.Get("name")
// 	description := r.Form.Get("description")
// 	configJSON := r.Form.Get("config")

// 	var config map[string]interface{}
// 	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
// 		h.logger.Errorf("Failed to unmarshal config: %v", err)
// 		http.Error(w, "Invalid config", http.StatusBadRequest)
// 		return
// 	}

// 	strat, err := models.NewStrategy(name, description, config)
// 	if err != nil {
// 		h.logger.Errorf("Failed to create strategy: %v", err)
// 		http.Error(w, "Failed to create strategy", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := h.strategyService.SaveStrategy(strat); err != nil {
// 		h.logger.Errorf("Failed to save strategy: %v", err)
// 		http.Error(w, "Failed to save strategy", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }

// func (h *StrategyHandler) PostUpdateStrategy(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	idStr := vars["id"]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		h.logger.Errorf("Invalid strategy ID: %v", err)
// 		http.Error(w, "Invalid strategy ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		h.logger.Errorf("Failed to parse form: %v", err)
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	name := r.Form.Get("name")
// 	description := r.Form.Get("description")
// 	configJSON := r.Form.Get("config")
// 	active := r.Form.Get("active") == "on"

// 	var config map[string]interface{}
// 	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
// 		h.logger.Errorf("Failed to unmarshal config: %v", err)
// 		http.Error(w, "Invalid config", http.StatusBadRequest)
// 		return
// 	}

// 	strat, err := h.strategyService.GetStrategyByID(id)
// 	if err != nil {
// 		h.logger.Errorf("Failed to get strategy by ID %d: %v", id, err)
// 		http.Error(w, "Failed to get strategy", http.StatusInternalServerError)
// 		return
// 	}

// 	strat.Name = name
// 	strat.Description = description
// 	strat.Active = active

// 	if err := strat.MarshalConfig(config); err != nil {
// 		h.logger.Errorf("Failed to marshal strategy config: %v", err)
// 		http.Error(w, "Failed to update strategy", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := h.strategyService.UpdateStrategy(strat); err != nil {
// 		h.logger.Errorf("Failed to update strategy: %v", err)
// 		http.Error(w, "Failed to update strategy", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }

// func (h *StrategyHandler) PostDeleteStrategy(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	idStr := vars["id"]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		h.logger.Errorf("Invalid strategy ID: %v", err)
// 		http.Error(w, "Invalid strategy ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := h.strategyService.DeleteStrategy(id); err != nil {
// 		h.logger.Errorf("Failed to delete strategy by ID %d: %v", id, err)
// 		http.Error(w, "Failed to delete strategy", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
