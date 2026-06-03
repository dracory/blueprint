package category_manager

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
		{"FieldTitle", FieldTitle, "title"},
		{"FieldDescription", FieldDescription, "description"},
		{"FieldStatus", FieldStatus, "status"},
		{"FieldParentID", FieldParentID, "parent_id"},
		{"FieldCreatedAt", FieldCreatedAt, "created_at"},
		{"FieldUpdatedAt", FieldUpdatedAt, "updated_at"},
		{"FieldCategories", FieldCategories, "categories"},
		{"FieldTotal", FieldTotal, "total"},
		{"FieldCategoryID", FieldCategoryID, "category_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}
