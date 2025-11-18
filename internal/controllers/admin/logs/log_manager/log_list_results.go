package log_manager

import (
	"project/internal/types"

	"github.com/dracory/logstore"
)

// logListFilters represents the filters, sorting, and paging for a log listing.
type logListFilters struct {
	FilterLevel            string
	FilterSearchMessage    string
	FilterSearchContext    string
	FilterSearchMessageNot string
	FilterSearchContextNot string
	FilterFrom             string
	FilterTo               string
	FilterSortBy           string
	FilterSortDirection    string
	FilterPage             int
	FilterPerPage          int
}

// logListResult represents the result of a log listing.
type logListResult struct {
	Logs    []logstore.LogInterface
	Total   int
	HasMore bool
}

// listLogs encapsulates the log listing logic so it can be reused by
// both controllers and LiveFlux components.
func listLogs(app types.AppInterface, f logListFilters) (logListResult, error) {
	if app == nil || app.GetLogStore() == nil {
		return logListResult{}, nil
	}

	if f.FilterPerPage <= 0 {
		f.FilterPerPage = 100
	}

	if f.FilterPage < 0 {
		f.FilterPage = 0
	}

	offset := f.FilterPage * f.FilterPerPage
	limit := f.FilterPerPage + 1 // lookahead item to detect if there is a next page

	// Build base filtered query (without limit/offset) for counting
	query := logstore.LogQuery()

	if f.FilterLevel != "" {
		query = query.SetLevel(f.FilterLevel)
	}

	if f.FilterSearchMessage != "" {
		query = query.SetMessageContains(f.FilterSearchMessage)
	}

	if f.FilterSearchContext != "" {
		query = query.SetContextContains(f.FilterSearchContext)
	}

	if f.FilterSearchMessageNot != "" {
		query = query.SetMessageNotContains(f.FilterSearchMessageNot)
	}

	if f.FilterSearchContextNot != "" {
		query = query.SetContextNotContains(f.FilterSearchContextNot)
	}

	if f.FilterFrom != "" {
		query = query.SetTimeGte(f.FilterFrom)
	}

	if f.FilterTo != "" {
		query = query.SetTimeLte(f.FilterTo)
	}

	if f.FilterSortBy != "" {
		query = query.SetOrderBy(f.FilterSortBy)
	}

	if f.FilterSortDirection != "" {
		query = query.SetOrderDirection(f.FilterSortDirection)
	}

	// Total count for pagination
	total := 0
	if n, err := app.GetLogStore().LogCount(query); err == nil {
		total = n
	}

	// Apply paging for the list query
	query = query.SetLimit(limit).SetOffset(offset)

	logs, err := app.GetLogStore().LogList(query)
	if err != nil {
		return logListResult{}, err
	}

	hasMore := false
	if len(logs) > f.FilterPerPage {
		hasMore = true
		logs = logs[:f.FilterPerPage]
	}

	return logListResult{
		Logs:    logs,
		Total:   total,
		HasMore: hasMore,
	}, nil
}
