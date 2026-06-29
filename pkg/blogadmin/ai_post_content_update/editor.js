const { createApp } = Vue;

const PostEditorApp = {
  data() {
    return {
      mounted: false,
      savingApply: false,
      savingClose: false,
      regenerating: null,
      error: '',
      success: '',
      backUrl: '',
      postId: '',
      title: '',
      summary: '',
      blocks: []
    };
  },

  computed: {
    saving() {
      return this.savingApply || this.savingClose;
    }
  },

  mounted() {
    this.postId = window.postEditorPostId || '';
    this.backUrl = window.postEditorBackUrl || '/admin/blog';
    this.fetchData();
    this.$nextTick(() => {
      this.initSortable();
      this.autoResizeAll();
    });
  },

  updated() {
    this.$nextTick(() => {
      this.autoResizeAll();
    });
  },

  methods: {
    blockLabel(type) {
      switch (type) {
        case 'h1': return 'Heading 1';
        case 'h2': return 'Heading 2';
        case 'code': return 'Code Block';
        default: return 'Paragraph';
      }
    },

    async fetchData() {
      try {
        const response = await fetch(urlPostEditorFetchData, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' }
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.title = data.data.title || '';
          this.summary = data.data.summary || '';
          this.blocks = data.data.blocks || [];
          this.mounted = true;
          this.$nextTick(() => {
            this.initSortable();
            this.autoResizeAll();
          });
        } else {
          this.error = data.message || 'Failed to load post data.';
          this.mounted = true;
        }
      } catch (err) {
        this.error = 'Failed to load post data.';
        this.mounted = true;
      }
    },

    initSortable() {
      const container = this.$refs.blocksContainer;
      if (!container || !window.Sortable || container._sortableInitialized) return;
      container._sortableInitialized = true;
      window.Sortable.create(container, {
        animation: 150,
        handle: '.block-drag-handle',
        onSort: () => {
          const cards = container.querySelectorAll('.block-card');
          const newOrder = [];
          cards.forEach(card => {
            const id = card.getAttribute('data-block-id');
            const block = this.blocks.find(b => b.id === id);
            if (block) newOrder.push(block);
          });
          this.blocks = newOrder;
        }
      });
    },

    autoResize(event) {
      const el = event.target;
      el.style.height = 'auto';
      el.style.height = el.scrollHeight + 'px';
    },

    autoResizeAll() {
      document.querySelectorAll('textarea.auto-resize-textarea').forEach(el => {
        el.style.height = 'auto';
        el.style.height = el.scrollHeight + 'px';
      });
    },

    deleteBlock(id) {
      this.blocks = this.blocks.filter(b => b.id !== id);
    },

    duplicateBlock(id) {
      const idx = this.blocks.findIndex(b => b.id === id);
      if (idx === -1) return;
      const original = this.blocks[idx];
      const clone = { ...original, id: 'blk_' + Math.random().toString(36).substr(2, 9) };
      this.blocks.splice(idx + 1, 0, clone);
    },

    async regenerateBlock(id) {
      this.regenerating = id;
      this.error = '';
      this.success = '';
      try {
        const response = await fetch(urlPostEditorRegenerate, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            block_id: id,
            blocks: this.blocks,
            title: this.title,
            summary: this.summary
          })
        });
        const data = await response.json();
        if (data.status === 'success') {
          const block = this.blocks.find(b => b.id === id);
          if (block) block.text = data.data.text;
          this.success = 'Block regenerated';
        } else {
          this.error = data.message || 'Failed to regenerate block.';
        }
      } catch (err) {
        this.error = 'Failed to regenerate block.';
      } finally {
        this.regenerating = null;
      }
    },

    async save(action) {
      if (action === 'apply') {
        this.savingApply = true;
      } else {
        this.savingClose = true;
      }
      this.error = '';
      this.success = '';

      try {
        const response = await fetch(urlPostEditorSave, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            action: action,
            title: this.title.trim(),
            summary: this.summary.trim(),
            blocks: this.blocks
          })
        });
        const data = await response.json();

        if (data.status === 'success') {
          if (action === 'save_close' && data.data?.redirect_url) {
            this.success = 'Post saved successfully';
            setTimeout(() => window.location.href = data.data.redirect_url, 1500);
          } else {
            this.success = data.message || 'Changes applied successfully';
          }
        } else {
          this.error = data.message || 'Failed to save post.';
        }
      } catch (err) {
        this.error = 'Failed to save post.';
      } finally {
        this.savingApply = false;
        this.savingClose = false;
      }
    }
  }
};

document.addEventListener('DOMContentLoaded', () => {
  const el = document.getElementById('post-editor-app');
  if (el) createApp(PostEditorApp).mount('#post-editor-app');
});
