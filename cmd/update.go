package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update BaseUI Core to the latest version",
	Long:  `Update BaseUI Core to the latest version. This command will update the core directory of your Base project to the latest version available on GitHub.`,
	Run:   updateBaseCore,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateBaseCore(cmd *cobra.Command, args []string) {
	fmt.Println("Updating BaseUI Core...")
	err := updateCore()
	if err != nil {
		fmt.Printf("Error updating BaseUI Core: %v\n", err)
		return
	}
	fmt.Println("BaseUI Core updated successfully.")
}

func updateCore() error {
	// Define the core repository URL
	coreRepoURL := "https://github.com/base-al/baseui-core.git"

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "base-core-update")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the core repository
	gitCmd := exec.Command("git", "clone", coreRepoURL, tempDir)
	if output, err := gitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone core repository: %v\n%s", err, output)
	}

	// Path to the core directory in the current project
	projectCoreDir := filepath.Join(".", "core")

	// Backup the current core directory
	backupDir := projectCoreDir + ".bak"
	if err := os.Rename(projectCoreDir, backupDir); err != nil {
		return fmt.Errorf("failed to backup current core directory: %v", err)
	}

	// Copy core files from temp directory to the project
	tempCoreDir := filepath.Join(tempDir, "core")
	err = filepath.Walk(tempCoreDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(tempCoreDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(projectCoreDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})

	if err != nil {
		// If there's an error, attempt to restore the backup
		os.RemoveAll(projectCoreDir)
		os.Rename(backupDir, projectCoreDir)
		return fmt.Errorf("failed to copy core files: %v", err)
	}

	// Run 'flutter pub get' in the project directory
	flutterCmd := exec.Command("flutter", "pub", "get")
	if output, err := flutterCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run 'flutter pub get': %v\n%s", err, output)
	}

	// Remove the backup directory
	os.RemoveAll(backupDir)

	fmt.Println("Core directory updated successfully.")
	return nil
}
