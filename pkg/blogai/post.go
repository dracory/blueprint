package blogai

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/dracory/customstore"
	"github.com/samber/lo"
)

type PostContentIntroduction struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}

type PostContentSection struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}

type PostContentConclusion struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}

type RecordPost struct {
	ID              string                  `json:"id,omitempty" validate:"omitempty"` // Unique identifier for the post
	Title           string                  `json:"title" validate:"required"`         // Title of the blog post
	Status          string                  `json:"status,omitempty"`                  // Status of the post (e.g., draft, published)
	Subtitle        string                  `json:"subtitle,omitempty"`                // Subtitle of the blog post
	Summary         string                  `json:"summary,omitempty"`                 // Brief summary of the post
	Introduction    PostContentIntroduction `json:"introduction"`                      // Introduction section of the post
	Sections        []PostContentSection    `json:"sections"`                          // Main content sections of the post
	Conclusion      PostContentConclusion   `json:"conclusion"`                        // Conclusion section of the post
	Keywords        []string                `json:"keywords,omitempty"`                // Keywords related to the post
	CreatedAt       string                  `json:"created_at,omitempty"`              // Timestamp when the post was created
	UpdatedAt       string                  `json:"updated_at,omitempty"`              // Timestamp when the post was last updated
	MetaDescription string                  `json:"meta_description,omitempty"`        // Meta description for SEO
	MetaKeywords    []string                `json:"meta_keywords,omitempty"`           // Meta keywords for SEO
	MetaTitle       string                  `json:"meta_title,omitempty"`              // Meta title for SEO
	Image           string                  `json:"image,omitempty"`                   // Generated image URL/path
}

func newRecordPostFromMap(data map[string]any) (RecordPost, error) {
	id := lo.
		IfF(lo.HasKey(data, "id"), func() string {
			return data["id"].(string)
		}).Else("")

	title := lo.
		IfF(lo.HasKey(data, "title"), func() string {
			return data["title"].(string)
		}).Else("")

	status := lo.
		IfF(lo.HasKey(data, "status"), func() string {
			return data["status"].(string)
		}).
		Else("")

	subtitle := lo.
		IfF(lo.HasKey(data, "subtitle"), func() string {
			return data["subtitle"].(string)
		}).
		Else("")

	summary := lo.
		IfF(lo.HasKey(data, "summary"), func() string {
			return data["summary"].(string)
		}).Else("")

	var introduction PostContentIntroduction
	if introData, ok := data["introduction"].(map[string]any); ok {
		title := ""
		if titleVal, ok := introData["title"].(string); ok {
			title = titleVal
		}

		var paragraphs []string
		if paragraphsData, ok := introData["paragraphs"].([]any); ok {
			paragraphs = lo.Map(paragraphsData, func(item any, _ int) string {
				return item.(string)
			})
		} else if contentVal, ok := introData["content"].(string); ok {
			// For backward compatibility with old format
			paragraphs = []string{contentVal}
		}

		introduction = PostContentIntroduction{
			Title:      title,
			Paragraphs: paragraphs,
		}
	}

	var sections []PostContentSection
	if sectionsData, ok := data["sections"].([]any); ok {
		for _, sectionData := range sectionsData {
			if sectionMap, ok := sectionData.(map[string]any); ok {
				title := ""
				if titleVal, ok := sectionMap["title"].(string); ok {
					title = titleVal
				}

				var paragraphs []string
				if paragraphsData, ok := sectionMap["paragraphs"].([]any); ok {
					paragraphs = lo.Map(paragraphsData, func(item any, _ int) string {
						return item.(string)
					})
				} else if contentVal, ok := sectionMap["content"].(string); ok {
					// For backward compatibility with old format
					paragraphs = []string{contentVal}
				}

				sections = append(sections, PostContentSection{
					Title:      title,
					Paragraphs: paragraphs,
				})
			}
		}
	}

	var conclusion PostContentConclusion
	if conclusionData, ok := data["conclusion"].(map[string]any); ok {
		title := ""
		if titleVal, ok := conclusionData["title"].(string); ok {
			title = titleVal
		}

		var paragraphs []string
		if paragraphsData, ok := conclusionData["paragraphs"].([]any); ok {
			paragraphs = lo.Map(paragraphsData, func(item any, _ int) string {
				return item.(string)
			})
		} else if contentVal, ok := conclusionData["content"].(string); ok {
			// For backward compatibility with old format
			paragraphs = []string{contentVal}
		}

		conclusion = PostContentConclusion{
			Title:      title,
			Paragraphs: paragraphs,
		}
	}

	keywords := lo.
		IfF(lo.HasKey(data, "keywords"), func() []string {
			keywords := data["keywords"].([]any)
			return lo.Map(keywords, func(item any, _ int) string {
				return item.(string)
			})
		}).
		Else([]string{})

	createdAt := lo.
		IfF(lo.HasKey(data, "created_at"), func() string {
			return data["created_at"].(string)
		}).
		Else(time.Now().UTC().Format(time.RFC3339))

	updatedAt := lo.
		IfF(lo.HasKey(data, "updated_at"), func() string {
			return data["updated_at"].(string)
		}).
		Else(time.Now().UTC().Format(time.RFC3339))

	metaTitle := lo.
		IfF(lo.HasKey(data, "meta_title"), func() string {
			return data["meta_title"].(string)
		}).
		Else(title) // Default to main title

	metaDescription := lo.
		IfF(lo.HasKey(data, "meta_description"), func() string {
			return data["meta_description"].(string)
		}).
		Else("")

	metaKeywords := lo.
		IfF(lo.HasKey(data, "meta_keywords"), func() []string {
			return lo.Map(data["meta_keywords"].([]any), func(item any, _ int) string {
				return item.(string)
			})
		}).
		Else([]string{})

	image := lo.
		IfF(lo.HasKey(data, "image"), func() string {
			return data["image"].(string)
		}).
		Else("")

	if id == "" {
		return RecordPost{}, errors.New("id is required")
	}

	if title == "" {
		return RecordPost{}, errors.New("title is required")
	}

	return RecordPost{
		ID:              id,
		Title:           title,
		Status:          status,
		Subtitle:        subtitle,
		Summary:         summary,
		Introduction:    introduction,
		Sections:        sections,
		Conclusion:      conclusion,
		Keywords:        keywords,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		MetaDescription: metaDescription,
		MetaKeywords:    metaKeywords,
		MetaTitle:       metaTitle,
		Image:           image,
	}, nil
}

func NewRecordPostFromCustomRecord(record customstore.RecordInterface) (RecordPost, error) {
	if record == nil {
		return RecordPost{}, errors.New("record is required")
	}

	payloadMap, err := record.PayloadMap()
	if err != nil {
		return RecordPost{}, err
	}

	payloadMap["id"] = record.ID()
	return newRecordPostFromMap(payloadMap)
}

// ToJSON converts the RecordPost to a JSON string
func (r RecordPost) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
