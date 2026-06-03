function initProductManager() {
  const { createApp } = Vue;

  createApp({
    data() {
      return {
        urls: window.productManagerUrls || {},
        products: [],
        selectedProducts: [],
        selectAll: false,
        loading: false,
        error: null,
        showFilterModal: false,
        showCreateModal: false,
        creating: false,
        newProduct: {
          title: ''
        },
        currentPage: 0,
        perPage: 10,
        totalProducts: 0,
        filters: {
          status: '',
          created_from: '',
          created_to: ''
        },
        sortByColumn: 'created_at',
        sortOrder: 'desc'
      };
    },
  computed: {
    totalPages() {
      return Math.ceil(this.totalProducts / this.perPage);
    },
    filterStatus() {
      const parts = [];
      if (this.filters.status) parts.push(`status: ${this.filters.status}`);
      if (this.filters.created_from) parts.push(`from: ${this.filters.created_from}`);
      if (this.filters.created_to) parts.push(`to: ${this.filters.created_to}`);
      
      if (parts.length === 0) return 'Showing all products';
      return 'Showing products with ' + parts.join(', ');
    },
    hasActiveFilters() {
      return this.filters.status !== '' ||
             this.filters.created_from !== '' ||
             this.filters.created_to !== '';
    }
  },
  mounted() {
    const urlParams = new URLSearchParams(window.location.search);
    this.filters.status = urlParams.get('status') || '';
    this.filters.created_from = urlParams.get('created_from') || '';
    this.filters.created_to = urlParams.get('created_to') || '';
    this.sortByColumn = urlParams.get('sort_by') || 'created_at';
    this.sortOrder = urlParams.get('sort_order') || 'desc';
    this.currentPage = parseInt(urlParams.get('page') || '0', 10);
    this.perPage = parseInt(urlParams.get('per_page') || '10', 10);
    this.loadProducts();
  },
  methods: {
    async loadProducts() {
      this.loading = true;
      this.error = null;
      try {
        const response = await fetch(this.urls.loadProducts, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            page: this.currentPage,
            per_page: this.perPage,
            status: this.filters.status,
            created_from: this.filters.created_from,
            created_to: this.filters.created_to,
            sort_by: this.sortByColumn,
            sort: this.sortOrder
          })
        });
        const result = await response.json();
        console.log('Products response:', result);
        if (result.status === 'success') {
          this.products = result.data.products;
          this.totalProducts = result.data.total || 0;
        } else {
          this.error = result.message || 'Failed to load products';
        }
      } catch (error) {
        console.error('Failed to load products:', error);
        this.error = 'Failed to load products';
      } finally {
        this.loading = false;
      }
    },
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
    goToPage(page) {
      if (page < 0) return;
      const maxPage = Math.ceil(this.totalProducts / this.perPage) - 1;
      if (page > maxPage) page = maxPage;
      this.currentPage = page;
      this.applyFilters();
    },
    changePerPage() {
      this.perPage = parseInt(this.perPage, 10);
      this.currentPage = 0;
      this.applyFilters();
    },
    openFilterModal() {
      this.showFilterModal = true;
    },
    closeFilterModal() {
      this.showFilterModal = false;
    },
    applyFilters() {
      const params = new URLSearchParams();
      if (this.filters.status) params.set('status', this.filters.status);
      if (this.filters.created_from) params.set('created_from', this.filters.created_from);
      if (this.filters.created_to) params.set('created_to', this.filters.created_to);
      params.set('page', this.currentPage);
      params.set('per_page', this.perPage);
      params.set('sort_order', this.sortOrder);
      params.set('sort_by', this.sortByColumn);

      const newUrl = `${window.location.pathname}?${params.toString()}`;
      window.history.pushState({}, '', newUrl);

      this.closeFilterModal();
      this.loadProducts();
    },
    clearFilters() {
      this.filters = {
        status: '',
        created_from: '',
        created_to: ''
      };
      this.currentPage = 0;
      this.applyFilters();
    },
    sortIcon(column) {
      if (this.sortByColumn !== column) return 'bi bi-arrow-down-up text-muted';
      return this.sortOrder === 'asc' ? 'bi bi-arrow-up' : 'bi bi-arrow-down';
    },
    async deleteProduct(productId) {
      Notiflix.Confirm.show(
        'Delete Product',
        'Are you sure you want to delete this product?',
        'Yes, delete it',
        'Cancel',
        async () => {
          try {
            const response = await fetch(this.urls.productDelete, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                product_id: productId
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.products = this.products.filter(p => p.id !== productId);
              Notiflix.Notify.success('Product deleted successfully', {
                position: 'right-top',
                timeout: 3000,
              });
            } else {
              Notiflix.Notify.failure(result.message || 'Failed to delete product', {
                position: 'right-top',
                timeout: 3000,
              });
            }
          } catch (error) {
            console.error('Failed to delete product:', error);
            Notiflix.Notify.failure('Failed to delete product', {
              position: 'right-top',
              timeout: 3000,
            });
          }
        },
        () => {
          // Cancel callback
        }
      );
    },
    async deleteSelectedProducts() {
      Notiflix.Confirm.show(
        'Delete Products',
        `Are you sure you want to delete ${this.selectedProducts.length} product(s)?`,
        'Yes, delete them',
        'Cancel',
        async () => {
          try {
            const response = await fetch(this.urls.productDeleteSelected, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({
                bulk_product_ids: this.selectedProducts
              })
            });
            const result = await response.json();
            if (result.status === 'success') {
              this.products = this.products.filter(p => !this.selectedProducts.includes(p.id));
              this.selectedProducts = [];
              this.selectAll = false;
              Notiflix.Notify.success('Products deleted successfully', {
                position: 'right-top',
                timeout: 3000,
              });
            } else {
              Notiflix.Notify.failure(result.message || 'Failed to delete products', {
                position: 'right-top',
                timeout: 3000,
              });
            }
          } catch (error) {
            console.error('Failed to delete products:', error);
            Notiflix.Notify.failure('Failed to delete products', {
              position: 'right-top',
              timeout: 3000,
            });
          }
        },
        () => {
          // Cancel callback
        }
      );
    },
    toggleSelectAll() {
      if (this.selectAll) {
        this.selectedProducts = this.products.map(p => p.id);
      } else {
        this.selectedProducts = [];
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
    },
    async createProduct() {
      if (!this.newProduct.title.trim()) {
        Notiflix.Notify.failure('Title is required', {
          position: 'right-top',
          timeout: 3000,
        });
        return;
      }

      this.creating = true;
      try {
        const response = await fetch(this.urls.createProduct, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            title: this.newProduct.title.trim()
          })
        });
        const result = await response.json();

        if (result.status === 'success') {
          Notiflix.Notify.success('Product created successfully', {
            position: 'right-top',
            timeout: 3000,
          });
          this.showCreateModal = false;
          this.newProduct.title = '';
          this.loadProducts();
        } else {
          Notiflix.Notify.failure(result.message || 'Failed to create product', {
            position: 'right-top',
            timeout: 3000,
          });
        }
      } catch (error) {
        console.error('Failed to create product:', error);
        Notiflix.Notify.failure('Failed to create product', {
          position: 'right-top',
          timeout: 3000,
        });
      } finally {
        this.creating = false;
      }
    }
  }
}).mount('#app');
}

document.addEventListener('DOMContentLoaded', function() {
  if (typeof Vue !== 'undefined') {
    initProductManager();
  } else {
    window.addEventListener('load', function() {
      if (typeof Vue !== 'undefined') {
        initProductManager();
      }
    });
  }
});
