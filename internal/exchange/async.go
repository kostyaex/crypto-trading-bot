package exchange

import (
	"sync"
	"time"
)

type AsyncManager struct {
	mu       sync.RWMutex
	results  map[CommandID]interface{}
	timeouts map[CommandID]time.Time
	ttl      time.Duration
}

func NewAsyncManager(ttl time.Duration) *AsyncManager {
	return &AsyncManager{
		results:  make(map[CommandID]interface{}),
		timeouts: make(map[CommandID]time.Time),
		ttl:      ttl,
	}
}

func (am *AsyncManager) StoreResult(cmdID CommandID, result interface{}) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.results[cmdID] = result
	am.timeouts[cmdID] = time.Now().Add(am.ttl)
}

func (am *AsyncManager) GetResult(cmdID CommandID) (interface{}, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if expiry, exists := am.timeouts[cmdID]; exists && time.Now().Before(expiry) {
		if result, ok := am.results[cmdID]; ok {
			return result, true
		}
	}
	return nil, false
}

func (am *AsyncManager) Cleanup() {
	am.mu.Lock()
	defer am.mu.Unlock()
	now := time.Now()
	for id, expiry := range am.timeouts {
		if now.After(expiry) {
			delete(am.results, id)
			delete(am.timeouts, id)
		}
	}
}
