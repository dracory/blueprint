const { createApp } = Vue;

createApp({
  data() {
    return {
      categories: [],
      selectedCategories: [],
      selectAll: false,
      loading: false,
      error: null
    };
  },
  mounted() {
    this.loadCategories();
  },
  methods: {
    async loadCategories() {
      this.loading = true;
      this.error = null;
      try {
        const response = await fetch(urlLoadCategories, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            page: 0,
            per_page: 10
          })
        });
        const result = await response.json();
        console.log('Categories response:', result);
        if (result.status === 'success') {
          this.categories = result.data.categories;
        } else {
          this.error = result.message || 'Failed to load categories';
        }
      } catch (error) {
        console.error('Failed to load categories:', error);
        this.error = 'Failed to load categories';
      } finally {
        this.loading = false;
      }
    },
    async deleteCategory(categoryId) {
      if (!confirm('Are you sure you want to delete this category?')) {
        return;
      }
      try {
        const response = await fetch(urlCategoryDelete, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            category_id: categoryId
          })
        });
        const result = await response.json();
        if (result.status === 'success') {
          this.categories = this.categories.filter(c => c.id !== categoryId);
          Swal.fire('Success', 'Category deleted successfully', 'success');
        } else {
          Swal.fire('Error', result.message || 'Failed to delete category', 'error');
        }
      } catch (error) {
        console.error('Failed to delete category:', error);
        Swal.fire('Error', 'Failed to delete category', 'error');
      }
    },
    async deleteSelectedCategories() {
      if (!confirm(`Are you sure you want to delete ${this.selectedCategories.length} category(ies)?`)) {
        return;
      }
      try {
        const response = await fetch(urlCategoryDeleteSelected, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            bulk_category_ids: this.selectedCategories
          })
        });
        const result = await response.json();
        if (result.status === 'success') {
          this.categories = this.categories.filter(c => !this.selectedCategories.includes(c.id));
          this.selectedCategories = [];
          this.selectAll = false;
          Swal.fire('Success', 'Categories deleted successfully', 'success');
        } else {
          Swal.fire('Error', result.message || 'Failed to delete categories', 'error');
        }
      } catch (error) {
        console.error('Failed to delete categories:', error);
        Swal.fire('Error', 'Failed to delete categories', 'error');
      }
    },
    toggleSelectAll() {
      if (this.selectAll) {
        this.selectedCategories = this.categories.map(c => c.id);
      } else {
        this.selectedCategories = [];
      }
    },
    getStatusBadgeClass(status) {
      switch (status) {
        case 'active':
          return 'bg-success';
        case 'inactive':
          return 'bg-secondary';
        case 'draft':
          return 'bg-warning';
        default:
          return 'bg-secondary';
      }
    },
    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      return date.toLocaleDateString();
    }
  }
}).mount('#app');
