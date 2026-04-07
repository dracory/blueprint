// Tags Component Vue.js Application
const { createApp } = Vue;

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', function() {
	// Ensure the tags-app element exists
	const tagsApp = document.getElementById('tags-app');
	if (!tagsApp) {
		console.error('Tags app container not found');
		return;
	}

	if (!urlTagsLoad || !urlTagsSave || !productId) {
		console.error('urlTagsLoad, urlTagsSave, or productId not found');
		return;
	}

	createApp({
		data() {
			return {
				productId: productId,
				urlTagsLoad: urlTagsLoad,
				urlTagsSave: urlTagsSave,
				tagItems: [],
				newTag: '',
				errorMessage: '',
				successMessage: '',
				loading: false
			}
		},
		methods: {
			/**
			 * Loads tags from the server via AJAX POST request
			 * Populates tagItems array with the response data
			 * Shows error message if loading fails
			 * @async
			 * @returns {Promise<void>}
			 */
			async loadTags() {
				this.loading = true;
				try {
					console.log('Loading tags from:', this.urlTagsLoad);
					console.log('Product ID:', this.productId);
					
					const response = await fetch(this.urlTagsLoad, {
						method: 'POST',
						headers: {
							'Accept': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						}
					});
					
					const result = await response.json();
					console.log('Response data:', result);
					if (result.status !== 'success') {
						throw new Error(result.message || 'Failed to load tags');
					}
					this.tagItems = result.data?.tags || [];
				} catch (error) {
					console.error('Error loading tags:', error);
					this.errorMessage = 'Failed to load tags: ' + error.message;
				} finally {
					this.loading = false;
				}
			},
			
			/**
			 * Saves tags to the server via AJAX POST request
			 * Sends current tagItems array as JSON payload
			 * Reloads tags after successful save and shows success message
			 * @async
			 * @returns {Promise<void>}
			 */
			async saveTags() {
				this.loading = true;
				try {
					console.log('Saving tags to:', this.urlTagsSave);
					console.log('Product ID:', this.productId);
					console.log('Tags to save:', this.tagItems);
					
					const response = await fetch(this.urlTagsSave, {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
							'Accept': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						},
						body: JSON.stringify({
							tags: this.tagItems
						})
					});
					
					const result = await response.json();
					console.log('Save response data:', result);
					if (result.status !== 'success') {
						throw new Error(result.message || 'Failed to save tags');
					}
					this.successMessage = result.message || 'Tags saved successfully';
					this.errorMessage = '';
					
					// Reload tags to show updated data
					await this.loadTags();
					
					// Clear success message after 3 seconds
					setTimeout(() => {
						this.successMessage = '';
					}, 3000);
					
				} catch (error) {
					console.error('Error saving tags:', error);
					this.errorMessage = 'Failed to save tags: ' + error.message;
					this.successMessage = '';
				} finally {
					this.loading = false;
				}
			},
			
			/**
			 * Adds a new tag to the local array
			 * Validates that tag is not empty and not duplicate
			 * Clears input field after successful addition
			 * @returns {void}
			 */
			addTag() {
				const tag = this.newTag.trim();
				
				if (!tag) {
					this.errorMessage = 'Please enter a tag';
					return;
				}
				
				// Check for duplicate tags
				if (this.tagItems.some(item => item.tag === tag)) {
					this.errorMessage = 'This tag already exists';
					return;
				}
				
				this.tagItems.push({
					id: Date.now().toString(),
					tag: tag
				});
				
				// Clear input field
				this.newTag = '';
				this.errorMessage = '';
			},
			
			/**
			 * Adds a quick tag from the preset buttons
			 * Checks for duplicates before adding
			 * @param {string} tag - The preset tag to add
			 * @returns {void}
			 */
			addQuickTag(tag) {
				// Check for duplicate tags
				if (this.tagItems.some(item => item.tag === tag)) {
					this.errorMessage = 'This tag already exists';
					setTimeout(() => {
						this.errorMessage = '';
					}, 2000);
					return;
				}
				
				this.tagItems.push({
					id: Date.now().toString(),
					tag: tag
				});
				
				this.errorMessage = '';
			},
			
			/**
			 * Removes a tag from the local array by index
			 * @param {number} index - The index of the tag to remove
			 * @returns {void}
			 */
			removeTag(index) {
				this.tagItems.splice(index, 1);
			}
		},
		
		async mounted() {
			await this.loadTags();
			
			// Listen for save events from external buttons
			document.getElementById('tags-save-btn-top')?.addEventListener('click', () => {
				this.saveTags();
			});
			document.getElementById('tags-save-btn-bottom')?.addEventListener('click', () => {
				this.saveTags();
			});
		}
	}).mount('#tags-app');
});
