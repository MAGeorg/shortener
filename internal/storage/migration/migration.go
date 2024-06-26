// пакет, реализующий работу с миграциями.
package migration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MAGeorg/shortener.git/internal/storage/migration/postgres"
)

// структура миграцмий, на данный момент содержит путь до одного файла с миграцими.
type Migration struct {
	Source string
	Flag   string
}

// Up runs an up migration.
func (m *Migration) Up(db *sql.DB) error {
	ctx := context.Background()
	return m.UpContext(ctx, db)
}

// UpContext runs an up migration.
func (m *Migration) UpContext(ctx context.Context, db *sql.DB) error {
	// проверка наличие схемы.
	if res := CheckExistScheme(ctx, db); !res {
		if err := m.run(ctx, db); err != nil {
			return err
		}
	}
	return nil
}

// Down runs a down migration.
func (m *Migration) Down(db *sql.DB) error {
	ctx := context.Background()
	return m.DownContext(ctx, db)
}

// DownContext runs a down migration.
func (m *Migration) DownContext(ctx context.Context, db *sql.DB) error {
	if err := m.run(ctx, db); err != nil {
		return err
	}
	return nil
}

func (m *Migration) run(ctx context.Context, db *sql.DB) error {
	switch m.Flag {
	case "sql":
		f, err := os.Open(m.Source)
		if err != nil {
			return fmt.Errorf("error %v: failed to open SQL migration file: %w", filepath.Base(m.Source), err)
		}
		defer f.Close()

		statements, err := ParseSQLMigration(f)
		if err != nil {
			return fmt.Errorf("error %v: failed to parse SQL migration file: %w", filepath.Base(m.Source), err)
		}
		fmt.Println("success read init migration scheme")

		start := time.Now()
		if err := runSQLMigration(ctx, db, statements[0]); err != nil {
			fmt.Printf("error run sql migration: %s", err.Error())
			return fmt.Errorf("error %v: failed to run SQL migration: %w", filepath.Base(m.Source), err)
		}
		finish := time.Since(start)

		//nolint:forbidigo // FP
		fmt.Printf("\nmigration execution time: %d ms\n", finish)

	case "go":
		queueUP := postgres.Up()

		start := time.Now()
		if err := runSQLMigration(ctx, db, queueUP); err != nil {
			fmt.Printf("error run sql migration: %s", err.Error())
			return fmt.Errorf("error up go-migration: failed to run SQL migration: %w", err)
		}
		finish := time.Since(start)

		//nolint:forbidigo // FP
		fmt.Printf("\nmigration execution time: %d ms\n", finish)
	}

	return nil
}

// проверка, что схемы еще нет.
func CheckExistScheme(ctx context.Context, db *sql.DB) bool {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false
	}

	query := `
	SELECT EXISTS (
		SELECT * FROM 
		shot_url
	);`

	res, err := tx.QueryContext(ctx, query)
	if err != nil {
		return false
	}
	defer res.Close()

	return res.Err() == nil
}
