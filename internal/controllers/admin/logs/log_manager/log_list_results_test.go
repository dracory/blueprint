package log_manager

import (
	"context"
	"errors"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/logstore"
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

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Logs) != 0 {
		t.Error("expected empty logs")
	}
	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
	if result.HasMore {
		t.Error("expected HasMore to be false")
	}
	_ = appNil
}

func TestListLogs_NilLogStore_ReturnsEmptyNoError(t *testing.T) {
	registry := testutils.Setup()

	result, err := listLogs(registry, logListFilters{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Logs) != 0 {
		t.Error("expected empty logs")
	}
	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
	if result.HasMore {
		t.Error("expected HasMore to be false")
	}
}

func TestListLogs_UsesDefaultsAndHasMoreTrimming(t *testing.T) {
	registry := testutils.Setup(testutils.WithLogStore(true))

	filters := logListFilters{
		FilterPerPage: 0,
		FilterPage:    -1,
	}

	// Seed 101 real log entries via the application's logger.
	for i := 0; i < 101; i++ {
		registry.GetLogger().Info("test log")
	}

	result, err := listLogs(registry, filters)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Logs) != 100 {
		t.Errorf("expected 100 logs, got %d", len(result.Logs))
	}
	if result.Total != 101 {
		t.Errorf("expected total 101, got %d", result.Total)
	}
	if !result.HasMore {
		t.Error("expected HasMore to be true")
	}
}

func TestListLogs_NoMorePagesWhenAtOrBelowPerPage(t *testing.T) {
	registry := testutils.Setup(testutils.WithLogStore(true))

	// Seed 100 real log entries via the application's logger.
	for i := 0; i < 100; i++ {
		registry.GetLogger().Info("test log")
	}

	// First page: should have 50 results and indicate there is a next page.
	filtersPage0 := logListFilters{
		FilterPerPage: 50,
		FilterPage:    0,
	}

	result, err := listLogs(registry, filtersPage0)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Logs) != 50 {
		t.Errorf("expected 50 logs, got %d", len(result.Logs))
	}
	if result.Total != 100 {
		t.Errorf("expected total 100, got %d", result.Total)
	}
	if !result.HasMore {
		t.Error("expected HasMore to be true for first page")
	}

	// Second page: also 50 results but no further pages available.
	filtersPage1 := logListFilters{
		FilterPerPage: 50,
		FilterPage:    1,
	}

	result, err = listLogs(registry, filtersPage1)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Logs) != 50 {
		t.Errorf("expected 50 logs, got %d", len(result.Logs))
	}
	if result.Total != 100 {
		t.Errorf("expected total 100, got %d", result.Total)
	}
	if result.HasMore {
		t.Error("expected HasMore to be false for last page")
	}
}

func TestListLogs_LogCountErrorDoesNotFailListing(t *testing.T) {
	registry := testutils.Setup(testutils.WithLogStore(true))

	logs := []logstore.LogInterface{nil, nil, nil}

	fakeStore := &fakeLogStore{
		logsToReturn: logs,
		countErr:     errors.New("count failed"),
	}

	registry.SetLogStore(fakeStore)

	filters := logListFilters{
		FilterPerPage: 10,
		FilterPage:    0,
	}

	result, err := listLogs(registry, filters)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Logs) != 3 {
		t.Errorf("expected 3 logs, got %d", len(result.Logs))
	}
	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
	if result.HasMore {
		t.Error("expected HasMore to be false")
	}
}

func TestListLogs_LogListErrorIsReturned(t *testing.T) {
	registry := testutils.Setup(testutils.WithLogStore(true))

	fakeStore := &fakeLogStore{
		listErr: errors.New("list failed"),
	}

	registry.SetLogStore(fakeStore)

	filters := logListFilters{
		FilterPerPage: 10,
		FilterPage:    0,
	}

	result, err := listLogs(registry, filters)

	if err == nil {
		t.Error("expected error")
	}
	if len(result.Logs) != 0 {
		t.Error("expected empty logs on error")
	}
}
