package main

import (
	"crypto-trading-bot/internal/generator"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	packageName := flag.String("package", "crypto-trading-bot", "Package name")
	modelName := flag.String("model", "", "Model name")
	fields := flag.String("fields", "", "Comma-separated list of fields in format 'FieldName:FieldType:JSONTag:SQLType'")

	flag.Parse()

	if *modelName == "" {
		fmt.Println("Model name is required")
		os.Exit(1)
	}

	if *fields == "" {
		fmt.Println("Fields are required")
		os.Exit(1)
	}

	fieldList := strings.Split(*fields, ",")
	var modelFields []generator.Field
	for _, field := range fieldList {
		parts := strings.Split(field, ":")
		if len(parts) != 4 {
			fmt.Printf("Invalid field format: %s\n", field)
			os.Exit(1)
		}
		modelFields = append(modelFields, generator.Field{
			FieldName: parts[0],
			FieldType: parts[1],
			JSONTag:   parts[2],
			SQLType:   parts[3],
		})
	}

	model := generator.Model{
		PackageName:    *packageName,
		ModelName:      *modelName,
		ModelNameLower: strings.ToLower(*modelName),
		Fields:         modelFields,
	}

	if err := generator.GenerateModel(model); err != nil {
		fmt.Printf("Failed to generate model: %v\n", err)
		os.Exit(1)
	}

	if err := generator.GenerateRepository(model); err != nil {
		fmt.Printf("Failed to generate repository: %v\n", err)
		os.Exit(1)
	}

	if err := generator.GenerateService(model); err != nil {
		fmt.Printf("Failed to generate service: %v\n", err)
		os.Exit(1)
	}

	if err := generator.GenerateCreateTableScript(model); err != nil {
		fmt.Printf("Failed to generate create table script: %v\n", err)
		os.Exit(1)
	}

	if err := generator.GenerateDropTableScript(model); err != nil {
		fmt.Printf("Failed to generate drop table script: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Code generation completed successfully")
}
