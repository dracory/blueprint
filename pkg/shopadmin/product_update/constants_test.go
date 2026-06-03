package product_update

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
		{"FieldPrice", FieldPrice, "price"},
		{"FieldProduct", FieldProduct, "product"},
		{"FieldProductID", FieldProductID, "product_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}
