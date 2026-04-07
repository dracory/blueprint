package log_manager

import (
	"context"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/logstore"
)

func TestLogTableComponent_Mount_SetsDefaultsAndCallsLoadLogs(t *testing.T) {
	registry := testutils.Setup()

	c := &logTableComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Page != 0 {
		t.Errorf("expected page 0, got %d", c.Page)
	}
	if c.PerPage != 100 {
		t.Errorf("expected perPage 100, got %d", c.PerPage)
	}
	if c.SortBy != logstore.COLUMN_TIME {
		t.Errorf("expected sortBy %s, got %s", logstore.COLUMN_TIME, c.SortBy)
	}
	if c.SortDirection != "desc" {
		t.Errorf("expected sortDirection desc, got %s", c.SortDirection)
	}
}

func TestLogTableComponent_HandleSort_TogglesDirectionAndUsesDefaultColumn(t *testing.T) {
	registry := testutils.Setup()
	c := &logTableComponent{registry: registry}
	ctx := context.Background()

	err := c.Handle(ctx, "sort", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.SortBy != logstore.COLUMN_TIME {
		t.Errorf("expected sortBy %s, got %s", logstore.COLUMN_TIME, c.SortBy)
	}
	if c.SortDirection != "asc" {
		t.Errorf("expected sortDirection asc, got %s", c.SortDirection)
	}

	err = c.Handle(ctx, "sort", url.Values{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.SortBy != logstore.COLUMN_TIME {
		t.Errorf("expected sortBy %s, got %s", logstore.COLUMN_TIME, c.SortBy)
	}
	if c.SortDirection != "desc" {
		t.Errorf("expected sortDirection desc, got %s", c.SortDirection)
	}
}

func TestLogTableComponent_LoadLogs_PopulatesFieldsFromListLogs(t *testing.T) {
	registry := testutils.Setup(testutils.WithLogStore(true))

	// Seed real log entries via the application's logger.
	for i := 0; i < 2; i++ {
		registry.GetLogger().Info("test log")
	}

	c := &logTableComponent{
		registry:      registry,
		Level:         "",
		SearchMessage: "",
		SearchContext: "",
		From:          "",
		To:            "",
		SortBy:        logstore.COLUMN_TIME,
		SortDirection: "desc",
		Page:          0,
		PerPage:       50,
	}

	err := c.loadLogs()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(c.Logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(c.Logs))
	}
	if c.Total != 2 {
		t.Errorf("expected total 2, got %d", c.Total)
	}
	if c.HasMore {
		t.Error("expected HasMore to be false")
	}
}
