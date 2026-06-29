const { createApp } = Vue;

const PostDetailsApp = {
  data() {
    return {
      loading: true,
      saving: false,
      regenerating: false,
      showAdvancedTools: false,
      showImagePicker: false,
      loadingMedia: false,
      mediaImages: [],
      postId: '',
      editorOptions: [
        // { value: 'blockarea', label: 'BlockArea' },
        // { value: 'blockeditor', label: 'BlockEditor' },
        { value: 'textarea', label: 'Text Area' },
        { value: 'htmlarea', label: 'HTML Area (Summernote)' },
        { value: 'markdown', label: 'Markdown (Text Area)' },
        { value: 'markdown_easymde', label: 'Markdown (EasyMDE)' },
        { value: 'markdown_codemirror', label: 'Markdown (CodeMirror)' },
        { value: 'html_codemirror', label: 'HTML (CodeMirror)' },
      ],
      form: {
        status: '',
        imageUrl: '',
        featured: '',
        publishedAt: '',
        editor: 'textarea', // Default to text area, will be updated from API
        memo: ''
      }
    };
  },

  mounted() {
    // Initialize postId from global variable
    if (typeof postId !== 'undefined') {
      this.postId = postId;
    }
    this.loadDetails();
  },

  methods: {
    async loadDetails() {
      this.loading = true;
      try {
        const response = await fetch(urlDetailsLoad, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ post_id: this.postId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.status = data.data?.status || '';
          this.form.imageUrl = data.data?.image_url || '';
          this.form.featured = data.data?.featured || '';
          this.form.publishedAt = data.data?.published_at || '';
          this.form.editor = data.data?.editor || '';
          this.form.memo = data.data?.memo || '';
          // Format publishedAt for datetime-local input
          if (this.form.publishedAt) {
            this.form.publishedAt = this.formatDateTimeForInput(this.form.publishedAt);
          }
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load details'
          });
        }
      } catch (error) {
        console.error('Error loading details:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to load details'
        });
      } finally {
        this.loading = false;
      }
    },

    async saveDetails() {
      if (this.saving) return;

      if (this.form.status === '') {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Status is required'
        });
        return;
      }

      this.saving = true;
      try {
        const response = await fetch(urlDetailsSave, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            post_id: this.postId,
            post_status: this.form.status,
            post_image_url: this.form.imageUrl,
            post_featured: this.form.featured,
            post_published_at: this.form.publishedAt,
            post_editor: this.form.editor,
            post_memo: this.form.memo
          })
        });

        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Success',
            text: 'Post saved successfully',
            position: 'top-end',
            timer: 3000,
            timerProgressBar: true,
            showConfirmButton: false
          });
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to save details'
          });
        }
      } catch (error) {
        console.error('Error saving details:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to save details'
        });
      } finally {
        this.saving = false;
      }
    },

    async confirmRegenerateImage() {
      if (this.regenerating) return;

      Swal.fire({
        title: 'Regenerate image?',
        text: 'This will replace the current image URL for this post.',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Yes, regenerate',
        cancelButtonText: 'Cancel'
      }).then((result) => {
        if (result.isConfirmed) {
          this.regenerateImage();
        }
      });
    },

    async regenerateImage() {
      this.regenerating = true;
      try {
        const response = await fetch(urlRegenerateImage, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ post_id: this.postId })
        });

        const data = await response.json();
        if (data.status === 'success') {
          this.form.imageUrl = data.data.image_url || '';
          Swal.fire({
            icon: 'success',
            title: 'Success',
            text: 'Image regenerated successfully',
            position: 'top-end',
            timer: 3000,
            timerProgressBar: true,
            showConfirmButton: false
          });
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to regenerate image'
          });
        }
      } catch (error) {
        console.error('Error regenerating image:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to regenerate image'
        });
      } finally {
        this.regenerating = false;
      }
    },

    toggleAdvancedTools() {
      this.showAdvancedTools = !this.showAdvancedTools;
    },

    async openImagePicker() {
      this.showImagePicker = true;
      this.loadingMedia = true;
      this.mediaImages = [];
      try {
        const response = await fetch(urlMediaLoad, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ post_id: this.postId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          const allFiles = data.data?.files || [];
          this.mediaImages = allFiles.filter(f => f.type && f.type.startsWith('image/'));
        }
      } catch (error) {
        console.error('Error loading media:', error);
      } finally {
        this.loadingMedia = false;
      }
    },

    selectImage(url) {
      this.form.imageUrl = url;
      this.showImagePicker = false;
    },

    formatDateTimeForInput(datetimeStr) {
      if (!datetimeStr) return '';
      // Convert from "2006-01-02 15:04:05" to "2006-01-02T15:04"
      const date = new Date(datetimeStr.replace(' ', 'T'));
      if (isNaN(date.getTime())) return '';
      
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      
      return `${year}-${month}-${day}T${hours}:${minutes}`;
    }
  }
};

// Mount the app immediately
const el = document.getElementById('post-details-app');
if (el) {
  createApp(PostDetailsApp).mount('#post-details-app');
}
