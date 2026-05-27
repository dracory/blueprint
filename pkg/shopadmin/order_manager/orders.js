const { createApp } = Vue;

createApp({
  data() {
    return {
      orders: [],
      loading: false,
      error: null
    };
  },
  mounted() {
    this.loadOrders();
  },
  methods: {
    async loadOrders() {
      this.loading = true;
      this.error = null;
      try {
        const response = await fetch(urlLoadOrders, {
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
        console.log('Orders response:', result);
        if (result.status === 'success') {
          this.orders = result.data.orders;
        } else {
          this.error = result.message || 'Failed to load orders';
        }
      } catch (error) {
        console.error('Failed to load orders:', error);
        this.error = 'Failed to load orders';
      } finally {
        this.loading = false;
      }
    },
    getStatusBadgeClass(status) {
      switch (status) {
        case 'completed':
          return 'bg-success';
        case 'pending':
          return 'bg-warning';
        case 'cancelled':
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
