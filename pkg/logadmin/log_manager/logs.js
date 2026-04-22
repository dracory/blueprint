const { createApp } = Vue;

/**
 * LogsApp is a Vue.js component for managing application logs.
 * It provides a table view with filtering, sorting, pagination,
 * and CRUD operations for logs.
 */
const LogsApp = {
  data() {
    return {
      // UI state
      loading: true,
      deleting: false,
      showFilterModal: false,
      showContextModal: false,

      // Log data
      logs: [],
      totalLogs: 0,
      hasMore: false,

      // Pagination
      currentPage: 0,
      perPage: 100,
      jumpToPage: 1,

      // Filters
      filters: {
        level: '',
        searchMessage: '',
        searchContext: '',
        searchMessageNot: '',
        searchContextNot: '',
        from: '',
        to: ''
      },

      // Sorting
      sortByColumn: 'time',
      sortOrder: 'desc',

      // Selection
      selectedLogs: [],
      selectAll: false,

      // Context modal
      contextLog: null
    };
  },

  computed: {
    /**
     * Returns total number of pages.
     */
    totalPages() {
      return Math.ceil(this.totalLogs / this.perPage);
    },

    /**
     * Returns a human-readable string describing the current filter state.
     */
    filterStatus() {
      const parts = [];
      if (this.filters.level) parts.push(`level: ${this.filters.level}`);
      if (this.filters.searchMessage) parts.push(`message: "${this.filters.searchMessage}"`);
      if (this.filters.searchContext) parts.push(`context: "${this.filters.searchContext}"`);
      if (this.filters.from) parts.push(`from: ${this.filters.from}`);
      if (this.filters.to) parts.push(`to: ${this.filters.to}`);
      
      if (parts.length === 0) return 'Showing all logs';
      return 'Showing logs with ' + parts.join(', ');
    },

    /**
     * Returns true if any filters are active.
     */
    hasActiveFilters() {
      return this.filters.level !== '' ||
             this.filters.searchMessage !== '' ||
             this.filters.searchContext !== '' ||
             this.filters.searchMessageNot !== '' ||
             this.filters.searchContextNot !== '' ||
             this.filters.from !== '' ||
             this.filters.to !== '';
    }
  },

  mounted() {
    // Read filters from URL parameters for shareable URLs
    const urlParams = new URLSearchParams(window.location.search);
    
    this.filters.level = urlParams.get('level') || '';
    this.filters.searchMessage = urlParams.get('search_message') || '';
    this.filters.searchContext = urlParams.get('search_context') || '';
    this.filters.searchMessageNot = urlParams.get('search_message_not') || '';
    this.filters.searchContextNot = urlParams.get('search_context_not') || '';
    this.filters.from = urlParams.get('from') || '';
    this.filters.to = urlParams.get('to') || '';
    this.sortByColumn = urlParams.get('sort_by') || 'time';
    this.sortOrder = urlParams.get('sort_order') || 'desc';
    this.currentPage = parseInt(urlParams.get('page') || '0', 10);
    this.perPage = parseInt(urlParams.get('per_page') || '100', 10);

    this.loadLogs();
  },

  methods: {
    /**
     * Loads logs from the server based on current filters, pagination, and sorting.
     */
    async loadLogs() {
      this.loading = true;
      try {
        const response = await fetch(urlLogsLoad, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            page: this.currentPage,
            per_page: this.perPage,
            level: this.filters.level,
            search_message: this.filters.searchMessage,
            search_context: this.filters.searchContext,
            search_message_not: this.filters.searchMessageNot,
            search_context_not: this.filters.searchContextNot,
            from: this.filters.from,
            to: this.filters.to,
            sort_by: this.sortByColumn,
            sort_order: this.sortOrder
          })
        });
        const data = await response.json();
        
        if (data.status === 'success') {
          this.logs = data.data?.logs || [];
          this.totalLogs = data.data?.total || 0;
          this.hasMore = data.data?.has_more || false;
          // Sync jumpToPage with currentPage to ensure UI consistency
          this.jumpToPage = this.currentPage + 1;
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load logs'
          });
        }
      } catch (error) {
        console.error('Error loading logs:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.message || 'Failed to load logs'
        });
      } finally {
        this.loading = false;
      }
    },

    /**
     * Sorts the logs by the specified column.
     */
    sortBy(column) {
      if (this.sortByColumn === column) {
        this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
      } else {
        this.sortByColumn = column;
        this.sortOrder = 'asc';
      }
      this.currentPage = 0;
      this.applyFilters();
    },

    /**
     * Navigates to the specified page number.
     */
    goToPage(page) {
      if (page < 0) return;
      const maxPage = Math.ceil(this.totalLogs / this.perPage) - 1;
      if (page > maxPage) page = maxPage;
      this.currentPage = page;
      this.jumpToPage = page + 1;
      this.applyFilters();
    },

    /**
     * Changes the per-page setting and reloads logs.
     */
    changePerPage() {
      this.perPage = parseInt(this.perPage, 10);
      this.currentPage = 0;
      this.jumpToPage = 1;
      this.applyFilters();
    },

    /**
     * Deletes the specified log after confirmation.
     */
    async deleteLog(log) {
      const result = await Swal.fire({
        icon: 'warning',
        title: 'Delete Log?',
        text: `Are you sure you want to delete this log entry?`,
        showCancelButton: true,
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonColor: '#dc3545'
      });

      if (!result.isConfirmed) return;

      try {
        const response = await fetch(urlLogDelete, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ log_id: log.id })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Deleted',
            text: 'Log deleted successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.selectedLogs = this.selectedLogs.filter(id => id !== log.id);
          this.loadLogs();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to delete log'
          });
        }
      } catch (error) {
        console.error('Error deleting log:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to delete log'
        });
      }
    },

    /**
     * Deletes selected logs after confirmation.
     */
    async deleteSelected() {
      const result = await Swal.fire({
        icon: 'warning',
        title: 'Delete Selected Logs?',
        text: `Are you sure you want to delete ${this.selectedLogs.length} log(s)?`,
        showCancelButton: true,
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonColor: '#dc3545'
      });

      if (!result.isConfirmed) return;

      this.deleting = true;
      try {
        const response = await fetch(urlLogDeleteSelected, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ bulk_log_ids: this.selectedLogs })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Deleted',
            text: 'Logs deleted successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.selectedLogs = [];
          this.selectAll = false;
          this.loadLogs();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to delete logs'
          });
        }
      } catch (error) {
        console.error('Error deleting logs:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to delete logs'
        });
      } finally {
        this.deleting = false;
      }
    },

    /**
     * Deletes all logs matching current filter criteria after confirmation.
     */
    async deleteAll() {
      const result = await Swal.fire({
        icon: 'warning',
        title: 'Delete All Logs?',
        text: `Are you sure you want to delete ALL ${this.totalLogs} log(s) matching the current filters? This cannot be undone.`,
        showCancelButton: true,
        confirmButtonText: 'Delete All',
        cancelButtonText: 'Cancel',
        confirmButtonColor: '#dc3545'
      });

      if (!result.isConfirmed) return;

      try {
        const response = await fetch(urlLogDeleteAll, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            level: this.filters.level,
            search_message: this.filters.searchMessage,
            search_context: this.filters.searchContext,
            search_message_not: this.filters.searchMessageNot,
            search_context_not: this.filters.searchContextNot,
            from: this.filters.from,
            to: this.filters.to
          })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Deleted',
            text: 'All logs deleted successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.selectedLogs = [];
          this.selectAll = false;
          this.loadLogs();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to delete logs'
          });
        }
      } catch (error) {
        console.error('Error deleting logs:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to delete logs'
        });
      }
    },

    /**
     * Shows the context modal for a log entry.
     */
    async showContext(log) {
      try {
        const response = await fetch(urlLogShowContext, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ log_id: log.id })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          this.contextLog = data.data?.log;
          this.showContextModal = true;
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load log context'
          });
        }
      } catch (error) {
        console.error('Error loading log context:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to load log context'
        });
      }
    },

    /**
     * Formats a timestamp to a human-readable format.
     */
    formatTime(timestamp) {
      if (!timestamp) return '-';
      const date = new Date(timestamp);
      const year = date.getFullYear();
      const dateTime = date.toLocaleString();
      return `${year}<br>${dateTime}`;
    },

    /**
     * Opens the filter modal.
     */
    openFilterModal() {
      this.showFilterModal = true;
    },

    /**
     * Closes the filter modal.
     */
    closeFilterModal() {
      this.showFilterModal = false;
    },

    /**
     * Closes the context modal.
     */
    closeContextModal() {
      this.showContextModal = false;
      this.contextLog = null;
    },

    /**
     * Applies the current filters and updates the URL.
     */
    applyFilters() {
      const params = new URLSearchParams();
      if (this.filters.level) params.set('level', this.filters.level);
      if (this.filters.searchMessage) params.set('search_message', this.filters.searchMessage);
      if (this.filters.searchContext) params.set('search_context', this.filters.searchContext);
      if (this.filters.searchMessageNot) params.set('search_message_not', this.filters.searchMessageNot);
      if (this.filters.searchContextNot) params.set('search_context_not', this.filters.searchContextNot);
      if (this.filters.from) params.set('from', this.filters.from);
      if (this.filters.to) params.set('to', this.filters.to);
      params.set('page', this.currentPage);
      params.set('per_page', this.perPage);
      params.set('sort_order', this.sortOrder);
      params.set('sort_by', this.sortByColumn);

      const newUrl = `${window.location.pathname}?${params.toString()}`;
      window.history.pushState({}, '', newUrl);

      this.closeFilterModal();
      this.loadLogs();
    },

    /**
     * Clears all filters and resets to default state.
     */
    clearFilters() {
      this.filters = {
        level: '',
        searchMessage: '',
        searchContext: '',
        searchMessageNot: '',
        searchContextNot: '',
        from: '',
        to: ''
      };
      this.currentPage = 0;
      this.applyFilters();
    },

    /**
     * Toggles select all logs on the current page.
     * Uses efficient approach to prevent freezing with large selections.
     */
    toggleSelectAll() {
      if (this.selectAll) {
        // Use a Set for better performance with large arrays
        const newSelection = new Set(this.selectedLogs);
        this.logs.forEach(log => newSelection.add(log.id));
        this.selectedLogs = Array.from(newSelection);
      } else {
        // Remove only current page logs from selection
        const currentPageIds = new Set(this.logs.map(log => log.id));
        this.selectedLogs = this.selectedLogs.filter(id => !currentPageIds.has(id));
      }
    }
  },

  watch: {
    selectedLogs(newVal) {
      // Only update selectAll if on current page
      const currentPageIds = new Set(this.logs.map(log => log.id));
      const selectedOnPage = newVal.filter(id => currentPageIds.has(id));
      this.selectAll = selectedOnPage.length === this.logs.length && this.logs.length > 0;
    }
  }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('logs-app');
  if (el) {
    createApp(LogsApp).mount('#logs-app');
  }
});
