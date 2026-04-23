const { createApp } = Vue;

/**
 * BlogPostsApp is a Vue.js component for managing blog posts.
 * It provides a table view with filtering, sorting, pagination,
 * and CRUD operations for posts.
 */
const BlogPostsApp = {
  data() {
    return {
      // UI state
      loading: true, // Whether posts are being loaded
      showCreateModal: false, // Whether the create post modal is visible
      showFilterModal: false, // Whether the filter modal is visible
      creating: false, // Whether a post is being created

      // Post data
      posts: [], // Array of post objects
      totalPosts: 0, // Total number of posts

      // Pagination
      currentPage: 0, // Current page number (0-indexed)
      perPage: 10, // Number of posts per page

      // Filters
      filters: {
        search: '', // Search query for title/content
        status: '', // Post status filter (draft, published, etc.)
        dateFrom: '', // Start date filter
        dateTo: '' // End date filter
      },

      // Sorting
      sortByColumn: 'created_at', // Column to sort by
      sortOrder: 'desc', // Sort order (asc or desc)

      // Create post form
      createForm: {
        title: '' // Post title for new post
      }
    };
  },

  computed: {
    /**
     * Returns the total number of pages based on total posts and per page setting.
     */
    totalPages() {
      return Math.ceil(this.totalPosts / this.perPage);
    },

    /**
     * Returns an array of visible page numbers for pagination.
     * Shows up to 5 pages centered around the current page.
     */
    visiblePages() {
      const pages = [];
      const start = Math.max(0, this.currentPage - 2);
      const end = Math.min(this.totalPages - 1, this.currentPage + 2);
      
      for (let i = start; i <= end; i++) {
        pages.push(i);
      }
      return pages;
    },

    /**
     * Returns a human-readable string describing the current filter state.
     */
    filterStatus() {
      const parts = [];
      if (this.filters.search) parts.push(`search: "${this.filters.search}"`);
      if (this.filters.status) parts.push(`status: ${this.filters.status}`);
      if (this.filters.dateFrom) parts.push(`from: ${this.filters.dateFrom}`);
      if (this.filters.dateTo) parts.push(`to: ${this.filters.dateTo}`);
      
      if (parts.length === 0) return 'Showing all posts';
      return 'Showing posts with ' + parts.join(', ');
    }
  },

  mounted() {
    // Read filters from URL parameters for shareable URLs
    const urlParams = new URLSearchParams(window.location.search);
    
    this.filters.search = urlParams.get('search') || '';
    this.filters.status = urlParams.get('status') || '';
    this.filters.dateFrom = urlParams.get('date_from') || '';
    this.filters.dateTo = urlParams.get('date_to') || '';
    this.sortByColumn = urlParams.get('sort_by') || 'created_at';
    this.sortOrder = urlParams.get('sort_order') || 'desc';
    this.currentPage = parseInt(urlParams.get('page') || '0', 10);
    this.perPage = parseInt(urlParams.get('per_page') || '10', 10);

    // Set default date range if not in URL
    if (!this.filters.dateFrom || !this.filters.dateTo) {
      const today = new Date();
      const lastYear = new Date(today);
      lastYear.setFullYear(today.getFullYear() - 1);
      
      this.filters.dateTo = today.toISOString().split('T')[0];
      this.filters.dateFrom = lastYear.toISOString().split('T')[0];
    }
    
    this.loadPosts();
  },

  methods: {
    /**
     * Loads posts from the server based on current filters, pagination, and sorting.
     * Makes a GET request to the posts load API endpoint.
     */
    async loadPosts() {
      this.loading = true;
      try {
        const params = new URLSearchParams({
          page: this.currentPage,
          per_page: this.perPage,
          search: this.filters.search,
          status: this.filters.status,
          date_from: this.filters.dateFrom,
          date_to: this.filters.dateTo,
          sort_order: this.sortOrder,
          sort_by: this.sortByColumn
        });

        const separator = urlPostsLoad.includes('?') ? '&' : '?';
        const url = `${urlPostsLoad}${separator}${params.toString()}`;
        
        const response = await fetch(url);
        const data = await response.json();
        
        if (data.status === 'success') {
          this.posts = data.data?.posts || [];
          this.totalPosts = data.data?.total || 0;
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load posts'
          });
        }
      } catch (error) {
        console.error('Error loading posts:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: error.message || 'Failed to load posts'
        });
      } finally {
        this.loading = false;
      }
    },

    /**
     * Sorts the posts by the specified column.
     * Toggles sort order if clicking the same column.
     */
    sortBy(column) {
      if (this.sortByColumn === column) {
        // Toggle sort order
        this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
      } else {
        // New column, default to asc
        this.sortByColumn = column;
        this.sortOrder = 'asc';
      }
      this.currentPage = 0; // Reset to first page when sorting
      this.applyFilters(); // This will update URL and reload
    },

    /**
     * Navigates to the specified page number.
     */
    goToPage(page) {
      if (page < 0 || page >= this.totalPages) return;
      this.currentPage = page;
      this.applyFilters(); // This will update URL and reload
    },

    /**
     * Deletes the specified post after confirmation.
     */
    async deletePost(post) {
      const result = await Swal.fire({
        icon: 'warning',
        title: 'Delete Post?',
        text: `Are you sure you want to delete "${post.title}"?`,
        showCancelButton: true,
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonColor: '#dc3545'
      });

      if (!result.isConfirmed) return;

      try {
        const response = await fetch(urlPostDelete, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ post_id: post.id })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Deleted',
            text: 'Post deleted successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.loadPosts();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to delete post'
          });
        }
      } catch (error) {
        console.error('Error deleting post:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to delete post'
        });
      }
    },

    /**
     * Formats a date string to a human-readable format (e.g., "15 Apr 2026").
     */
    formatDate(dateString) {
      if (!dateString) return '-';
      const date = new Date(dateString);
      const options = { day: 'numeric', month: 'short', year: 'numeric' };
      return date.toLocaleDateString('en-GB', options);
    },

    /**
     * Returns the URL for viewing a post on the website.
     */
    getWebsitePostUrl(postId, slug) {
      return `/blog/post/${postId}/${slug}`;
    },

    /**
     * Returns the URL for the AI post content update page.
     */
    getAiPostContentUrl(postId) {
      return urlAiPostContentUpdate.replace('POST_ID_PLACEHOLDER', postId);
    },

    /**
     * Returns the URL for the post update/edit page.
     */
    getPostUpdateUrl(postId) {
      return urlPostUpdate.replace('POST_ID_PLACEHOLDER', postId);
    },

    /**
     * Opens the create post modal and resets the form.
     */
    openCreateModal() {
      this.createForm.title = '';
      this.showCreateModal = true;
    },

    /**
     * Closes the create post modal and resets the form.
     */
    closeCreateModal() {
      this.showCreateModal = false;
      this.createForm.title = '';
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
     * Applies the current filters and updates the URL.
     * This makes the filter state shareable via URL.
     */
    applyFilters() {
      // Update URL with filter parameters using GET
      const urlParams = new URLSearchParams(window.location.search);
      
      // Preserve controller parameter
      const controller = urlParams.get('controller');
      
      const params = new URLSearchParams();
      if (controller) params.set('controller', controller);
      if (this.filters.search) params.set('search', this.filters.search);
      if (this.filters.status) params.set('status', this.filters.status);
      if (this.filters.dateFrom) params.set('date_from', this.filters.dateFrom);
      if (this.filters.dateTo) params.set('date_to', this.filters.dateTo);
      params.set('page', 0); // Reset to first page when applying filters
      params.set('per_page', this.perPage);
      params.set('sort_order', this.sortOrder);
      params.set('sort_by', this.sortByColumn);

      // Update URL without reloading
      const newUrl = `${window.location.pathname}?${params.toString()}`;
      window.history.pushState({}, '', newUrl);

      this.closeFilterModal();
      this.loadPosts();
    },

    /**
     * Clears all filters and resets to default state.
     */
    clearFilters() {
      this.filters = {
        search: '',
        status: '',
        dateFrom: '',
        dateTo: ''
      };

      // Reset date range to default
      const today = new Date();
      const lastYear = new Date(today);
      lastYear.setFullYear(today.getFullYear() - 1);
      
      this.filters.dateTo = today.toISOString().split('T')[0];
      this.filters.dateFrom = lastYear.toISOString().split('T')[0];

      this.applyFilters();
    },

    /**
     * Creates a new post with the title from the form.
     * Opens the edit page in a new tab on success.
     */
    async createPost() {
      if (!this.createForm.title) return;

      this.creating = true;
      try {
        const response = await fetch(urlPostCreate, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            title: this.createForm.title
          })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Success',
            text: 'Post created successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.closeCreateModal();
          // Navigate to edit page
          window.open(urlPostUpdate.replace('POST_ID_PLACEHOLDER', data.data.id), '_blank');
          this.loadPosts();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to create post'
          });
        }
      } catch (error) {
        console.error('Error creating post:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to create post'
        });
      } finally {
        this.creating = false;
      }
    }
  }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('blog-posts-app');
  if (el) {
    createApp(BlogPostsApp).mount('#blog-posts-app');
  }
});
