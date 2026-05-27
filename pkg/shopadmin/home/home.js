const { createApp } = Vue;

createApp({
  data() {
    return {
      stats: {
        product_count: 0,
        category_count: 0,
        order_count: 0
      },
      urlProducts: urlProducts,
      urlCategories: urlCategories,
      urlDiscounts: urlDiscounts,
      urlOrders: urlOrders
    };
  },
  mounted() {
    this.loadStats();
  },
  methods: {
    async loadStats() {
      try {
        const response = await fetch(urlLoadStats);
        const result = await response.json();
        if (result.status === 'success') {
          this.stats = result.data;
        }
      } catch (error) {
        console.error('Failed to load stats:', error);
      }
    }
  }
}).mount('#app');
