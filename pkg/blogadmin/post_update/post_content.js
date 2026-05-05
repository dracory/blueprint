const { createApp } = Vue;

const PostContentApp = {
  data() {
    return {
      loading: true,
      saving: false,
      postId: '',
      editor: '',
      summernote: null,
      easyMDE: null,
      codeMirror: null,
      codeMirrorHTML: null,
      blockArea: null,
      form: {
        title: '',
        summary: '',
        content: ''
      }
    };
  },

  mounted() {
    console.log('PostContentApp mounted');
    // Initialize postId from global variable
    if (typeof postId !== 'undefined') {
      this.postId = postId;
      console.log('PostId set from global variable:', this.postId);
    }
    this.loadContent();
  },

  beforeUnmount() {
    console.log('PostContentApp beforeUnmount - cleaning up editors');
    this.cleanupEditors();
  },

  methods: {
    /**
     * Loads post content from the server and initializes the appropriate editor
     */
    async loadContent() {
      console.log('loadContent called for post:', this.postId);
      this.loading = true;
      try {
        const response = await fetch(urlContentLoad, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ post_id: this.postId })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.form.title = data.data?.title || '';
          this.form.summary = data.data?.summary || '';
          this.form.content = data.data?.content || '';
          // Normalize editor value to lowercase
          const newEditor = (data.data?.editor || '').toLowerCase();
          console.log('Loaded editor from server:', newEditor);

          // Cleanup previous editor if switching
          if (this.editor && this.editor !== newEditor) {
            console.log('Switching editor from', this.editor, 'to', newEditor);
            this.cleanupEditors();
          }
          this.editor = newEditor;

          // Initialize editor after data is loaded
          this.$nextTick(() => {
            this.initializeEditor();
          });
        } else {
          Swal.fire({
            icon: 'error',
            title: 'Error',
            text: data.message || 'Failed to load content'
          });
        }
      } catch (error) {
        console.error('Error loading content:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to load content'
        });
      } finally {
        this.loading = false;
      }
    },

    /**
     * Saves post content to the server
     */
    async saveContent() {
      console.log('saveContent called');
      if (this.saving) return;

      if (this.form.title === '') {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Title is required'
        });
        return;
      }

      // Sync content from editor if needed
      console.log('Syncing content from editor before save...');
      this.syncContentFromEditor();
      console.log('Content synced, current content length:', this.form.content?.length);

      this.saving = true;
      try {
        const response = await fetch(urlContentSave, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            post_id: this.postId,
            post_title: this.form.title,
            post_summary: this.form.summary,
            post_content: this.form.content
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
            text: data.message || 'Failed to save content'
          });
        }
      } catch (error) {
        console.error('Error saving content:', error);
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Failed to save content'
        });
      } finally {
        this.saving = false;
      }
    },

    /**
     * Initializes the appropriate editor based on the editor type
     */
    initializeEditor() {
      console.log('initializeEditor called for editor:', this.editor);

      switch (this.editor) {
        case 'htmlarea':
          this.initializeEditorSummernote();
          break;
        case 'html_codemirror':
          this.initializeEditorCodeMirrorHTML();
          break;
        case 'markdown_easymde':
          this.initializeEditorEasyMDE();
          break;
        case 'markdown_codemirror':
          this.initializeEditorCodeMirror();
          break;
        case 'markdown':
          this.initializeEditorMarkdown();
          break;
        case 'blockarea':
          this.initializeEditorBlockArea();
          break;
        case 'blockeditor':
          this.initializeEditorBlockEditor();
          break;
        case 'textarea':
          this.initializeEditorTextarea();
          break;
        default:
          console.log('Unknown editor type:', this.editor);
      }
    },

    /**
     * Initializes Summernote WYSIWYG editor (for htmlarea editor type)
     */
    initializeEditorSummernote() {
      console.log('Initializing HTML Area (Summernote) editor...');
      if (typeof $ === 'undefined' || typeof $.fn.summernote === 'undefined') {
        console.log('Summernote not loaded yet, retrying in 300ms...');
        setTimeout(() => this.initializeEditorSummernote(), 300);
        return;
      }

      try {
        const summernoteContainer = document.getElementById('summernote-container');
        if (!summernoteContainer) {
          console.log('Summernote container not found, retrying in 300ms...');
          setTimeout(() => this.initializeEditorSummernote(), 300);
          return;
        }

        $(summernoteContainer).summernote({
          height: 400,
          placeholder: 'Compose your post content...',
          callbacks: {
            onChange: (contents) => {
              this.form.content = contents;
            }
          }
        });

        // Set initial content
        if (this.form.content) {
          $(summernoteContainer).summernote('code', this.form.content);
        }

        console.log('HTML Area (Summernote) initialized successfully');
      } catch (error) {
        console.error('Error initializing HTML Area (Summernote):', error);
      }
    },

    /**
     * Initializes EasyMDE markdown editor
     */
    initializeEditorEasyMDE() {
      console.log('Initializing EasyMDE...');
      const textAreaId = 'easymde-content';
      const textArea = document.getElementById(textAreaId);

      if (!textArea) {
        console.log('EasyMDE textarea not found, ID:', textAreaId, 'retrying in 200ms...');
        setTimeout(() => this.initializeEditorEasyMDE(), 200);
        return;
      }

      this.initializeEasyMDEInstance(textArea);
    },

    /**
     * Initializes CodeMirror markdown editor
     */
    initializeEditorCodeMirror() {
      console.log('Initializing CodeMirror...');
      const textAreaId = 'codemirror-content';
      const textArea = document.getElementById(textAreaId);

      if (!textArea) {
        console.log('CodeMirror textarea not found, ID:', textAreaId, 'retrying in 200ms...');
        setTimeout(() => this.initializeEditorCodeMirror(), 200);
        return;
      }

      if (typeof CodeMirror === 'undefined') {
        console.log('CodeMirror not loaded yet, retrying in 300ms...');
        setTimeout(() => this.initializeEditorCodeMirror(), 300);
        return;
      }

      setTimeout(() => {
        try {
          // Avoid re-initializing if already attached
          if (this.codeMirror) {
            console.log('CodeMirror already initialized, skipping');
            return;
          }

          console.log('Creating new CodeMirror instance...');
          const editor = CodeMirror.fromTextArea(textArea, {
            mode: 'markdown',
            lineNumbers: true,
            lineWrapping: true,
            theme: 'default'
          });
          this.codeMirror = editor;

          // Set initial content
          if (this.form.content) {
            console.log('Setting initial content for CodeMirror, length:', this.form.content.length);
            editor.setValue(this.form.content);
          }

          // Sync changes back to form
          editor.on('change', () => {
            this.form.content = editor.getValue();
          });

          console.log('CodeMirror initialized successfully');
        } catch (error) {
          console.error('Error initializing CodeMirror:', error);
        }
      }, 100);
    },

    /**
     * Initializes CodeMirror HTML editor
     */
    initializeEditorCodeMirrorHTML() {
      console.log('Initializing CodeMirror HTML...');
      const textAreaId = 'codemirror-html-content';
      const textArea = document.getElementById(textAreaId);

      if (!textArea) {
        console.log('CodeMirror HTML textarea not found, ID:', textAreaId, 'retrying in 200ms...');
        setTimeout(() => this.initializeEditorCodeMirrorHTML(), 200);
        return;
      }

      if (typeof CodeMirror === 'undefined') {
        console.log('CodeMirror not loaded yet, retrying in 300ms...');
        setTimeout(() => this.initializeEditorCodeMirrorHTML(), 300);
        return;
      }

      setTimeout(() => {
        try {
          // Avoid re-initializing if already attached
          if (this.codeMirrorHTML) {
            console.log('CodeMirror HTML already initialized, skipping');
            return;
          }

          console.log('Creating new CodeMirror HTML instance...');
          const editor = CodeMirror.fromTextArea(textArea, {
            mode: 'xml',
            lineNumbers: true,
            lineWrapping: true,
            theme: 'default'
          });
          this.codeMirrorHTML = editor;

          // Set initial content
          if (this.form.content) {
            console.log('Setting initial content for CodeMirror HTML, length:', this.form.content.length);
            editor.setValue(this.form.content);
          }

          // Sync changes back to form
          editor.on('change', () => {
            this.form.content = editor.getValue();
          });

          console.log('CodeMirror HTML initialized successfully');
        } catch (error) {
          console.error('Error initializing CodeMirror HTML:', error);
        }
      }, 100);
    },

    /**
     * Initializes plain textarea with auto-resize for markdown content
     */
    initializeEditorMarkdown() {
      console.log('Initializing Markdown editor with auto-resize...');
      const textAreaId = 'markdown-content';
      const textArea = document.getElementById(textAreaId);

      if (!textArea) {
        console.log('Markdown textarea not found, ID:', textAreaId, 'retrying in 200ms...');
        setTimeout(() => this.initializeEditorMarkdown(), 200);
        return;
      }

      setTimeout(() => {
        // Set initial content if available
        if (this.form.content) {
          textArea.value = this.form.content;
          console.log('Set initial content for Markdown editor');
        }

        const autoResize = () => {
          textArea.style.height = 'auto';
          textArea.style.height = textArea.scrollHeight + 'px';
        };
        autoResize();
        textArea.addEventListener('input', autoResize);
        console.log('Markdown editor initialized successfully');
      }, 100);
    },

    /**
     * Initializes BlockArea editor
     */
    initializeEditorBlockArea() {
      console.log('Initializing BlockArea...');
      const textAreaId = 'blockarea-content';
      const textArea = document.getElementById(textAreaId);

      if (!textArea) {
        console.log('BlockArea textarea not found, ID:', textAreaId, 'retrying in 200ms...');
        setTimeout(() => this.initializeEditorBlockArea(), 200);
        return;
      }

      if (typeof BlockArea === 'undefined') {
        console.log('BlockArea not loaded yet, retrying in 300ms...');
        setTimeout(() => this.initializeEditorBlockArea(), 300);
        return;
      }

      setTimeout(() => {
        const blockArea = new BlockArea(textArea.id);
        blockArea.setParentId(this.postId);
        blockArea.registerBlock(BlockAreaHeading);
        blockArea.registerBlock(BlockAreaText);
        blockArea.registerBlock(BlockAreaImage);
        blockArea.registerBlock(BlockAreaCode);
        blockArea.registerBlock(BlockAreaRawHtml);

        // Set initial content if available
        if (this.form.content) {
          textArea.value = this.form.content;
        }

        blockArea.init();
        this.blockArea = blockArea;

        // Sync changes back to form
        blockArea.on('change', () => {
          this.form.content = textArea.value;
        });
        console.log('BlockArea initialized successfully');
      }, 500);
    },

    /**
     * Initializes BlockEditor
     */
    initializeEditorBlockEditor() {
      console.log('Initializing BlockEditor...');
      if (typeof BlockEditor === 'undefined') {
        console.log('BlockEditor not loaded yet, retrying in 300ms...');
        setTimeout(() => this.initializeEditorBlockEditor(), 300);
        return;
      }

      const value = this.form.content || '[]';
      const be = new BlockEditor({
        name: 'post_content',
        value: value,
        handleEndpoint: urlBlockEditorHandle,
        blockDefinitions: window.BlockEditorDefinitions || []
      });
      be.init();
      console.log('BlockEditor initialized successfully');
    },

    /**
     * Initializes plain textarea editor
     */
    initializeEditorTextarea() {
      console.log('Initializing plain textarea editor...');
      const textArea = document.querySelector('textarea[name="post_content"]');

      if (!textArea) {
        console.log('Textarea not found, retrying in 200ms...');
        setTimeout(() => this.initializeEditorTextarea(), 200);
        return;
      }

      // Set initial content if available
      if (this.form.content) {
        textArea.value = this.form.content;
      }

      console.log('Plain textarea editor initialized successfully');
    },

    /**
     * Cleans up all editor instances to prevent memory leaks
     */
    cleanupEditors() {
      console.log('cleanupEditors called');
      // Cleanup Summernote (HTML Area)
      if (this.summernote && typeof $ !== 'undefined') {
        console.log('Cleaning up HTML Area (Summernote)...');
        try {
          $('#summernote-container').summernote('destroy');
          console.log('HTML Area (Summernote) cleaned up successfully');
        } catch (e) {
          console.error('Error destroying HTML Area (Summernote):', e);
        }
        this.summernote = null;
      }

      // Cleanup EasyMDE
      if (this.easyMDE) {
        console.log('Cleaning up EasyMDE...');
        try {
          this.easyMDE.toTextArea();
          console.log('EasyMDE cleaned up successfully');
          this.easyMDE = null;
        } catch (e) {
          console.error('Error destroying EasyMDE:', e);
        }
      }

      // Cleanup CodeMirror (Markdown)
      if (this.codeMirror) {
        console.log('Cleaning up CodeMirror (Markdown)...');
        try {
          this.codeMirror.toTextArea();
          console.log('CodeMirror (Markdown) cleaned up successfully');
          this.codeMirror = null;
        } catch (e) {
          console.error('Error destroying CodeMirror (Markdown):', e);
        }
      }

      // Cleanup CodeMirror (HTML)
      if (this.codeMirrorHTML) {
        console.log('Cleaning up CodeMirror (HTML)...');
        try {
          this.codeMirrorHTML.toTextArea();
          console.log('CodeMirror (HTML) cleaned up successfully');
          this.codeMirrorHTML = null;
        } catch (e) {
          console.error('Error destroying CodeMirror (HTML):', e);
        }
      }

      // Cleanup BlockArea
      if (this.blockArea) {
        console.log('Cleaning up BlockArea...');
        this.blockArea = null;
        console.log('BlockArea cleaned up successfully');
      }
    },

    /**
     * Creates and configures the EasyMDE instance
     * @param {HTMLTextAreaElement} textArea - The textarea element to attach EasyMDE to
     */
    initializeEasyMDEInstance(textArea) {
      console.log('initializeEasyMDEInstance called with textarea:', textArea?.id);
      setTimeout(() => {
        if (typeof EasyMDE === 'undefined') {
          console.log('EasyMDE not loaded yet, retrying in 300ms...');
          setTimeout(() => this.initializeEasyMDEInstance(textArea), 300);
          return;
        }

        try {
          // Avoid re-initializing if already attached
          if (this.easyMDE) {
            console.log('EasyMDE already initialized, skipping');
            return;
          }

          console.log('Creating new EasyMDE instance...');
          const easyMDE = new EasyMDE({
            element: textArea,
          });
          this.easyMDE = easyMDE;

          // Set initial content
          if (this.form.content) {
            console.log('Setting initial content for EasyMDE, length:', this.form.content.length);
            easyMDE.value(this.form.content);
          }

          // Keep textarea value in sync with EasyMDE content
          const form = textArea.closest('form');
          const syncToTextarea = () => {
            if (!this.easyMDE) return;
            textArea.value = this.easyMDE.value();
            this.form.content = textArea.value;
          };

          // Sync on editor changes
          easyMDE.codemirror.on('change', syncToTextarea);

          // Ensure sync right before form submit
          if (form) {
            form.addEventListener('submit', syncToTextarea);
          }

          console.log('EasyMDE initialized successfully');
        } catch (error) {
          console.error('Error initializing EasyMDE:', error);
        }
      }, 100);
    },

    /**
     * Syncs content from the active editor to the form data
     */
    syncContentFromEditor() {
      console.log('syncContentFromEditor called for editor:', this.editor);

      // Sync from Summernote
      if (this.editor === 'htmlarea' && typeof $ !== 'undefined') {
        console.log('Syncing from HTML Area (Summernote)...');
        const snContainer = $('#summernote-container');
        if (snContainer.length > 0 && typeof $.fn.summernote !== 'undefined') {
          this.form.content = snContainer.summernote('code');
          console.log('HTML Area (Summernote) content synced, length:', this.form.content?.length);
        }
      }

      // Sync from EasyMDE
      if (this.editor === 'markdown_easymde' && this.easyMDE) {
        console.log('Syncing from EasyMDE...');
        this.form.content = this.easyMDE.value();
        console.log('EasyMDE content synced, length:', this.form.content?.length);
      }

      // Sync from CodeMirror (Markdown)
      if (this.editor === 'markdown_codemirror' && this.codeMirror) {
        console.log('Syncing from CodeMirror (Markdown)...');
        this.form.content = this.codeMirror.getValue();
        console.log('CodeMirror (Markdown) content synced, length:', this.form.content?.length);
      }

      // Sync from CodeMirror (HTML)
      if (this.editor === 'html_codemirror' && this.codeMirrorHTML) {
        console.log('Syncing from CodeMirror (HTML)...');
        this.form.content = this.codeMirrorHTML.getValue();
        console.log('CodeMirror (HTML) content synced, length:', this.form.content?.length);
      }

      // Sync from textareas
      if (['textarea', 'markdown', 'blockarea'].includes(this.editor)) {
        console.log('Syncing from textarea for editor:', this.editor);
        let textArea = null;
        if (this.editor === 'blockarea') {
          textArea = document.getElementById('blockarea-content');
        } else if (this.editor === 'markdown') {
          textArea = document.getElementById('markdown-content');
        } else {
          textArea = document.querySelector('textarea[name="post_content"]');
        }
        if (textArea) {
          this.form.content = textArea.value;
          console.log('Textarea content synced, length:', this.form.content?.length);
        }
      }

      // Sync from BlockEditor
      if (this.editor === 'blockeditor') {
        console.log('BlockEditor handles sync internally');
        // BlockEditor should handle this internally
      }
    }
  }
};

// Mount the app immediately
const el = document.getElementById('post-content-app');
if (el) {
  console.log('Mounting PostContentApp on #post-content-app');
  createApp(PostContentApp).mount('#post-content-app');
  console.log('PostContentApp mounted successfully');
} else {
  console.error('Could not find #post-content-app element');
}
