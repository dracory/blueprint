package user_manager

import (
	"testing"
)

func TestFieldConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"FieldID", FieldID, "id"},
		{"FieldFirstName", FieldFirstName, "first_name"},
		{"FieldLastName", FieldLastName, "last_name"},
		{"FieldEmail", FieldEmail, "email"},
		{"FieldStatus", FieldStatus, "status"},
		{"FieldCreatedAt", FieldCreatedAt, "created_at"},
		{"FieldUpdatedAt", FieldUpdatedAt, "updated_at"},
		{"FieldUsers", FieldUsers, "users"},
		{"FieldTotal", FieldTotal, "total"},
		{"FieldUserID", FieldUserID, "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}
