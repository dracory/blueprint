const { createApp } = Vue;

createApp({
  data() {
    return {
      discounts: [],
      selectedDiscounts: [],
      selectAll: false,
      loading: false,
      error: null
    };
  },
  mounted() {
    this.loadDiscounts();
  },
  methods: {
    async loadDiscounts() {
      this.loading = true;
      this.error = null;
      try {
        const response = await fetch(urlLoadDiscounts, {
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
        console.log('Discounts response:', result);
        if (result.status === 'success') {
          this.discounts = result.data.discounts;
        } else {
          this.error = result.message || 'Failed to load discounts';
        }
      } catch (error) {
        console.error('Failed to load discounts:', error);
        this.error = 'Failed to load discounts';
      } finally {
        this.loading = false;
      }
    },
    async deleteDiscount(discountId) {
      if (!confirm('Are you sure you want to delete this discount?')) {
        return;
      }
      try {
        const response = await fetch(urlDiscountDelete, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            discount_id: discountId
          })
        });
        const result = await response.json();
        if (result.status === 'success') {
          this.discounts = this.discounts.filter(d => d.id !== discountId);
          Swal.fire('Success', 'Discount deleted successfully', 'success');
        } else {
          Swal.fire('Error', result.message || 'Failed to delete discount', 'error');
        }
      } catch (error) {
        console.error('Failed to delete discount:', error);
        Swal.fire('Error', 'Failed to delete discount', 'error');
      }
    },
    async deleteSelectedDiscounts() {
      if (!confirm(`Are you sure you want to delete ${this.selectedDiscounts.length} discount(s)?`)) {
        return;
      }
      try {
        const response = await fetch(urlDiscountDeleteSelected, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            bulk_discount_ids: this.selectedDiscounts
          })
        });
        const result = await response.json();
        if (result.status === 'success') {
          this.discounts = this.discounts.filter(d => !this.selectedDiscounts.includes(d.id));
          this.selectedDiscounts = [];
          this.selectAll = false;
          Swal.fire('Success', 'Discounts deleted successfully', 'success');
        } else {
          Swal.fire('Error', result.message || 'Failed to delete discounts', 'error');
        }
      } catch (error) {
        console.error('Failed to delete discounts:', error);
        Swal.fire('Error', 'Failed to delete discounts', 'error');
      }
    },
    toggleSelectAll() {
      if (this.selectAll) {
        this.selectedDiscounts = this.discounts.map(d => d.id);
      } else {
        this.selectedDiscounts = [];
      }
    },
    getStatusBadgeClass(status) {
      switch (status) {
        case 'active':
          return 'bg-success';
        case 'inactive':
          return 'bg-secondary';
        case 'expired':
          return 'bg-danger';
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
