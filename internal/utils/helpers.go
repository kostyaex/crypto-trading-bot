package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const TimeFormatForHuman = "2006-01-02 15:04:05"

// Чтобы обрезать путь до указанного каталога
func TrimPathToDir(fullPath string, targetDir string) (string, error) {
	// // Нормализация пути
	// fullPath, err := filepath.Abs(fullPath)
	// if err != nil {
	// 	return "", err
	// }

	// // Нормализация пути к целевому каталогу
	// targetDir, err = filepath.Abs(targetDir)
	// if err != nil {
	// 	return "", err
	// }

	if filepath.Base(fullPath) == targetDir {
		return fullPath, nil
	}

	// Обрезка пути до целевого каталога
	for {
		parentDir := filepath.Dir(fullPath)
		if parentDir == "/" {
			// Если мы достигли корня, то не можем обрезать дальше
			break
		}

		if filepath.Base(parentDir) == targetDir {
			return parentDir, nil
		}

		fullPath = parentDir
	}

	return "", fmt.Errorf("не удалось обрезать путь до каталога %s", targetDir)
}

func FileList(filePath string, namePrefix string) ([]string, error) {
	names := make([]string, 0)

	files, err := os.ReadDir(filePath)

	if err != nil {
		return names, err
	}

	// Сортируем по дате в обратном порядке
	sort.Slice(files, func(i int, j int) bool {
		fileI, _ := files[i].Info()
		fileJ, _ := files[j].Info()
		return fileI.ModTime().After(fileJ.ModTime())
	})

	for _, file := range files {
		name := fmt.Sprintf("%s%s", namePrefix, file.Name())
		fmt.Println(name)
		names = append(names, name)
	}

	return names, nil
}

// -------------------------------------------------------------
// Для работы с временем

// Возвращает начало временного интервала для заданного времени и типа интервала
func GetIntervalBounds(t time.Time, interval string) (start, end, nextStart time.Time, err error) {

	// Рассчитываем начало интервала
	if interval == "1M" {
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		nextStart = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
		end = nextStart.Add(-time.Nanosecond) // Конец интервала
		return
	} else if interval == "1d" {
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
		nextDay := t.Add(24 * time.Hour)
		nextStart = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, t.Location())
		return
	} else if interval == "1w" {
		start, end, nextStart = getWeekBounds(t)
		return
	}

	// Парсинг интервала
	duration, err := parseInterval(interval)
	if err != nil {
		return
	}
	start = t.Truncate(duration)
	end = start.Add(duration).Add(-time.Nanosecond) // Конец интервала
	nextStart = start.Add(duration)

	return
}

// определить границы недели
func getWeekBounds(t time.Time) (start, end, nextStart time.Time) {
	year, week := t.ISOWeek()
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
	daysSinceEpoch := int(firstDay.Weekday()) - 1
	if daysSinceEpoch < 0 {
		daysSinceEpoch += 7
	}
	start = firstDay.AddDate(0, 0, (week-1)*7-daysSinceEpoch)
	end = start.Add(7 * 24 * time.Hour).Add(-time.Nanosecond)
	nextStart = start.Add(7 * 24 * time.Hour)
	return
}

// Парсинг строки интервала в time.Duration
func parseInterval(interval string) (time.Duration, error) {
	// Словарь интервалов
	intervals := map[string]time.Duration{
		"1s":  time.Second,
		"1m":  time.Minute,
		"3m":  3 * time.Minute,
		"5m":  5 * time.Minute,
		"15m": 15 * time.Minute,
		"30m": 30 * time.Minute,
		"1h":  time.Hour,
		"4h":  4 * time.Hour,
		"6h":  6 * time.Hour,
		"8h":  8 * time.Hour,
		"12h": 12 * time.Hour,
		"1d":  24 * time.Hour,
		"3d":  72 * time.Hour,
		"1w":  168 * time.Hour,
		//"1M":  30 * 24 * time.Hour, // Примерное значение для месяца
	}

	// Проверка допустимых значений
	if duration, ok := intervals[interval]; ok {
		return duration, nil
	}
	return 0, fmt.Errorf("недопустимый интервал: %s", interval)
}

// -------------------------------------------------------------

// Функция для разбиения данных из канала на блоки с перекрытием
func SplitChannelWithOverlap[T any](inputChan <-chan T, blockSize int, overlap int, outputChan chan<- []T) {
	var buffer []T // Буфер для хранения текущего блока
	var block []T  // Текущий блок для отправки

	for {
		// Читаем один элемент из входного канала
		element, ok := <-inputChan
		if !ok {
			// Если канал закрыт, проверяем, есть ли оставшиеся данные в буфере
			if len(buffer) >= overlap {
				// Отправляем оставшийся блок (всю оставшуюся часть буфера)
				//outputChan <- buffer[len(buffer)-min(len(buffer), blockSize):]

				// Если канал закрыт, проверяем, есть ли оставшиеся данные в буфере
				// Но мы не отправляем остатки, если они не составляют полный блок
				close(outputChan) // Закрываем выходной канал
				return
			}
			close(outputChan) // Закрываем выходной канал
			return
		}

		// Добавляем элемент в буфер
		buffer = append(buffer, element)

		// Проверяем, достаточно ли элементов для формирования нового блока
		if len(buffer) >= blockSize {
			// Собираем текущий блок
			block = make([]T, blockSize)
			copy(block, buffer[len(buffer)-blockSize:])

			// Отправляем блок в выходной канал
			outputChan <- block

			// Обновляем буфер, оставляя только последние `overlap` элементов
			buffer = buffer[len(buffer)-overlap:]
		}
	}
}
