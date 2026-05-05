(function() {
  const { createApp } = Vue;

  const PostVersioningApp = {
    data() {
      return {
        loading: true,
        restoring: false,
        postId: '',
        urlVersionsLoad: '',
        urlVersionDetail: '',
        urlVersionRestore: '',
        versioningEnabled: false,
        versions: [],
        detailVersion: null,
        detailLoading: false,
        detailAttributes: [],
        detailError: '',
        selectedAttributes: [],
        errorMessage: '',
        successMessage: ''
      };
    },

    mounted() {
      if (window.postVersioningConfig) {
        this.postId = window.postVersioningConfig.postId || '';
        this.urlVersionsLoad = window.postVersioningConfig.urlVersionsLoad || '';
        this.urlVersionDetail = window.postVersioningConfig.urlVersionDetail || '';
        this.urlVersionRestore = window.postVersioningConfig.urlVersionRestore || '';
      }

      // Load versions when modal is opened
      const modalEl = document.getElementById('versionHistoryModal');
      if (modalEl) {
        modalEl.addEventListener('shown.bs.modal', () => {
          this.loadVersions();
        });

        // Reset detail view when modal is closed
        modalEl.addEventListener('hidden.bs.modal', () => {
          this.closeDetail();
        });
      }

      // Wire up footer buttons
      const backBtn = document.getElementById('versionHistoryModalBackBtn');
      const restoreBtn = document.getElementById('versionHistoryModalRestoreBtn');
      if (backBtn) {
        backBtn.addEventListener('click', () => this.closeDetail());
      }
      if (restoreBtn) {
        restoreBtn.addEventListener('click', () => this.restoreSelectedAttributes());
      }
    },

    watch: {
      detailVersion(newVal) {
        // Update modal title
        const titleEl = document.getElementById('versionHistoryModalLabel');
        if (titleEl) {
          titleEl.textContent = newVal ? this.formatDateTime(newVal.created_at) : 'Post Revisions';
        }

        // Show/hide footer buttons
        const backBtn = document.getElementById('versionHistoryModalBackBtn');
        const restoreBtn = document.getElementById('versionHistoryModalRestoreBtn');
        if (backBtn) {
          backBtn.classList.toggle('d-none', !newVal);
        }
        if (restoreBtn) {
          restoreBtn.classList.toggle('d-none', !newVal);
        }
      },
      restoring(newVal) {
        // Update restore button disabled state
        const restoreBtn = document.getElementById('versionHistoryModalRestoreBtn');
        if (restoreBtn) {
          restoreBtn.disabled = newVal || this.selectedAttributes.length === 0;
        }
      },
      selectedAttributes(newVal) {
        // Update restore button disabled state
        const restoreBtn = document.getElementById('versionHistoryModalRestoreBtn');
        if (restoreBtn) {
          restoreBtn.disabled = this.restoring || newVal.length === 0;
        }
      }
    },

    methods: {
      async loadVersions() {
        this.loading = true;
        this.errorMessage = '';
        this.successMessage = '';
        try {
          const response = await fetch(this.urlVersionsLoad, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ post_id: this.postId })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.versioningEnabled = data.data?.versioning_enabled || false;
            this.versions = data.data?.versions || [];
          } else {
            this.errorMessage = data.message || 'Failed to load versions';
          }
        } catch (error) {
          console.error('Error loading versions:', error);
          this.errorMessage = 'Failed to load versions';
        } finally {
          this.loading = false;
        }
      },

      async openVersionDetail(version) {
        this.detailVersion = version;
        this.detailLoading = true;
        this.detailError = '';
        this.selectedAttributes = [];

        try {
          const response = await fetch(this.urlVersionDetail, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ version_id: version.id })
          });
          const data = await response.json();
          if (data.status === 'success') {
            this.detailAttributes = data.data?.attributes || [];
          } else {
            this.detailError = data.message || 'Failed to load version detail';
          }
        } catch (error) {
          console.error('Error loading version detail:', error);
          this.detailError = 'Failed to load version detail';
        } finally {
          this.detailLoading = false;
        }
      },

      closeDetail() {
        this.detailVersion = null;
        this.detailAttributes = [];
        this.selectedAttributes = [];
        this.detailError = '';
      },

      async restoreSelectedAttributes() {
        if (this.restoring || this.selectedAttributes.length === 0) return;

        this.restoring = true;
        this.errorMessage = '';
        this.successMessage = '';

        try {
          const response = await fetch(this.urlVersionRestore, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              post_id: this.postId,
              version_id: this.detailVersion.id,
              attributes: this.selectedAttributes
            })
          });

          const data = await response.json();
          if (data.status === 'success') {
            this.successMessage = 'Selected attributes restored successfully';
            this.selectedAttributes = [];
            // Reload versions after restore
            await this.loadVersions();

            setTimeout(() => {
              this.successMessage = '';
            }, 3000);
          } else {
            this.errorMessage = data.message || 'Failed to restore attributes';
          }
        } catch (error) {
          console.error('Error restoring attributes:', error);
          this.errorMessage = 'Failed to restore attributes';
        } finally {
          this.restoring = false;
        }
      },

      formatDateTime(datetimeStr) {
        if (!datetimeStr) return 'Unknown';
        const parsedTime = new Date(datetimeStr);
        if (isNaN(parsedTime.getTime())) return 'Unknown';
        return parsedTime.toLocaleString('en-US', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
          hour12: false
        });
      },

      formatRelativeTime(datetimeStr) {
        if (!datetimeStr) return 'Unknown';
        const parsedTime = new Date(datetimeStr);
        if (isNaN(parsedTime.getTime())) return 'Unknown';

        const now = new Date();
        const diffMs = now - parsedTime;
        const diffSec = Math.floor(diffMs / 1000);
        const diffMin = Math.floor(diffSec / 60);
        const diffHour = Math.floor(diffMin / 60);
        const diffDay = Math.floor(diffHour / 24);

        if (diffSec < 60) {
          return 'just now';
        } else if (diffMin < 60) {
          return diffMin + ' minute' + (diffMin === 1 ? '' : 's') + ' ago';
        } else if (diffHour < 24) {
          return diffHour + ' hour' + (diffHour === 1 ? '' : 's') + ' ago';
        } else if (diffDay < 30) {
          return diffDay + ' day' + (diffDay === 1 ? '' : 's') + ' ago';
        } else if (diffDay < 365) {
          const months = Math.floor(diffDay / 30);
          return months + ' month' + (months === 1 ? '' : 's') + ' ago';
        } else {
          const years = Math.floor(diffDay / 365);
          return years + ' year' + (years === 1 ? '' : 's') + ' ago';
        }
      }
    }
  };

  const el = document.getElementById('post-versioning-app');
  if (el) {
    createApp(PostVersioningApp).mount('#post-versioning-app');
  }
})();
