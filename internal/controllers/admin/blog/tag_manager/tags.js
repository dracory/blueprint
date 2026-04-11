const { createApp } = Vue;

const BlogTagsApp = {
  data() {
    return {
      loading: true,
      tags: [],
      showModal: false,
      editingTag: null,
      saving: false,
      form: {
        name: '',
        slug: ''
      },
      showPostsModal: false,
      currentTag: null,
      loadingPosts: false,
      posts: [],
      postsMessage: ''
    };
  },

  mounted() {
    this.loadTags();
  },

  methods: {
    async loadTags() {
      this.loading = true;
      try {
        const response = await fetch(urlTagsLoad);
        const data = await response.json();
        if (data.status === 'success') {
          this.tags = data.data?.tags || [];
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load tags'
          });
        }
      } catch (error) {
        console.error('Error loading tags:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to load tags'
        });
      } finally {
        this.loading = false;
      }
    },

    openCreateModal() {
      this.editingTag = null;
      this.form = {
        name: '',
        slug: ''
      };
      this.showModal = true;
    },

    editTag(tag) {
      this.editingTag = tag;
      this.form = {
        name: tag.name,
        slug: tag.slug
      };
      this.showModal = true;
    },

    closeModal() {
      this.showModal = false;
      this.editingTag = null;
      this.form = { name: '', slug: '' };
    },

    autoGenerateSlug() {
      if (!this.editingTag && this.form.name && !this.form.slug) {
        this.form.slug = this.slugify(this.form.name);
      }
    },

    slugify(text) {
      return text
        .toLowerCase()
        .trim()
        .replace(/[^\w\s-]/g, '')
        .replace(/[\s_-]+/g, '-')
        .replace(/^-+|-+$/g, '');
    },

    async saveTag() {
      if (!this.form.name) return;

      this.saving = true;
      try {
        const url = this.editingTag 
          ? urlTagUpdate.replace('TAG_ID_PLACEHOLDER', this.editingTag.id)
          : urlTagCreate;
        
        const method = this.editingTag ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
          method: method,
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            name: this.form.name,
            slug: this.form.slug || this.slugify(this.form.name)
          })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Success',
            text: this.editingTag ? 'Tag updated successfully' : 'Tag created successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.closeModal();
          this.loadTags();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to save tag'
          });
        }
      } catch (error) {
        console.error('Error saving tag:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to save tag'
        });
      } finally {
        this.saving = false;
      }
    },

    async deleteTag(tag) {
      const result = await Swal.fire({
        icon: 'warning',
        title: 'Delete Tag?',
        text: `Are you sure you want to delete "${tag.name}"?`,
        showCancelButton: true,
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonColor: '#dc3545'
      });

      if (!result.isConfirmed) return;

      try {
        const response = await fetch(urlTagDelete, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ tag_id: tag.id })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Deleted',
            text: 'Tag deleted successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.loadTags();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to delete tag'
          });
        }
      } catch (error) {
        console.error('Error deleting tag:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to delete tag'
        });
      }
    },

    async viewTagPosts(tag) {
      this.currentTag = tag;
      this.showPostsModal = true;
      this.loadingPosts = true;
      this.posts = [];
      this.postsMessage = '';

      try {
        const url = urlTagPostsLoad.replace('TAG_ID_PLACEHOLDER', tag.id);
        const response = await fetch(url);
        const data = await response.json();

        if (data.status === 'success') {
          this.posts = data.data?.posts || [];
          this.postsMessage = data.data?.message || '';
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load posts'
          });
          this.postsMessage = 'Failed to load posts.';
        }
      } catch (error) {
        console.error('Error loading tag posts:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to load posts'
        });
        this.postsMessage = 'Failed to load posts.';
      } finally {
        this.loadingPosts = false;
      }
    },

    closePostsModal() {
      this.showPostsModal = false;
      this.currentTag = null;
      this.posts = [];
      this.postsMessage = '';
    }
  }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('blog-tags-app');
  if (el) {
    createApp(BlogTagsApp).mount('#blog-tags-app');
  }
});
