package aipostgenerator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/wf"
	"github.com/dromara/carbon/v2"
)

func (c *AiPostGeneratorController) onGeneratePost(r *http.Request) string {
	recordPostID := req.GetStringTrimmed(r, "record_post_id")
	if recordPostID == "" {
		return shared.ErrorPopup("Record Post ID is required").ToHTML()
	}

	initialData := map[string]any{
		"record_post_id": recordPostID,
		"post_title":     "",
	}

	stepGetPostDetails := wf.NewStep()
	stepGetPostDetails.SetName("Get Post Details")
	stepGetPostDetails.SetHandler(c.stepHandlerGetPostDetails)

	stepGeneratePost := wf.NewStep()
	stepGeneratePost.SetName("Generate Post Content")
	stepGeneratePost.SetHandler(c.stepHandlerGeneratePost)

	stepSavePost := wf.NewStep()
	stepSavePost.SetName("Save Blog Post")
	stepSavePost.SetHandler(c.stepHandlerSavePost)

	pipeline := wf.NewPipeline()
	pipeline.RunnableAdd(stepGetPostDetails)
	pipeline.RunnableAdd(stepGeneratePost)
	pipeline.RunnableAdd(stepSavePost)

	_, _, err := pipeline.Run(r.Context(), initialData)

	if err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error generating post: %s", err.Error())).ToHTML()
	}

	return hb.Swal(hb.SwalOptions{
		Title:            "Success",
		Text:             "Blog post generated successfully! Redirecting to editor...",
		Icon:             "success",
		Timer:            3000,
		TimerProgressBar: true,
		RedirectURL:      shared.NewLinks().AiPostEditor(map[string]string{"id": recordPostID}),
		RedirectSeconds:  3,
	}).ToHTML()
}

func (c *AiPostGeneratorController) stepHandlerGetPostDetails(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	recordPostID, ok := data["record_post_id"].(string)
	if !ok || recordPostID == "" {
		return ctx, data, errors.New("record_post_id is missing from pipeline data")
	}

	customStore := c.registry.GetCustomStore()
	if customStore == nil {
		return ctx, data, errors.New("custom store not configured")
	}

	record, err := customStore.RecordFindByID(recordPostID)
	if err != nil {
		return ctx, data, errors.New("failed to find post record with ID " + recordPostID + ": " + err.Error())
	}

	if record == nil {
		return ctx, data, errors.New("post record not found")
	}

	blogAiPost, err := blogai.NewRecordPostFromCustomRecord(record)
	if err != nil {
		return ctx, data, fmt.Errorf("failed to parse post record: %w", err)
	}

	data["post_title"] = blogAiPost.Title

	return ctx, data, nil
}

func (c *AiPostGeneratorController) stepHandlerGeneratePost(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	recordPostID, ok := data["record_post_id"].(string)
	if !ok || recordPostID == "" {
		return ctx, data, errors.New("record_post_id is missing from pipeline data")
	}

	recordPostTitle, ok := data["post_title"].(string)
	if !ok || recordPostTitle == "" {
		return ctx, data, errors.New("post_title is missing from pipeline data")
	}

	blogAgent := blogai.NewBlogWriterAgent(c.registry.GetLogger())
	if blogAgent == nil {
		return ctx, data, errors.New("failed to initialize LLM engine")
	}

	model, err := shared.LlmEngine(c.registry)
	if err != nil {
		return ctx, data, fmt.Errorf("failed to initialize LLM engine: %w", err)
	}

	blogAiPost, err := blogAgent.WriteBlogPost(model, recordPostTitle)
	if err != nil {
		return ctx, data, fmt.Errorf("failed to generate blog post: %w", err)
	}

	blogAiPost.ID = recordPostID
	blogAiPost.Status = blogai.POST_STATUS_DRAFT
	blogAiPost.CreatedAt = carbon.Now().ToDateTimeString(carbon.UTC)
	blogAiPost.UpdatedAt = carbon.Now().ToDateTimeString(carbon.UTC)

	data["blogai_post"] = blogAiPost

	return ctx, data, nil
}

func (c *AiPostGeneratorController) stepHandlerSavePost(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	recordPostID, ok := data["record_post_id"].(string)
	if !ok || recordPostID == "" {
		return ctx, data, errors.New("record_post_id is missing from pipeline data")
	}

	blogAiPost, okContent := data["blogai_post"].(blogai.RecordPost)
	if !okContent {
		return ctx, data, errors.New("blogai_post is missing from pipeline data for saving")
	}

	customStore := c.registry.GetCustomStore()
	if customStore == nil {
		return ctx, data, errors.New("custom store not configured")
	}

	record, err := customStore.RecordFindByID(recordPostID)
	if err != nil {
		return ctx, data, errors.New("failed to find post record with ID " + recordPostID + ": " + err.Error())
	}

	if record == nil {
		return ctx, data, errors.New("post record not found")
	}

	record.SetPayload(blogAiPost.ToJSON())

	if err := customStore.RecordUpdate(record); err != nil {
		return ctx, data, fmt.Errorf("failed to update post record: %w", err)
	}

	return ctx, data, nil
}
