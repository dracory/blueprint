function initMetadataApp() {
    if (typeof Vue === 'undefined') {
        setTimeout(initMetadataApp, 100);
        return;
    }

    const { createApp } = Vue;

    createApp({
        data() {
            return {
                loading: false,
                metadataItems: [],
                originalMetadataItems: [],
                newKey: '',
                newValue: '',
                showModal: false,
                successMessage: '',
                errorMessage: '',
                productID: productId
            };
        },
        computed: {
            hasChanges() {
                if (this.metadataItems.length !== this.originalMetadataItems.length) {
                    return true;
                }
                for (let i = 0; i < this.metadataItems.length; i++) {
                    if (this.metadataItems[i].key !== this.originalMetadataItems[i].key ||
                        this.metadataItems[i].value !== this.originalMetadataItems[i].value) {
                        return true;
                    }
                }
                return false;
            }
        },
        methods: {
            async loadMetadata() {
                this.loading = true;
                try {
                    const response = await fetch(urlMetadataLoad, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ product_id: this.productID })
                    });
                    
                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const data = await response.json();
                        this.metadataItems = data.metadata || [];
                        // Sort by key name
                        this.metadataItems.sort((a, b) => a.key.localeCompare(b.key));
                        // Deep copy to track changes
                        this.originalMetadataItems = this.metadataItems.map(item => ({
                            id: item.id,
                            key: item.key,
                            value: item.value
                        }));
                    } else {
                        const text = await response.text();
                        console.error('Non-JSON response:', text);
                        throw new Error('Invalid response format: ' + contentType);
                    }
                } catch (error) {
                    console.error('Error loading metadata:', error);
                    this.errorMessage = 'Failed to load metadata: ' + error.message;
                } finally {
                    this.loading = false;
                }
            },

            async saveMetadata() {
                this.loading = true;
                try {
                    const response = await fetch(urlMetadataSave, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            product_id: this.productID,
                            metadata: this.metadataItems
                        })
                    });
                    
                    const data = await response.json();
                    if (data.status === 'success') {
                        this.successMessage = data.message || 'Metadata saved successfully';
                        this.errorMessage = '';
                        // Deep copy to track changes
                        this.originalMetadataItems = this.metadataItems.map(item => ({
                            id: item.id,
                            key: item.key,
                            value: item.value
                        }));

                        setTimeout(() => {
                            this.successMessage = '';
                        }, 3000);
                    } else {
                        this.errorMessage = data.message || 'Failed to save metadata';
                        this.successMessage = '';
                    }
                } catch (error) {
                    console.error('Error saving metadata:', error);
                    this.errorMessage = 'Failed to save metadata: ' + error.message;
                    this.successMessage = '';
                } finally {
                    this.loading = false;
                }
            },

            addItem() {
                const key = this.newKey.trim();
                const value = this.newValue.trim();

                if (!key) {
                    this.errorMessage = 'Please enter a key for the metadata';
                    return;
                }

                // Check for duplicate keys
                if (this.metadataItems.some(item => item.key === key)) {
                    this.errorMessage = 'A metadata item with this key already exists';
                    return;
                }

                this.metadataItems.push({
                    id: Date.now().toString(),
                    key: key,
                    value: value
                });

                // Sort by key name
                this.metadataItems.sort((a, b) => a.key.localeCompare(b.key));

                // Clear input fields and close modal
                this.newKey = '';
                this.newValue = '';
                this.errorMessage = '';
                this.showModal = false;
            },

            removeItem(index) {
                this.metadataItems.splice(index, 1);
            }
        },
        mounted() {
            this.loadMetadata();
        }
    }).mount('#metadata-app');
}

// Initialize the metadata app when DOM is ready
document.addEventListener('DOMContentLoaded', initMetadataApp);
