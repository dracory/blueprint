package user_update

import (
	"testing"
)

func TestFieldConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"FieldStatus", FieldStatus, "status"},
		{"FieldRole", FieldRole, "role"},
		{"FieldFirstName", FieldFirstName, "first_name"},
		{"FieldLastName", FieldLastName, "last_name"},
		{"FieldEmail", FieldEmail, "email"},
		{"FieldBusinessName", FieldBusinessName, "business_name"},
		{"FieldPhone", FieldPhone, "phone"},
		{"FieldCountry", FieldCountry, "country"},
		{"FieldTimezone", FieldTimezone, "timezone"},
		{"FieldMemo", FieldMemo, "memo"},
		{"FieldStatusField", FieldStatusField, "field_status"},
		{"FieldCountries", FieldCountries, "countries"},
		{"FieldTimezones", FieldTimezones, "timezones"},
		{"FieldIsoCode2", FieldIsoCode2, "iso_code_2"},
		{"FieldName", FieldName, "name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}
