package utils

import (
	"testing"
)

// Вспомогательная функция для сравнения двух [][]int
func equal(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// Тестовая функция
func TestSplitChannelWithOverlap(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		blockSize int
		overlap   int
		expected  [][]int
	}{
		{
			name:      "Basic case with numbers",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			blockSize: 3,
			overlap:   2,
			expected: [][]int{
				{0, 1, 2},
				{1, 2, 3},
				{2, 3, 4},
				{3, 4, 5},
				{4, 5, 6},
				{5, 6, 7},
				{6, 7, 8},
				{7, 8, 9},
			},
		},
		{
			name:      "No overlap",
			input:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			blockSize: 3,
			overlap:   0,
			expected: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем входной канал и заполняем его данными
			inputChan := make(chan int)
			go func() {
				defer close(inputChan)
				for _, d := range tt.input {
					inputChan <- d
				}
			}()

			// Создаем выходной канал
			outputChan := make(chan []int)

			// Запускаем функцию для обработки канала
			go SplitChannelWithOverlap(inputChan, tt.blockSize, tt.overlap, outputChan)

			// Собираем результаты из выходного канала
			var result [][]int
			for block := range outputChan {
				result = append(result, block)
			}

			// Проверяем, что полученные блоки соответствуют ожидаемым
			if !equal(result, tt.expected) {
				t.Errorf("Test %s failed:\nGot: %v\nExpected: %v", tt.name, result, tt.expected)
			}
		})
	}
}
