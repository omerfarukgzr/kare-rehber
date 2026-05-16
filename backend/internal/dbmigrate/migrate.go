package dbmigrate

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var migrationsFS embed.FS

// New, embedded migrate instance ile beraber döner.
// Çağıran dispose etmek için Close() çağırmalı.
func New(databaseURL string) (*migrate.Migrate, error) {
	d, err := iofs.New(migrationsFS, "sql")
	if err != nil {
		return nil, fmt.Errorf("iofs source: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("migrate init: %w", err)
	}
	return m, nil
}

// RunUp embedded migration'ları en sona kadar uygular. ErrNoChange normal kabul edilir.
func RunUp(databaseURL string) error {
	m, err := New(databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}
