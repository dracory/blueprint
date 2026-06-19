package migrations

import "testing"

func TestTableCustomCreate_InterfaceMethods(t *testing.T) {
	t.Skip("Enable if custom migrations are used")

	migration := &TableCustomCreate{}

	if migration.Signature() != "2026_03_22_0001_table_custom_create" {
		t.Errorf("Expected signature '2026_03_22_0001_table_custom_create', got '%s'", migration.Signature())
	}

	if migration.Description() != "Example: Create custom table with indexes" {
		t.Errorf("Expected description 'Example: Create custom table with indexes', got '%s'", migration.Description())
	}
}
