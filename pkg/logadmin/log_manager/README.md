# Log Manager Subcontroller

The log manager provides a web interface for viewing, filtering, and deleting application logs.

## Features

- **Filter Logs**: Filter by level (trace, debug, info, warning, error)
- **Search**: Search within message and context fields
- **Exclude**: Exclude logs matching certain message/context patterns
- **Date Range**: Filter logs by time range
- **Sorting**: Sort by any column (time, level, etc.)
- **Pagination**: Paginated view with configurable per-page setting
- **Bulk Actions**: Delete selected logs or delete all logs matching filters
- **Context Viewer**: View detailed context for individual log entries
- **Vue.js Interface**: Modern reactive UI with no page reloads

## Known Limitations

### Context Field Database Retrieval

**Issue**: The log list API (`handleLoadLogs`) still retrieves the `context` field from the database even though it is not sent to the frontend. This is due to a limitation in the `logstore.LogQuery()` API which does not support field selection.

**Impact**: 
- Database queries retrieve all fields including potentially large context data
- Memory is used to store context data that is immediately discarded
- Network bandwidth between database and application is wasted on unused data

**Current Optimization**:
- Context is not sent to the frontend in JSON responses
- Context is only fetched on-demand when user clicks "View Context"
- This reduces frontend payload size but doesn't reduce database query cost

**Future Improvement**:
To fully optimize this, the `logstore` library would need to be enhanced to support field selection, for example:
```go
query = logstore.LogQuery().SelectFields("id", "time", "level", "message")
```

This would allow the database query to only retrieve the fields needed for the list view, significantly improving performance for logs with large context data.

## API Endpoints

### load-logs
**Method**: POST  
**Purpose**: Load logs with filtering, sorting, and pagination  
**Body**: JSON with filter, sort, and pagination parameters  
**Response**: JSON with logs array, total count, and has_more flag

### delete-log
**Method**: POST  
**Purpose**: Delete a single log entry  
**Body**: JSON with `log_id`  
**Response**: Success/error status

### delete-selected
**Method**: POST  
**Purpose**: Delete multiple log entries  
**Body**: JSON with `bulk_log_ids` array  
**Response**: Success/error status

### delete-all
**Method**: POST  
**Purpose**: Delete all logs matching current filter criteria  
**Body**: JSON with filter criteria  
**Response**: Success/error status with deleted_count

### show-context
**Method**: POST  
**Purpose**: Retrieve full log entry including context  
**Body**: JSON with `log_id`  
**Response**: JSON with complete log data including context

## Data Structures

### logListFilters
Used for filtering and pagination:
- `FilterLevel`: Log level filter
- `FilterSearchMessage`: Search within message
- `FilterSearchContext`: Search within context
- `FilterSearchMessageNot`: Exclude message pattern
- `FilterSearchContextNot`: Exclude context pattern
- `FilterFrom`: Start date/time
- `FilterTo`: End date/time
- `FilterSortBy`: Sort column
- `FilterSortDirection`: Sort direction (asc/desc)
- `FilterPage`: Page number
- `FilterPerPage`: Items per page

### logListResult
Result of log listing operation:
- `Logs`: Array of log entries
- `Total`: Total count matching filters
- `HasMore`: Whether more pages exist

## Vue.js Component

The log manager uses an embedded Vue.js component (`logs.html` + `logs.js`) for the frontend interface:

- **State Management**: Reactive state for filters, pagination, and selection
- **Shareable URLs**: Filter and pagination state preserved in URL parameters
- **SweetAlert2**: Used for confirmations and alerts
- **Bootstrap 5**: Styling with dark mode support

## Testing

Run tests with:
```bash
go test ./pkg/logadmin/log_manager/...
```

Tests cover:
- Vue app rendering
- Load logs API endpoint
- ListLogs function with various filter combinations
