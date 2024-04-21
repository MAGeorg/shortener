package migration

import (
	"context"
	"database/sql"
	"fmt"
)

func runSQLMigration(ctx context.Context, db *sql.DB, query string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, query); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to execute SQL query %q\n%w", query, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
