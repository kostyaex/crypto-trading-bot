package utils

import (
	"fmt"
	"path/filepath"
	"time"
)

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
