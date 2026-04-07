// Metadata Component Vue.js Application
const { createApp } = Vue;

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', function() {
	// Ensure the metadata-app element exists
	const metadataApp = document.getElementById('metadata-app');
	if (!metadataApp) {
		console.error('Metadata app container not found');
		return;
	}

	if (!urlMetasLoad || !urlMetasSave || !productId) {
		console.error('urlMetasLoad, urlMetasSave, or productId not found');
		return;
	}

	createApp({
		data() {
			return {
				productId: productId,
				urlMetasLoad: urlMetasLoad,
				urlMetasSave: urlMetasSave,
				metadataItems: [],
				newKey: '',
				newValue: '',
				errorMessage: '',
				successMessage: '',
				loading: false
			}
		},
		methods: {
			/**
			 * Loads metadata from the server via AJAX POST request
			 * Populates metadataItems array with the response data
			 * Shows error message if loading fails
			 * @async
			 * @returns {Promise<void>}
			 */
			async loadMetadata() {
				this.loading = true;
				try {
					console.log('Loading metadata from:', this.urlMetasLoad);
					console.log('Product ID:', this.productId);
					
					const response = await fetch(this.urlMetasLoad, {
						method: 'POST',
						headers: {
							'Accept': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						}
					});
					
					const result = await response.json();
					console.log('Response data:', result);
					if (result.status !== 'success') {
						throw new Error(result.message || 'Failed to load metadata');
					}
					this.metadataItems = result.data?.metadata || [];
				} catch (error) {
					console.error('Error loading metadata:', error);
					this.errorMessage = 'Failed to load metadata: ' + error.message;
				} finally {
					this.loading = false;
				}
			},
			
			/**
			 * Saves metadata to the server via AJAX POST request
			 * Sends current metadataItems array as JSON payload
			 * Reloads metadata after successful save and shows success message
			 * @async
			 * @returns {Promise<void>}
			 */
			async saveMetadata() {
				this.loading = true;
				try {
					console.log('Saving metadata to:', this.urlMetasSave);
					console.log('Product ID:', this.productId);
					console.log('Metadata to save:', this.metadataItems);
					
					const response = await fetch(this.urlMetasSave, {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
							'Accept': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						},
						body: JSON.stringify({
							metadata: this.metadataItems
						})
					});
					
					const result = await response.json();
					console.log('Save response data:', result);
					if (result.status !== 'success') {
						throw new Error(result.message || 'Failed to save metadata');
					}
					this.successMessage = result.message || 'Metadata saved successfully';
					this.errorMessage = '';
					
					// Reload metadata to show updated data
					await this.loadMetadata();
					
					// Clear success message after 3 seconds
					setTimeout(() => {
						this.successMessage = '';
					}, 3000);
					
				} catch (error) {
					console.error('Error saving metadata:', error);
					this.errorMessage = 'Failed to save metadata: ' + error.message;
					this.successMessage = '';
				} finally {
					this.loading = false;
				}
			},
			
			/**
			 * Adds a new metadata item to the local array
			 * Validates that key is not empty and not duplicate
			 * Clears input fields after successful addition
			 * @returns {void}
			 */
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
				
				// Clear input fields
				this.newKey = '';
				this.newValue = '';
				this.errorMessage = '';
			},
			
			/**
			 * Removes a metadata item from the local array by index
			 * @param {number} index - The index of the item to remove
			 * @returns {void}
			 */
			removeItem(index) {
				this.metadataItems.splice(index, 1);
			}
		},
		
		async mounted() {
			await this.loadMetadata();
			
			// Listen for save events from external buttons
			document.getElementById('metadata-save-btn-top')?.addEventListener('click', () => {
				this.saveMetadata();
			});
			document.getElementById('metadata-save-btn-bottom')?.addEventListener('click', () => {
				this.saveMetadata();
			});
		}
	}).mount('#metadata-app');
});
