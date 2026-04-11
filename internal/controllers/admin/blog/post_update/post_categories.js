const { createApp } = Vue;

const PostCategoriesApp = {
  data() {
    return {
      loading: true,
      categories: [],
      saving: false
    };
  },

  mounted() {
    this.loadCategories();
  },

  methods: {
    async loadCategories() {
      this.loading = true;
      try {
        const response = await fetch(urlCategoriesLoad);
        const data = await response.json();
        if (data.status === 'success') {
          this.categories = data.data?.categories || [];
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load categories'
          });
        }
      } catch (error) {
        console.error('Error loading categories:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to load categories'
        });
      } finally {
        this.loading = false;
      }
    },

    async toggleCategory(category) {
      if (this.saving) return;

      this.saving = true;
      try {
        const url = category.assigned ? urlCategoryRemove : urlCategoryAdd;
        
        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ category_id: category.id })
        });

        const data = await response.json();
        if (data.status === 'success') {
          // Toggle the assigned state locally for immediate feedback
          category.assigned = !category.assigned;
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to update category'
          });
          // Reload to restore correct state
          this.loadCategories();
        }
      } catch (error) {
        console.error('Error toggling category:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to update category'
        });
        this.loadCategories();
      } finally {
        this.saving = false;
      }
    }
  }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('post-categories-app');
  if (el) {
    createApp(PostCategoriesApp).mount('#post-categories-app');
  }
});
