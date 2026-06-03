package discount_manager

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
		{"FieldCode", FieldCode, "code"},
		{"FieldAmount", FieldAmount, "amount"},
		{"FieldStatus", FieldStatus, "status"},
		{"FieldCreatedAt", FieldCreatedAt, "created_at"},
		{"FieldUpdatedAt", FieldUpdatedAt, "updated_at"},
		{"FieldDiscounts", FieldDiscounts, "discounts"},
		{"FieldTotal", FieldTotal, "total"},
		{"FieldDiscountID", FieldDiscountID, "discount_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}
