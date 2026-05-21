package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreBlogMigrate)(nil)

type StoreBlogMigrate struct {
	registry RegistryInterface
}

func (m *StoreBlogMigrate) ID() string {
	return "2026_03_21_0002_store_blog_migrate"
}

func (m *StoreBlogMigrate) Description() string {
	return "Run blog store MigrateUp to create blog tables"
}

func (m *StoreBlogMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetBlogStore()
	if store == nil {
		return errors.New("blog store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreBlogMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	return nil
}

func (m *StoreBlogMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:02:00", "UTC").StdTime()
}
