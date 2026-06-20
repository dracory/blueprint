package migrations

import (
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*TableCustomCreate)(nil)

// TableCustomCreate is an example template for custom schema migrations.
// Copy this file and adapt it to create your own table migrations.
// Register new migrations in registry.go getSQLMigrations().
type TableCustomCreate struct {
	migrator.BaseMigration
}

func (m *TableCustomCreate) Signature() string {
	return "2026_03_22_0001_table_custom_create"
}

func (m *TableCustomCreate) Description() string {
	return "Example: Create custom table with indexes"
}

func (m *TableCustomCreate) Up() error {
	if m.GetSchema().HasTable("custom_example") {
		return nil
	}

	return m.GetSchema().Create("custom_example", func(blueprint contractsschema.Blueprint) {
		blueprint.ID()
		blueprint.String("name")
		blueprint.String("email")
		blueprint.Unique("email")
		blueprint.String("status")
		blueprint.Timestamps()
	})
}

func (m *TableCustomCreate) Down() error {
	return m.GetSchema().DropIfExists("custom_example")
}
