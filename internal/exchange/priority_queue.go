// priority_queue.go
package exchange

import (
	"container/heap"
	"sync"
	"time"
)

// Record — обобщённая структура записи с временной меткой
type Record[T any] struct {
	Timestamp time.Time
	Data      T
}

// RecordHeap — минимальная куча по Timestamp (самые ранние — первыми)
type RecordHeap[T any] []*Record[T]

func (h RecordHeap[T]) Len() int           { return len(h) }
func (h RecordHeap[T]) Less(i, j int) bool { return h[i].Timestamp.Before(h[j].Timestamp) }
func (h RecordHeap[T]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *RecordHeap[T]) Push(x interface{}) {
	*h = append(*h, x.(*Record[T]))
}

func (h *RecordHeap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// PriorityQueueManager — потокобезопасный приоритетный менеджер
type PriorityQueueManager[T any] struct {
	mu   sync.Mutex
	heap *RecordHeap[T]
}

// NewPriorityQueueManager создаёт новый менеджер
func NewPriorityQueueManager[T any]() *PriorityQueueManager[T] {
	h := &RecordHeap[T]{}
	heap.Init(h)
	return &PriorityQueueManager[T]{
		heap: h,
	}
}

// PushBatch добавляет пакет записей. Потокобезопасно.
func (pqm *PriorityQueueManager[T]) PushBatch(records ...*Record[T]) {
	if len(records) == 0 {
		return
	}

	pqm.mu.Lock()
	defer pqm.mu.Unlock()

	for _, rec := range records {
		heap.Push(pqm.heap, rec)
	}
}

// PopOne извлекает одну запись с наименьшим Timestamp.
// Возвращает запись и true, если она есть; nil и false — если очередь пуста.
// НЕ БЛОКИРУЕТСЯ.
func (pqm *PriorityQueueManager[T]) PopOne() (*Record[T], bool) {
	pqm.mu.Lock()
	defer pqm.mu.Unlock()

	if pqm.heap.Len() == 0 {
		return nil, false
	}

	item := heap.Pop(pqm.heap).(*Record[T])
	return item, true
}

// Len возвращает текущее количество элементов в очереди
func (pqm *PriorityQueueManager[T]) Len() int {
	pqm.mu.Lock()
	defer pqm.mu.Unlock()
	return pqm.heap.Len()
}
