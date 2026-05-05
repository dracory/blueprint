const { createApp } = Vue;

const PostSEOApp = {
  data() {
    return {
      loading: true,
      saving: false,
      postId: '',
      newOldSlug: '',
      form: {
        slug: '',
        canonicalUrl: '',
        metaDescription: '',
        metaKeywords: '',
        metaRobots: '',
        oldSlugs: []
      }
    };
  },

  mounted() {
    // Initialize postId from global variable
    if (typeof postId !== 'undefined') {
      this.postId = postId;
    }
    this.loadSEO();
  },

  methods: {
    async loadSEO() {
      this.loading = true;
      try {
        const response = await fetch(urlSEOLoad, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ post_id: this.postId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.slug = data.data?.slug || '';
          this.form.canonicalUrl = data.data?.canonical_url || '';
          this.form.metaDescription = data.data?.meta_description || '';
          this.form.metaKeywords = data.data?.meta_keywords || '';
          this.form.metaRobots = data.data?.meta_robots || '';
          this.form.oldSlugs = data.data?.old_slugs || [];
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load SEO data'
          });
        }
      } catch (error) {
        console.error('Error loading SEO:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to load SEO data'
        });
      } finally {
        this.loading = false;
      }
    },

    async saveSEO() {
      if (this.saving) return;

      this.saving = true;
      try {
        const response = await fetch(urlSEOSave, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            post_id: this.postId,
            post_slug: this.form.slug,
            post_canonical_url: this.form.canonicalUrl,
            post_meta_description: this.form.metaDescription,
            post_meta_keywords: this.form.metaKeywords,
            post_meta_robots: this.form.metaRobots,
            post_old_slugs: this.form.oldSlugs
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
            text: data.message || 'Failed to save SEO data'
          });
        }
      } catch (error) {
        console.error('Error saving SEO:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to save SEO data'
        });
      } finally {
        this.saving = false;
      }
    },

    addOldSlug() {
      const slug = this.newOldSlug.trim();
      if (!slug) return;

      // Check for duplicates
      if (this.form.oldSlugs.includes(slug)) {
        Swal.fire({
          icon: 'warning',
          title: 'Duplicate',
          text: 'This slug already exists in the old slugs list',
          timer: 2000,
          timerProgressBar: true,
          showConfirmButton: false
        });
        return;
      }

      this.form.oldSlugs.push(slug);
      this.newOldSlug = '';
    },

    removeOldSlug(index) {
      this.form.oldSlugs.splice(index, 1);
    }
  }
};

// Mount the app immediately
const el = document.getElementById('post-seo-app');
if (el) {
  createApp(PostSEOApp).mount('#post-seo-app');
}
