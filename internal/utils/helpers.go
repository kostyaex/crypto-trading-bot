package utils

import (
	"fmt"
	"path/filepath"
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
