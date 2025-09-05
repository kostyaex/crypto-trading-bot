package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// Field представляет поле модели.
type Field struct {
	FieldName string
	FieldType string
	JSONTag   string
	SQLType   string
}

// Model представляет модель для генерации.
type Model struct {
	PackageName    string
	ModelName      string
	ModelNameLower string
	Fields         []Field
}

// GenerateModel генерирует код модели.
func GenerateModel(model Model) error {
	tmpl, err := template.ParseFiles("templates/model.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse model template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, model); err != nil {
		return fmt.Errorf("failed to execute model template: %v", err)
	}

	modelFilePath := filepath.Join("internal", "models", fmt.Sprintf("%s.go", model.ModelNameLower))
	if err := os.WriteFile(modelFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write model file: %v", err)
	}

	return nil
}

// GenerateRepository генерирует код репозитория.
func GenerateRepository(model Model) error {
	tmpl, err := template.New("repository.tmpl").Funcs(template.FuncMap{
		"ToSnakeCase": toSnakeCase,
	}).ParseFiles("templates/repository.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse repository template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, model); err != nil {
		return fmt.Errorf("failed to execute repository template: %v", err)
	}

	repoFilePath := filepath.Join("internal", "repositories", fmt.Sprintf("%s_repo.go", model.ModelNameLower))
	if err := os.WriteFile(repoFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write repository file: %v", err)
	}

	return nil
}

// GenerateService генерирует код сервиса.
func GenerateService(model Model) error {
	tmpl, err := template.ParseFiles("templates/service.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse service template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, model); err != nil {
		return fmt.Errorf("failed to execute service template: %v", err)
	}

	serviceFilePath := filepath.Join("internal", "services", fmt.Sprintf("%s_service.go", model.ModelNameLower))
	if err := os.WriteFile(serviceFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %v", err)
	}

	return nil
}

// GenerateCreateTableScript генерирует скрипт создания таблицы.
func GenerateCreateTableScript(model Model) error {
	tmpl, err := template.New("create_table.tmpl").Funcs(template.FuncMap{
		"ToSnakeCase": toSnakeCase,
	}).ParseFiles("templates/create_table.tmpl")
	//tmpl, err := template.ParseFiles("templates/create_table.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse create table template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, model); err != nil {
		return fmt.Errorf("failed to execute create table template: %v", err)
	}

	createFilePath := filepath.Join("migrations", fmt.Sprintf("0000001_create_%s_table.sql", model.ModelNameLower))
	if err := os.WriteFile(createFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write create table script: %v", err)
	}

	return nil
}

// GenerateDropTableScript генерирует скрипт удаления таблицы.
func GenerateDropTableScript(model Model) error {
	tmpl, err := template.New("drop_table.tmpl").Funcs(template.FuncMap{
		"ToSnakeCase": toSnakeCase,
	}).ParseFiles("templates/drop_table.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse drop table template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, model); err != nil {
		return fmt.Errorf("failed to execute drop table template: %v", err)
	}

	dropFilePath := filepath.Join("migrations", fmt.Sprintf("0000001_drop_%s_table.sql", model.ModelNameLower))
	if err := os.WriteFile(dropFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write drop table script: %v", err)
	}

	return nil
}

// toSnakeCase преобразует имя поля в snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (unicode.IsUpper(r) || unicode.IsDigit(r)) && (!unicode.IsUpper(rune(s[i-1])) || !unicode.IsDigit(rune(s[i-1]))) {
			result.WriteString("_")
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
