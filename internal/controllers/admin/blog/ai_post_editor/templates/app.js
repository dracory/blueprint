const postEditorApp = {
    data() {
        return {
            loading: true,
            saving: false,
            loadingSections: {
                introduction: false,
                conclusion: false,
                summary: false,
                metas: false
            },
            loadingImage: false,
            post: null, // will be loaded dynamically
        }
    },
    async mounted() {
        this.loading = false;
        await this.loadPost();
    },
    methods: {
        /**
         * Loads the blog post data from the backend and sets up the editor state.
         * Shows an error if the load fails.
         * @returns {Promise<void>}
         */
        async loadPost() {
            try {
                const data = await apiPost({
                     action: "load_post",
                     ts: new Date().getTime(),
                });
                if (data.status !== 'success') {
                    this.showError(data.message || data.error || 'Failed to load post');
                    return;
                }
                this.post = data.data;
            } catch (error) {
                this.showError('Failed to load post: ' + error);
            }
        },

        /**
         * Regenerates a section (introduction, conclusion, or dynamic section) via the backend.
         * Updates the section in-place on success.
         * @param {string} section - Section identifier ('introduction', 'conclusion', or 'section_{index}')
         * @returns {Promise<void>}
         */
        async regenerateSection(section) {
            // Set loading state for this section
            if (section.startsWith('section_')) {
                const index = section.split('_')[1];
                this.loadingSections[section] = true;
            } else {
                this.loadingSections[section] = true;
            }

            try {
                const data = await apiPost({
                    action: 'regenerate_section',
                    section: section,
                    id: this.post.id
                });
                if (data.status !== 'success') {
                    this.showError(data.message || data.error || 'Failed to regenerate section');
                    return;
                }
                const sectionObj = data.data && data.data.section ? data.data.section : null;
                if (!sectionObj) {
                    this.showError('No section returned from server');
                    return;
                }
                // Update the section content
                if (section === 'introduction') {
                    this.post.introduction = sectionObj;
                } else if (section === 'conclusion') {
                    this.post.conclusion = sectionObj;
                } else if (section.startsWith('section_')) {
                    const index = parseInt(section.split('_')[1]);
                    if (index >= 0 && index < this.post.sections.length) {
                        this.post.sections[index] = sectionObj;
                    }
                }
            } catch (error) {
                this.showError('Failed to regenerate section: ' + error);
            } finally {
                // Clear loading state
                this.loadingSections[section] = false;
            }
        },
        
        /**
         * Regenerates the featured image for the post by requesting the backend.
         * Updates the imageUrl on success.
         * @returns {Promise<void>}
         */
        async regenerateImage() {
            this.loadingImage = true;
            try {
                const data = await apiPost({
                    action: 'regenerate_image',
                    id: this.post.id
                });
                if (data.status !== 'success') {
                    this.showError(data.message || data.error || 'Failed to generate image');
                    return;
                }
                this.post.image = data.data.image;
            } catch (error) {
                this.showError('Failed to generate image: ' + error);
            } finally {
                this.loadingImage = false;
            }
        },
        
        /**
         * Saves the blog post as a finalized post by sending all data to the backend.
         * Shows an error if saving fails.
         * @returns {Promise<void>}
         */
        async savePost() {
            this.saving = true;
            try {
                const response = await apiPost({
                    action: 'create_final_post',
                    id: this.post.id,
                    post: JSON.stringify(this.post),
                });
                if (response.status !== 'success') {
                    this.showError(response.message || response.error || 'Failed to save post');
                    return;
                }
                // If HTML is returned (legacy), insert it
                if (typeof response === 'string') {
                    document.body.insertAdjacentHTML('beforeend', response);
                } else if (response.data && response.data.redirect) {
                    // Handle redirect URL from JSON response
                    window.location.href = response.data.redirect;
                } else {
                    // Show success message if no redirect
                    this.showSuccess('Post saved successfully!');
                }
            } catch (error) {
                this.showError('Failed to save post: ' + error);
            } finally {
                this.saving = false;
            }
        },

        /**
         * Saves the blog post record as a draft by sending all data to the backend.
         * Shows an error if saving fails.
         * @returns {Promise<void>}
         */
        async saveDraft() {
            try {
                const response = await apiPost({
                    action: 'save_draft',
                    id: this.post.id,
                    post: JSON.stringify(this.post),
                });
                if (response.status !== 'success') {
                    this.showError(response.message || response.error || 'Failed to save draft');
                    return;
                }
                if (response.data && response.data.redirect) {
                    window.location = response.data.redirect;
                } else {
                    this.showSuccess('Draft saved successfully!');
                }
            } catch (error) {
                this.showError('Failed to save draft: ' + error);
            }
        },

        /**
         * Shows an error message using SweetAlert2.
         * @param {string} message - The error message to display.
         */
        showError(message) {
            Swal.fire({
                title: 'Error',
                text: message,
                icon: 'error',
                confirmButtonText: 'OK'
            });
        },

        /**
         * Shows a success message using SweetAlert2.
         * @param {string} message - The success message to display.
         */
        showSuccess(message) {
            Swal.fire({
                title: 'Success',
                text: message,
                icon: 'success',
                confirmButtonText: 'OK'
            });
        },

        /**
         * Adds a new empty paragraph to the specified section.
         * Handles migration from old format (content -> paragraphs).
         * @param {string} section - Section type ('introduction', 'conclusion', 'section')
         * @param {number} sectionIndex - Index of the section (for dynamic sections)
         */
        addParagraph(section, sectionIndex) {
            if (section === 'introduction') {
                if (!this.post.introduction.paragraphs) {
                    // Handle migration from old format
                    if (this.post.introduction.content) {
                        this.post.introduction.paragraphs = [this.post.introduction.content];
                        delete this.post.introduction.content;
                    } else {
                        this.post.introduction.paragraphs = [];
                    }
                }
                this.post.introduction.paragraphs.push('');
            } else if (section === 'conclusion') {
                if (!this.post.conclusion.paragraphs) {
                    // Handle migration from old format
                    if (this.post.conclusion.content) {
                        this.post.conclusion.paragraphs = [this.post.conclusion.content];
                        delete this.post.conclusion.content;
                    } else {
                        this.post.conclusion.paragraphs = [];
                    }
                }
                this.post.conclusion.paragraphs.push('');
            } else if (section === 'section') {
                if (!this.post.sections[sectionIndex].paragraphs) {
                    // Handle migration from old format
                    if (this.post.sections[sectionIndex].content) {
                        this.post.sections[sectionIndex].paragraphs = [this.post.sections[sectionIndex].content];
                        delete this.post.sections[sectionIndex].content;
                    } else {
                        this.post.sections[sectionIndex].paragraphs = [];
                    }
                }
                this.post.sections[sectionIndex].paragraphs.push('');
            }
        },

        /**
         * Regenerates a single paragraph in a section by requesting the backend.
         * Updates the paragraph in-place on success.
         * @param {string} sectionType - Section type ('introduction', 'conclusion', 'section')
         * @param {number|null} sectionIndex - Index of the section (for dynamic sections)
         * @param {number} pIndex - Index of the paragraph to regenerate
         */
        async regenerateParagraph(sectionType, sectionIndex, pIndex) {
            try {
                const payload = {
                    id: this.post.id,
                    action: 'regenerate_paragraph',
                    section_type: sectionType,
                    paragraph_index: pIndex
                };
                if (sectionIndex !== null && sectionIndex !== undefined) payload.section_index = sectionIndex;
                const data = await apiPost(payload);
                if (data.status !== 'success') {
                    this.showError(data.message || data.error || 'Failed to regenerate paragraph');
                    return;
                }
                const newParagraph = data.data && data.data.paragraph ? data.data.paragraph : null;
                if (!newParagraph) {
                    this.showError('No paragraph returned from server');
                    return;
                }
                if (sectionType === 'introduction') {
                    this.post.introduction.paragraphs.splice(pIndex, 1, newParagraph);
                } else if (sectionType === 'conclusion') {
                    this.post.conclusion.paragraphs.splice(pIndex, 1, newParagraph);
                } else if (sectionType === 'section') {
                    this.post.sections[sectionIndex].paragraphs.splice(pIndex, 1, newParagraph);
                }
            } catch (error) {
                this.showError('Failed to regenerate paragraph: ' + error);
            }
        },
        /**
         * Removes a paragraph from a section.
         * @param {string} sectionType - Section type ('introduction', 'conclusion', 'section')
         * @param {number|null} sectionIndex - Index of the section (for dynamic sections)
         * @param {number} pIndex - Index of the paragraph to remove
         */
        removeParagraph(sectionType, sectionIndex, pIndex) {
            if (sectionType === 'introduction') {
                this.post.introduction.paragraphs.splice(pIndex, 1);
            } else if (sectionType === 'conclusion') {
                this.post.conclusion.paragraphs.splice(pIndex, 1);
            } else if (sectionType === 'section') {
                this.post.sections[sectionIndex].paragraphs.splice(pIndex, 1);
            }
        },

        /**
         * Regenerates the summary for the post using the backend.
         * Updates the summary field on success.
         * @returns {Promise<void>}
         */
        async regenerateSummary() {
            this.loadingSections.summary = true;
            try {
                const data = await apiPost({
                    action: 'regenerate_summary',
                    id: this.post.id
                });
                if (data.status !== 'success') {
                    this.showError(data.message || data.error || 'Failed to regenerate summary');
                    return;
                }
                this.post.summary = data.data && data.data.summary ? data.data.summary : '';
            } catch (error) {
                this.showError('Failed to regenerate summary: ' + error);
            } finally {
                this.loadingSections.summary = false;
            }
        },

        /**
         * Regenerates the meta title, description, and keywords for the post using the backend.
         * Updates the meta fields on success.
         * @returns {Promise<void>}
         */
        async regenerateMetas() {
            this.loadingSections.metas = true;
            try {
                const data = await apiPost({
                    action: 'regenerate_metas',
                    id: this.post.id
                });
                if (data.status !== 'success') {
                    this.showError(data.message || data.error || 'Failed to regenerate meta information');
                    return;
                }
                if (data.data) {
                    this.post.metaTitle = data.data.metaTitle || '';
                    this.post.metaDescription = data.data.metaDescription || '';
                    // Convert comma-separated keywords string to array
                    this.post.metaKeywords = data.data.metaKeywords
                        ? data.data.metaKeywords.split(',').map(k => k.trim()).filter(Boolean)
                        : [];
                }
            } catch (error) {
                this.showError('Failed to regenerate meta information: ' + error);
            } finally {
                this.loadingSections.metas = false;
            }
        },

        /**
         * Handler for the Regenerate Summary button. Calls regenerateSummary and shows a success message.
         * @returns {Promise<void>}
         */
        async onRegenerateSummary() {
            await this.regenerateSummary();
            this.showSuccess('Summary regenerated successfully!');
        },

        /**
         * Handler for the Regenerate Metas button. Calls regenerateMetas and shows a success message.
         * @returns {Promise<void>}
         */
        async onRegenerateMetas() {
            await this.regenerateMetas();
            this.showSuccess('Meta information regenerated successfully!');
        },
    },
    computed: {
        keywordsString: {
            get() {
                return (this.post && this.post.metaKeywords) ? this.post.metaKeywords.join(', ') : '';
            },
            set(val) {
                if (this.post) {
                    this.post.metaKeywords = val.split(',').map(k => k.trim()).filter(Boolean);
                }
            }
        }
    }
};

// Initialize Vue app after the page is loaded
document.addEventListener('DOMContentLoaded', function() {
    Vue.createApp(postEditorApp).mount('#post-editor-app');
});

async function apiPost(body) {
    body.id = '{{ id }}'

    params = new URLSearchParams(body);
    const response = await fetch('{{ url }}', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: params
    });

    return await response.json();
}