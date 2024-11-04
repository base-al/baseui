package utils

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
)

var PluralizeClient *pluralize.Client

func init() {
	PluralizeClient = pluralize.NewClient()
}

//go:embed templates/*
var TemplateFS embed.FS

type FieldStruct struct {
	Name           string
	Type           string
	JSONName       string
	DBName         string
	AssociatedType string
	PluralType     string
	Relationship   string
	IsSort         bool // Add this field
}

// Update GenerateFieldStructs to handle sort fields
func GenerateFieldStructs(fields []string) []FieldStruct {
	var fieldStructs []FieldStruct

	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) >= 2 {
			name := ToPascalCase(parts[0])
			fieldType := parts[1]
			jsonName := ToSnakeCase(parts[0])
			dbName := ToSnakeCase(parts[0])
			var associatedType, pluralType, relationship string
			isSort := false

			// Check if this is a sort field
			if fieldType == "sort" {
				isSort = true
				fieldType = "int"
			}

			dartType := GetDartType(fieldType)

			fieldStructs = append(fieldStructs, FieldStruct{
				Name:           name,
				Type:           dartType,
				JSONName:       jsonName,
				DBName:         dbName,
				AssociatedType: associatedType,
				PluralType:     pluralType,
				Relationship:   relationship,
				IsSort:         isSort,
			})
		}
	}

	return fieldStructs
}

// Add helper function to check for sort fields
func HasSortField(fields []FieldStruct) bool {
	for _, field := range fields {
		if field.IsSort {
			return true
		}
	}
	return false
}

// Update GenerateFileFromTemplate to include sort information
func GenerateFileFromTemplate(dir, filename, templateFile string, data map[string]interface{}) error {
	tmplContent, err := TemplateFS.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("error reading template %s: %v", templateFile, err)
	}

	funcMap := template.FuncMap{
		"toLower":      ToLower,
		"toUpper":      ToUpper,
		"toTitle":      ToTitle,
		"toCamelCase":  ToCamelCase,
		"toTitleCase":  ToTitle,
		"toPascalCase": ToPascalCase,
		"toSnakeCase":  ToSnakeCase,
		"toPlural":     ToPlural,
		"getInputType": GetInputType,
		"getDartType":  GetDartType,
		"hasSortField": HasSortField,
	}

	// Add sort information to template data
	if fields, ok := data["Fields"].([]FieldStruct); ok {
		data["HasSort"] = HasSortField(fields)
	}

	tmpl, err := template.New(filepath.Base(templateFile)).Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("error parsing template %s: %v", templateFile, err)
	}

	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template for %s: %v", filename, err)
	}

	return nil
}

func ParseTemplate(name, content string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"toLower":      ToLower,
		"toUpper":      ToUpper,
		"toTitle":      ToTitle,
		"toCamelCase":  ToCamelCase,
		"toTitleCase":  ToTitle,
		"toPascalCase": ToPascalCase,
		"toSnakeCase":  ToSnakeCase,
		"toPlural":     ToPlural,
		"getInputType": GetInputType,
		"getDartType":  GetDartType,
	}

	return template.New(name).Funcs(funcMap).Parse(content)
}
