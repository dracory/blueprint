package log_manager

import (
	"context"
	"errors"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/logstore"
	"github.com/stretchr/testify/assert"
)

type fakeLogStore struct {
	logstore.StoreInterface

	logsToReturn []logstore.LogInterface
	count        int
	listErr      error
	countErr     error
}

func (s *fakeLogStore) LogList(ctx context.Context, q logstore.LogQueryInterface) ([]logstore.LogInterface, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.logsToReturn, nil
}

func (s *fakeLogStore) LogCount(ctx context.Context, q logstore.LogQueryInterface) (int, error) {
	if s.countErr != nil {
		return 0, s.countErr
	}
	return s.count, nil
}

func TestListLogs_NilApp_ReturnsEmptyNoError(t *testing.T) {
	var appNil logstore.StoreInterface

	result, err := listLogs(nil, logListFilters{})

	assert.NoError(t, err)
	assert.Empty(t, result.Logs)
	assert.Equal(t, 0, result.Total)
	assert.False(t, result.HasMore)
	_ = appNil
}

func TestListLogs_NilLogStore_ReturnsEmptyNoError(t *testing.T) {
	app := testutils.Setup()

	result, err := listLogs(app, logListFilters{})

	assert.NoError(t, err)
	assert.Empty(t, result.Logs)
	assert.Equal(t, 0, result.Total)
	assert.False(t, result.HasMore)
}

func TestListLogs_UsesDefaultsAndHasMoreTrimming(t *testing.T) {
	app := testutils.Setup(testutils.WithLogStore(true))

	filters := logListFilters{
		FilterPerPage: 0,
		FilterPage:    -1,
	}

	// Seed 101 real log entries via the application's logger.
	for i := 0; i < 101; i++ {
		app.GetLogger().Info("test log")
	}

	result, err := listLogs(app, filters)

	assert.NoError(t, err)
	assert.Equal(t, 100, len(result.Logs))
	assert.Equal(t, 101, result.Total)
	assert.True(t, result.HasMore)
}

func TestListLogs_NoMorePagesWhenAtOrBelowPerPage(t *testing.T) {
	app := testutils.Setup(testutils.WithLogStore(true))

	// Seed 100 real log entries via the application's logger.
	for i := 0; i < 100; i++ {
		app.GetLogger().Info("test log")
	}

	// First page: should have 50 results and indicate there is a next page.
	filtersPage0 := logListFilters{
		FilterPerPage: 50,
		FilterPage:    0,
	}

	result, err := listLogs(app, filtersPage0)

	assert.NoError(t, err)
	assert.Equal(t, 50, len(result.Logs))
	assert.Equal(t, 100, result.Total)
	assert.True(t, result.HasMore)

	// Second page: also 50 results but no further pages available.
	filtersPage1 := logListFilters{
		FilterPerPage: 50,
		FilterPage:    1,
	}

	result, err = listLogs(app, filtersPage1)

	assert.NoError(t, err)
	assert.Equal(t, 50, len(result.Logs))
	assert.Equal(t, 100, result.Total)
	assert.False(t, result.HasMore)
}

func TestListLogs_LogCountErrorDoesNotFailListing(t *testing.T) {
	app := testutils.Setup(testutils.WithLogStore(true))

	logs := []logstore.LogInterface{nil, nil, nil}

	fakeStore := &fakeLogStore{
		logsToReturn: logs,
		countErr:     errors.New("count failed"),
	}

	app.SetLogStore(fakeStore)

	filters := logListFilters{
		FilterPerPage: 10,
		FilterPage:    0,
	}

	result, err := listLogs(app, filters)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(result.Logs))
	assert.Equal(t, 0, result.Total)
	assert.False(t, result.HasMore)
}

func TestListLogs_LogListErrorIsReturned(t *testing.T) {
	app := testutils.Setup(testutils.WithLogStore(true))

	fakeStore := &fakeLogStore{
		listErr: errors.New("list failed"),
	}

	app.SetLogStore(fakeStore)

	filters := logListFilters{
		FilterPerPage: 10,
		FilterPage:    0,
	}

	result, err := listLogs(app, filters)

	assert.Error(t, err)
	assert.Empty(t, result.Logs)
}
