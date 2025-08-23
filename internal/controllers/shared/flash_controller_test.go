package shared

import (
	"net/http"
	"net/url"
	"project/internal/helpers"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/str"
	"github.com/dracory/test"
)

func TestFlash(t *testing.T) {
	app := testutils.Setup()

	body, response, err := test.CallStringEndpoint(http.MethodPost, NewFlashController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"type":    {"success"},
			"message": {"Authentication Provider Error. Once is required field"},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`The message is no longer available`,
		`<a href="/">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

func TestFlashMessage_Info(t *testing.T) {
	app := testutils.Setup()

	infoUrl := helpers.ToFlashInfoURL(app.GetCacheStore(), "This is an info message", "/testbackendpoint", 5)

	flashMessageID := str.RightFrom(infoUrl, `/flash?message_id=`)

	body, response, err := test.CallStringEndpoint(http.MethodPost, NewFlashController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"message_id": {flashMessageID},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`<div class="alert alert-info">`,
		`This is an info message`,
		`<a href="/testbackendpoint">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

func TestFlashMessage_Error(t *testing.T) {
	app := testutils.Setup()

	errorUrl := helpers.ToFlashErrorURL(app.GetCacheStore(), "This is an error message", "/testbackendpoint", 5)

	flashMessageID := str.RightFrom(errorUrl, `/flash?message_id=`)

	body, response, err := test.CallStringEndpoint(http.MethodPost, NewFlashController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"message_id": {flashMessageID},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`<div class="alert alert-danger">`,
		`This is an error message`,
		`<a href="/testbackendpoint">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

func TestFlashMessage_Success(t *testing.T) {
	app := testutils.Setup()

	successUrl := helpers.ToFlashSuccessURL(app.GetCacheStore(), "This is a success message", "/testbackendpoint", 5)

	flashMessageID := str.RightFrom(successUrl, `/flash?message_id=`)

	body, response, err := test.CallStringEndpoint(http.MethodPost, NewFlashController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"message_id": {flashMessageID},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`<div class="alert alert-success">`,
		`This is a success message`,
		`<a href="/testbackendpoint">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

func TestFlashMessage_Warning(t *testing.T) {
	app := testutils.Setup()

	warningUrl := helpers.ToFlashWarningURL(app.GetCacheStore(), "This is a warning message", "/testbackendpoint", 5)

	flashMessageID := str.RightFrom(warningUrl, `/flash?message_id=`)

	body, response, err := test.CallStringEndpoint(http.MethodPost, NewFlashController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"message_id": {flashMessageID},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`<div class="alert alert-warning">`,
		`This is a warning message`,
		`<a href="/testbackendpoint">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

func TestFlashMessage_Get(t *testing.T) {
	app := testutils.Setup()

	infoUrl := helpers.ToFlashInfoURL(app.GetCacheStore(), "This is an info message", "/testbackendpoint", 5)

	flashMessageID := str.RightFrom(infoUrl, `/flash?message_id=`)

	body, response, err := test.CallStringEndpoint(http.MethodGet, NewFlashController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"message_id": {flashMessageID},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`<div class="alert alert-info">`,
		`This is an info message`,
		`<a href="/testbackendpoint">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

func TestFlashMessage_Delete(t *testing.T) {
	app := testutils.Setup()

	infoUrl := helpers.ToFlashInfoURL(app.GetCacheStore(), "This is an info message", "/testbackendpoint", 5)

	flashMessageID := str.RightFrom(infoUrl, `/flash?message_id=`)

	body, response, err := test.CallStringEndpoint(http.MethodDelete, NewFlashController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"message_id": {flashMessageID},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`<div class="alert alert-info">`,
		`This is an info message`,
		`<a href="/testbackendpoint">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

func TestFlashMessage_Post(t *testing.T) {
	app := testutils.Setup()

	infoUrl := helpers.ToFlashInfoURL(app.GetCacheStore(), "This is an info message", "/testbackendpoint", 5)

	flashMessageID := str.RightFrom(infoUrl, `/flash?message_id=`)

	body, response, err := test.CallStringEndpoint(http.MethodPost, NewFlashController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"message_id": {flashMessageID},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`<div class="alert alert-info">`,
		`This is an info message`,
		`<a href="/testbackendpoint">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}

func TestFlashMessage_Put(t *testing.T) {
	app := testutils.Setup()

	infoUrl := helpers.ToFlashInfoURL(app.GetCacheStore(), "This is an info message", "/testbackendpoint", 5)

	flashMessageID := str.RightFrom(infoUrl, `/flash?message_id=`)

	body, response, err := test.CallStringEndpoint(http.MethodPut, NewFlashController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"message_id": {flashMessageID},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatal(`Response MUST be `, http.StatusOK, ` but was: `, response.StatusCode)
	}

	expecteds := []string{
		`<div class="alert alert-info">`,
		`This is an info message`,
		`<a href="/testbackendpoint">Click here to continue</a>`,
	}

	for _, expected := range expecteds {
		if !strings.Contains(body, expected) {
			t.Fatal(`Response MUST contain`, expected, ` but was `, body)
		}
	}
}
