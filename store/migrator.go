package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"itsfriday/server/version"
)

//go:embed migration
var migrationFS embed.FS

const (
	LatestSchemaFileName = "LATEST.sql"
)

func (s *Store) Migrate(ctx context.Context) error {
	migrationHistoryList, err := s.driver.FindMigrationHistoryList(ctx, &FindMigrationHistory{})
	if err != nil || len(migrationHistoryList) == 0 {
		filePath := s.getMigrationBasePath() + LatestSchemaFileName
		bytes, err := migrationFS.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read latest schema file: %w", err)
		}
		schemaVersion, err := s.GetCurrentSchemaVersion()
		if err != nil {
			return fmt.Errorf("failed to get current schema version: %w", err)
		}
	
		tx, err := s.driver.GetDB().Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback()
		if err := s.execute(ctx, tx, string(bytes)); err != nil {
			return fmt.Errorf("failed to execute SQL file %s, err %w", filePath, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		if _, err := s.driver.UpsertMigrationHistory(ctx, &UpsertMigrationHistory{
			Version: schemaVersion,
		}); err != nil {
			return fmt.Errorf("failed to upsert migration history: %w", err)
		}
	}

    return nil
}

func (s *Store) getMigrationBasePath() string {
	return fmt.Sprintf("migration/%s/", s.Profile.Driver)
}

func (*Store) execute(ctx context.Context, tx *sql.Tx, stmt string) error {
	if _, err := tx.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("failed to execute statement: %s", err)
	}
	return nil
}

func (s *Store) GetCurrentSchemaVersion() (string, error) {
	currentVersion := version.GetCurrentVersion(s.Profile.Mode)
	minorVersion := version.GetMinorVersion(currentVersion)
	return fmt.Sprintf("%s.0", minorVersion), nil
}
