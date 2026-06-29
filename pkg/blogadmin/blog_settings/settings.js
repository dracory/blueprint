const { createApp } = Vue;

const BlogSettingsApp = {
  data() {
    return {
      mounted: false,
      savingApply: false,
      savingClose: false,
      isEnvOverride: false,
      infoMessage: '',
      returnUrl: '',
      form: {
        blogTopic: ''
      }
    };
  },

  computed: {
    saving() {
      return this.savingApply || this.savingClose;
    }
  },

  mounted() {
    this.returnUrl = window.blogSettingsReturnUrl || '/admin/blog';
    this.fetchData();
  },

  methods: {
    async fetchData() {
      try {
        const response = await fetch(urlBlogSettingsFetchData, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' }
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.blogTopic = data.data.blog_topic || '';
          this.isEnvOverride = data.data.is_env_override || false;
          this.infoMessage = data.data.info_message || '';
          this.mounted = true;
        } else {
          Swal.fire({ icon: 'error', title: 'Oops...', text: data.message || 'Failed to load settings.' });
          this.mounted = true;
        }
      } catch (error) {
        console.error('Error loading blog settings:', error);
        Swal.fire({ icon: 'error', title: 'Oops...', text: 'Failed to load settings.' });
        this.mounted = true;
      }
    },

    async save(action) {
      if (action === 'apply') {
        this.savingApply = true;
      } else {
        this.savingClose = true;
      }

      try {
        const response = await fetch(urlBlogSettingsSubmit, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            action: action,
            blog_topic: this.form.blogTopic.trim()
          })
        });
        const data = await response.json();

        if (data.status === 'success') {
          if (action === 'save_close' && data.data?.redirect_url) {
            Swal.fire({ icon: 'success', title: 'Saved', text: data.message || 'Settings saved.', timer: 1500, showConfirmButton: false });
            setTimeout(() => window.location.href = data.data.redirect_url, 1500);
          } else {
            Swal.fire({ icon: 'success', title: 'Saved', text: data.message || 'Settings saved.', timer: 3000, toast: true, position: 'top-end', showConfirmButton: false });
          }
        } else {
          Swal.fire({ icon: 'error', title: 'Oops...', text: data.message || 'Failed to save settings.' });
        }
      } catch (error) {
        console.error('Error saving blog settings:', error);
        Swal.fire({ icon: 'error', title: 'Oops...', text: 'Failed to save settings.' });
      } finally {
        this.savingApply = false;
        this.savingClose = false;
      }
    }
  }
};

document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('blog-settings-app');
  if (el) createApp(BlogSettingsApp).mount('#blog-settings-app');
});
