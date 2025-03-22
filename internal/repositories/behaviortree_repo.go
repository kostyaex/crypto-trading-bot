package repositories

import (
	"crypto-trading-bot/internal/utils"
	"database/sql"
	"encoding/json"
)

type BehaviorTreeRepository interface {
	GetBehaviorTreeState(strategyID int) (map[string]interface{}, error)
	SaveBehaviorTreeState(strategyID int, state map[string]interface{}) error
}

type behaviorTreeRepository struct {
	db     *DB
	logger *utils.Logger
}

func NewBehaviorTreeRepositoryRepository(db *DB, logger *utils.Logger) BehaviorTreeRepository {
	return &behaviorTreeRepository{db: db, logger: logger}
}

func (r *behaviorTreeRepository) GetBehaviorTreeState(strategyID int) (map[string]interface{}, error) {
	query := `
        SELECT state
        FROM behavior_trees
        WHERE strategy_id = $1
        ORDER BY last_executed DESC
        LIMIT 1;
    `

	var stateJSON json.RawMessage
	err := r.db.QueryRow(query, strategyID).Scan(&stateJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warnf("No behavior tree state found for strategy ID %d", strategyID)
			return nil, nil
		}
		r.logger.Errorf("Failed to get behavior tree state for strategy ID %d: %v", strategyID, err)
		return nil, err
	}

	var state map[string]interface{}
	if err := json.Unmarshal(stateJSON, &state); err != nil {
		r.logger.Errorf("Failed to unmarshal behavior tree state for strategy ID %d: %v", strategyID, err)
		return nil, err
	}

	r.logger.Infof("Behavior tree state retrieved for strategy ID %d: %v", strategyID, state)
	return state, nil
}

func (r *behaviorTreeRepository) SaveBehaviorTreeState(strategyID int, state map[string]interface{}) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		r.logger.Errorf("Failed to marshal behavior tree state for strategy ID %d: %v", strategyID, err)
		return err
	}

	query := `
        INSERT INTO behavior_trees (strategy_id, state, last_executed)
        VALUES ($1, $2, NOW());
    `

	_, err = r.db.Exec(query, strategyID, stateJSON)
	if err != nil {
		r.logger.Errorf("Failed to save behavior tree state for strategy ID %d: %v", strategyID, err)
		return err
	}

	r.logger.Infof("Behavior tree state saved for strategy ID %d: %v", strategyID, state)
	return nil
}
