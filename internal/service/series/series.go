package series

import (
	"crypto-trading-bot/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type Point struct {
	Value      float64            `json:"value"`
	Weight     float64            `json:"weight"`
	Time       time.Time          `json:"time"`
	MarketData *models.MarketData `json:"market_data"` // основание по которому формировалась точка
}

type Series struct {
	Points []Point `json:"points"`
}

// Возвращает первую точку
func (s *Series) First() *Point {
	if len(s.Points) == 0 {
		return nil
	}
	return &s.Points[0]
}

// Возвращает последнюю точку
func (s *Series) Last() *Point {
	if len(s.Points) == 0 {
		return nil
	}
	return &s.Points[len(s.Points)-1]
}

// Формирует строку для вывода в лог
func (s *Series) String() string {
	if len(s.Points) == 0 {
		return "[]"
	}

	if len(s.Points) == 1 {
		first := s.Points[0].MarketData
		return fmt.Sprintf("[%s $%.2f]", first.Timestamp.Format(time.DateTime), first.ClusterPrice)
	}

	first := s.Points[0].MarketData
	last := s.Points[len(s.Points)-1].MarketData

	return fmt.Sprintf("(%d)[%s $%.2f - %s $%.2f]",
		len(s.Points),
		first.Timestamp.Format(time.DateTime), first.ClusterPrice,
		last.Timestamp.Format(time.DateTime), last.ClusterPrice)
}

func SaveSeries(series []Series, filePath string) {
	// Кодируем массив структур в JSON
	jsonData, err := json.MarshalIndent(series, "", "	")
	if err != nil {
		fmt.Println("Ошибка кодирования в JSON:", err)
		return
	}

	// Записываем JSON в файл
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Printf("Ошибка записи в файл:%s\n", err)
		return
	}

	fmt.Printf("JSON записан в файл %s\n", filePath)
}

func SeriesDumpList(filePath string, namePrefix string) ([]string, error) {
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
