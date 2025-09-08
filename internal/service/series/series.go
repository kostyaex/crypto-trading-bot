package series

import (
	"crypto-trading-bot/internal/types"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func SaveSeries(series []types.Series, filePath string) {
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
