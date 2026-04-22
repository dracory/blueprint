package log_manager

import (
	"embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/logadmin/shared"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/logstore"
	"github.com/dracory/req"
	"github.com/dracory/sb"
)

//go:embed *.html
//go:embed *.js
var logsFiles embed.FS

const (
	actionLoadLogs          = "load-logs"
	actionLogDelete         = "delete-log"
	actionLogDeleteSelected = "delete-selected"
	actionLogDeleteAll      = "delete-all"
	actionLogShowContext    = "show-context"
)

// == CONTROLLER ==============================================================

type logManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewLogManagerController(registry registry.RegistryInterface) *logManagerController {
	return &logManagerController{registry: registry}
}

func (controller *logManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadLogs:
		return controller.handleLoadLogs(w, r)
	case actionLogDelete:
		return controller.handleLogDelete(w, r)
	case actionLogDeleteSelected:
		return controller.handleLogDeleteSelected(w, r)
	case actionLogDeleteAll:
		return controller.handleLogDeleteAll(w, r)
	case actionLogShowContext:
		return controller.handleLogShowContext(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *logManagerController) renderPage(r *http.Request) string {
	if controller.registry.GetLogStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "Log store is not initialized", links.Admin().Home(), 10)
	}

	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Home(), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Logs", URL: links.Admin().Logs()},
	})

	heading := hb.Heading1().HTML("Log Manager")

	htmlContent, err := logsFiles.ReadFile("logs.html")
	if err != nil {
		slog.Error("Failed to read logs HTML template", "error", err)
		return hb.Div().HTML("Error loading logs component").ToHTML()
	}

	jsContent, err := logsFiles.ReadFile("logs.js")
	if err != nil {
		slog.Error("Failed to read logs JavaScript file", "error", err)
		return hb.Div().HTML("Error loading logs component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const urlLogsLoad = '` + shared.NewLinks("/admin/logs").LogManager(map[string]string{"action": actionLoadLogs}) + `';
		const urlLogDelete = '` + shared.NewLinks("/admin/logs").LogManager(map[string]string{"action": actionLogDelete}) + `';
		const urlLogDeleteSelected = '` + shared.NewLinks("/admin/logs").LogManager(map[string]string{"action": actionLogDeleteSelected}) + `';
		const urlLogDeleteAll = '` + shared.NewLinks("/admin/logs").LogManager(map[string]string{"action": actionLogDeleteAll}) + `';
		const urlLogShowContext = '` + shared.NewLinks("/admin/logs").LogManager(map[string]string{"action": actionLogShowContext}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueContainer)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Log Manager",
		Content: content,
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *logManagerController) handleLoadLogs(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	logStore := controller.registry.GetLogStore()
	if logStore == nil {
		return api.Error("Log store not available").ToString()
	}

	// Parse request body
	var reqBody struct {
		Page             int    `json:"page"`
		PerPage          int    `json:"per_page"`
		SortOrder        string `json:"sort_order"`
		SortBy           string `json:"sort_by"`
		Level            string `json:"level"`
		SearchMessage    string `json:"search_message"`
		SearchContext    string `json:"search_context"`
		SearchMessageNot string `json:"search_message_not"`
		SearchContextNot string `json:"search_context_not"`
		From             string `json:"from"`
		To               string `json:"to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	// Helper functions for trimmed values with defaults (similar to req.GetStringTrimmedOr)
	getInt := func(val int, defaultVal int) int {
		if val == 0 {
			return defaultVal
		}
		return val
	}

	getStringTrimmed := func(val string, defaultVal string) string {
		val = strings.TrimSpace(val)
		if val == "" {
			return defaultVal
		}
		return val
	}

	page := getInt(reqBody.Page, 0)
	perPage := getInt(reqBody.PerPage, 100)
	sortOrder := getStringTrimmed(reqBody.SortOrder, sb.DESC)
	sortBy := getStringTrimmed(reqBody.SortBy, logstore.COLUMN_TIME)
	level := getStringTrimmed(reqBody.Level, "")
	searchMessage := getStringTrimmed(reqBody.SearchMessage, "")
	searchContext := getStringTrimmed(reqBody.SearchContext, "")
	searchMessageNot := getStringTrimmed(reqBody.SearchMessageNot, "")
	searchContextNot := getStringTrimmed(reqBody.SearchContextNot, "")
	from := getStringTrimmed(reqBody.From, "")
	to := getStringTrimmed(reqBody.To, "")

	// Build query with column selection to exclude context (performance optimization)
	query := logstore.LogQuery().
		SetColumns([]string{
			logstore.COLUMN_ID,
			logstore.COLUMN_TIME,
			logstore.COLUMN_LEVEL,
			logstore.COLUMN_MESSAGE,
		})

	if level != "" {
		query = query.SetLevel(level)
	}
	if searchMessage != "" {
		query = query.SetMessageContains(searchMessage)
	}
	if searchContext != "" {
		query = query.SetContextContains(searchContext)
	}
	if searchMessageNot != "" {
		query = query.SetMessageNotContains(searchMessageNot)
	}
	if searchContextNot != "" {
		query = query.SetContextNotContains(searchContextNot)
	}
	if from != "" {
		query = query.SetTimeGte(from)
	}
	if to != "" {
		query = query.SetTimeLte(to)
	}
	if sortBy != "" {
		query = query.SetOrderBy(sortBy)
	}
	if sortOrder != "" {
		query = query.SetOrderDirection(sortOrder)
	}

	// Total count
	total, err := logStore.LogCount(ctx, query)
	if err != nil {
		slog.Error("Failed to get logs count", "error", err)
		return api.Error("Failed to get logs count").ToString()
	}

	// Apply pagination
	offset := page * perPage
	limit := perPage + 1 // lookahead item to detect if there is a next page
	query = query.SetLimit(limit).SetOffset(offset)

	logs, err := logStore.LogList(ctx, query)
	if err != nil {
		slog.Error("Failed to load logs", "error", err)
		return api.Error("Failed to load logs").ToString()
	}

	hasMore := false
	if len(logs) > perPage {
		hasMore = true
		logs = logs[:perPage]
	}

	// Convert to JSON-friendly format (exclude context as it can be large)
	logList := []map[string]any{}
	for _, log := range logs {
		logList = append(logList, map[string]any{
			"id":      log.GetID(),
			"time":    log.GetTime(),
			"level":   log.GetLevel(),
			"message": log.GetMessage(),
		})
	}

	return api.SuccessWithData("Logs loaded successfully", map[string]any{
		"logs":     logList,
		"total":    total,
		"has_more": hasMore,
	}).ToString()
}

func (controller *logManagerController) handleLogDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	logStore := controller.registry.GetLogStore()
	if logStore == nil {
		return api.Error("Log store not available").ToString()
	}

	// Parse request body
	var reqBody struct {
		LogID string `json:"log_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.LogID == "" {
		return api.Error("Log ID is required").ToString()
	}

	if err := logStore.LogDeleteByID(ctx, reqBody.LogID); err != nil {
		slog.Error("Failed to delete log", "error", err)
		return api.Error("Failed to delete log").ToString()
	}

	return api.Success("Log deleted successfully").ToString()
}

func (controller *logManagerController) handleLogDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	logStore := controller.registry.GetLogStore()
	if logStore == nil {
		return api.Error("Log store not available").ToString()
	}

	// Parse request body
	var reqBody struct {
		BulkLogIDs []string `json:"bulk_log_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if len(reqBody.BulkLogIDs) == 0 {
		return api.Error("No log IDs provided").ToString()
	}

	// Use bulk delete for better performance
	if err := logStore.LogDeleteByIDs(ctx, reqBody.BulkLogIDs); err != nil {
		slog.Error("Failed to delete logs", "error", err)
		return api.Error("Failed to delete logs").ToString()
	}

	return api.Success("Logs deleted successfully").ToString()
}

func (controller *logManagerController) handleLogDeleteAll(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	logStore := controller.registry.GetLogStore()
	if logStore == nil {
		return api.Error("Log store not available").ToString()
	}

	// Parse request body with filter criteria
	var reqBody struct {
		Level            string `json:"level"`
		SearchMessage    string `json:"search_message"`
		SearchContext    string `json:"search_context"`
		SearchMessageNot string `json:"search_message_not"`
		SearchContextNot string `json:"search_context_not"`
		From             string `json:"from"`
		To               string `json:"to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	// Build query to find all logs matching criteria
	query := logstore.LogQuery()

	if reqBody.Level != "" {
		query = query.SetLevel(reqBody.Level)
	}
	if reqBody.SearchMessage != "" {
		query = query.SetMessageContains(reqBody.SearchMessage)
	}
	if reqBody.SearchContext != "" {
		query = query.SetContextContains(reqBody.SearchContext)
	}
	if reqBody.SearchMessageNot != "" {
		query = query.SetMessageNotContains(reqBody.SearchMessageNot)
	}
	if reqBody.SearchContextNot != "" {
		query = query.SetContextNotContains(reqBody.SearchContextNot)
	}
	if reqBody.From != "" {
		query = query.SetTimeGte(reqBody.From)
	}
	if reqBody.To != "" {
		query = query.SetTimeLte(reqBody.To)
	}

	// Get all matching logs (no limit)
	logs, err := logStore.LogList(ctx, query)
	if err != nil {
		slog.Error("Failed to load logs for deletion", "error", err)
		return api.Error("Failed to load logs").ToString()
	}

	// Delete all matching logs
	for _, log := range logs {
		_ = logStore.LogDeleteByID(ctx, log.GetID())
	}

	return api.SuccessWithData("All logs deleted successfully", map[string]any{
		"deleted_count": len(logs),
	}).ToString()
}

func (controller *logManagerController) handleLogShowContext(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	logStore := controller.registry.GetLogStore()
	if logStore == nil {
		return api.Error("Log store not available").ToString()
	}

	// Parse request body
	var reqBody struct {
		LogID string `json:"log_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqBody.LogID == "" {
		return api.Error("Log ID is required").ToString()
	}

	logEntry, err := logStore.LogFindByID(ctx, reqBody.LogID)
	if err != nil || logEntry == nil {
		return api.Error("Log not found").ToString()
	}

	return api.SuccessWithData("Log context loaded successfully", map[string]any{
		"log": map[string]any{
			"id":      logEntry.GetID(),
			"time":    logEntry.GetTime(),
			"level":   logEntry.GetLevel(),
			"message": logEntry.GetMessage(),
			"context": logEntry.GetContext(),
		},
	}).ToString()
}
