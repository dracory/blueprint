package aipostcontentupdate

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/llm"
	"github.com/dracory/uid"
	"github.com/samber/lo"
)

// == COMPONENT ===============================================================

type formAiPostContentUpdate struct {
	liveflux.Base
	App        types.RegistryInterface
	Post       *blogstore.Post
	Blocks     []Block
	Error      string
	Success    string
	RedirectTo string
}

const (
	actionReorder         = "reorder"
	actionRegenerateBlock = "regenerate_block"
	actionApply           = "apply"
	actionSaveClose       = "save_close"
	actionDeleteBlock     = "delete_block"
	actionDuplicateBlock  = "duplicate_block"
)

// == CONSTRUCTOR =============================================================

func NewFormAiPostContentUpdate(app types.RegistryInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&formAiPostContentUpdate{})
	if err != nil {
		log.Println(err)
		return nil
	}
	if c, ok := inst.(*formAiPostContentUpdate); ok {
		c.App = app
	}
	return inst
}

// == PUBLIC METHODS ==========================================================

func (c *formAiPostContentUpdate) GetKind() string {
	return "admin_blog_ai_post_content_update_form_component"
}

func (c *formAiPostContentUpdate) Mount(ctx context.Context, params map[string]string) error {
	if c.App == nil {
		c.Error = "Application not initialized"
		return nil
	}

	postID := strings.TrimSpace(params["post_id"])
	if postID == "" {
		c.Error = "Post ID is required"
		return nil
	}

	blogStore := c.App.GetBlogStore()
	if blogStore == nil {
		c.Error = "Blog store is not configured"
		return nil
	}

	post, err := blogStore.PostFindByID(ctx, postID)
	if err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("AI content editor: failed to load post", "error", err.Error())
		}
		c.Error = "Post not found"
		return nil
	}
	if post == nil {
		c.Error = "Post not found"
		return nil
	}

	c.Post = post
	c.Blocks = MarkdownToBlocks(post.Content())
	c.Error = ""

	return nil
}

func (c *formAiPostContentUpdate) onDeleteBlock(ctx context.Context, data url.Values) error {
	// Ensure server-side order matches client-side sorted order
	c.Blocks = c.orderedBlocksFromData(data)
	id := strings.TrimSpace(data.Get("block_id"))
	if id == "" {
		return nil
	}
	var blocks []Block
	for _, b := range c.Blocks {
		if b.ID == id {
			continue
		}
		blocks = append(blocks, b)
	}
	c.Blocks = blocks
	c.Success = "Block deleted"
	c.Error = ""
	c.RedirectTo = ""
	return nil
}

func (c *formAiPostContentUpdate) onDuplicateBlock(ctx context.Context, data url.Values) error {
	// Ensure server-side order matches client-side sorted order
	c.Blocks = c.orderedBlocksFromData(data)
	id := strings.TrimSpace(data.Get("block_id"))
	if id == "" {
		return nil
	}
	var blocks []Block
	for _, b := range c.Blocks {
		blocks = append(blocks, b)
		if b.ID == id {
			clone := b
			clone.ID = uid.HumanUid()
			blocks = append(blocks, clone)
		}
	}
	c.Blocks = blocks
	c.Success = "Block duplicated"
	c.Error = ""
	c.RedirectTo = ""
	return nil
}

// orderedBlocksFromData returns c.Blocks reordered according to the
// comma-separated block_order (list of block IDs) present in data.
// If block_order is empty or invalid, it falls back to the current order.
func (c *formAiPostContentUpdate) orderedBlocksFromData(data url.Values) []Block {
	if data == nil {
		return c.Blocks
	}

	orderCSV := strings.TrimSpace(data.Get("block_order"))
	if orderCSV == "" {
		return c.Blocks
	}

	idOrder := strings.Split(orderCSV, ",")
	byID := make(map[string]Block, len(c.Blocks))
	for _, b := range c.Blocks {
		byID[b.ID] = b
	}

	var ordered []Block
	for _, id := range idOrder {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if b, ok := byID[id]; ok {
			ordered = append(ordered, b)
		}
	}

	if len(ordered) == 0 {
		return c.Blocks
	}

	return ordered
}

func (c *formAiPostContentUpdate) Handle(ctx context.Context, action string, data url.Values) error {
	if data == nil {
		data = url.Values{}
	}

	switch action {
	case actionReorder:
		// Update in-memory order only; no persistence
		c.Blocks = c.orderedBlocksFromData(data)
		return nil
	case actionRegenerateBlock:
		return c.onRegenerateBlock(ctx, data)
	case actionApply:
		return c.onSave(ctx, "Changes applied successfully", false, data)
	case actionSaveClose:
		return c.onSave(ctx, "Post saved successfully", true, data)
	case actionDeleteBlock:
		return c.onDeleteBlock(ctx, data)
	case actionDuplicateBlock:
		return c.onDuplicateBlock(ctx, data)
	default:
		return nil
	}
}

func (c *formAiPostContentUpdate) onRegenerateBlock(ctx context.Context, data url.Values) error {
	// Ensure blocks are in the same order as the client
	c.Blocks = c.orderedBlocksFromData(data)
	id := strings.TrimSpace(data.Get("block_id"))
	if id == "" {
		return nil
	}

	_, idx, found := lo.FindIndexOf(c.Blocks, func(b Block) bool {
		return b.ID == id
	})
	if !found || idx < 0 || idx >= len(c.Blocks) {
		return nil
	}

	block := c.Blocks[idx]
	if strings.TrimSpace(block.Text) == "" {
		return nil
	}

	// Build full-post context where this block is replaced by a marker
	contextBlocks := make([]Block, len(c.Blocks))
	copy(contextBlocks, c.Blocks)
	contextBlocks[idx].Text = "== BLOCK TO REPLACE HERE =="
	contextMarkdown := BlocksToMarkdown(contextBlocks)

	systemPrompt := `You are an expert blog editor.
You will receive the FULL blog post as markdown, where exactly ONE block is replaced
with the marker "== BLOCK TO REPLACE HERE ==".

Your task is to REWRITE that missing block to improve clarity, style, and readability,
while preserving the original meaning and fitting naturally into the surrounding content.

IMPORTANT:
- You MUST significantly rephrase the original block. Do NOT return the same text.
- Do NOT copy long spans verbatim from the original block content.
- If the original block is a heading, return a short, strong heading.
- If it is a paragraph, return one or more paragraphs of body text.

Return ONLY the rewritten text for that block as markdown, with no additional explanations.`

	// userPrompt := fmt.Sprintf("Post title: %s\nBlock type: %s\n\nFull post markdown with marker:\n%s\n\nOriginal block content (for meaning only; DO NOT copy it verbatim):\n%s", c.Post.Title(), block.Type, contextMarkdown, block.Text)

	userPrompt := []string{}
	userPrompt = append(userPrompt, "Post title: "+c.Post.Title())
	userPrompt = append(userPrompt, "Block type: "+string(block.Type))
	userPrompt = append(userPrompt, "Full post markdown with marker:\n"+contextMarkdown)
	userPrompt = append(userPrompt, "Original block content to be regenerated (for meaning only; MUST BE REPHRASED SIGNIFICANTLY; DO NOT copy it verbatim):\n"+block.Text)

	fmt.Println("System Prompt:", systemPrompt)
	fmt.Println("User Prompt:", strings.Join(userPrompt, "\n\n"))

	engine, err := shared.LlmEngine(c.App)
	if err != nil || engine == nil {
		c.Error = "Failed to initialize LLM engine. Please try again later."
		c.Success = ""
		c.RedirectTo = ""
		return nil
	}
	resp, err := engine.Generate(systemPrompt, strings.Join(userPrompt, "\n\n"), llm.LlmOptions{
		MaxTokens:    512,
		Temperature:  0.7,
		OutputFormat: llm.OutputFormatText,
		// Verbose:      true,
	})
	if err != nil {
		c.Error = "Failed to regenerate block content. Please try again later."
		c.Success = ""
		c.RedirectTo = ""
		return nil
	}

	newText := strings.TrimSpace(resp)
	if newText == "" || newText == block.Text {
		c.Error = "LLM did not provide a meaningful rewrite for this block."
		c.Success = ""
		c.RedirectTo = ""
		return nil
	}

	c.Blocks[idx].Text = newText
	c.Error = ""
	c.Success = "Block regenerated"
	c.RedirectTo = ""

	return nil
}

func (c *formAiPostContentUpdate) Render(ctx context.Context) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Dashboard", URL: links.Admin().Home()},
		{Name: "Blog", URL: links.Admin().Blog()},
		{Name: "Post Manager", URL: shared.NewLinks().PostManager()},
		{Name: "Edit Post", URL: shared.NewLinks().PostUpdate(map[string]string{"post_id": c.PostID()})},
		{Name: "AI Content Editor"},
	})

	header := hb.Heading1().HTML("AI Post Content Editor")
	backButton := hb.A().
		Class("btn btn-secondary").
		Href(shared.NewLinks().PostUpdate(map[string]string{"post_id": c.PostID()})).
		HTML("← Back to Post")

	content := hb.Div()

	alerts := hb.Div()
	if c.Error != "" {
		alerts = alerts.Child(hb.SwalError(hb.SwalOptions{
			Text:             c.Error,
			Timer:            5000,
			TimerProgressBar: true,
			Position:         "top-end",
		}))
	}
	if c.Success != "" {
		// If RedirectTo is set, show success with redirect; otherwise just a toast
		if c.RedirectTo != "" {
			alerts = alerts.Child(hb.SwalSuccess(hb.SwalOptions{
				Text:             c.Success,
				RedirectURL:      c.RedirectTo,
				RedirectSeconds:  2,
				Timer:            2000,
				TimerProgressBar: true,
				Position:         "top-end",
			}))
		} else {
			alerts = alerts.Child(hb.SwalSuccess(hb.SwalOptions{
				Text:             c.Success,
				Timer:            3000,
				TimerProgressBar: true,
				Position:         "top-end",
			}))
		}
	}

	if c.Error != "" && c.Post == nil {
		return c.wrapPage(breadcrumbs, header, backButton, alerts)
	}

	// Title and summary fields
	titleInput := hb.Input().
		Class("form-control").
		Name("title").
		Value(c.Post.Title())

	summaryTextarea := hb.TextArea().
		Class("form-control auto-resize-textarea").
		Name("summary").
		Attr("rows", "3").
		Attr("oninput", "this.style.height='auto';this.style.height=this.scrollHeight+'px';").
		Text(c.Post.Summary())

	titleGroup := hb.Div().
		Class("mb-3").
		Child(hb.Label().Class("form-label fw-semibold").Text("Title")).
		Child(titleInput)

	summaryGroup := hb.Div().
		Class("mb-0").
		Child(hb.Label().Class("form-label fw-semibold").Text("Summary")).
		Child(summaryTextarea)

	titleCardHeader := hb.Div().
		Class("card-header d-flex justify-content-between align-items-center").
		Child(hb.Span().Class("fw-semibold").Text("Title & Summary"))

	titleCardBody := hb.Div().
		Class("card-body").
		Child(titleGroup).
		Child(summaryGroup)

	content = content.
		Child(hb.Div().
			Class("card mb-3").
			Child(titleCardHeader).
			Child(titleCardBody))

	// Blocks
	blocksContainer := hb.Div().ID("BlocksContainer")
	for _, block := range c.Blocks {
		blockName := fmt.Sprintf("block_%s", block.ID)
		blockTypeName := fmt.Sprintf("block_type_%s", block.ID)

		label := hb.Span().Class("fw-semibold")
		switch block.Type {
		case BlockHeader1:
			label = label.Text("Heading 1")
		case BlockHeader2:
			label = label.Text("Heading 2")
		case BlockCode:
			label = label.Text("Code Block")
		default:
			label = label.Text("Paragraph")
		}

		// hidden type field
		hiddenType := hb.Input().
			Type("hidden").
			Name(blockTypeName).
			Value(string(block.Type))

		var field hb.TagInterface
		if block.Type == BlockHeader1 || block.Type == BlockHeader2 {
			field = hb.Input().
				Class("form-control").
				Name(blockName).
				Value(block.Text)
		} else {
			rows := "3"
			if block.Type == BlockCode {
				rows = "6"
			}
			field = hb.TextArea().
				Class("form-control font-monospace auto-resize-textarea").
				Name(blockName).
				Attr("rows", rows).
				Attr("oninput", "this.style.height='auto';this.style.height=this.scrollHeight+'px';").
				Text(block.Text)
		}

		dragHandle := hb.Span().
			Class("block-drag-handle me-2").
			Style("cursor: move;").
			Child(hb.I().Class("bi bi-arrows-move"))

		// Per-block actions: regenerate, duplicate and delete
		regenerateButton := hb.Button().
			Type("submit").
			Class("btn btn-outline-info btn-sm").
			Name("block_id").
			Value(block.ID).
			Attr(liveflux.DataFluxAction, actionRegenerateBlock).
			Child(hb.I().Class("bi bi-stars"))

		duplicateButton := hb.Button().
			Type("submit").
			Class("btn btn-outline-secondary btn-sm").
			Name("block_id").
			Value(block.ID).
			Attr(liveflux.DataFluxAction, actionDuplicateBlock).
			Child(hb.I().Class("bi bi-files"))

		deleteButton := hb.Button().
			Type("submit").
			Class("btn btn-outline-danger btn-sm").
			Name("block_id").
			Value(block.ID).
			Attr(liveflux.DataFluxAction, actionDeleteBlock).
			Child(hb.I().Class("bi bi-trash"))

		headerActions := hb.Div().
			Class("btn-group btn-group-sm").
			Child(regenerateButton).
			Child(duplicateButton).
			Child(deleteButton)

		cardHeader := hb.Div().
			Class("card-header d-flex justify-content-between align-items-center").
			Child(hb.Div().Class("d-flex align-items-center gap-2").
				Child(dragHandle).
				Child(label)).
			Child(headerActions)

		cardBody := hb.Div().
			Class("card-body").
			Child(hiddenType).
			Child(field)

		card := hb.Div().
			Class("card mb-3 block-card").
			Attr("data-block-id", block.ID).
			Child(cardHeader).
			Child(cardBody)

		blocksContainer = blocksContainer.Child(card)
	}

	content = content.Child(blocksContainer)

	// Action buttons
	toolbar := hb.Div().Class("d-flex justify-content-end gap-2 mb-3").
		Child(hb.Button().
			Type("submit").
			Class("btn btn-outline-secondary").
			Attr(liveflux.DataFluxAction, actionApply).
			Attr(liveflux.DataFluxIndicator, ".apply-spinner-inner").
			Child(hb.I().Class("bi bi-check2 me-2")).
			Child(hb.Span().Text("Apply")).
			Child(hb.Div().
				Style("display: none").
				Class("apply-spinner-inner ms-2").
				Class("spinner-border spinner-border-sm text-secondary")),
		).
		Child(hb.Button().
			Type("submit").
			Class("btn btn-primary").
			Attr(liveflux.DataFluxAction, actionSaveClose).
			Attr(liveflux.DataFluxIndicator, ".save-close-spinner-inner").
			Child(hb.I().Class("bi bi-check2-all me-2")).
			Child(hb.Span().Text("Save & Close")).
			Child(hb.Div().
				Style("display: none").
				Class("save-close-spinner-inner ms-2").
				Class("spinner-border spinner-border-sm text-secondary")))

	orderInput := hb.Input().
		Type("hidden").
		ID("BlocksOrderInput").
		Name("block_order")

	reorderButton := hb.Button().
		Type("submit").
		ID("BlocksReorderAction").
		Style("display:none").
		Attr(liveflux.DataFluxAction, actionReorder)

	autoResizeScript := hb.Script(`
document.querySelectorAll('textarea.auto-resize-textarea').forEach(function(t) {
  t.style.height = 'auto';
  t.style.height = t.scrollHeight + 'px';
});
`)

	sortableScript := hb.Script(`
(function() {
  var container = document.getElementById('BlocksContainer');
  var orderInput = document.getElementById('BlocksOrderInput');
  var reorderButton = document.getElementById('BlocksReorderAction');
  if (!container || !orderInput || !window.Sortable) return;

  // Prevent multiple initializations on the same container
  if (container._livefluxSortableInitialized) {
    return;
  }
  container._livefluxSortableInitialized = true;

  function updateOrder() {
    var ids = [];
    container.querySelectorAll('.block-card').forEach(function(card) {
      var id = card.getAttribute('data-block-id');
      if (id) { ids.push(id); }
    });
    orderInput.value = ids.join(',');
  }

  updateOrder();

  Sortable.create(container, {
    animation: 150,
    handle: '.block-drag-handle',
    onSort: function() {
      updateOrder();
      if (reorderButton) {
        reorderButton.click();
      }
    }
  });
})();
`)

	form := hb.Form().
		Child(alerts).
		Child(toolbar).
		Child(content).
		Child(orderInput).
		Child(reorderButton).
		Child(autoResizeScript).
		Child(sortableScript)

	return c.wrapPage(breadcrumbs, header, backButton, form)
}

func (c *formAiPostContentUpdate) onSave(ctx context.Context, successMessage string, redirect bool, data url.Values) error {
	if c.App == nil || c.Post == nil {
		c.Error = "Post is not loaded"
		c.Success = ""
		return nil
	}

	// Update title and summary
	title := strings.TrimSpace(data.Get("title"))
	summary := strings.TrimSpace(data.Get("summary"))

	if title != "" {
		c.Post.SetTitle(title)
	}
	if summary != "" {
		c.Post.SetSummary(summary)
	}

	// Determine block order from Sortable and rebuild blocks from form
	ordered := c.orderedBlocksFromData(data)

	// Rebuild blocks from form in the new order using ID-based field names
	blocks := make([]Block, len(ordered))
	for i, base := range ordered {
		id := base.ID
		typeKey := fmt.Sprintf("block_type_%s", id)
		textKey := fmt.Sprintf("block_%s", id)

		typeVal := BlockType(strings.TrimSpace(data.Get(typeKey)))
		if typeVal == "" {
			typeVal = base.Type
		}
		textVal := data.Get(textKey)

		blocks[i] = Block{
			ID:   id,
			Type: typeVal,
			Text: textVal,
		}
	}

	markdown := BlocksToMarkdown(blocks)
	c.Post.SetContent(markdown)
	// preserve existing status, only updating content/title/summary
	c.Post.SetEditor(blogstore.POST_EDITOR_MARKDOWN)

	if err := c.App.GetBlogStore().PostUpdate(ctx, c.Post); err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("AI content editor: failed to save post", "error", err.Error())
		}
		c.Error = "Failed to save post. Please try again later."
		c.Success = ""
		return nil
	}

	c.Blocks = blocks
	c.Error = ""
	c.Success = successMessage
	if redirect {
		c.RedirectTo = shared.NewLinks().PostManager()
	} else {
		c.RedirectTo = ""
	}

	return nil
}

func (c *formAiPostContentUpdate) PostID() string {
	if c.Post == nil {
		return ""
	}
	return c.Post.ID()
}

func (c *formAiPostContentUpdate) wrapPage(breadcrumbs hb.TagInterface, header hb.TagInterface, backButton hb.TagInterface, content hb.TagInterface) hb.TagInterface {
	return c.Root(
		hb.Div().
			Class("container py-4").
			Child(breadcrumbs).
			Child(
				hb.Div().
					Class("d-flex justify-content-between align-items-center mb-3").
					Child(header).
					Child(backButton),
			).
			Child(content),
	)
}

func init() {
	if err := liveflux.Register(&formAiPostContentUpdate{}); err != nil {
		log.Printf("Failed to register formAiPostContentUpdate component: %v", err)
	}
}
