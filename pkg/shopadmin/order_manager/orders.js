const { createApp } = Vue;

const OrdersApp = {
    data() {
        return {
            // UI state
            loading: true,
            showFilterModal: false,

            // Order data
            orders: [],
            totalOrders: 0,

            // Pagination
            currentPage: 0,
            perPage: 10,
            jumpToPage: 1,

            // Filters
            filters: {
                status: '',
                customer_name: '',
                customer_email: '',
                order_id: '',
                created_from: '',
                created_to: ''
            },

            // Sorting
            sortByColumn: 'created_at',
            sortOrder: 'desc'
        };
    },

    computed: {
        /**
         * Returns total number of pages.
         */
        totalPages() {
            return Math.ceil(this.totalOrders / this.perPage);
        },

        /**
         * Returns a human-readable string describing the current filter state.
         */
        filterStatus() {
            const parts = [];
            if (this.filters.status) parts.push(`status: ${this.filters.status}`);
            if (this.filters.customer_name) parts.push(`customer name: "${this.filters.customer_name}"`);
            if (this.filters.customer_email) parts.push(`customer email: "${this.filters.customer_email}"`);
            if (this.filters.order_id) parts.push(`order id: "${this.filters.order_id}"`);
            if (this.filters.created_from) parts.push(`from: ${this.filters.created_from}`);
            if (this.filters.created_to) parts.push(`to: ${this.filters.created_to}`);
            
            if (parts.length === 0) return 'Showing all orders';
            return 'Showing orders with ' + parts.join(', ');
        },

        /**
         * Returns true if any filters are active.
         */
        hasActiveFilters() {
            return this.filters.status !== '' ||
                   this.filters.customer_name !== '' ||
                   this.filters.customer_email !== '' ||
                   this.filters.order_id !== '' ||
                   this.filters.created_from !== '' ||
                   this.filters.created_to !== '';
        }
    },

    mounted() {
        // Read filters from URL parameters for shareable URLs
        const urlParams = new URLSearchParams(window.location.search);
        
        this.filters.status = urlParams.get('status') || '';
        this.filters.customer_name = urlParams.get('customer_name') || '';
        this.filters.customer_email = urlParams.get('customer_email') || '';
        this.filters.order_id = urlParams.get('order_id') || '';
        this.filters.created_from = urlParams.get('created_from') || '';
        this.filters.created_to = urlParams.get('created_to') || '';
        this.sortByColumn = urlParams.get('sort_by') || 'created_at';
        this.sortOrder = urlParams.get('sort_order') || 'desc';
        this.currentPage = parseInt(urlParams.get('page') || '0', 10);
        this.perPage = parseInt(urlParams.get('per_page') || '10', 10);

        this.loadOrders();
    },
    methods: {
        /**
         * Loads orders from the server based on current filters, pagination, and sorting.
         */
        async loadOrders() {
            this.loading = true;
            try {
                const response = await fetch(urlLoadOrders, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        page: this.currentPage,
                        per_page: this.perPage,
                        status: this.filters.status,
                        customer_name: this.filters.customer_name,
                        customer_email: this.filters.customer_email,
                        order_id: this.filters.order_id,
                        created_from: this.filters.created_from,
                        created_to: this.filters.created_to,
                        sort_by: this.sortByColumn,
                        sort_order: this.sortOrder
                    })
                });
                const data = await response.json();

                if (data.status === 'success') {
                    this.orders = data.data?.orders || [];
                    this.totalOrders = data.data?.total || 0;
                    // Sync jumpToPage with currentPage to ensure UI consistency
                    this.jumpToPage = this.currentPage + 1;
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to load orders'
                    });
                }
            } catch (error) {
                console.error('Error loading orders:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: error.message || 'Failed to load orders'
                });
            } finally {
                this.loading = false;
            }
        },
        /**
         * Sorts the orders by the specified column.
         */
        sortBy(column) {
            if (this.sortByColumn === column) {
                this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
            } else {
                this.sortByColumn = column;
                this.sortOrder = 'asc';
            }
            
            // For computed fields (customer_name, customer_email), sort client-side
            if (column === 'customer_name' || column === 'customer_email') {
                this.sortOrdersClientSide();
            } else {
                this.currentPage = 0;
                this.applyFilters();
            }
        },

        /**
         * Sorts orders client-side for computed fields.
         */
        sortOrdersClientSide() {
            const column = this.sortByColumn;
            const order = this.sortOrder === 'asc' ? 1 : -1;
            
            this.orders.sort((a, b) => {
                const aVal = (a[column] || '').toLowerCase();
                const bVal = (b[column] || '').toLowerCase();
                
                if (aVal < bVal) return -1 * order;
                if (aVal > bVal) return 1 * order;
                return 0;
            });
        },

        /**
         * Navigates to the specified page number.
         */
        goToPage(page) {
            if (page < 0) return;
            const maxPage = Math.ceil(this.totalOrders / this.perPage) - 1;
            if (page > maxPage) page = maxPage;
            this.currentPage = page;
            this.jumpToPage = page + 1;
            this.applyFilters();
        },

        /**
         * Changes the per-page setting and reloads orders.
         */
        changePerPage() {
            this.perPage = parseInt(this.perPage, 10);
            this.currentPage = 0;
            this.jumpToPage = 1;
            this.applyFilters();
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
         */
        applyFilters() {
            const params = new URLSearchParams();
            if (this.filters.status) params.set('status', this.filters.status);
            if (this.filters.customer_name) params.set('customer_name', this.filters.customer_name);
            if (this.filters.customer_email) params.set('customer_email', this.filters.customer_email);
            if (this.filters.order_id) params.set('order_id', this.filters.order_id);
            if (this.filters.created_from) params.set('created_from', this.filters.created_from);
            if (this.filters.created_to) params.set('created_to', this.filters.created_to);
            params.set('page', this.currentPage);
            params.set('per_page', this.perPage);
            params.set('sort_order', this.sortOrder);
            params.set('sort_by', this.sortByColumn);

            const newUrl = `${window.location.pathname}?${params.toString()}`;
            window.history.pushState({}, '', newUrl);

            this.closeFilterModal();
            this.loadOrders();
        },

        /**
         * Clears all filters and resets to default state.
         */
        clearFilters() {
            this.filters = {
                status: '',
                customer_name: '',
                customer_email: '',
                order_id: '',
                created_from: '',
                created_to: ''
            };
            this.currentPage = 0;
            this.applyFilters();
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
            
            const day = date.getDate().toString().padStart(2, '0');
            const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
            const month = months[date.getMonth()];
            const year = date.getFullYear();
            
            const hours = date.getHours().toString().padStart(2, '0');
            const minutes = date.getMinutes().toString().padStart(2, '0');
            
            return `${day} ${month} ${year}<br><span class="small text-muted">${hours}:${minutes}</span>`;
        }
    }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    const el = document.getElementById('app');
    if (el) {
        createApp(OrdersApp).mount('#app');
    }
});
