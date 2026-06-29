package admin

import (
	"testing"
)

func TestDiscountFormValidationRule(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		status       string
		code         string
		discountType string
		startsAt     string
		endsAt       string
		wantPass     bool
		wantMessage  string
	}{
		{
			name:         "all fields valid",
			title:        "Summer Sale",
			status:       "active",
			code:         "SUMMER20",
			discountType: "percent",
			startsAt:     "2026-06-01 00:00:00",
			endsAt:       "2026-08-31 23:59:59",
			wantPass:     true,
		},
		{
			name:         "missing title",
			title:        "",
			status:       "active",
			code:         "SUMMER20",
			discountType: "percent",
			startsAt:     "2026-06-01 00:00:00",
			endsAt:       "2026-08-31 23:59:59",
			wantPass:     false,
			wantMessage:  "title is required",
		},
		{
			name:         "missing status",
			title:        "Summer Sale",
			status:       "",
			code:         "SUMMER20",
			discountType: "percent",
			startsAt:     "2026-06-01 00:00:00",
			endsAt:       "2026-08-31 23:59:59",
			wantPass:     false,
			wantMessage:  "status is required",
		},
		{
			name:         "missing code",
			title:        "Summer Sale",
			status:       "active",
			code:         "",
			discountType: "percent",
			startsAt:     "2026-06-01 00:00:00",
			endsAt:       "2026-08-31 23:59:59",
			wantPass:     false,
			wantMessage:  "code is required",
		},
		{
			name:         "missing discount type",
			title:        "Summer Sale",
			status:       "active",
			code:         "SUMMER20",
			discountType: "",
			startsAt:     "2026-06-01 00:00:00",
			endsAt:       "2026-08-31 23:59:59",
			wantPass:     false,
			wantMessage:  "discount type is required",
		},
		{
			name:         "missing starts_at",
			title:        "Summer Sale",
			status:       "active",
			code:         "SUMMER20",
			discountType: "percent",
			startsAt:     "",
			endsAt:       "2026-08-31 23:59:59",
			wantPass:     false,
			wantMessage:  "starts_at is required",
		},
		{
			name:         "missing ends_at",
			title:        "Summer Sale",
			status:       "active",
			code:         "SUMMER20",
			discountType: "percent",
			startsAt:     "2026-06-01 00:00:00",
			endsAt:       "",
			wantPass:     false,
			wantMessage:  "ends_at is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewDiscountFormValidationRule(DiscountFormData{
				Title:        tt.title,
				Status:       tt.status,
				Code:         tt.code,
				DiscountType: tt.discountType,
				StartsAt:     tt.startsAt,
				EndsAt:       tt.endsAt,
			})
			if rule.Passes() != tt.wantPass {
				t.Errorf("Passes() = %v, want %v", rule.Passes(), tt.wantPass)
			}
			if rule.Message() != tt.wantMessage {
				t.Errorf("Message() = %q, want %q", rule.Message(), tt.wantMessage)
			}
		})
	}
}
