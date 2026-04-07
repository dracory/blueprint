document.addEventListener('DOMContentLoaded', function() {
  const { createApp } = Vue;

  const app = createApp({
    data() {
      return {
        categories: [],
        loading: false,
        showCreateModal: false,
        newCategory: {
          title: '',
          description: '',
          status: 'draft',
          parent_id: ''
        }
      };
    },
    mounted() {
      this.loadCategories();
    },
    methods: {
      async loadCategories() {
        try {
          this.loading = true;
          const response = await fetch(urlCategoriesLoad);
          const data = await response.json();
          
          if (data.status === 'success') {
            this.categories = this.buildCategoryTree(data.data.categories);
          } else {
            Swal.fire('Error', data.message || 'Failed to load categories', 'error');
          }
        } catch (error) {
          console.error('Error loading categories:', error);
          Swal.fire('Error', 'Failed to load categories', 'error');
        } finally {
          this.loading = false;
        }
      },
      buildCategoryTree(categories) {
        const categoryMap = {};
        const roots = [];
        
        // First pass: create map
        categories.forEach(cat => {
          categoryMap[cat.id] = { ...cat, children: [] };
        });
        
        // Second pass: build tree
        categories.forEach(cat => {
          if (cat.parent_id && categoryMap[cat.parent_id]) {
            categoryMap[cat.parent_id].children.push(categoryMap[cat.id]);
          } else {
            roots.push(categoryMap[cat.id]);
          }
        });
        
        return roots;
      },
      openCreateModal(parentId = null) {
        this.newCategory.parent_id = parentId || '';
        this.showCreateModal = true;
      },
      closeCreateModal() {
        this.showCreateModal = false;
        this.newCategory = {
          title: '',
          description: '',
          status: 'draft',
          parent_id: ''
        };
      },
      async createCategory() {
        if (!this.newCategory.title) {
          Swal.fire('Error', 'Title is required', 'error');
          return;
        }
        
        try {
          const formData = new FormData();
          formData.append('title', this.newCategory.title);
          formData.append('description', this.newCategory.description);
          formData.append('status', this.newCategory.status);
          formData.append('parent_id', this.newCategory.parent_id);
          
          const response = await fetch(urlCategoryCreate, {
            method: 'POST',
            body: formData
          });
          
          const html = await response.text();
          
          if (html.includes('Category created successfully')) {
            Swal.fire('Success', 'Category created successfully', 'success');
            this.closeCreateModal();
            this.loadCategories();
          } else {
            Swal.fire('Error', 'Failed to create category', 'error');
          }
        } catch (error) {
          console.error('Error creating category:', error);
          Swal.fire('Error', 'Failed to create category', 'error');
        }
      },
      editCategory(categoryId) {
        window.location.href = urlCategoryUpdate.replace('category_id=', 'category_id=' + categoryId);
      },
      async deleteCategory(categoryId) {
        const result = await Swal.fire({
          title: 'Are you sure?',
          text: 'This will delete the category and all its subcategories',
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: 'Yes, delete it!',
          cancelButtonText: 'No, cancel!'
        });
        
        if (result.isConfirmed) {
          try {
            const formData = new FormData();
            formData.append('category_id', categoryId);
            
            const response = await fetch(urlCategoryDelete, {
              method: 'POST',
              body: formData
            });
            
            const data = await response.json();
            
            if (data.success) {
              Swal.fire('Deleted!', 'Category has been deleted.', 'success');
              this.loadCategories();
            } else {
              Swal.fire('Error', data.message || 'Failed to delete category', 'error');
            }
          } catch (error) {
            console.error('Error deleting category:', error);
            Swal.fire('Error', 'Failed to delete category', 'error');
          }
        }
      },
      getCategoryClass(status) {
        switch(status) {
          case 'active': return 'text-success';
          case 'inactive': return 'text-danger';
          default: return 'text-warning';
        }
      }
    }
  });

  app.component('category-tree-item', {
    props: ['category'],
    template: `
      <div class="category-item mb-2">
        <div class="card">
          <div class="card-body d-flex justify-content-between align-items-center">
            <div>
              <h5 :class="getCategoryClass(category.status)">
                <i class="bi bi-folder me-2"></i>{{ category.title }}
              </h5>
              <small class="text-muted">{{ category.description || 'No description' }}</small>
            </div>
            <div>
              <button class="btn btn-sm btn-outline-primary me-1" @click="$emit('edit', category.id)">
                <i class="bi bi-pencil"></i>
              </button>
              <button class="btn btn-sm btn-outline-secondary me-1" @click="$emit('add-child', category.id)">
                <i class="bi bi-plus"></i>
              </button>
              <button class="btn btn-sm btn-outline-danger" @click="$emit('delete', category.id)">
                <i class="bi bi-trash"></i>
              </button>
            </div>
          </div>
        </div>
        <div v-if="category.children && category.children.length > 0" class="ms-4 mt-2">
          <category-tree-item 
            v-for="child in category.children" 
            :key="child.id" 
            :category="child"
            @edit="$emit('edit', $event)"
            @delete="$emit('delete', $event)"
            @add-child="$emit('add-child', $event)">
          </category-tree-item>
        </div>
      </div>
    `,
    methods: {
      getCategoryClass(status) {
        switch(status) {
          case 'active': return 'text-success';
          case 'inactive': return 'text-danger';
          default: return 'text-warning';
        }
      }
    }
  });

  app.mount('#categories-wrapper');
});
