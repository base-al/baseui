package cmd

import (
	"base/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:     "generate [name] [field:type...]",
	Aliases: []string{"g"},
	Short:   "Generate a new Flutter UI module",
	Long: `Generate a new Flutter UI module with the specified name and fields.
			This command creates a complete module structure including:
			- Models for data representation
			- Controllers for business logic
			- Services for API communication
			- Views for user interface
			- Bindings for dependency injection
			- Routes for navigation`,
	Args: cobra.MinimumNArgs(1),
	Run:  generateUIModule,
}

func init() {
	rootCmd.AddCommand(guiCmd)
}

func generateUIModule(cmd *cobra.Command, args []string) {
	singularName := args[0]
	fields := args[1:]

	fmt.Println("\n=== Starting Module Generation ===")

	// Convert names to appropriate cases
	dirName := utils.ToSnakeCase(singularName)
	pluralName := utils.ToPascalCase(utils.ToPlural(singularName))
	pluralDirName := utils.ToSnakeCase(pluralName)
	structName := utils.ToPascalCase(singularName)

	fmt.Printf("Names configured:\n")
	fmt.Printf("- dirName: %s\n", dirName)
	fmt.Printf("- pluralName: %s\n", pluralName)
	fmt.Printf("- pluralDirName: %s\n", pluralDirName)
	fmt.Printf("- structName: %s\n", structName)

	// Base module directory
	baseDir := filepath.Join("lib", "modules", pluralDirName)
	fmt.Printf("\nBase directory: %s\n", baseDir)

	// Create module directory structure
	dirs := []string{
		filepath.Join(baseDir, "bindings"),
		filepath.Join(baseDir, "controllers"),
		filepath.Join(baseDir, "models"),
		filepath.Join(baseDir, "services"),
		filepath.Join(baseDir, "views"),
		filepath.Join(baseDir, "views", "mode"),
	}

	fmt.Println("\nCreating directories:")
	for _, dir := range dirs {
		fmt.Printf("- Creating: %s\n", dir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Process fields into FieldStruct
	processedFields := utils.GenerateFieldStructs(fields)
	fmt.Printf("\nProcessed %d fields\n", len(processedFields))
	hasSort := utils.HasSortField(processedFields)

	// Template data
	data := map[string]interface{}{
		"StructName":            structName,
		"PluralName":            pluralName,
		"LowerStructName":       utils.ToCamelCase(singularName),
		"LowerPluralStructName": utils.ToSnakeCase(pluralName),
		"TableName":             utils.ToSnakeCase(pluralName),
		"RouteName":             utils.ToSnakeCase(pluralName),
		"PackageName":           dirName,
		"Fields":                processedFields,
		"HasSort":               hasSort,
	}

	// Generate files using Flutter templates
	files := map[string]string{
		"routes.dart":                                            "routes.tmpl",
		filepath.Join("models", dirName+".dart"):                 "model.tmpl",
		filepath.Join("controllers", dirName+"_controller.dart"): "controller.tmpl",
		filepath.Join("services", dirName+"_service.dart"):       "service.tmpl",
		filepath.Join("bindings", dirName+"_binding.dart"):       "binding.tmpl",
		filepath.Join("views", "create.dart"):                    "create.tmpl",
		filepath.Join("views", "edit.dart"):                      "edit.tmpl",
		filepath.Join("views", "show.dart"):                      "show.tmpl",
		filepath.Join("views", "index.dart"):                     "index.tmpl",
		filepath.Join("views", "mode", "list.dart"):              "list.tmpl",
		filepath.Join("views", "mode", "grid.dart"):              "grid.tmpl",
	}
	if hasSort {
		files[filepath.Join("views", "mode", "sort.dart")] = "sort.tmpl"
	}
	fmt.Println("\nGenerating files:")
	for filePath, templateName := range files {
		fullPath := filepath.Join(baseDir, filePath)
		templatePath := "templates/" + templateName
		fmt.Printf("\nProcessing file: %s\n", filePath)
		fmt.Printf("- Template: %s\n", templatePath)
		fmt.Printf("- Full path: %s\n", fullPath)

		err := utils.GenerateFileFromTemplate(
			filepath.Join(baseDir, filepath.Dir(filePath)),
			filepath.Base(filePath),
			templatePath,
			data,
		)
		if err != nil {
			fmt.Printf("Error generating file %s: %v\n", filePath, err)
			return
		}
	}

	// Update main routes.dart
	fmt.Println("\nUpdating routes...")
	if err := utils.UpdateFlutterRoutes(singularName, pluralName); err != nil {
		fmt.Printf("Error updating routes.dart: %v\n", err)
		return
	}

	fmt.Printf("\n=== Module Generation Complete ===\n")
	fmt.Printf("Flutter module %s generated successfully in %s\n", singularName, baseDir)
	fmt.Printf("Fields: %v\n", fields)
	fmt.Println("\nRoutes have been updated in lib/app/routes.dart")
}
