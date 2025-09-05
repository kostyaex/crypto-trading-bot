package trader

import (
	"context"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/trading/dispatcher"
	"sync"
)

type Runner struct {
	Strategy *models.Strategy
	Pipeline *Pipeline
	cancel   context.CancelFunc // Для остановки
	running  bool
	doneCh   chan struct{}
	Err      error
	mu       sync.RWMutex
}

func NewRunner(Strategy *models.Strategy, Pipeline *Pipeline) *Runner {
	return &Runner{
		Strategy: Strategy,
		Pipeline: Pipeline,
		doneCh:   make(chan struct{}),
	}
}

func (r *Runner) Start(parentCtx context.Context) {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return
	}
	r.running = true
	r.mu.Unlock()

	// Создаем отдельный контекст, для независимого завершения своего pipeline
	ctx, cancel := context.WithCancel(parentCtx)
	r.cancel = cancel

	// Запускаем пайплайн в горутине
	go func() {
		defer r.onFinish()
		err := r.Pipeline.Run(ctx) // ← ctx.Done() остановит цикл
		if err != nil && err != context.Canceled {
			r.mu.Lock()
			r.Err = err
			r.mu.Unlock()
		}
	}()
}

func (r *Runner) Stop() {
	r.mu.RLock()
	cancel := r.cancel
	r.mu.RUnlock()

	if cancel != nil {
		cancel() // ← это остановит Pipeline
	}
}

func (r *Runner) onFinish() {
	r.mu.Lock()
	r.running = false
	r.cancel = nil
	r.mu.Unlock()

	r.doneCh <- struct{}{} // Отправляем сигнал о завершении
}

func (r *Runner) IsRunning() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}

// Канал получения сигнала о завершении
func (r *Runner) Done() <-chan struct{} {
	return r.doneCh
}

// Обновить правила
func (r *Runner) UpdateDispatcher(Dispatcher *dispatcher.Dispatcher) {

	// Обновляем пайплайн (без остановки!)
	if r.IsRunning() && r.Pipeline != nil {
		r.mu.Lock()
		r.Pipeline.UpdateDispatcher(Dispatcher)
		r.mu.Unlock()
	}
}
