// Package models holds example neat ORM models.
//
// This is an EXAMPLE: blueprint's convention is the dracory stores
// (entitystore, userstore, ...), so nothing here is wired up.
//
// Use this folder while you are still exploring a domain and don't yet
// know the final shape of a reusable store package: prototype with a
// plain model + neat's query builder (see dracory/pico), then convert to
// a store package once the fields and queries settle.
package models

import (
	"time"

	"github.com/dracory/neat"
	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/support/uid"
)

const (
	NoteTableName = "note"

	NOTE_STATUS_ACTIVE   = "active"
	NOTE_STATUS_INACTIVE = "inactive"
)

// Note is an example ORM model. Datetime fields are stored as
// ISO 8601 strings.
type Note struct {
	orm.ShortID
	Status        string `json:"status" db:"status"`
	Title         string `json:"title" db:"title"`
	Content       string `json:"content" db:"content"`
	CreatedAt     string `json:"created_at" db:"created_at"`
	UpdatedAt     string `json:"updated_at" db:"updated_at"`
	SoftDeletedAt string `json:"soft_deleted_at" db:"soft_deleted_at"`
}

// TableName returns the database table name for Note.
func (Note) TableName() string {
	return NoteTableName
}

// NewNote creates a new Note instance with a generated short ID, active
// status, and datetime fields initialized to current UTC
// (CreatedAt/UpdatedAt) or the neat.NullDateTime sentinel (SoftDeletedAt).
func NewNote() *Note {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	return &Note{
		ShortID:       orm.ShortID{ID: uid.GenerateShortID()},
		Status:        NOTE_STATUS_ACTIVE,
		CreatedAt:     now,
		UpdatedAt:     now,
		SoftDeletedAt: neat.NullDateTime,
	}
}

// SetID sets the ID and returns the pointer for chaining.
func (n *Note) SetID(id string) *Note {
	n.ID = id
	return n
}

// SetStatus sets the status and returns the pointer for chaining.
func (n *Note) SetStatus(status string) *Note {
	n.Status = status
	return n
}

// SetTitle sets the title and returns the pointer for chaining.
func (n *Note) SetTitle(title string) *Note {
	n.Title = title
	return n
}

// SetContent sets the content and returns the pointer for chaining.
func (n *Note) SetContent(content string) *Note {
	n.Content = content
	return n
}

// SetCreatedAt sets the created_at datetime string and returns the pointer for chaining.
func (n *Note) SetCreatedAt(s string) *Note {
	n.CreatedAt = s
	return n
}

// SetUpdatedAt sets the updated_at datetime string and returns the pointer for chaining.
func (n *Note) SetUpdatedAt(s string) *Note {
	n.UpdatedAt = s
	return n
}
