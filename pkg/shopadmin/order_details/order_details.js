const { createApp } = Vue;

const OrderDetailsApp = {
    data() {
        return {
            // UI state
            loading: true,

            // Order data
            order: null
        };
    },
    mounted() {
        this.loadOrderDetails();
    },
    methods: {
        async loadOrderDetails() {
            this.loading = true;
            try {
                const response = await fetch(urlLoadOrderDetails, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        order_id: ORDER_ID
                    })
                });
                const data = await response.json();

                if (data.status === 'success') {
                    this.order = data.data?.order || null;
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to load order details'
                    });
                }
            } catch (error) {
                console.error('Error loading order details:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.message || 'Failed to load order details'
                });
            } finally {
                this.loading = false;
            }
        },
        formatDate(dateString) {
            if (!dateString) return '-';
            const date = new Date(dateString);
            const day = date.getDate().toString().padStart(2, '0');
            const month = date.toLocaleString('en-GB', { month: 'short' });
            const year = date.getFullYear();
            const hours = date.getHours().toString().padStart(2, '0');
            const minutes = date.getMinutes().toString().padStart(2, '0');
            return `${day} ${month} ${year}<br><small class="text-muted">${hours}:${minutes}</small>`;
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
        }
    }
};

createApp(OrderDetailsApp).mount('#app');
