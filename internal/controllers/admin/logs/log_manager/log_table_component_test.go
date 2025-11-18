package log_manager

import (
	"context"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/logstore"
	"github.com/stretchr/testify/assert"
)

func TestLogTableComponent_Mount_SetsDefaultsAndCallsLoadLogs(t *testing.T) {
	app := testutils.Setup()

	c := &logTableComponent{App: app}

	err := c.Mount(context.Background(), map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, 0, c.Page)
	assert.Equal(t, 100, c.PerPage)
	assert.Equal(t, logstore.COLUMN_TIME, c.SortBy)
	assert.Equal(t, "desc", c.SortDirection)
}

func TestLogTableComponent_HandleSort_TogglesDirectionAndUsesDefaultColumn(t *testing.T) {
	app := testutils.Setup()
	c := &logTableComponent{App: app}
	ctx := context.Background()

	err := c.Handle(ctx, "sort", nil)
	assert.NoError(t, err)
	assert.Equal(t, logstore.COLUMN_TIME, c.SortBy)
	assert.Equal(t, "asc", c.SortDirection)

	err = c.Handle(ctx, "sort", url.Values{})
	assert.NoError(t, err)
	assert.Equal(t, logstore.COLUMN_TIME, c.SortBy)
	assert.Equal(t, "desc", c.SortDirection)
}

func TestLogTableComponent_LoadLogs_PopulatesFieldsFromListLogs(t *testing.T) {
	app := testutils.Setup(testutils.WithLogStore(true))

	logs := []logstore.LogInterface{nil, nil}
	fakeStore := &fakeLogStore{
		logsToReturn: logs,
		count:        2,
	}

	app.SetLogStore(fakeStore)

	c := &logTableComponent{
		App:           app,
		Level:         "info",
		SearchMessage: "foo",
		SearchContext: "bar",
		From:          "2024-01-01",
		To:            "2024-12-31",
		SortBy:        logstore.COLUMN_TIME,
		SortDirection: "desc",
		Page:          1,
		PerPage:       50,
	}

	err := c.loadLogs()

	assert.NoError(t, err)
	assert.Equal(t, logs, c.Logs)
	assert.Equal(t, 2, c.Total)
	assert.False(t, c.HasMore)
}
