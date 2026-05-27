package log_manager

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/logstore"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestLogManagerController_Functional(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithLogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	controller := NewLogManagerController(registry)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	// Create some logs
	logStore := registry.GetLogStore()
	l1 := logstore.NewLog()
	l1.SetMessage("Test Message 1")
	l1.SetLevel("info")
	logStore.LogCreate(ctx, l1)

	l2 := logstore.NewLog()
	l2.SetMessage("Test Message 2")
	l2.SetLevel("error")
	logStore.LogCreate(ctx, l2)

	t.Run("renderPage", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/logs", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		assert.Contains(t, resp, "Log Manager")
	})

	t.Run("handleLoadLogs", func(t *testing.T) {
		loadData := map[string]any{
			"page":     0,
			"per_page": 10,
		}
		body, _ := json.Marshal(loadData)
		req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLoadLogs, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		assert.Contains(t, resp, "success")
		assert.Contains(t, resp, "Test Message 1")
		assert.Contains(t, resp, "Test Message 2")
	})

	t.Run("handleLogShowContext", func(t *testing.T) {
		logs, _ := logStore.LogList(ctx, logstore.LogQuery())
		logID := logs[0].GetID()

		showData := map[string]string{
			"log_id": logID,
		}
		body, _ := json.Marshal(showData)
		req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLogShowContext, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		assert.Contains(t, resp, "success")
		assert.Contains(t, resp, "context")
	})

	t.Run("handleLogDelete", func(t *testing.T) {
		logs, _ := logStore.LogList(ctx, logstore.LogQuery())
		logID := logs[0].GetID()

		deleteData := map[string]string{
			"log_id": logID,
		}
		body, _ := json.Marshal(deleteData)
		req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLogDelete, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		assert.Contains(t, resp, "success")

		// Verify deletion
		l, _ := logStore.LogFindByID(ctx, logID)
		assert.Nil(t, l)
	})

	t.Run("handleLogDeleteSelected", func(t *testing.T) {
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
		assert.Contains(t, resp, "success")

		// Verify deletion
		lFound, _ := logStore.LogFindByID(ctx, logID)
		assert.Nil(t, lFound)
	})

	t.Run("handleLogDeleteAll", func(t *testing.T) {
		// Create another log
		l3 := logstore.NewLog()
		l3.SetMessage("Delete me")
		logStore.LogCreate(ctx, l3)

		deleteData := map[string]string{}
		body, _ := json.Marshal(deleteData)
		req := httptest.NewRequest(http.MethodPost, "/admin/logs?action="+actionLogDeleteAll, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		assert.Contains(t, resp, "success")

		// Verify deletion
		count, _ := logStore.LogCount(ctx, logstore.LogQuery())
		assert.Equal(t, int(0), int(count))
	})
}
