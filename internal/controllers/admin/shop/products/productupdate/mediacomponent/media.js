// Media Component Vue.js Application
const { createApp } = Vue;

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', function() {
	// Ensure the media-app element exists
	const mediaApp = document.getElementById('media-app');
	if (!mediaApp) {
		console.error('Media app container not found');
		return;
	}

	if (!urlMediaLoad || !urlMediaSave || !productId) {
		console.error('urlMediaLoad, urlMediaSave, or productId not found');
		return;
	}

	createApp({
		data() {
			return {
				productId: productId,
				urlMediaLoad: urlMediaLoad,
				urlMediaSave: urlMediaSave,
				mediaItems: [],
				newMediaUrl: '',
				errorMessage: '',
				successMessage: '',
				loading: false,
				draggedIndex: null,
				dragOverIndex: null
			}
		},
		methods: {
			/**
			 * Loads media from the server via AJAX POST request
			 * Populates mediaItems array with the response data
			 * Shows error message if loading fails
			 * @async
			 * @returns {Promise<void>}
			 */
			async loadMedia() {
				this.loading = true;
				try {
					console.log('Loading media from:', this.urlMediaLoad);
					console.log('Product ID:', this.productId);
					
					const response = await fetch(this.urlMediaLoad, {
						method: 'POST',
						headers: {
							'Accept': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						}
					});
					
					const result = await response.json();
					console.log('Response data:', result);
					if (result.status !== 'success') {
						throw new Error(result.message || 'Failed to load media');
					}
					this.mediaItems = result.data?.media || [];
				} catch (error) {
					console.error('Error loading media:', error);
					this.errorMessage = 'Failed to load media: ' + error.message;
				} finally {
					this.loading = false;
				}
			},
			
			/**
			 * Saves media to the server via AJAX POST request
			 * Sends current mediaItems array as JSON payload
			 * Reloads media after successful save and shows success message
			 * @async
			 * @returns {Promise<void>}
			 */
			async saveMedia() {
				this.loading = true;
				try {
					console.log('Saving media to:', this.urlMediaSave);
					console.log('Product ID:', this.productId);
					console.log('Media to save:', this.mediaItems);
					
					const response = await fetch(this.urlMediaSave, {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
							'Accept': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						},
						body: JSON.stringify({
							media: this.mediaItems
						})
					});
					
					const result = await response.json();
					console.log('Save response data:', result);
					if (result.status !== 'success') {
						throw new Error(result.message || 'Failed to save media');
					}
					this.successMessage = result.message || 'Media saved successfully';
					this.errorMessage = '';
					
					// Reload media to show updated data
					await this.loadMedia();
					
					// Clear success message after 3 seconds
					setTimeout(() => {
						this.successMessage = '';
					}, 3000);
					
				} catch (error) {
					console.error('Error saving media:', error);
					this.errorMessage = 'Failed to save media: ' + error.message;
					this.successMessage = '';
				} finally {
					this.loading = false;
				}
			},
			
			/**
			 * Adds a new media item to the local array
			 * Validates that URL is not empty
			 * Clears input field after successful addition and auto-saves
			 * @returns {void}
			 */
			async addMedia() {
				const url = this.newMediaUrl.trim();
				
				if (!url) {
					this.errorMessage = 'Please enter a media URL';
					return;
				}
				
				// Extract filename from URL
				let fileName = '';
				try {
					const urlObj = new URL(url);
					const pathParts = urlObj.pathname.split('/');
					fileName = pathParts[pathParts.length - 1] || '';
				} catch (e) {
					// If URL parsing fails, try simple string manipulation
					const parts = url.split('/');
					fileName = parts[parts.length - 1] || '';
				}
				
				this.mediaItems.push({
					id: Date.now().toString(),
					fileName: fileName,
					url: url,
					isMain: this.mediaItems.length === 0 // First media is main
				});
				
				// Clear input field
				this.newMediaUrl = '';
				this.errorMessage = '';
				
				// Auto-save
				await this.saveMedia();
			},
			
			/**
			 * Sets a media item as the main image
			 * Clears isMain flag from all other items
			 * Auto-saves after change
			 * @param {number} index - The index of the item to set as main
			 * @returns {void}
			 */
			async setAsMain(index) {
				// Clear all isMain flags
				this.mediaItems.forEach(item => item.isMain = false);
				// Set the selected item as main
				this.mediaItems[index].isMain = true;
				// Auto-save the change
				await this.saveMedia();
			},
			
			/**
			 * Removes a media item from the local array by index
			 * Updates isMain flag for remaining items
			 * @param {number} index - The index of the item to remove
			 * @returns {void}
			 */
			async removeMedia(index) {
				this.mediaItems.splice(index, 1);
				
				// Update main flag - first item becomes main if exists
				if (this.mediaItems.length > 0) {
					this.mediaItems.forEach((item, i) => {
						item.isMain = i === 0;
					});
				}
				
				// Auto-save
				await this.saveMedia();
			},

			/**
			 * Starts dragging a media item
			 * @param {DragEvent} event - The drag event
			 * @param {number} index - The index of the item being dragged
			 */
			dragStart(event, index) {
				this.draggedIndex = index;
				event.dataTransfer.effectAllowed = 'move';
				event.dataTransfer.setData('text/plain', index);
			},

			/**
			 * Handles drag over another item
			 * @param {DragEvent} event - The drag event
			 * @param {number} index - The index of the item being dragged over
			 */
			dragOver(event, index) {
				event.preventDefault();
				event.dataTransfer.dropEffect = 'move';
				if (this.draggedIndex !== null && this.draggedIndex !== index) {
					this.dragOverIndex = index;
				}
			},

			/**
			 * Handles drag enter on an item
			 * @param {DragEvent} event - The drag event
			 * @param {number} index - The index of the item
			 */
			dragEnter(event, index) {
				event.preventDefault();
				if (this.draggedIndex !== null && this.draggedIndex !== index) {
					this.dragOverIndex = index;
				}
			},

			/**
			 * Handles drag leave from an item
			 * @param {DragEvent} event - The drag event
			 * @param {number} index - The index of the item
			 */
			dragLeave(event, index) {
				if (this.dragOverIndex === index) {
					this.dragOverIndex = null;
				}
			},

			/**
			 * Handles dropping an item to reorder
			 * @param {DragEvent} event - The drop event
			 * @param {number} dropIndex - The index where item is dropped
			 */
			drop(event, dropIndex) {
				event.preventDefault();
				if (this.draggedIndex === null || this.draggedIndex === dropIndex) {
					this.draggedIndex = null;
					this.dragOverIndex = null;
					return;
				}

				// Reorder items
				const draggedItem = this.mediaItems[this.draggedIndex];
				this.mediaItems.splice(this.draggedIndex, 1);
				this.mediaItems.splice(dropIndex, 0, draggedItem);

				// Update sequence for all items
				this.mediaItems.forEach((item, i) => {
					item.sequence = i;
				});

				// Reset drag state
				this.draggedIndex = null;
				this.dragOverIndex = null;

				// Auto-save the new order
				this.saveMedia();
			}
		},
		
		async mounted() {
			await this.loadMedia();
			
			// Listen for save events from external buttons
			document.getElementById('media-save-btn-top')?.addEventListener('click', () => {
				this.saveMedia();
			});
			document.getElementById('media-save-btn-bottom')?.addEventListener('click', () => {
				this.saveMedia();
			});
		}
	}).mount('#media-app');
});
