const { createApp } = Vue;

createApp({
  data() {
    return {
      details: {
        title: '',
        description: '',
        status: 'draft',
        parent_id: ''
      },
      categories: [],
      loading: false
    };
  },
  mounted() {
    this.loadDetails();
    this.loadCategories();
    this.setupSaveButtons();
  },
  methods: {
    async loadDetails() {
      try {
        this.loading = true;
        const response = await fetch(urlDetailsLoad);
        const data = await response.json();
        
        if (data.status === 'success') {
          this.details = data.data.details;
        } else {
          Swal.fire('Error', data.message || 'Failed to load category details', 'error');
        }
      } catch (error) {
        console.error('Error loading details:', error);
        Swal.fire('Error', 'Failed to load category details', 'error');
      } finally {
        this.loading = false;
      }
    },
    async loadCategories() {
      try {
        const response = await fetch(urlCategoriesList);
        const data = await response.json();
        
        if (data.status === 'success') {
          this.categories = data.data.categories;
        } else {
          console.error('Error loading categories:', data.message);
        }
      } catch (error) {
        console.error('Error loading categories:', error);
      }
    },
    async saveDetails() {
      try {
        this.loading = true;
        const response = await fetch(urlDetailsSave, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(this.details)
        });
        const data = await response.json();
        
        if (data.status === 'success') {
          Swal.fire('Success', 'Category saved successfully', 'success');
        } else {
          Swal.fire('Error', data.message || 'Failed to save category', 'error');
        }
      } catch (error) {
        console.error('Error saving details:', error);
        Swal.fire('Error', 'Failed to save category', 'error');
      } finally {
        this.loading = false;
      }
    },
    setupSaveButtons() {
      const saveTopBtn = document.getElementById('details-save-btn-top');
      const saveBottomBtn = document.getElementById('details-save-btn-bottom');
      
      if (saveTopBtn) {
        saveTopBtn.addEventListener('click', () => this.saveDetails());
      }
      if (saveBottomBtn) {
        saveBottomBtn.addEventListener('click', () => this.saveDetails());
      }
    }
  }
}).mount('#details-wrapper');
