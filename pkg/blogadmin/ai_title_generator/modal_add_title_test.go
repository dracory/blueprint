package aititlegenerator

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/testutils"
)

// TestOnAddTitleModal_BasicStructure tests modal generation structure
func TestOnAddTitleModal_BasicStructure(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitleModal(req)

	// Verify modal structure
	if !strings.Contains(result, "ModalAiAddTitle") {
		t.Error("Should contain modal ID")
	}
	if !strings.Contains(result, "modal") {
		t.Error("Should contain modal classes")
	}
	if !strings.Contains(result, "Add Custom Title") {
		t.Error("Should contain modal title")
	}
}

// TestOnAddTitleModal_FormElements tests form elements
func TestOnAddTitleModal_FormElements(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitleModal(req)

	// Verify form elements
	if !strings.Contains(result, "custom_title") {
		t.Error("Should contain custom_title input name")
	}
	if !strings.Contains(result, "Enter a custom blog title") {
		t.Error("Should contain placeholder text")
	}
	if !strings.Contains(result, "form-label") {
		t.Error("Should contain form label styling")
	}
	if !strings.Contains(result, "Approved titles move to the AI post generator") {
		t.Error("Should contain help text")
	}
}

// TestOnAddTitleModal_Buttons tests button elements
func TestOnAddTitleModal_Buttons(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitleModal(req)

	// Verify buttons
	if !strings.Contains(result, "Save Title") {
		t.Error("Should contain submit button text")
	}
	if !strings.Contains(result, "Close") {
		t.Error("Should contain close button")
	}
	if !strings.Contains(result, "btn-primary") {
		t.Error("Should contain primary button class")
	}
	if !strings.Contains(result, "btn-secondary") {
		t.Error("Should contain secondary button class")
	}
}

// TestOnAddTitleModal_HTMXAttributes tests HTMX integration
func TestOnAddTitleModal_HTMXAttributes(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitleModal(req)

	// Verify HTMX attributes
	if !strings.Contains(result, "hx-post") {
		t.Error("Should contain hx-post attribute")
	}
	if !strings.Contains(result, `hx-target="body"`) {
		t.Error("Should contain hx-target='body'")
	}
	if !strings.Contains(result, `hx-swap="beforeend"`) {
		t.Error("Should contain hx-swap='beforeend'")
	}
	if !strings.Contains(result, "hx-indicator") {
		t.Error("Should contain hx-indicator attribute")
	}
	if !strings.Contains(result, "spinner-border") {
		t.Error("Should contain spinner for loading state")
	}
}

// TestOnAddTitleModal_WithCustomTitle tests pre-populated title value
func TestOnAddTitleModal_WithCustomTitle(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "custom_title=Pre-filled+Title",
		},
	}

	result := c.onAddTitleModal(req)

	// The value should be pre-populated
	if !strings.Contains(result, "Pre-filled Title") {
		t.Error("Should contain pre-filled title value")
	}
}

// TestOnAddTitleModal_JavaScriptFunction tests modal close script
func TestOnAddTitleModal_JavaScriptFunction(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitleModal(req)

	// Verify JavaScript close function
	if !strings.Contains(result, "closeModal") {
		t.Error("Should contain closeModal function")
	}
	if !strings.Contains(result, "ModalBackdrop") {
		t.Error("Should contain backdrop class reference")
	}
	if !strings.Contains(result, "getElementById") {
		t.Error("Should contain getElementById for DOM manipulation")
	}
}

// TestOnAddTitleModal_Backdrop tests modal backdrop
func TestOnAddTitleModal_Backdrop(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitleModal(req)

	// Verify backdrop
	if !strings.Contains(result, "modal-backdrop") {
		t.Error("Should contain modal-backdrop class")
	}
	if !strings.Contains(result, "fade show") {
		t.Error("Should contain fade show classes")
	}
}

// TestOnAddTitleModal_ActionParameter tests action parameter in form
func TestOnAddTitleModal_ActionParameter(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitleModal(req)

	// Verify action parameter
	if !strings.Contains(result, ACTION_ADD_TITLE) {
		t.Error("Should contain add_title action")
	}
}

// TestOnAddTitleModal_RequiredAttribute tests required field validation
func TestOnAddTitleModal_RequiredAttribute(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitleModal(req)

	// Verify required attribute
	if !strings.Contains(result, `required="required"`) {
		t.Error("Should contain required attribute on input")
	}
}
