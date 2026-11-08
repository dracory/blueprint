package models

import (
	"testing"
	"time"

	"github.com/dracory/neat"
)

func TestNoteTableName(t *testing.T) {
	note := Note{}
	if note.TableName() != "note" {
		t.Errorf("TableName() = %q, want %q", note.TableName(), "note")
	}
}

func TestNewNote(t *testing.T) {
	note := NewNote()

	if note.ID == "" {
		t.Error("expected generated ID")
	}
	if note.Status != NOTE_STATUS_ACTIVE {
		t.Errorf("Status = %q, want %q", note.Status, NOTE_STATUS_ACTIVE)
	}
	if note.CreatedAt == "" || note.UpdatedAt == "" {
		t.Error("expected CreatedAt and UpdatedAt to be initialized")
	}
	if _, err := time.Parse("2006-01-02 15:04:05", note.CreatedAt); err != nil {
		t.Errorf("CreatedAt %q is not a valid datetime: %v", note.CreatedAt, err)
	}
	if note.SoftDeletedAt != neat.NullDateTime {
		t.Errorf("SoftDeletedAt = %q, want %q", note.SoftDeletedAt, neat.NullDateTime)
	}
}

func TestNoteSettersChain(t *testing.T) {
	note := NewNote().
		SetID("custom-id").
		SetStatus(NOTE_STATUS_INACTIVE).
		SetTitle("Title").
		SetContent("Content").
		SetCreatedAt("2020-01-01 00:00:00").
		SetUpdatedAt("2020-01-02 00:00:00")

	if note.ID != "custom-id" {
		t.Errorf("ID = %q, want %q", note.ID, "custom-id")
	}
	if note.Status != NOTE_STATUS_INACTIVE {
		t.Errorf("Status = %q, want %q", note.Status, NOTE_STATUS_INACTIVE)
	}
	if note.Title != "Title" {
		t.Errorf("Title = %q, want %q", note.Title, "Title")
	}
	if note.Content != "Content" {
		t.Errorf("Content = %q, want %q", note.Content, "Content")
	}
	if note.CreatedAt != "2020-01-01 00:00:00" {
		t.Errorf("CreatedAt = %q", note.CreatedAt)
	}
	if note.UpdatedAt != "2020-01-02 00:00:00" {
		t.Errorf("UpdatedAt = %q", note.UpdatedAt)
	}
}

func TestNewNoteUniqueIDs(t *testing.T) {
	if NewNote().ID == NewNote().ID {
		t.Error("expected unique IDs across NewNote calls")
	}
}
