const { createApp } = Vue;

const TitleSettingsApp = {
  data() {
    return {
      mounted: false,
      showModal: false,
      savingApply: false,
      savingClose: false,
      infoMessage: '',
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
    this.fetchData();
  },

  methods: {
    async fetchData() {
      try {
        const response = await fetch(urlTitleSettingsFetchData, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' }
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.blogTopic = data.data.blog_topic || '';
          this.infoMessage = data.data.info_message || '';
          this.mounted = true;
        }
      } catch (error) {
        console.error('Error loading settings:', error);
      }
    },

    async save(action) {
      if (action === 'apply') {
        this.savingApply = true;
      } else {
        this.savingClose = true;
      }

      try {
        const response = await fetch(urlTitleSettingsSubmit, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            action: action,
            blog_topic: this.form.blogTopic.trim()
          })
        });
        const data = await response.json();

        if (data.status === 'success') {
          if (action === 'save_close') {
            this.showModal = false;
            if (data.data?.redirect_url) {
              window.location.href = data.data.redirect_url;
              return;
            }
          }
          Swal.fire({ icon: 'success', title: 'Saved', text: data.message || 'Settings saved.', timer: 2000, toast: true, position: 'top-end', showConfirmButton: false });
        } else {
          Swal.fire({ icon: 'error', title: 'Oops...', text: data.message || 'Failed to save settings.' });
        }
      } catch (error) {
        console.error('Error saving settings:', error);
        Swal.fire({ icon: 'error', title: 'Oops...', text: 'Failed to save settings.' });
      } finally {
        this.savingApply = false;
        this.savingClose = false;
      }
    }
  }
};

document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('title-settings-app');
  if (el) createApp(TitleSettingsApp).mount('#title-settings-app');
});
