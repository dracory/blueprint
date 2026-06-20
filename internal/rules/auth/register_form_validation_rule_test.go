package authrules

import (
	"testing"
)

func TestRegisterFormValidationRule_AllFieldsValid_Passes(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Country:   "US",
		Timezone:  "America/New_York",
	})
	if r.Fails() {
		t.Errorf("expected rule to pass with all valid fields, got: %s", r.Message())
	}
}

func TestRegisterFormValidationRule_MissingFirstName_Fails(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		LastName: "Doe",
		Country:  "US",
		Timezone: "America/New_York",
	})
	if r.Passes() {
		t.Error("expected rule to fail when first name is missing")
	}
	if r.Message() != "First name is required field" {
		t.Errorf("expected message 'First name is required field', got: %s", r.Message())
	}
}

func TestRegisterFormValidationRule_MissingLastName_Fails(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		FirstName: "John",
		Country:   "US",
		Timezone:  "America/New_York",
	})
	if r.Passes() {
		t.Error("expected rule to fail when last name is missing")
	}
	if r.Message() != "Last name is required field" {
		t.Errorf("expected message 'Last name is required field', got: %s", r.Message())
	}
}

func TestRegisterFormValidationRule_MissingCountry_Fails(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		FirstName: "John",
		LastName:  "Doe",
		Timezone:  "America/New_York",
	})
	if r.Passes() {
		t.Error("expected rule to fail when country is missing")
	}
	if r.Message() != "Country is required field" {
		t.Errorf("expected message 'Country is required field', got: %s", r.Message())
	}
}

func TestRegisterFormValidationRule_MissingTimezone_Fails(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		FirstName: "John",
		LastName:  "Doe",
		Country:   "US",
	})
	if r.Passes() {
		t.Error("expected rule to fail when timezone is missing")
	}
	if r.Message() != "Timezone is required field" {
		t.Errorf("expected message 'Timezone is required field', got: %s", r.Message())
	}
}

func TestRegisterFormValidationRule_PasswordTooShort_Fails(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		FirstName: "John",
		LastName:  "Doe",
		Country:   "US",
		Timezone:  "America/New_York",
		Password:  "short",
	})
	if r.Passes() {
		t.Error("expected rule to fail when password is too short")
	}
	if r.Message() != "Password must be at least 8 characters" {
		t.Errorf("expected message 'Password must be at least 8 characters', got: %s", r.Message())
	}
}

func TestRegisterFormValidationRule_PasswordMismatch_Fails(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		FirstName:       "John",
		LastName:        "Doe",
		Country:         "US",
		Timezone:        "America/New_York",
		Password:        "password123",
		PasswordConfirm: "different456",
	})
	if r.Passes() {
		t.Error("expected rule to fail when passwords do not match")
	}
	if r.Message() != "Passwords do not match" {
		t.Errorf("expected message 'Passwords do not match', got: %s", r.Message())
	}
}

func TestRegisterFormValidationRule_ValidWithPassword_Passes(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		FirstName:       "John",
		LastName:        "Doe",
		Country:         "US",
		Timezone:        "America/New_York",
		Password:        "securePass123",
		PasswordConfirm: "securePass123",
	})
	if r.Fails() {
		t.Errorf("expected rule to pass with valid password, got: %s", r.Message())
	}
}

func TestRegisterFormValidationRule_EmptyPassword_Skips_PasswordValidation(t *testing.T) {
	r := NewRegisterFormValidationRule(RegisterFormData{
		FirstName: "John",
		LastName:  "Doe",
		Country:   "US",
		Timezone:  "America/New_York",
	})
	if r.Fails() {
		t.Errorf("expected rule to pass when password is not provided, got: %s", r.Message())
	}
}
