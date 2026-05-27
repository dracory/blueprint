const { createApp } = Vue;

createApp({
  data() {
    return {
      urlUpdateProduct: urlUpdateProduct,
      products: [],
      selectedProducts: [],
      selectAll: false,
      loading: false,
      error: null
    };
  },
  mounted() {
    this.loadProducts();
  },
  methods: {
    async loadProducts() {
      this.loading = true;
      this.error = null;
      try {
        const response = await fetch(urlLoadProducts, {
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
        console.log('Products response:', result);
        if (result.status === 'success') {
          this.products = result.data.products;
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
    async deleteProduct(productId) {
      if (!confirm('Are you sure you want to delete this product?')) {
        return;
      }
      try {
        const response = await fetch(urlProductDelete, {
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
          Swal.fire('Success', 'Product deleted successfully', 'success');
        } else {
          Swal.fire('Error', result.message || 'Failed to delete product', 'error');
        }
      } catch (error) {
        console.error('Failed to delete product:', error);
        Swal.fire('Error', 'Failed to delete product', 'error');
      }
    },
    async deleteSelectedProducts() {
      if (!confirm(`Are you sure you want to delete ${this.selectedProducts.length} product(s)?`)) {
        return;
      }
      try {
        const response = await fetch(urlProductDeleteSelected, {
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
          Swal.fire('Success', 'Products deleted successfully', 'success');
        } else {
          Swal.fire('Error', result.message || 'Failed to delete products', 'error');
        }
      } catch (error) {
        console.error('Failed to delete products:', error);
        Swal.fire('Error', 'Failed to delete products', 'error');
      }
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
    }
  }
}).mount('#app');
