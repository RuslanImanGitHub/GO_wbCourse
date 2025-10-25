package migrations

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// Migration представляет одну миграцию
type Migration struct {
	Version string
	Up      string
	Down    string
}

// Migrator управляет миграциями
type Migrator struct {
	db         *sql.DB
	migrations map[string]Migration
	fsys       fs.FS // Используем fs.FS для работы с встроенными файлами
}

// NewMigrator создает новый экземпляр мигратора
func NewMigrator(db *sql.DB, fsys fs.FS) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make(map[string]Migration),
		fsys:       fsys,
	}
}

// LoadMigrations загружает миграции из файловой системы
func (m *Migrator) LoadMigrations() error {
	return fs.WalkDir(m.fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Парсим имя файла: 001_create_table.up.sql
		baseName := filepath.Base(path)
		parts := strings.Split(baseName, "_")
		if len(parts) < 3 {
			return nil // Пропускаем файлы с неправильным форматом
		}

		version := parts[0]
		direction := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".sql")
		direction = strings.TrimSuffix(direction, ".up")
		direction = strings.TrimSuffix(direction, ".down")

		// Читаем содержимое файла
		content, err := fs.ReadFile(m.fsys, path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		migration, exists := m.migrations[version]
		if !exists {
			migration = Migration{Version: version}
		}

		if strings.Contains(path, ".up.sql") {
			migration.Up = string(content)
		} else if strings.Contains(path, ".down.sql") {
			migration.Down = string(content)
		}

		m.migrations[version] = migration
		return nil
	})
}

// Init создает таблицу для отслеживания миграций
func (m *Migrator) Init() error {
	query := `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `
	_, err := m.db.Exec(query)
	return err
}

// GetAppliedMigrations возвращает список примененных миграций
func (m *Migrator) GetAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

// MigrateUp применяет все непримененные миграции
func (m *Migrator) MigrateUp() error {
	if err := m.Init(); err != nil {
		return err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Сортируем миграции по версии
	var versions []string
	for version := range m.migrations {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	for _, version := range versions {
		if applied[version] {
			continue
		}

		migration := m.migrations[version]
		if migration.Up == "" {
			return fmt.Errorf("no UP migration for version %s", version)
		}

		// Выполняем в транзакции
		tx, err := m.db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", version, err)
		}

		// Записываем в историю
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		log.Printf("Applied migration: %s", version)
	}

	log.Printf("All migrations applied successfully")
	return nil
}

// MigrateDown откатывает миграции
func (m *Migrator) MigrateDown(targetVersion string) error {
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Сортируем версии в обратном порядке
	var versions []string
	for version := range applied {
		versions = append(versions, version)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))

	for _, version := range versions {
		migration, exists := m.migrations[version]
		if !exists || migration.Down == "" {
			return fmt.Errorf("no DOWN migration for version %s", version)
		}

		// Выполняем в транзакции
		tx, err := m.db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(migration.Down); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to revert migration %s: %w", version, err)
		}

		if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		log.Printf("Reverted migration: %s", version)

		if version == targetVersion {
			break
		}
	}

	return nil
}
