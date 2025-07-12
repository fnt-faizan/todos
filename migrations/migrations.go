package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations runs all SQL migration files in the migrations directory
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Get all SQL files in migrations directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}

	// Sort files by name (assuming they're numbered)
	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, filepath.Join(migrationsDir, file.Name()))
		}
	}
	sort.Strings(sqlFiles)

	// Execute each migration file
	for _, file := range sqlFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %w", file, err)
		}

		sqlStatements := string(content)
		_, err = db.Exec(sqlStatements)
		if err != nil {
			return fmt.Errorf("error executing migration %s: %w", file, err)
		}
		log.Printf("Successfully executed migration: %s", file)
	}

	return nil
}
