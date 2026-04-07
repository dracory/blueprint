// Details Component Vue.js Application
const { createApp } = Vue;

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', function() {
	// Ensure the details-app element exists
	const detailsApp = document.getElementById('details-app');
	if (!detailsApp) {
		console.error('Details app container not found');
		return;
	}

	if (!urlDetailsLoad || !urlDetailsSave || !productId) {
		console.error('urlDetailsLoad, urlDetailsSave, or productId not found');
		return;
	}

	createApp({
		components: {
			QuillEditor
		},
		data() {
			return {
				productId: productId,
				urlDetailsLoad: urlDetailsLoad,
				urlDetailsSave: urlDetailsSave,
				details: {
					id: '',
					title: '',
					description: '',
					price: '',
					quantity: '',
					status: 'draft'
				},
				errorMessage: '',
				successMessage: '',
				loading: false
			}
		},
		methods: {
			/**
			 * Loads product details from the server via AJAX POST request
			 * @async
			 * @returns {Promise<void>}
			 */
			async loadDetails() {
				this.loading = true;
				try {
					console.log('Loading details from:', this.urlDetailsLoad);
					
					const response = await fetch(this.urlDetailsLoad, {
						method: 'POST',
						headers: {
							'Accept': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						}
					});
					
					const result = await response.json();
					if (result.status !== 'success') {
						throw new Error(result.message || 'Failed to load details');
					}
					if (result.data?.details) {
						this.details = result.data.details;
					}
				} catch (error) {
					console.error('Error loading details:', error);
					this.errorMessage = 'Failed to load details: ' + error.message;
				} finally {
					this.loading = false;
				}
			},
			
			/**
			 * Saves product details to the server via AJAX POST request
			 * @async
			 * @returns {Promise<void>}
			 */
			async saveDetails() {
				this.loading = true;
				try {
					console.log('Saving details to:', this.urlDetailsSave);
					
					const response = await fetch(this.urlDetailsSave, {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
							'Accept': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						},
						body: JSON.stringify({
							details: this.details
						})
					});
					
					const result = await response.json();
					if (result.status !== 'success') {
						throw new Error(result.message || 'Failed to save details');
					}
					this.successMessage = result.message || 'Product details saved successfully';
					this.errorMessage = '';
					
					setTimeout(() => {
						this.successMessage = '';
					}, 3000);
					
				} catch (error) {
					console.error('Error saving details:', error);
					this.errorMessage = 'Failed to save details: ' + error.message;
					this.successMessage = '';
				} finally {
					this.loading = false;
				}
			}
		},
		
		async mounted() {
			await this.loadDetails();
			
			// Listen for save events from external buttons
			document.getElementById('details-save-btn-top')?.addEventListener('click', () => {
				this.saveDetails();
			});
			document.getElementById('details-save-btn-bottom')?.addEventListener('click', () => {
				this.saveDetails();
			});
		}
	}).mount('#details-app');
});
