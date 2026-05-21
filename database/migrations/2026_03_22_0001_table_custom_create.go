package migrations

import (
	"context"
	"database/sql"
	"time"

	"github.com/dracory/migrate"

	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*TableCustomCreate)(nil)

// TableCustomCreate is an example template for custom SQL migrations.
// Copy this file and adapt it to create your own table migrations.
// Uncomment and modify the Table/Column definitions below to match your schema.
// Register new migrations in registry.go getSQLMigrations().
type TableCustomCreate struct{}

func (m *TableCustomCreate) ID() string {
	return "2026_03_22_0001_table_custom_create"
}

func (m *TableCustomCreate) Description() string {
	return "Example: Create custom table with indexes"
}

func (m *TableCustomCreate) Up(ctx context.Context, tx *sql.Tx) error {
	// dialect := database.DatabaseType(tx)

	// tableCreateSql, err := sb.NewBuilder(dialect).
	// 	Table(models.UserTableName).
	// 	Column(sb.Column{
	// 		Name:       models.UserColumnID,
	// 		Type:       sb.COLUMN_TYPE_STRING,
	// 		Length:     40,
	// 		PrimaryKey: true,
	// 	}).
	// 	Column(sb.Column{
	// 		Name:     models.UserColumnEmail,
	// 		Type:     sb.COLUMN_TYPE_STRING,
	// 		Length:   255,
	// 		Unique:   true,
	// 		Nullable: false,
	// 	}).
	// 	Column(sb.Column{
	// 		Name:     models.ColumnName,
	// 		Type:     sb.COLUMN_TYPE_STRING,
	// 		Length:   255,
	// 		Nullable: false,
	// 	}).
	// 	Column(sb.Column{
	// 		Name:     models.UserColumnPasswordHash,
	// 		Type:     sb.COLUMN_TYPE_STRING,
	// 		Nullable: false,
	// 	}).
	// 	Column(sb.Column{
	// 		Name:     models.UserColumnCreatedAt,
	// 		Type:     sb.COLUMN_TYPE_DATETIME,
	// 		Nullable: false,
	// 	}).
	// 	Column(sb.Column{
	// 		Name:     models.UserColumnUpdatedAt,
	// 		Type:     sb.COLUMN_TYPE_DATETIME,
	// 		Nullable: false,
	// 	}).
	// 	Column(sb.Column{
	// 		Name:     models.UserColumnSoftDeleted,
	// 		Type:     sb.COLUMN_TYPE_DATETIME,
	// 		Nullable: false,
	// 		Default:  models.UserSoftDeletedMaxTime, // Default to MAX_DATETIME for active users
	// 	}).
	// 	Create()

	// if err != nil {
	// 	return err
	// }

	// _, err = tx.Exec(tableCreateSql)
	// if err != nil {
	// 	return err
	// }

	// // Create indexes using sb builder
	// indexSQL1, err := sb.NewBuilder(dialect).
	// 	Table("users").
	// 	CreateIndex("idx_users_email", "email")
	// if err != nil {
	// 	return err
	// }

	// indexSQL2, err := sb.NewBuilder(dialect).
	// 	Table("users").
	// 	CreateIndex("idx_users_created_at", "created_at")
	// if err != nil {
	// 	return err
	// }

	// // Execute indexes
	// for _, indexSQL := range []string{indexSQL1, indexSQL2} {
	// 	if _, err := tx.Exec(indexSQL); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (m *TableCustomCreate) Down(ctx context.Context, tx *sql.Tx) error {
	dialect := database.DatabaseType(tx)
	tableDropSql, err := sb.NewBuilder(dialect).
		Table("users").
		Drop()
	if err != nil {
		return err
	}

	_, err = tx.Exec(tableDropSql)
	return err
}

func (m *TableCustomCreate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-22 00:01:00", "UTC").StdTime()
}
