package blogai

import (
	"strings"

	"github.com/dracory/llm"
	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
)

type TitleGeneratorAgentV1 struct{}

func NewTitleGeneratorAgentV1() *TitleGeneratorAgentV1 {
	return &TitleGeneratorAgentV1{}
}

// GenerateTitles generates 10 titles for the given topic.
func (t *TitleGeneratorAgentV1) GenerateTitles(ai llm.LlmInterface, topic string, existingTitles []string) ([]string, error) {
	systemPrompt := `
You are an experienced blog post writer. Your task is to generate 10 titles for the given topic.

The following titles are already existing:
` + strings.Join(existingTitles, "\n") + `

You MUST NOT repeat any of the existing titles.

Expected JSON structure:
{
	"titles": [
		"Title 1",
		"Title 2",
		"Title 3",
		"Title 4",
		"Title 5",
		"Title 6",
		"Title 7",
		"Title 8",
		"Title 9",
		"Title 10"
	]
}
`
	userMessage := `
Your task is to generate 10 titles for the given topic.

Topic: ` + topic + `

`

	response, err := ai.Generate(systemPrompt, userMessage, llm.LlmOptions{})
	if err != nil {
		return nil, err
	}

	_, _ = cfmt.Successln("Response: ", response)

	titlesResponse := gjson.Get(response, "titles").Array()

	titles := lo.Map(titlesResponse, func(title gjson.Result, _ int) string {
		return title.String()
	})

	titles = lo.Filter(titles, func(title string, _ int) bool {
		return !lo.Contains(existingTitles, title)
	})

	return titles, nil
}
