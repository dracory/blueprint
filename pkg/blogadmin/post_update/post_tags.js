const { createApp } = Vue;

const PostTagsApp = {
  data() {
    return {
      loading: true,
      tags: [],
      filteredTags: [],
      searchQuery: '',
      saving: false
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
          this.filteredTags = [...this.tags];
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

    filterTags() {
      const query = this.searchQuery.toLowerCase();
      this.filteredTags = this.tags.filter(tag => 
        tag.name.toLowerCase().includes(query) || 
        tag.slug.toLowerCase().includes(query)
      );
    },

    async toggleTag(tag) {
      if (this.saving) return;

      this.saving = true;
      try {
        const url = tag.assigned ? urlTagRemove : urlTagAdd;
        
        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ tag_id: tag.id })
        });

        const data = await response.json();
        if (data.status === 'success') {
          // Toggle the assigned state locally for immediate feedback
          tag.assigned = !tag.assigned;
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to update tag'
          });
          // Reload to restore correct state
          this.loadTags();
        }
      } catch (error) {
        console.error('Error toggling tag:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to update tag'
        });
        this.loadTags();
      } finally {
        this.saving = false;
      }
    }
  }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('post-tags-app');
  if (el) {
    createApp(PostTagsApp).mount('#post-tags-app');
  }
});
