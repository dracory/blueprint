package models

import (
	"context"
	"time"

	"project/internal/platform/database"

	basedb "github.com/dracory/base/database"
)

// Model represents a base model with common functionality
type Model struct {
	DB *database.Database
}

// NewModel creates a new Model instance
func NewModel(db *database.Database) *Model {
	return &Model{
		DB: db,
	}
}

// Entity represents a base entity with common fields
type Entity struct {
	ID        string     `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// Store represents a data store interface
type Store interface {
	// Common store methods
	FindByID(ctx context.Context, id string) (interface{}, error)
	FindAll(ctx context.Context) ([]interface{}, error)
	Create(ctx context.Context, entity interface{}) error
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id string) error
}

// GetQueryableContext returns a QueryableContext for the database
func (m *Model) GetQueryableContext(ctx context.Context) basedb.QueryableContext {
	return basedb.Context(ctx, m.DB.GetDB())
}

// WithTransaction executes a function within a database transaction
func (m *Model) WithTransaction(ctx context.Context, fn func(basedb.QueryableContext) error) error {
	// Begin a transaction
	tx, err := m.DB.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Create a transaction context
	txCtx := basedb.Context(ctx, tx)

	// Defer a rollback in case anything fails
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Execute the function
	err = fn(txCtx)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
