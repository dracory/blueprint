function initDetailsApp() {
    console.log('initDetailsApp called');
    console.log('Global productId:', typeof productId !== 'undefined' ? productId : 'undefined');
    console.log('Global urlDetailsLoad:', typeof urlDetailsLoad !== 'undefined' ? urlDetailsLoad : 'undefined');
    console.log('Global urlDetailsSave:', typeof urlDetailsSave !== 'undefined' ? urlDetailsSave : 'undefined');
    
    if (typeof Vue === 'undefined') {
        console.log('Vue not defined, waiting...');
        setTimeout(initDetailsApp, 100);
        return;
    }

    console.log('Vue is defined, creating app');
    const { createApp } = Vue;

    createApp({
        data() {
            return {
                loading: false,
                productID: typeof productId !== 'undefined' ? productId : '',
                status: '',
                title: '',
                description: '',
                price: '',
                quantity: '',
                memo: '',
                successMessage: '',
                errorMessage: ''
            };
        },
        watch: {
            description(newVal, oldVal) {
                console.log('Watch triggered - description changed from:', oldVal, 'to:', newVal);
                if (typeof $ !== 'undefined' && typeof $.fn.summernote !== 'undefined') {
                    this.$nextTick(() => {
                        console.log('Calling initSummernote from watch');
                        this.initSummernote();
                    });
                }
            },
            loading(newVal) {
                console.log('Loading changed to:', newVal);
                if (!newVal) {
                    console.log('Loading is false, checking Summernote');
                    if (typeof $ !== 'undefined' && typeof $.fn.summernote !== 'undefined') {
                        console.log('Summernote exists:', $('#product_description').next('.note-editor').length);
                        if ($('#product_description').next('.note-editor').length === 0) {
                            console.log('Summernote gone, re-initializing');
                            this.initSummernote();
                        }
                    }
                }
            }
        },
        methods: {
            async loadDetails() {
                console.log('loadDetails called, productID:', this.productID);
                this.loading = true;
                try {
                    const response = await fetch(urlDetailsLoad, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ product_id: this.productID })
                    });
                    
                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        const data = await response.json();
                        if (data.status === 'success') {
                            this.status = data.data.status || '';
                            this.title = data.data.title || '';
                            this.description = data.data.description || '';
                            this.price = data.data.price || '';
                            this.quantity = data.data.quantity || '';
                            this.memo = data.data.memo || '';
                            
                            // Initialize Summernote after data is loaded
                            this.$nextTick(() => {
                                console.log('Calling initSummernote from loadDetails');
                                this.initSummernote();
                            });
                        } else {
                            this.errorMessage = data.message || 'Failed to load details';
                        }
                    } else {
                        const text = await response.text();
                        console.error('Non-JSON response:', text);
                        throw new Error('Invalid response format: ' + contentType);
                    }
                } catch (error) {
                    console.error('Error loading details:', error);
                    this.errorMessage = 'Failed to load details: ' + error.message;
                } finally {
                    this.loading = false;
                }
            },

            async saveDetails() {
                // Validate form
                if (!this.status) {
                    this.errorMessage = 'Status is required';
                    return;
                }
                if (!this.title) {
                    this.errorMessage = 'Title is required';
                    return;
                }
                if (!this.price) {
                    this.errorMessage = 'Price is required';
                    return;
                }
                if (!this.quantity) {
                    this.errorMessage = 'Quantity is required';
                    return;
                }
                if (isNaN(parseFloat(this.price))) {
                    this.errorMessage = 'Price must be numeric';
                    return;
                }
                if (parseFloat(this.price) < 0) {
                    this.errorMessage = 'Price cannot be negative';
                    return;
                }
                if (isNaN(parseInt(this.quantity))) {
                    this.errorMessage = 'Quantity must be numeric';
                    return;
                }
                if (parseInt(this.quantity) < 0) {
                    this.errorMessage = 'Quantity cannot be negative';
                    return;
                }

                // Sync Summernote content before saving
                if (typeof $ !== 'undefined' && typeof $.fn.summernote !== 'undefined') {
                    this.description = $('#product_description').summernote('code');
                }

                this.loading = true;
                try {
                    const response = await fetch(urlDetailsSave, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            product_id: this.productID,
                            status: this.status,
                            title: this.title,
                            description: this.description,
                            price: this.price,
                            quantity: this.quantity,
                            memo: this.memo
                        })
                    });
                    
                    const data = await response.json();
                    console.log('Save response:', data);
                    if (data.status === 'success') {
                        this.successMessage = data.message || 'Details saved successfully';
                        this.errorMessage = '';
                        console.log('After save success, description:', this.description);
                        console.log('Checking if Summernote exists:', $('#product_description').next('.note-editor').length);
                        
                        setTimeout(() => {
                            this.successMessage = '';
                        }, 3000);
                    } else {
                        this.errorMessage = data.message || 'Failed to save details';
                        this.successMessage = '';
                    }
                } catch (error) {
                    console.error('Error saving details:', error);
                    this.errorMessage = 'Failed to save details: ' + error.message;
                    this.successMessage = '';
                } finally {
                    console.log('Save finally block, loading:', this.loading);
                    this.loading = false;
                    console.log('After loading set to false');
                    this.initSummernote();
                }
            },

            initSummernote() {
                console.log('initSummernote called, description:', this.description);
                console.log('jQuery defined:', typeof $ !== 'undefined');
                console.log('Summernote defined:', typeof $.fn.summernote !== 'undefined');
                console.log('Vue ref container:', this.$refs.descriptionContainer);
                
                if (typeof $ !== 'undefined' && typeof $.fn.summernote !== 'undefined') {
                    // Use Vue ref to get container
                    const container = this.$refs.descriptionContainer;
                    
                    if (!container) {
                        console.log('Container ref not found');
                        return;
                    }
                    
                    console.log('Container found via ref');
                    
                    // Check if textarea already exists in container
                    const $container = $(container);
                    const $existingTextarea = $container.find('#product_description');
                    
                    console.log('Existing textarea:', $existingTextarea.length);
                    
                    if ($existingTextarea.length === 0) {
                        // Create textarea and container dynamically
                        console.log('Creating textarea dynamically');
                        const textareaHtml = '<div id="description-editor-container"><textarea id="product_description" class="form-control" rows="10" placeholder="Enter product description..."></textarea></div>';
                        $container.html(textareaHtml);
                    }
                    
                    const $el = $('#product_description');
                    console.log('Element found:', $el.length);
                    console.log('Existing editor:', $el.next('.note-editor').length);
                    
                    // Only initialize if not already initialized
                    if ($el.next('.note-editor').length === 0) {
                        console.log('Initializing new Summernote instance');
                        $el.summernote({
                            height: 300,
                            placeholder: 'Enter product description...',
                            callbacks: {
                                onChange: (contents, $editable) => {
                                    console.log('Summernote onChange triggered');
                                    this.description = contents;
                                }
                            }
                        });
                        
                        console.log('Summernote initialized, setting content');
                        // Set content after initialization
                        if (this.description) {
                            console.log('Setting initial content');
                            $el.summernote('code', this.description);
                        }
                        console.log('Content set, checking if editor exists:', $el.next('.note-editor').length);
                    } else {
                        // Already initialized, just update the code
                        console.log('Summernote already initialized, updating code');
                        if (this.description) {
                            $el.summernote('code', this.description);
                        }
                    }
                } else {
                    console.log('Waiting for Summernote to load...');
                    setTimeout(() => this.initSummernote(), 200);
                }
            },
        },
        mounted() {
            this.loadDetails();
            // Expose app instance globally for Summernote callback
            window.detailsApp = this;
        }
    }).mount('#details-app');
}

// Initialize the details app when DOM is ready
console.log('DOMContentLoaded listener added');
document.addEventListener('DOMContentLoaded', initDetailsApp);
