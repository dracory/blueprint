const { createApp } = Vue;

/**
 * FilesApp is a Vue.js component for managing files and directories.
 * It provides a table view with navigation, upload, rename, delete,
 * and bulk operations.
 */
const FilesApp = {
  data() {
    return {
      // UI state
      loading: true,
      uploading: false,
      creating: false,
      deleting: false,
      renaming: false,
      cloning: false,
      moving: false,
      modal: '',
      uploadProgress: 0,

      // Current directory
      currentDirectory: '',
      parentDirectory: '',

      // Data
      directories: [],
      files: [],

      // Selection
      selectedItems: [],
      selectAll: false,

      // Upload
      filesToUpload: [],

      // Modals state
      createDirName: '',
      renameOldName: '',
      renameNewName: '',
      renameType: 'file',
      renameItemPath: '',
      deleteItemName: '',
      deleteItemType: 'file',
      deleteItemPath: '',
      cloneFileName: '',
      cloneNewName: '',
      viewFileUrl: '',
      moveDestinationDir: '',
      moveDestinations: [],

      // Opener messaging
      openerArgs: {}
    };
  },

  computed: {
    /**
     * Returns the number of selected items.
     */
    selectedCount() {
      return this.selectedItems.length;
    },

    /**
     * Returns whether this window was opened by another window (file picker mode).
     */
    hasOpener() {
      return typeof window !== 'undefined' && window.opener !== null;
    }
  },

  mounted() {
    // Read current directory from URL
    const urlParams = new URLSearchParams(window.location.search);
    this.currentDirectory = urlParams.get('current_dir') || '';

    // Inter-window messaging setup
    this.setupOpenerMessaging();

    this.loadFiles();
  },

  methods: {
    /**
     * Loads files and directories from the server.
     */
    async loadFiles() {
      this.loading = true;
      try {
        const formData = new FormData();
        formData.append('action', 'load-files');
        formData.append('current_dir', this.currentDirectory);

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.directories = data.data?.directories || [];
          this.files = data.data?.files || [];
          this.parentDirectory = data.data?.parent_directory || '';
          this.selectedItems = [];
          this.selectAll = false;
        } else {
          this.showNotification(data.message || 'Failed to load files', 'error');
        }
      } catch (error) {
        console.error('Error loading files:', error);
        this.showNotification('Failed to load files', 'error');
      } finally {
        this.loading = false;
      }
    },

    /**
     * Navigates to a directory.
     */
    navigateTo(dir) {
      this.currentDirectory = dir;
      const params = new URLSearchParams();
      if (dir) params.set('current_dir', dir);
      const newUrl = `${window.location.pathname}?${params.toString()}`;
      window.history.pushState({}, '', newUrl);
      this.loadFiles();
    },

    /**
     * Uploads a single file.
     */
    async uploadFile() {
      const input = this.$refs.uploadFileInput;
      if (!input || !input.files || input.files.length === 0) {
        this.showNotification('Please select a file to upload', 'error');
        return;
      }

      this.uploading = true;
      try {
        const formData = new FormData();
        formData.append('action', 'file_upload');
        formData.append('current_dir', this.currentDirectory);
        formData.append('upload_file', input.files[0]);

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.showNotification(data.message || 'File uploaded successfully', 'success');
          this.closeModal();
          this.loadFiles();
        } else {
          this.showNotification(data.message || 'Failed to upload file', 'error');
        }
      } catch (error) {
        console.error('Error uploading file:', error);
        this.showNotification('Failed to upload file', 'error');
      } finally {
        this.uploading = false;
      }
    },

    /**
     * Handles drag-and-drop files.
     */
    handleDrop(event) {
      const files = event.dataTransfer.files;
      this.handleFiles(files);
    },

    /**
     * Handles file selection for multiple upload.
     */
    handleFiles(files) {
      for (const file of files) {
        this.filesToUpload.push(file);
      }
    },

    /**
     * Removes a file from the upload queue.
     */
    removeFile(index) {
      this.filesToUpload.splice(index, 1);
    },

    /**
     * Uploads multiple files.
     */
    async uploadMultipleFiles() {
      if (this.filesToUpload.length === 0) {
        this.showNotification('Please select at least one file to upload', 'error');
        return;
      }

      this.uploading = true;
      let successCount = 0;
      let errorCount = 0;

      for (let i = 0; i < this.filesToUpload.length; i++) {
        const file = this.filesToUpload[i];
        this.uploadProgress = Math.round(((i + 1) / this.filesToUpload.length) * 100);

        const formData = new FormData();
        formData.append('action', 'file_upload');
        formData.append('current_dir', this.currentDirectory);
        formData.append('upload_file', file);

        try {
          const response = await fetch(urlFileManager, {
            method: 'POST',
            body: formData
          });
          const data = await response.json();
          if (data.status === 'success') {
            successCount++;
          } else {
            errorCount++;
            console.error('Upload error for ' + file.name + ':', data.message);
          }
        } catch (error) {
          errorCount++;
          console.error('Upload error for ' + file.name + ':', error);
        }
      }

      this.uploadProgress = 100;

      if (errorCount === 0) {
        this.showNotification('All ' + successCount + ' files uploaded successfully', 'success');
        this.filesToUpload = [];
        this.closeModal();
        this.loadFiles();
      } else {
        this.showNotification('Uploaded ' + successCount + ' files, ' + errorCount + ' failed', 'error');
      }

      this.uploading = false;
      this.uploadProgress = 0;
    },

    /**
     * Creates a new directory.
     */
    async createDirectory() {
      if (!this.createDirName.trim()) {
        this.showNotification('Directory name is required', 'error');
        return;
      }

      this.creating = true;
      try {
        const formData = new FormData();
        formData.append('action', 'directory_create');
        formData.append('current_dir', this.currentDirectory);
        formData.append('create_dir', this.createDirName.trim());

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.showNotification(data.message || 'Directory created successfully', 'success');
          this.createDirName = '';
          this.closeModal();
          this.loadFiles();
        } else {
          this.showNotification(data.message || 'Failed to create directory', 'error');
        }
      } catch (error) {
        console.error('Error creating directory:', error);
        this.showNotification('Failed to create directory', 'error');
      } finally {
        this.creating = false;
      }
    },

    /**
     * Deletes a directory.
     */
    async deleteDirectory() {
      this.deleting = true;
      try {
        const formData = new FormData();
        formData.append('action', 'directory_delete');
        formData.append('current_dir', this.currentDirectory);
        formData.append('delete_dir', this.deleteItemName);

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.showNotification(data.message || 'Directory deleted successfully', 'success');
          this.closeModal();
          this.loadFiles();
        } else {
          this.showNotification(data.message || 'Failed to delete directory', 'error');
        }
      } catch (error) {
        console.error('Error deleting directory:', error);
        this.showNotification('Failed to delete directory', 'error');
      } finally {
        this.deleting = false;
      }
    },

    /**
     * Deletes a file.
     */
    async deleteFile() {
      this.deleting = true;
      try {
        const formData = new FormData();
        formData.append('action', 'file_delete');
        formData.append('current_dir', this.currentDirectory);
        formData.append('delete_file', this.deleteItemName);

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.showNotification(data.message || 'File deleted successfully', 'success');
          this.closeModal();
          this.loadFiles();
        } else {
          this.showNotification(data.message || 'Failed to delete file', 'error');
        }
      } catch (error) {
        console.error('Error deleting file:', error);
        this.showNotification('Failed to delete file', 'error');
      } finally {
        this.deleting = false;
      }
    },

    /**
     * Renames a file or directory.
     */
    async renameItem() {
      if (!this.renameNewName.trim()) {
        this.showNotification('New name is required', 'error');
        return;
      }

      this.renaming = true;
      try {
        const formData = new FormData();
        formData.append('action', 'file_rename');
        formData.append('current_dir', this.currentDirectory);
        formData.append('rename_file', this.renameOldName);
        formData.append('new_file', this.renameNewName.trim());

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.showNotification(data.message || 'Renamed successfully', 'success');
          this.closeModal();
          this.loadFiles();
        } else {
          this.showNotification(data.message || 'Failed to rename', 'error');
        }
      } catch (error) {
        console.error('Error renaming:', error);
        this.showNotification('Failed to rename', 'error');
      } finally {
        this.renaming = false;
      }
    },

    /**
     * Moves selected items to a destination directory.
     */
    async moveSelected() {
      if (this.selectedItems.length === 0) {
        this.showNotification('Please select at least one item to move', 'error');
        return;
      }

      this.moving = true;
      try {
        const formData = new FormData();
        formData.append('action', 'bulk_move');
        formData.append('current_dir', this.currentDirectory);
        formData.append('destination_dir', this.moveDestinationDir);
        formData.append('selected_items', JSON.stringify(this.selectedItems));

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.showNotification(data.message || 'Items moved successfully', 'success');
          this.closeModal();
          this.loadFiles();
        } else {
          this.showNotification(data.message || 'Failed to move items', 'error');
        }
      } catch (error) {
        console.error('Error moving items:', error);
        this.showNotification('Failed to move items', 'error');
      } finally {
        this.moving = false;
      }
    },

    /**
     * Deletes selected items.
     */
    async deleteSelected() {
      if (this.selectedItems.length === 0) {
        this.showNotification('Please select at least one item to delete', 'error');
        return;
      }

      this.deleting = true;
      try {
        const formData = new FormData();
        formData.append('action', 'bulk_delete');
        formData.append('current_dir', this.currentDirectory);
        formData.append('selected_items', JSON.stringify(this.selectedItems));

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.showNotification(data.message || 'Items deleted successfully', 'success');
          this.closeModal();
          this.loadFiles();
        } else {
          this.showNotification(data.message || 'Failed to delete items', 'error');
        }
      } catch (error) {
        console.error('Error deleting items:', error);
        this.showNotification('Failed to delete items', 'error');
      } finally {
        this.deleting = false;
      }
    },

    /**
     * Fetches move destinations for the selected items.
     */
    async fetchMoveDestinations() {
      try {
        const formData = new FormData();
        formData.append('action', 'get_move_destinations');
        formData.append('current_dir', this.currentDirectory);
        formData.append('selected_items', JSON.stringify(this.selectedItems));

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.moveDestinations = data.data?.directories || [];
        } else {
          this.showNotification(data.message || 'Failed to load destinations', 'error');
        }
      } catch (error) {
        console.error('Error loading destinations:', error);
        this.showNotification('Failed to load destinations', 'error');
      }
    },

    /**
     * Opens a modal by name.
     */
    openModal(name) {
      this.modal = name;
      if (name === 'moveSelected') {
        this.fetchMoveDestinations();
      }
    },

    /**
     * Closes the current modal.
     */
    closeModal() {
      this.modal = '';
      this.createDirName = '';
      this.renameNewName = '';
      this.renameOldName = '';
      this.renameType = 'file';
      this.deleteItemName = '';
      this.deleteItemType = 'file';
      this.viewFileUrl = '';
      this.cloneFileName = '';
      this.cloneNewName = '';
      this.moveDestinationDir = '';
      this.moveDestinations = [];
      this.filesToUpload = [];
      this.uploadProgress = 0;
    },

    /**
     * Shows the delete directory modal.
     */
    showDeleteDirectoryModal(dir) {
      this.deleteItemName = dir.Name;
      this.deleteItemType = 'directory';
      this.openModal('deleteDirectory');
    },

    /**
     * Shows the delete file modal.
     */
    showDeleteFileModal(file) {
      this.deleteItemName = file.Name;
      this.deleteItemType = 'file';
      this.openModal('deleteFile');
    },

    /**
     * Shows the rename modal.
     */
    showRenameModal(item, type) {
      this.renameOldName = item.Name;
      this.renameNewName = item.Name;
      this.renameType = type;
      this.openModal('rename');
    },

    /**
     * Shows the view file modal.
     */
    showViewModal(file) {
      this.viewFileUrl = file.URL;
      this.openModal('viewFile');
    },

    /**
     * Shows the clone file modal.
     */
    showCloneModal(file) {
      this.cloneFileName = file.Name;
      this.cloneNewName = this.generateCopyName(file.Name);
      this.openModal('cloneFile');
    },

    /**
     * Generates a copy filename by appending _copy before extension.
     */
    generateCopyName(filename) {
      const lastDot = filename.lastIndexOf('.');
      if (lastDot === -1) {
        return filename + '_copy';
      }
      const base = filename.substring(0, lastDot);
      const ext = filename.substring(lastDot);
      return base + '_copy' + ext;
    },

    /**
     * Clones a file by duplicating it.
     */
    async cloneFile() {
      if (!this.cloneNewName.trim()) {
        this.showNotification('New file name is required', 'error');
        return;
      }

      this.cloning = true;
      try {
        const formData = new FormData();
        formData.append('action', 'file_clone');
        formData.append('current_dir', this.currentDirectory);
        formData.append('clone_file', this.cloneFileName);
        formData.append('new_file', this.cloneNewName.trim());

        const response = await fetch(urlFileManager, {
          method: 'POST',
          body: formData
        });
        const data = await response.json();

        if (data.status === 'success') {
          this.showNotification(data.message || 'File cloned successfully', 'success');
          this.closeModal();
          this.loadFiles();
        } else {
          this.showNotification(data.message || 'Failed to clone file', 'error');
        }
      } catch (error) {
        console.error('Error cloning file:', error);
        this.showNotification('Failed to clone file', 'error');
      } finally {
        this.cloning = false;
      }
    },

    /**
     * Toggles select all on the current page.
     */
    toggleSelectAll() {
      if (this.selectAll) {
        const allItems = [];
        this.directories.forEach(dir => {
          allItems.push({ path: dir.Path, type: 'directory', name: dir.Name });
        });
        this.files.forEach(file => {
          allItems.push({ path: file.Path, type: 'file', name: file.Name });
        });
        this.selectedItems = allItems;
      } else {
        this.selectedItems = [];
      }
    },

    /**
     * Selects all items.
     */
    selectAllItems() {
      this.selectAll = true;
      this.toggleSelectAll();
    },

    /**
     * Unselects all items.
     */
    unselectAllItems() {
      this.selectAll = false;
      this.toggleSelectAll();
    },

    /**
     * Sends a message to the opener window when a file is selected.
     */
    fileSelectedUrl(selectedFileUrl) {
      if (window.opener === null) {
        return true;
      }
      window.opener.postMessage({ msg: 'media-manager-file-selected', url: selectedFileUrl, args: this.openerArgs }, '*');
      window.close();
    },

    /**
     * Sets up inter-window messaging.
     */
    setupOpenerMessaging() {
      const messageReceived = (event) => {
        this.openerArgs = event.data;
        // Show select buttons when message is received
        const selectButtons = document.querySelectorAll('.btn-select');
        selectButtons.forEach(btn => btn.style.display = 'inline-block');
      };

      window.addEventListener('message', messageReceived, false);

      if (window.opener !== null) {
        window.opener.postMessage({ msg: 'media-manager-loaded' }, '*');
      } else {
        // Hide select buttons when no opener
        const selectButtons = document.querySelectorAll('.btn-select');
        selectButtons.forEach(btn => btn.style.display = 'none');
      }
    },

    /**
     * Checks if a filename has an image extension.
     */
    isImage(filename) {
      const ext = filename.split('.').pop().toLowerCase();
      return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'bmp', 'ico'].includes(ext);
    },

    /**
     * Formats file size to human readable format.
     */
    formatFileSize(bytes) {
      if (bytes < 1024) return bytes + ' B';
      const k = 1024;
      const sizes = ['KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
    },

    /**
     * Shows a notification toast.
     */
    showNotification(message, type) {
      const notification = document.createElement('div');
      const alertClass = type === 'success' ? 'success' : 'danger';
      notification.className = 'alert alert-' + alertClass + ' alert-dismissible fade show position-fixed';
      notification.style.cssText = 'top: 20px; right: 20px; z-index: 9999; min-width: 300px;';
      notification.innerHTML =
        '<div>' + message + '</div>' +
        '<button type="button" class="btn-close" data-bs-dismiss="alert"></button>';

      document.body.appendChild(notification);
      setTimeout(() => {
        if (notification.parentNode) {
          notification.parentNode.removeChild(notification);
        }
      }, 3000);
    }
  },

  watch: {
    selectedItems(newVal) {
      const allItemsCount = this.directories.length + this.files.length;
      this.selectAll = newVal.length === allItemsCount && allItemsCount > 0;
    }
  }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('files-app');
  if (el) {
    createApp(FilesApp).mount('#files-app');
  }
});
