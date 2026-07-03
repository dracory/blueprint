function initTagsApp() {
    if (typeof Vue === 'undefined') {
        setTimeout(initTagsApp, 100);
        return;
    }

    const { createApp } = Vue;

    createApp({
        data() {
            return {
                loading: false,
                tags: [],
                originalTags: [],
                modalTagInput: '',
                showModal: false,
                successMessage: '',
                errorMessage: '',
                productID: productId
            };
        },
        computed: {
            hasChanges() {
                if (this.tags.length !== this.originalTags.length) {
                    return true;
                }
                const sortedTags = [...this.tags].sort();
                const sortedOriginal = [...this.originalTags].sort();
                for (let i = 0; i < sortedTags.length; i++) {
                    if (sortedTags[i] !== sortedOriginal[i]) {
                        return true;
                    }
                }
                return false;
            }
        },
        methods: {
            async loadTags() {
                this.loading = true;
                try {
                    const response = await fetch(urlTagsLoad, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ product_id: this.productID })
                    });

                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const data = await response.json();
                        this.tags = data.tags || [];
                        this.originalTags = [...this.tags];
                    } else {
                        const text = await response.text();
                        console.error('Non-JSON response:', text);
                        throw new Error('Invalid response format: ' + contentType);
                    }
                } catch (error) {
                    console.error('Error loading tags:', error);
                    this.errorMessage = 'Failed to load tags: ' + error.message;
                } finally {
                    this.loading = false;
                }
            },

            async saveTags() {
                this.loading = true;
                try {
                    const response = await fetch(urlTagsSave, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            product_id: this.productID,
                            tags: this.tags
                        })
                    });

                    const data = await response.json();
                    if (data.status === 'success') {
                        this.successMessage = data.message || 'Tags saved successfully';
                        this.errorMessage = '';
                        this.originalTags = [...this.tags];

                        setTimeout(() => {
                            this.successMessage = '';
                        }, 3000);
                    } else {
                        this.errorMessage = data.message || 'Failed to save tags';
                        this.successMessage = '';
                    }
                } catch (error) {
                    console.error('Error saving tags:', error);
                    this.errorMessage = 'Failed to save tags: ' + error.message;
                    this.successMessage = '';
                } finally {
                    this.loading = false;
                }
            },

            addQuickTag(tag) {
                if (!this.tags.includes(tag)) {
                    this.tags.push(tag);
                    this.errorMessage = '';
                    this.showModal = false;
                }
            },

            addModalTag() {
                const tag = this.modalTagInput.trim();
                if (!tag) return;

                if (!this.tags.includes(tag)) {
                    this.tags.push(tag);
                    this.modalTagInput = '';
                    this.errorMessage = '';
                    this.showModal = false;
                }
            },

            removeTag(index) {
                this.tags.splice(index, 1);
            }
        },
        mounted() {
            this.loadTags();
        }
    }).mount('#tags-app');
}

document.addEventListener('DOMContentLoaded', initTagsApp);
