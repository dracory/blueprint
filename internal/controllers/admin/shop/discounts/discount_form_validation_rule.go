package admin

import (
	"github.com/dracory/rule"
)

// DiscountFormData holds the discount form data for validation
type DiscountFormData struct {
	Title        string
	Status       string
	Code         string
	DiscountType string
	StartsAt     string
	EndsAt       string
}

// DiscountFormValidationRule validates discount form fields
type DiscountFormValidationRule struct {
	rule.Rule
	message string
}

// NewDiscountFormValidationRule creates a new DiscountFormValidationRule
func NewDiscountFormValidationRule(data DiscountFormData) *DiscountFormValidationRule {
	r := &DiscountFormValidationRule{}
	r.SetContext(discountFormContext{
		title:        data.Title,
		status:       data.Status,
		code:         data.Code,
		discountType: data.DiscountType,
		startsAt:     data.StartsAt,
		endsAt:       data.EndsAt,
	})
	r.SetCondition(r.passes)
	return r
}

// Message returns the validation error message
func (r *DiscountFormValidationRule) Message() string {
	return r.message
}

// passes validates all required discount fields
func (r *DiscountFormValidationRule) passes(ctx any) bool {
	context := ctx.(discountFormContext)

	if context.title == "" {
		r.message = "title is required"
		return false
	}
	if context.status == "" {
		r.message = "status is required"
		return false
	}
	if context.code == "" {
		r.message = "code is required"
		return false
	}
	if context.discountType == "" {
		r.message = "discount type is required"
		return false
	}
	if context.startsAt == "" {
		r.message = "starts_at is required"
		return false
	}
	if context.endsAt == "" {
		r.message = "ends_at is required"
		return false
	}

	r.message = ""
	return true
}

type discountFormContext struct {
	title        string
	status       string
	code         string
	discountType string
	startsAt     string
	endsAt       string
}
