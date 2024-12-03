package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/gertd/go-pluralize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var pluralizeClient *pluralize.Client

func init() {
	pluralizeClient = pluralize.NewClient()
}

func GetDartType(t string) string {
	switch t {
	case "int", "int64", "uint", "uint64":
		return "int"
	case "string", "text":
		return "String"
	case "datetime", "time":
		return "DateTime"
	case "float", "float64", "double":
		return "double"
	case "bool", "boolean":
		return "bool"
	case "list":
		return "List<dynamic>"
	case "map":
		return "Map<String, dynamic>"
	case "binary", "blob":
		return "Uint8List"
	case "decimal":
		return "Decimal"
	case "json":
		return "Map<String, dynamic>"
	default:
		return t
	}
}

func GetInputType(dartType string) string {
	switch dartType {
	case "int":
		return "number"
	case "double":
		return "number"
	case "bool":
		return "checkbox"
	case "DateTime":
		return "datetime-local"
	default:
		return "text"
	}
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToSlug(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func ToTitle(s string) string {
	return cases.Title(language.Und).String(s)
}

func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func ToCamelCase(s string) string {
	s = ToPascalCase(s)
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func ToPascalCase(s string) string {
	words := splitIntoWords(s)
	for i, word := range words {
		words[i] = cases.Title(language.Und).String(word)
	}
	return strings.Join(words, "")
}

func ToPlural(s string) string {
	return pluralizeClient.Plural(s)
}

func splitIntoWords(s string) []string {
	var words []string
	var currentWord strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(rune(s[i-1])) || unicode.IsLower(r)) {
			words = append(words, currentWord.String())
			currentWord.Reset()
		}
		if r == '_' || r == ' ' || r == '-' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		} else {
			currentWord.WriteRune(r)
		}
	}
	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}
	return words
}
func UpdateFlutterRoutes(singularName, pluralName string) error {
	// Update main routes.dart
	mainRoutesFile := "lib/app/routes.dart"
	content, err := os.ReadFile(mainRoutesFile)
	if err != nil {
		return err
	}

	pluralSnake := ToSnakeCase(pluralName)

	// Add import for module routes
	importStr := fmt.Sprintf(`import '../modules/%s/routes.dart';`, pluralSnake)
	content = addFlutterImport(content, importStr)

	// Add module routes to AppPages
	routeStr := fmt.Sprintf(`...%sRoutes.routes,`, pluralName)
	content = addRoutesToAppPages(content, routeStr)

	if err := os.WriteFile(mainRoutesFile, content, 0644); err != nil {
		return fmt.Errorf("error writing routes file: %v", err)
	}

	// Update navbar.dart
	navbarFile := "lib/app/navbar.dart"
	navContent, err := os.ReadFile(navbarFile)
	if err != nil {
		return err
	}

	// Add import statements
	navContent = addFlutterImport(navContent, fmt.Sprintf(`import '../modules/%s/routes.dart';`, pluralSnake))

	// Add nav item
	navItem := fmt.Sprintf(`  NavLink(
    icon: Icon(Icons.list),
    label: '%s',
    path: %sRoute.list,
  ),`, pluralName, pluralName)

	navContent = addNavItem(navContent, navItem)

	return os.WriteFile(navbarFile, navContent, 0644)
}

func addNavItem(content []byte, navItem string) []byte {
	if bytes.Contains(content, []byte(navItem)) {
		return content
	}

	// Find end of destinations array
	marker := []byte("];")
	if idx := bytes.LastIndex(content, marker); idx != -1 {
		// Insert before the closing bracket
		return append(content[:idx], append([]byte(navItem+"\n"), content[idx:]...)...)
	}

	return content
}
func addRoutesToAppPages(content []byte, routeStr string) []byte {
	if bytes.Contains(content, []byte(routeStr)) {
		return content
	}

	// Find "// MODULE PAGES" marker
	marker := []byte("// MODULE PAGES")
	if idx := bytes.Index(content, marker); idx != -1 {
		// Insert at the marker position
		insertStr := []byte(routeStr + "\n")
		return append(content[:idx], append(insertStr, content[idx:]...)...)
	}

	return content
}

func addFlutterImport(content []byte, importStr string) []byte {
	if bytes.Contains(content, []byte(importStr)) {
		return content
	}

	// Find "// MODULE IMPORTS" marker
	marker := []byte("// MODULE IMPORTS")
	if idx := bytes.Index(content, marker); idx != -1 {
		// Insert before the marker
		insertStr := []byte(importStr + "\n")
		return append(content[:idx], append(insertStr, content[idx:]...)...)
	}

	return content
}
