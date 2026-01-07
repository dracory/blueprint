package aititlegenerator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/blogstore"
	"github.com/dracory/customstore"
	"github.com/dracory/wf"
	"github.com/samber/lo"
)

const (
	dataKeyExistingTitles  = "existing_titles"
	dataKeyGeneratedTitles = "generated_titles"
)

func (c *AiTitleGeneratorController) onGenerateTitles(r *http.Request) string {
	stepFindExistingTitles := wf.NewStep()
	stepFindExistingTitles.SetName("Find Existing Titles")
	stepFindExistingTitles.SetHandler(c.stepHandlerFindExistingTitles)

	stepGenerateTitles := wf.NewStep()
	stepGenerateTitles.SetName("Generate Titles")
	stepGenerateTitles.SetHandler(c.stepHandlerGenerateTitles)

	stepSaveTitles := wf.NewStep()
	stepSaveTitles.SetName("Save Generated Titles")
	stepSaveTitles.SetHandler(c.stepHandlerSaveTitles)

	pipeline := wf.NewPipeline()
	pipeline.RunnableAdd(stepFindExistingTitles)
	pipeline.RunnableAdd(stepGenerateTitles)
	pipeline.RunnableAdd(stepSaveTitles)

	initialData := map[string]any{}
	_, processedData, err := pipeline.Run(r.Context(), initialData)

	if err != nil {
		return shared.ErrorPopup(fmt.Sprintf("Error generating titles: %s", err.Error())).ToHTML()
	}

	generatedCount := 0
	if processedData != nil {
		if generated, ok := processedData[dataKeyGeneratedTitles]; ok {
			if titles, ok := generated.([]string); ok {
				generatedCount = len(titles)
			}
		}
	}

	successMessage := "Titles generated successfully"
	if generatedCount > 0 {
		suffix := "s"
		if generatedCount == 1 {
			suffix = ""
		}
		successMessage = fmt.Sprintf("Generated %d new title%s successfully", generatedCount, suffix)
	}

	return shared.SuccessPopupWithRedirect(
		successMessage,
		shared.NewLinks().AiTitleGenerator(),
		5,
	).ToHTML()
}

// stepHandlerFindExistingTitles finds the existing titles
//
// Business logic:
// 1. Loads all stored titles from the custom store
// 2. Loads all blog post titles from the blog store
// 3. Merges the two lists into a single slice
// 4. Adds the merged list to the data map as dataKeyExistingTitles
func (c *AiTitleGeneratorController) stepHandlerFindExistingTitles(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	records, err := c.registry.GetCustomStore().RecordList(customstore.RecordQuery().
		SetType(blogai.POST_RECORD_TYPE))
	if err != nil {
		return ctx, data, fmt.Errorf("failed to fetch titles: %w", err)
	}

	postList, err := c.registry.GetBlogStore().PostList(ctx, blogstore.PostQueryOptions{})
	if err != nil {
		return ctx, data, fmt.Errorf("failed to fetch blog list: %w", err)
	}

	existingTitles := []string{}
	for _, record := range records {
		recordPost, err := blogai.NewRecordPostFromCustomRecord(record)
		if err != nil {
			continue
		}
		existingTitles = append(existingTitles, recordPost.Title)
	}

	for _, post := range postList {
		existingTitles = append(existingTitles, post.Title())
	}

	data[dataKeyExistingTitles] = existingTitles
	return ctx, data, nil
}

// stepHandlerGenerateTitles generates the titles
//
// Business logic:
// 1. Checks if dataKeyExistingTitles is present in the data map
// 2. If not, returns an error
// 3. Initializes the LLM engine
// 4. Hydrates a TitleGeneratorAgent
// 5. Generates titles while avoiding duplicates
// 6. Adds the generated titles to the data map as "generated_titles"
func (c *AiTitleGeneratorController) stepHandlerGenerateTitles(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	if !lo.HasKey(data, dataKeyExistingTitles) {
		return ctx, data, errors.New("existing titles is required field")
	}

	existingTitles := data[dataKeyExistingTitles].([]string)

	blogTopic, err := c.loadBlogTopic(ctx)
	if err != nil {
		return ctx, data, fmt.Errorf("failed to load blog topic: %w", err)
	}

	model, err := shared.LlmEngine(c.registry)
	if err != nil {
		return ctx, data, fmt.Errorf("failed to initialize LLM engine: %w", err)
	}

	titleAgent := blogai.NewTitleGeneratorAgent()
	titles, err := titleAgent.GenerateTitles(model, blogTopic, existingTitles)
	if err != nil {
		return ctx, data, fmt.Errorf("failed to generate titles: %w", err)
	}

	data[dataKeyGeneratedTitles] = titles
	return ctx, data, nil
}

func (c *AiTitleGeneratorController) loadBlogTopic(ctx context.Context) (string, error) {
	store := c.registry.GetSettingStore()
	if store == nil {
		return "", errors.New("setting store is not configured")
	}

	topic, err := store.Get(ctx, SETTING_KEY_BLOG_TOPIC, "")
	if err != nil {
		return "", fmt.Errorf("failed to read blog topic: %w", err)
	}

	topic = strings.TrimSpace(topic)
	if topic == "" {
		return "", errors.New("blog topic is empty")
	}

	return topic, nil
}

func (c *AiTitleGeneratorController) stepHandlerSaveTitles(ctx context.Context, data map[string]any) (context.Context, map[string]any, error) {
	titles, ok := data[dataKeyGeneratedTitles].([]string)
	if !ok || len(titles) == 0 {
		return ctx, data, errors.New("no titles generated")
	}

	for _, title := range titles {
		record := customstore.NewRecord(blogai.POST_RECORD_TYPE)
		recordID := record.ID()

		record.SetPayloadMap(map[string]interface{}{
			"id":         recordID,
			"title":      title,
			"status":     blogai.POST_STATUS_PENDING,
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		})

		if err := c.registry.GetCustomStore().RecordCreate(record); err != nil {
			return ctx, data, fmt.Errorf("failed to save title '%s': %w", title, err)
		}
	}

	return ctx, data, nil
}
