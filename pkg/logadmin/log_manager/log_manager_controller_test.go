package log_manager

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"project/internal/config"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/logstore"
	"github.com/dracory/test"
)

func TestLogManagerController_RenderPage(t *testing.T) {
	app := testutils.Setup(
		testutils.WithLogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewLogManagerController(app)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	// Create some logs
	logStore := app.GetLogStore()
	l1 := logstore.NewLog()
	l1.SetMessage("Test Message 1")
	l1.SetLevel("info")
	logStore.LogCreate(ctx, l1)

	l2 := logstore.NewLog()
	l2.SetMessage("Test Message 2")
	l2.SetLevel("error")
	logStore.LogCreate(ctx, l2)

	req := httptest.NewRequest(http.MethodGet, "/admin/logs", nil).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "Log Manager") {
		t.Error("expected Log Manager in response")
	}
}

func TestLogManagerController_HandleLoadLogs(t *testing.T) {
	app := testutils.Setup(
		testutils.WithLogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewLogManagerController(app)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	// Create some logs
	logStore := app.GetLogStore()
	l1 := logstore.NewLog()
	l1.SetMessage("Test Message 1")
	l1.SetLevel("info")
	logStore.LogCreate(ctx, l1)

	l2 := logstore.NewLog()
	l2.SetMessage("Test Message 2")
	l2.SetLevel("error")
	logStore.LogCreate(ctx, l2)

	loadData := map[string]any{
		"page":     0,
		"per_page": 10,
	}
	body, _ := json.Marshal(loadData)
	req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLoadLogs, bytes.NewBuffer(body)).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "success") {
		t.Error("expected success in response")
	}
	if !strings.Contains(resp, "Test Message 1") {
		t.Error("expected Test Message 1 in response")
	}
	if !strings.Contains(resp, "Test Message 2") {
		t.Error("expected Test Message 2 in response")
	}
}

func TestLogManagerController_HandleLogShowContext(t *testing.T) {
	app := testutils.Setup(
		testutils.WithLogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewLogManagerController(app)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	// Create some logs
	logStore := app.GetLogStore()
	l1 := logstore.NewLog()
	l1.SetMessage("Test Message 1")
	l1.SetLevel("info")
	logStore.LogCreate(ctx, l1)

	l2 := logstore.NewLog()
	l2.SetMessage("Test Message 2")
	l2.SetLevel("error")
	logStore.LogCreate(ctx, l2)

	logs, _ := logStore.LogList(ctx, logstore.LogQuery())
	logID := logs[0].GetID()

	showData := map[string]string{
		"log_id": logID,
	}
	body, _ := json.Marshal(showData)
	req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLogShowContext, bytes.NewBuffer(body)).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "success") {
		t.Error("expected success in response")
	}
	if !strings.Contains(resp, "context") {
		t.Error("expected context in response")
	}
}

func TestLogManagerController_HandleLogDelete(t *testing.T) {
	app := testutils.Setup(
		testutils.WithLogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewLogManagerController(app)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	// Create some logs
	logStore := app.GetLogStore()
	l1 := logstore.NewLog()
	l1.SetMessage("Test Message 1")
	l1.SetLevel("info")
	logStore.LogCreate(ctx, l1)

	l2 := logstore.NewLog()
	l2.SetMessage("Test Message 2")
	l2.SetLevel("error")
	logStore.LogCreate(ctx, l2)

	logs, _ := logStore.LogList(ctx, logstore.LogQuery())
	logID := logs[0].GetID()

	deleteData := map[string]string{
		"log_id": logID,
	}
	body, _ := json.Marshal(deleteData)
	req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLogDelete, bytes.NewBuffer(body)).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "success") {
		t.Error("expected success in response")
	}

	// Verify deletion
	l, _ := logStore.LogFindByID(ctx, logID)
	if l != nil {
		t.Error("expected log to be nil after deletion")
	}
}

func TestLogManagerController_HandleLogDeleteSelected(t *testing.T) {
	app := testutils.Setup(
		testutils.WithLogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewLogManagerController(app)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	// Create some logs
	logStore := app.GetLogStore()
	l1 := logstore.NewLog()
	l1.SetMessage("Test Message 1")
	l1.SetLevel("info")
	logStore.LogCreate(ctx, l1)

	l2 := logstore.NewLog()
	l2.SetMessage("Test Message 2")
	l2.SetLevel("error")
	logStore.LogCreate(ctx, l2)

	// Create a log first since previous one was deleted
	l := logstore.NewLog()
	l.SetMessage("Delete selected")
	logStore.LogCreate(ctx, l)

	logs, _ := logStore.LogList(ctx, logstore.LogQuery())
	var logID string
	for _, log := range logs {
		if log.GetMessage() == "Delete selected" {
			logID = log.GetID()
			break
		}
	}

	deleteData := map[string][]string{
		"bulk_log_ids": {logID},
	}
	body, _ := json.Marshal(deleteData)
	req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLogDeleteSelected, bytes.NewBuffer(body)).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "success") {
		t.Error("expected success in response")
	}

	// Verify deletion
	lFound, _ := logStore.LogFindByID(ctx, logID)
	if lFound != nil {
		t.Error("expected log to be nil after deletion")
	}
}

func TestLogManagerController_HandleLogDeleteAll(t *testing.T) {
	app := testutils.Setup(
		testutils.WithLogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewLogManagerController(app)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	// Create some logs
	logStore := app.GetLogStore()
	l1 := logstore.NewLog()
	l1.SetMessage("Test Message 1")
	l1.SetLevel("info")
	logStore.LogCreate(ctx, l1)

	l2 := logstore.NewLog()
	l2.SetMessage("Test Message 2")
	l2.SetLevel("error")
	logStore.LogCreate(ctx, l2)

	// Create another log
	l3 := logstore.NewLog()
	l3.SetMessage("Delete me")
	logStore.LogCreate(ctx, l3)

	deleteData := map[string]string{}
	body, _ := json.Marshal(deleteData)
	req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLogDeleteAll, bytes.NewBuffer(body)).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "success") {
		t.Error("expected success in response")
	}

	// Verify deletion
	count, _ := logStore.LogCount(ctx, logstore.LogQuery())
	if int(count) != 0 {
		t.Errorf("expected 0 logs, got %d", count)
	}
}

func TestLogManagerController_RendersVueApp(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithLogStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	html, resp, err := test.CallStringEndpoint(http.MethodGet, NewLogManagerController(app).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	// Page should render Vue app mount point
	if !strings.Contains(html, "logs-app") {
		t.Error("expected logs-app in HTML")
	}
	// Page should include Vue CDN
	if !strings.Contains(html, "vue.global.js") {
		t.Error("expected Vue CDN script in HTML")
	}
	// Page should include SweetAlert2
	if !strings.Contains(html, "sweetalert2") {
		t.Error("expected SweetAlert2 in HTML")
	}
}

func TestLogManagerController_LoadLogsAction(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithLogStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test load-logs action
	queryParams := url.Values{}
	queryParams.Set("action", "load-logs")

	requestBody := map[string]any{
		"page":       0,
		"per_page":   100,
		"sort_order": "desc",
		"sort_by":    "time",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	body, resp, err := test.CallStringEndpoint(http.MethodPost, NewLogManagerController(app).Handler, test.NewRequestOptions{
		GetValues: queryParams,
		Body:      string(bodyBytes),
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	// Should return JSON response
	if !strings.Contains(body, `"status"`) {
		t.Error("expected JSON response with status field")
	}
	if !strings.Contains(body, `"logs"`) {
		t.Error("expected JSON response with logs field")
	}
}
