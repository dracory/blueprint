package blogai

import (
	"fmt"
	"strings"

	"github.com/dracory/llm"
	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
)

type TitleGeneratorAgent struct{}

func NewTitleGeneratorAgent() *TitleGeneratorAgent {
	return &TitleGeneratorAgent{}
}

// GenerateTitles generates 10-15 optimized titles for the given topic.
func (t *TitleGeneratorAgent) GenerateTitles(ai llm.LlmInterface, topic string, existingTitles []string) ([]string, error) {
	existingTitlesText := "None provided."
	if len(existingTitles) > 0 {
		existingTitlesText = "- " + strings.Join(existingTitles, "\n- ")
	}

	systemPrompt := fmt.Sprintf(`
You are an expert blog title generator capable of crafting high-performing titles across any niche. Each request includes a topic (for example: legal, contract review, sports, technology). Tailor tone, terminology, and audience focus to that topic.

## Your Objective
Generate compelling, SEO-optimized blog post titles that:
- Capture readers' attention immediately
- Clearly communicate the article's value proposition
- Reflect the motivations and pain points of audiences who care about the topic
- Avoid duplicating any provided titles

## Title Style Guidelines
1. Be specific with concrete numbers, timeframes, outcomes, or audiences.
2. Use persuasive power words appropriate for the topic (e.g. "essential", "critical", "hidden", "ultimate", "proven", "complete").
3. Include action verbs that imply movement or decision-making ("review", "avoid", "negotiate", "protect", "master").
4. Address common questions, mistakes, or challenges relevant to the topic's audience.
5. Promise clear, practical value or insight.

## Title Formats to Mix In
- How-to guides ("How to [Achieve Outcome] in [Timeframe]")
- Lists ("7 [Topic] Mistakes That Cost You [Resource]")
- Questions ("Is Your [Topic Element] Ready for [Scenario]?")
- Warnings ("Never [Action] Without Checking These [Number] Things")
- Comparisons ("[Option A] vs [Option B]: Which Fits Your [Need]?")
- Ultimate guides ("The Complete Guide to [Topic Focus]")
- Checklists or frameworks ("The [Number]-Point Checklist for [Outcome]")

## Angles to Explore
- Foundational best practices and step-by-step walkthroughs
- Common pitfalls, risks, or red flags
- Key terminology or concept breakdowns
- Industry- or audience-specific considerations
- Cost/time-saving strategies and negotiation tips
- Emerging trends, expert insights, or FAQs

## Existing Titles to Avoid
%s

## Output Requirements
- Generate between 10 and 15 unique titles tailored to the topic.
- Provide a one-sentence description for each title highlighting the article's value.
- Ensure titles are distinct from one another and from the existing titles list.
- Respond ONLY with valid JSON matching this schema:
{
  "titles": [
    {
      "title": "string",
      "description": "string"
    }
  ]
}
- Do not include any additional commentary, explanations, or markdown outside the JSON.
`, existingTitlesText)

	userMessage := fmt.Sprintf(`
Generate 10-15 compelling, SEO-optimized blog post titles (with descriptions) about the following topic while following all constraints:

Topic: %s
`, topic)

	response, err := ai.Generate(systemPrompt, userMessage, llm.LlmOptions{})
	if err != nil {
		return nil, err
	}

	cfmt.Successln("Response: ", response)

	titlesResponse := gjson.Get(response, "titles").Array()

	titles := lo.Map(titlesResponse, func(item gjson.Result, _ int) string {
		if item.IsObject() {
			return strings.TrimSpace(item.Get("title").String())
		}
		return strings.TrimSpace(item.String())
	})

	titles = lo.Filter(titles, func(title string, _ int) bool {
		return title != ""
	})

	existingSet := map[string]struct{}{}
	for _, existingTitle := range existingTitles {
		normalized := strings.ToLower(strings.TrimSpace(existingTitle))
		if normalized == "" {
			continue
		}
		existingSet[normalized] = struct{}{}
	}

	seen := map[string]struct{}{}
	filtered := make([]string, 0, len(titles))
	for _, title := range titles {
		normalized := strings.ToLower(strings.TrimSpace(title))
		if normalized == "" {
			continue
		}
		if _, exists := existingSet[normalized]; exists {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		filtered = append(filtered, title)
	}

	return filtered, nil
}
