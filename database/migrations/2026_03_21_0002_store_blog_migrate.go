package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreBlogMigrate)(nil)

type StoreBlogMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreBlogMigrate) Signature() string {
	return "2026_03_21_0002_store_blog_migrate"
}

func (m *StoreBlogMigrate) Description() string {
	return "Run blog store MigrateUp to create blog tables"
}

func (m *StoreBlogMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetBlogStore()
	if store == nil {
		return errors.New("blog store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreBlogMigrate) Down() error {
	store := m.app.GetBlogStore()
	if store == nil {
		return errors.New("blog store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

