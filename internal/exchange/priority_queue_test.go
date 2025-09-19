// priority_queue_test.go
package exchange

import (
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"
)

// Пример пользовательского типа данных
type TestData struct {
	ID   int
	Name string
}

func TestPushBatchAndPopOne_Ordering(t *testing.T) {
	pqm := NewPriorityQueueManager[TestData]()

	now := time.Now()
	records := []*Record[TestData]{
		{Timestamp: now.Add(5 * time.Second), Data: TestData{ID: 3, Name: "поздний"}},
		{Timestamp: now.Add(1 * time.Second), Data: TestData{ID: 1, Name: "ранний"}},
		{Timestamp: now.Add(3 * time.Second), Data: TestData{ID: 2, Name: "средний"}},
	}

	pqm.PushBatch(records...)

	// Должны извлекаться в порядке возрастания Timestamp
	expectedOrder := []int{1, 2, 3}
	for i, expectedID := range expectedOrder {
		rec, ok := pqm.PopOne()
		if !ok {
			t.Fatalf("ожидалась запись %d, но очередь пуста на шаге %d", expectedID, i+1)
		}
		if rec.Data.ID != expectedID {
			t.Errorf("ожидался ID=%d, получили ID=%d", expectedID, rec.Data.ID)
		}
	}

	// Очередь должна быть пуста
	if rec, ok := pqm.PopOne(); ok {
		t.Errorf("ожидалась пустая очередь, но получена запись: %+v", rec)
	}
}

func TestPopOne_EmptyQueue(t *testing.T) {
	pqm := NewPriorityQueueManager[TestData]()

	rec, ok := pqm.PopOne()
	if ok {
		t.Errorf("ожидалось отсутствие записи, но получена: %+v", rec)
	}
	if rec != nil {
		t.Errorf("ожидался nil, но получен: %+v", rec)
	}
}

func TestLen(t *testing.T) {
	pqm := NewPriorityQueueManager[TestData]()

	if pqm.Len() != 0 {
		t.Errorf("ожидалась длина 0, получено: %d", pqm.Len())
	}

	now := time.Now()
	records := []*Record[TestData]{
		{Timestamp: now, Data: TestData{ID: 1}},
		{Timestamp: now.Add(1 * time.Second), Data: TestData{ID: 2}},
	}

	pqm.PushBatch(records...)

	if pqm.Len() != 2 {
		t.Errorf("ожидалась длина 2, получено: %d", pqm.Len())
	}

	pqm.PopOne() // извлекаем одну

	if pqm.Len() != 1 {
		t.Errorf("ожидалась длина 1, получено: %d", pqm.Len())
	}
}

func TestConcurrentAccess(t *testing.T) {
	pqm := NewPriorityQueueManager[int]()

	var wg sync.WaitGroup
	numWorkers := 10
	recordsPerWorker := 100

	// Горутины-писатели
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			records := make([]*Record[int], recordsPerWorker)
			baseTime := time.Now().Add(time.Duration(workerID) * time.Millisecond)
			for i := 0; i < recordsPerWorker; i++ {
				records[i] = &Record[int]{
					Timestamp: baseTime.Add(time.Duration(i) * time.Millisecond),
					Data:      workerID*1000 + i,
				}
			}
			pqm.PushBatch(records...)
		}(w)
	}

	time.Sleep(3 * time.Second)

	// Горутина-читатель
	readResults := make(chan int, numWorkers*recordsPerWorker)
	go func() {
		defer close(readResults)
		for {
			rec, ok := pqm.PopOne()
			if !ok {
				// Проверяем, может, просто ещё не все записи добавлены?
				time.Sleep(10 * time.Millisecond)
				continue
			}
			readResults <- rec.Data
			if len(readResults) == numWorkers*recordsPerWorker {
				break
			}
		}
	}()

	time.Sleep(3 * time.Second)

	wg.Wait()

	// Собираем все результаты
	var results []int
	for r := range readResults {
		results = append(results, r)
	}

	// Проверяем, что прочитали всё
	if len(results) != numWorkers*recordsPerWorker {
		t.Errorf("ожидалось %d записей, получено %d", numWorkers*recordsPerWorker, len(results))
	}

	// Проверяем, что порядок соблюден (сортируем по времени вручную для сравнения)
	// Создадим ожидаемый порядок: все записи должны быть отсортированы по времени
	expected := make([]int, 0, len(results))
	for w := 0; w < numWorkers; w++ {
		//baseTime := time.Now().Add(time.Duration(w) * time.Millisecond)
		for i := 0; i < recordsPerWorker; i++ {
			expected = append(expected, w*1000+i)
		}
	}

	// Но так как время у всех пересекается — отсортируем ожидаемый массив по "виртуальному времени"
	sort.Slice(expected, func(i, j int) bool {
		wi, ii := expected[i]/1000, expected[i]%1000
		wj, ij := expected[j]/1000, expected[j]%1000
		ti := time.Now().Add(time.Duration(wi) * time.Millisecond).Add(time.Duration(ii) * time.Millisecond)
		tj := time.Now().Add(time.Duration(wj) * time.Millisecond).Add(time.Duration(ij) * time.Millisecond)
		return ti.Before(tj)
	})

	if len(results) != len(expected) {
		t.Fatalf("разная длина: results=%d, expected=%d", len(results), len(expected))
	}

	for i := 0; i < len(results); i++ {
		if results[i] != expected[i] {
			t.Errorf("на позиции %d ожидалось %d, получено %d", i, expected[i], results[i])
		}
	}
}

func TestPushBatch_EmptyBatch(t *testing.T) {
	pqm := NewPriorityQueueManager[string]()

	// Не должно паниковать
	pqm.PushBatch(nil)

	pqm.PushBatch([]*Record[string]{}...)

	if pqm.Len() != 0 {
		t.Errorf("ожидалась пустая очередь после пустых батчей")
	}

	// Добавляем реальные данные
	pqm.PushBatch([]*Record[string]{
		{Timestamp: time.Now(), Data: "test"},
	}...)

	if pqm.Len() != 1 {
		t.Errorf("ожидалась длина 1, получено %d", pqm.Len())
	}
}

func BenchmarkPushBatchAndPopOne(b *testing.B) {
	pqm := NewPriorityQueueManager[int]()
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Пушим 100 записей
		records := make([]*Record[int], 100)
		for j := 0; j < 100; j++ {
			records[j] = &Record[int]{
				Timestamp: now.Add(time.Duration(j) * time.Millisecond),
				Data:      j,
			}
		}
		pqm.PushBatch(records...)

		// Извлекаем 100 записей
		for j := 0; j < 100; j++ {
			_, ok := pqm.PopOne()
			if !ok {
				b.Fatal("не удалось извлечь запись")
			}
		}
	}
}

func ExamplePriorityQueueManager() {
	pqm := NewPriorityQueueManager[string]()

	// Добавляем записи
	now := time.Now()
	pqm.PushBatch([]*Record[string]{
		{Timestamp: now.Add(10 * time.Second), Data: "третий"},
		{Timestamp: now, Data: "первый"},
		{Timestamp: now.Add(5 * time.Second), Data: "второй"},
	}...)

	// Извлекаем по одной — в порядке возрастания времени
	for i := 0; i < 3; i++ {
		if rec, ok := pqm.PopOne(); ok {
			fmt.Printf("Извлечено: %s\n", rec.Data)
		}
	}
	// Output:
	// Извлечено: первый
	// Извлечено: второй
	// Извлечено: третий
}
