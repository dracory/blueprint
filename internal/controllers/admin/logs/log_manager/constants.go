package log_manager

// Log filter query parameter names
const FILTER_LEVEL = "level"
const FILTER_SEARCH_MESSAGE = "search_message"
const FILTER_SEARCH_CONTEXT = "search_context"
const FILTER_SEARCH_MESSAGE_NOT = "search_message_not"
const FILTER_SEARCH_CONTEXT_NOT = "search_context_not"
const FILTER_FROM = "from"
const FILTER_TO = "to"

// Log table query parameter names
const SORT_BY = "sort_by"
const SORT_DIRECTION = "sort_direction"

// Pagination query parameter names
const PAGE = "page"
const PER_PAGE = "per_page"

// Log actions and parameters
const ACTION_DELETE = "delete"
const ACTION_SHOW_CONTEXT = "show_context"
const ACTION_CLOSE_CONTEXT = "close_context"
const ACTION_DELETE_SELECTED = "delete_selected"
const PARAM_LOG_ID = "log_id"
const PARAM_BULK_LOG_IDS = "bulk_log_ids"
