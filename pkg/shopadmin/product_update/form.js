const { createApp } = Vue;

createApp({
    data() {
        return {
            loading: true,
            saving: false,
            action: '',
            redirectTo: '',
            productId: '',
            returnUrl: '',
            form: {
                title: '',
                description: '',
                status: '',
                price: ''
            }
        };
    },
    mounted() {
        this.productId = productID;
        this.returnUrl = '/admin/shop?controller=products';
        this.loadProduct();
    },
    methods: {
        async loadProduct() {
            this.loading = true;
            this.errorMessage = '';
            try {
                const response = await fetch(urlLoadProduct, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        product_id: this.productId
                    })
                });
                const result = await response.json();

                if (result.status === 'success') {
                    const d = result.data.product;
                    this.form.title = d.title || '';
                    this.form.description = d.description || '';
                    this.form.status = d.status || '';
                    this.form.price = d.price || '';
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to load product', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                }
            } catch (err) {
                console.error('Error loading product:', err);
                Notiflix.Notify.failure('Failed to load product', {
                    position: 'right-top',
                    timeout: 3000,
                });
            } finally {
                this.loading = false;
            }
        },
        async save(actionType) {
            this.action = actionType;
            this.redirectTo = '';

            if (!this.form.title.trim()) {
                Notiflix.Notify.failure('Title is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            if (!this.form.status) {
                Notiflix.Notify.failure('Status is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }
            if (!this.form.price) {
                Notiflix.Notify.failure('Price is required', {
                    position: 'right-top',
                    timeout: 3000,
                });
                return;
            }

            this.saving = true;
            try {
                const payload = {
                    product_id: this.productId,
                    title: this.form.title.trim(),
                    description: this.form.description.trim(),
                    status: this.form.status,
                    price: parseFloat(this.form.price)
                };

                const response = await fetch(urlUpdateProduct, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });
                const result = await response.json();

                if (result.status === 'success') {
                    if (actionType === 'save') {
                        this.redirectTo = this.returnUrl;
                        Notiflix.Notify.success('Product saved successfully', {
                            position: 'right-top',
                            timeout: 3000,
                        });
                        setTimeout(() => {
                            if (this.redirectTo) {
                                window.location.href = this.redirectTo;
                            }
                        }, 3000);
                    } else {
                        Notiflix.Notify.success('Product saved successfully', {
                            position: 'right-top',
                            timeout: 3000,
                        });
                        await this.loadProduct();
                    }
                } else {
                    Notiflix.Notify.failure(result.message || 'Failed to save product', {
                        position: 'right-top',
                        timeout: 3000,
                    });
                }
            } catch (err) {
                console.error('Error saving product:', err);
                Notiflix.Notify.failure('Failed to save product', {
                    position: 'right-top',
                    timeout: 3000,
                });
            } finally {
                this.saving = false;
            }
        }
    }
}).mount('#app-product-update');
