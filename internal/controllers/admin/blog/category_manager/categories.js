const { createApp } = Vue;

const BlogCategoriesApp = {
  data() {
    return {
      loading: true,
      categories: [],
      showModal: false,
      editingCategory: null,
      saving: false,
      form: {
        name: '',
        slug: '',
        description: ''
      }
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
          this.$nextTick(() => {
            this.initSortable();
          });
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

    initSortable() {
      const list = document.querySelector('.sortable-list');
      if (!list) {
        setTimeout(() => this.initSortable(), 200);
        return;
      }

      if (typeof Sortable === 'undefined') {
        setTimeout(() => this.initSortable(), 100);
        return;
      }

      // Destroy existing instance if present
      if (this.sortable) {
        this.sortable.destroy();
      }

      this.sortable = new Sortable(list, {
        handle: '.drag-handle',
        animation: 150,
        ghostClass: 'sortable-ghost',
        chosenClass: 'sortable-chosen',
        dragClass: 'sortable-drag',
        forceFallback: true,
        fallbackClass: 'sortable-fallback',
        fallbackTolerance: 3,
        onEnd: (evt) => {
          this.handleReorder(evt.oldIndex, evt.newIndex);
        }
      });
    },

    async handleReorder(oldIndex, newIndex) {
      if (oldIndex === newIndex) return;

      // Reorder locally first for immediate feedback
      const movedItem = this.categories.splice(oldIndex, 1)[0];
      this.categories.splice(newIndex, 0, movedItem);

      // Send new order to server
      try {
        const categoryIds = this.categories.map(c => c.id);
        const response = await fetch(urlCategoriesReorder, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ category_ids: categoryIds })
        });
        const data = await response.json();
        if (data.status !== 'success') {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to reorder categories'
          });
          // Reload to restore original order
          this.loadCategories();
        }
      } catch (error) {
        console.error('Error reordering categories:', error);
        this.loadCategories();
      }
    },

    openCreateModal() {
      this.editingCategory = null;
      this.form = {
        name: '',
        slug: '',
        description: ''
      };
      this.showModal = true;
    },

    editCategory(category) {
      this.editingCategory = category;
      this.form = {
        name: category.name,
        slug: category.slug,
        description: category.description || ''
      };
      this.showModal = true;
    },

    closeModal() {
      this.showModal = false;
      this.editingCategory = null;
      this.form = { name: '', slug: '', description: '' };
    },

    autoGenerateSlug() {
      if (!this.editingCategory && this.form.name && !this.form.slug) {
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

    async saveCategory() {
      if (!this.form.name) return;

      this.saving = true;
      try {
        const url = this.editingCategory 
          ? urlCategoryUpdate.replace('CATEGORY_ID_PLACEHOLDER', this.editingCategory.id)
          : urlCategoryCreate;
        
        const method = this.editingCategory ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
          method: method,
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            name: this.form.name,
            slug: this.form.slug || this.slugify(this.form.name),
            description: this.form.description
          })
        });

        const data = await response.json();
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Success',
            text: this.editingCategory ? 'Category updated successfully' : 'Category created successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.closeModal();
          this.loadCategories();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to save category'
          });
        }
      } catch (error) {
        console.error('Error saving category:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to save category'
        });
      } finally {
        this.saving = false;
      }
    },

    async deleteCategory(category) {
      const result = await Swal.fire({
        icon: 'warning',
        title: 'Delete Category?',
        text: `Are you sure you want to delete "${category.name}"?`,
        showCancelButton: true,
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonColor: '#dc3545'
      });

      if (!result.isConfirmed) return;

      try {
        const response = await fetch(urlCategoryDelete, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ category_id: category.id })
        });

        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire({
            icon: 'success',
            title: 'Deleted',
            text: 'Category deleted successfully',
            timer: 1500,
            showConfirmButton: false
          });
          this.loadCategories();
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to delete category'
          });
        }
      } catch (error) {
        console.error('Error deleting category:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to delete category'
        });
      }
    }
  }
};

// Mount the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('blog-categories-app');
  if (el) {
    createApp(BlogCategoriesApp).mount('#blog-categories-app');
  }
});
