function initMediaApp() {
    if (typeof Vue === 'undefined') {
        setTimeout(initMediaApp, 100);
        return;
    }

    const { createApp } = Vue;

    createApp({
    data() {
        return {
            loading: false,
            uploading: false,
            uploadProgress: 0,
            isDragOver: false,
            mediaItems: [],
            newMediaUrl: '',
            newMediaFileName: '',
            showAddModal: false,
            showEditModal: false,
            editIndex: null,
            editForm: {
                url: '',
                fileName: ''
            },
            draggedIndex: null,
            dragOverIndex: null
        };
    },
    methods: {
        async loadMedia() {
            this.loading = true;
            try {
                const response = await fetch(urlMediaLoad, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ product_id: productId })
                });
                const data = await response.json();
                if (data.status === 'success') {
                    this.mediaItems = data.data.media || [];
                    // Extract fileName from URL if empty
                    this.mediaItems.forEach(item => {
                        if (!item.fileName && item.url) {
                            item.fileName = this.getFileNameFromUrl(item.url);
                        }
                    });
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to load media'
                    });
                }
            } catch (error) {
                console.error('Error loading media:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to load media'
                });
            } finally {
                this.loading = false;
            }
        },

        async saveMedia() {
            this.loading = true;
            try {
                const response = await fetch(urlMediaSave, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        product_id: productId,
                        media: this.mediaItems
                    })
                });
                const data = await response.json();
                if (data.status === 'success') {
                    Swal.fire({
                        icon: 'success',
                        title: 'Success',
                        text: 'Media saved successfully',
                        timer: 2000,
                        timerProgressBar: true,
                        showConfirmButton: false
                    });
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: data.message || 'Failed to save media'
                    });
                }
            } catch (error) {
                console.error('Error saving media:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to save media'
                });
            } finally {
                this.loading = false;
            }
        },

        addMedia() {
            if (!this.newMediaUrl.trim()) return;

            const mediaType = this.determineMediaType(this.newMediaUrl);
            this.mediaItems.push({
                id: 'new_' + Date.now(),
                url: this.newMediaUrl.trim(),
                fileName: this.newMediaFileName.trim() || this.getFileNameFromUrl(this.newMediaUrl),
                type: mediaType,
                isMain: this.mediaItems.length === 0
            });
            // Clear input fields and close modal
            this.newMediaUrl = '';
            this.newMediaFileName = '';
            this.showAddModal = false;

            // Auto-save after adding
            this.saveMedia();
        },

        handleFileSelect(event) {
            const files = event.target.files;
            if (files && files.length > 0) {
                this.uploadFiles(files);
            }
            event.target.value = '';
        },

        handleDrop(event) {
            this.isDragOver = false;
            const files = event.dataTransfer.files;
            if (files && files.length > 0) {
                this.uploadFiles(files);
            }
        },

        async uploadFiles(files) {
            this.uploading = true;
            this.uploadProgress = 0;

            const formData = new FormData();
            formData.append('product_id', productId);
            for (let i = 0; i < files.length; i++) {
                formData.append('files[]', files[i]);
            }

            try {
                const xhr = new XMLHttpRequest();
                xhr.open('POST', urlMediaUpload);

                xhr.upload.onprogress = (e) => {
                    if (e.lengthComputable) {
                        this.uploadProgress = Math.round((e.loaded / e.total) * 100);
                    }
                };

                xhr.onload = () => {
                    try {
                        const data = JSON.parse(xhr.responseText);
                        if (data.status === 'success') {
                            Swal.fire({
                                icon: 'success',
                                title: 'Success',
                                text: 'Files uploaded successfully',
                                position: 'top-end',
                                timer: 3000,
                                timerProgressBar: true,
                                showConfirmButton: false
                            });
                            this.loadMedia();
                        } else {
                            Swal.fire({
                                icon: 'error',
                                title: 'Error',
                                text: data.message || 'Failed to upload files'
                            });
                        }
                    } catch (e) {
                        Swal.fire({
                            icon: 'error',
                            title: 'Error',
                            text: 'Failed to parse upload response'
                        });
                    }
                    this.uploading = false;
                    if (this.showAddModal) {
                        this.showAddModal = false;
                    }
                };

                xhr.onerror = () => {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error',
                        text: 'Failed to upload files'
                    });
                    this.uploading = false;
                };

                xhr.send(formData);
            } catch (error) {
                console.error('Error uploading files:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Error',
                    text: 'Failed to upload files'
                });
                this.uploading = false;
            }
        },

        removeMedia(index) {
            this.mediaItems.splice(index, 1);
            
            // Update isMain flag if needed
            if (this.mediaItems.length > 0) {
                this.mediaItems[0].isMain = true;
            }
            
            // Auto-save after removing
            this.saveMedia();
        },

        determineMediaType(url) {
            if (!url) return 'file';
            const extension = url.split('.').pop().toLowerCase();
            const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'];
            const videoExtensions = ['mp4', 'webm', 'ogg', 'mov'];
            const audioExtensions = ['mp3', 'wav', 'ogg'];

            if (imageExtensions.includes(extension)) return 'image';
            if (videoExtensions.includes(extension)) return 'video';
            if (audioExtensions.includes(extension)) return 'audio';
            return 'file';
        },

        getFileNameFromUrl(url) {
            if (!url) return '';
            const parts = url.split('/');
            return parts[parts.length - 1] || 'Unnamed';
        },

        dragStart(event, index) {
            this.draggedIndex = index;
            event.dataTransfer.effectAllowed = 'move';
            event.dataTransfer.setData('text/plain', index);
        },

        dragOver(event, index) {
            event.preventDefault();
            event.dataTransfer.dropEffect = 'move';
            if (this.draggedIndex !== null && this.draggedIndex !== index) {
                this.dragOverIndex = index;
            }
        },

        dragEnter(event, index) {
            event.preventDefault();
            if (this.draggedIndex !== null && this.draggedIndex !== index) {
                this.dragOverIndex = index;
            }
        },

        dragLeave(event, index) {
            if (this.dragOverIndex === index) {
                this.dragOverIndex = null;
            }
        },

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
        },

        openUrl(url) {
            if (url) {
                window.open(url, '_blank');
            }
        },

        openEditModal(index) {
            this.editIndex = index;
            this.editForm = {
                url: this.mediaItems[index].url,
                fileName: this.mediaItems[index].fileName
            };
            this.showEditModal = true;
        },

        closeEditModal() {
            this.showEditModal = false;
            this.editIndex = null;
            this.editForm = {
                url: '',
                fileName: ''
            };
        },

        async saveEdit() {
            if (!this.editForm.url.trim()) {
                this.errorMessage = 'URL cannot be empty';
                return;
            }

            // Update the media item
            this.mediaItems[this.editIndex].url = this.editForm.url.trim();
            // If file name is cleared, extract from URL
            if (!this.editForm.fileName.trim()) {
                this.mediaItems[this.editIndex].fileName = this.getFileNameFromUrl(this.editForm.url);
            } else {
                this.mediaItems[this.editIndex].fileName = this.editForm.fileName.trim();
            }
            
            // Close modal
            this.closeEditModal();

            // Auto-save
            await this.saveMedia();
        }
    },
    mounted() {
        this.loadMedia();
    }
}).mount('#media-app');
}

// Initialize the media app when DOM is ready
document.addEventListener('DOMContentLoaded', initMediaApp);
