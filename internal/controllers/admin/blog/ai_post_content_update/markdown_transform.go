package aipostcontentupdate

import (
	"strings"

	"project/pkg/blogai"

	"github.com/dracory/uid"
	"github.com/samber/lo"
)

var (
	introKeywords      = []string{"introduction", "intro", "getting started", "overview"}
	conclusionKeywords = []string{"conclusion", "summary", "final thoughts", "wrap up", "wrapping up", "closing thoughts"}
)

type parsedSection struct {
	title string
	lines []string
}

// == FLAT BLOCK MODEL ==========================================================

type BlockType string

const (
	BlockHeader1   BlockType = "h1"
	BlockHeader2   BlockType = "h2"
	BlockParagraph BlockType = "p"
	BlockCode      BlockType = "code"
)

type Block struct {
	ID   string
	Type BlockType
	Text string
}

// MarkdownToBlocks converts markdown into a flat list of sequential blocks.
// It preserves order and does not attempt to infer introduction/conclusion.
func MarkdownToBlocks(markdown string) []Block {
	lines := strings.Split(markdown, "\n")
	var (
		blocks    []Block
		paraBuf   []string
		inCode    bool
		codeFence string
	)

	flushParagraph := func() {
		if len(paraBuf) == 0 {
			return
		}
		text := normalizeParagraph(strings.Join(paraBuf, "\n"))
		if strings.TrimSpace(text) != "" {
			blocks = append(blocks, Block{ID: uid.HumanUid(), Type: BlockParagraph, Text: text})
		}
		paraBuf = nil
	}

	var codeBuf []string
	flushCode := func() {
		if len(codeBuf) == 0 {
			return
		}
		text := strings.Join(codeBuf, "\n")
		blocks = append(blocks, Block{ID: uid.HumanUid(), Type: BlockCode, Text: text})
		codeBuf = nil
	}

	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		trimmed := strings.TrimSpace(line)

		if fence, ok := isCodeFence(trimmed); ok {
			if !inCode {
				// starting code block
				flushParagraph()
				inCode = true
				codeFence = fence
				continue
			}
			// potential closing fence
			if fence == codeFence {
				inCode = false
				flushCode()
				codeFence = ""
				continue
			}
		}

		if inCode {
			codeBuf = append(codeBuf, line)
			continue
		}

		// Outside code blocks
		if trimmed == "" {
			flushParagraph()
			continue
		}

		if strings.HasPrefix(trimmed, "# ") {
			flushParagraph()
			text := strings.TrimSpace(trimmed[2:])
			blocks = append(blocks, Block{ID: uid.HumanUid(), Type: BlockHeader1, Text: text})
			continue
		}

		if headingLevel(trimmed) >= 2 {
			flushParagraph()
			sectionTitle := strings.TrimSpace(trimmed[headingPrefixLength(trimmed):])
			blocks = append(blocks, Block{ID: uid.HumanUid(), Type: BlockHeader2, Text: sectionTitle})
			continue
		}

		paraBuf = append(paraBuf, line)
	}

	flushParagraph()
	flushCode()

	return blocks
}

// BlocksToMarkdown converts a sequence of blocks back to markdown.
func BlocksToMarkdown(blocks []Block) string {
	var builder strings.Builder

	for _, block := range blocks {
		text := strings.TrimRight(block.Text, "\n")
		if strings.TrimSpace(text) == "" && block.Type != BlockCode {
			continue
		}

		switch block.Type {
		case BlockHeader1:
			builder.WriteString("# ")
			builder.WriteString(strings.TrimSpace(text))
			builder.WriteString("\n\n")
		case BlockHeader2:
			builder.WriteString("## ")
			builder.WriteString(strings.TrimSpace(text))
			builder.WriteString("\n\n")
		case BlockParagraph:
			cleaned := normalizeParagraph(text)
			if strings.TrimSpace(cleaned) == "" {
				continue
			}
			builder.WriteString(cleaned)
			builder.WriteString("\n\n")
		case BlockCode:
			builder.WriteString("```")
			builder.WriteString("\n")
			builder.WriteString(text)
			builder.WriteString("\n```")
			builder.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(builder.String())
}

// MarkdownToRecordPost converts markdown content into a blogai.RecordPost structure that mirrors
// the block-based editor expectations (introduction, sections, conclusion).
// fallbackTitle is used when the markdown content does not contain a level-1 heading.
func MarkdownToRecordPost(markdown, fallbackTitle string) blogai.RecordPost {
	title, preamble, sections := parseMarkdown(markdown)

	record := blogai.RecordPost{
		Title:        lo.Ternary(strings.TrimSpace(title) != "", strings.TrimSpace(title), strings.TrimSpace(fallbackTitle)),
		Introduction: blogai.PostContentIntroduction{Title: "Introduction"},
		Sections:     []blogai.PostContentSection{},
		Conclusion:   blogai.PostContentConclusion{Title: "Conclusion"},
	}

	preambleParagraphs := paragraphsFromLines(preamble)
	if len(preambleParagraphs) > 0 {
		record.Introduction.Paragraphs = preambleParagraphs
	}

	if len(record.Introduction.Paragraphs) == 0 && len(sections) > 0 && isIntroductionTitle(sections[0].title) {
		record.Introduction.Title = sections[0].title
		record.Introduction.Paragraphs = paragraphsFromLines(sections[0].lines)
		sections = sections[1:]
	}

	if len(sections) > 0 && isConclusionTitle(sections[len(sections)-1].title) {
		last := sections[len(sections)-1]
		record.Conclusion.Title = last.title
		record.Conclusion.Paragraphs = paragraphsFromLines(last.lines)
		sections = sections[:len(sections)-1]
	}

	for _, section := range sections {
		paragraphs := paragraphsFromLines(section.lines)
		if len(paragraphs) == 0 && strings.TrimSpace(section.title) == "" {
			continue
		}
		record.Sections = append(record.Sections, blogai.PostContentSection{
			Title:      section.title,
			Paragraphs: paragraphs,
		})
	}

	if len(record.Introduction.Paragraphs) == 0 {
		// Prefer to promote the first section to the introduction only when we
		// don't already have explicit conclusion content. This ensures that in
		// documents where there is a clear conclusion section, earlier sections
		// remain as main sections rather than being reclassified as the intro.
		if len(record.Sections) > 0 && len(record.Conclusion.Paragraphs) == 0 {
			first := record.Sections[0]
			record.Introduction.Title = lo.Ternary(strings.TrimSpace(first.Title) != "", first.Title, record.Introduction.Title)
			record.Introduction.Paragraphs = append([]string{}, first.Paragraphs...)
			record.Sections = record.Sections[1:]
		} else if len(record.Sections) == 0 && len(record.Conclusion.Paragraphs) > 0 {
			// If there are no sections at all but we do have conclusion content,
			// use it as an introduction and clear the conclusion paragraphs.
			record.Introduction.Paragraphs = append([]string{}, record.Conclusion.Paragraphs...)
			record.Conclusion.Paragraphs = []string{}
		}
	}

	if strings.TrimSpace(record.Introduction.Title) == "" {
		record.Introduction.Title = "Introduction"
	}

	if strings.TrimSpace(record.Conclusion.Title) == "" {
		record.Conclusion.Title = "Conclusion"
	}

	return record
}

// RecordPostToMarkdown converts a RecordPost back to markdown text.
func RecordPostToMarkdown(record blogai.RecordPost) string {
	var builder strings.Builder

	title := strings.TrimSpace(record.Title)
	if title != "" {
		builder.WriteString("# ")
		builder.WriteString(title)
		builder.WriteString("\n\n")
	}

	writeSection(&builder, record.Introduction.Title, record.Introduction.Paragraphs)

	for _, section := range record.Sections {
		// Skip sections that have no meaningful title and whose paragraphs are
		// all effectively empty. This matches the expectation that placeholder
		// or empty sections are not rendered in the final markdown.
		hasContent := false
		for _, p := range section.Paragraphs {
			if strings.TrimSpace(normalizeParagraph(p)) != "" {
				hasContent = true
				break
			}
		}
		if !hasContent && strings.TrimSpace(section.Title) == "" {
			continue
		}

		writeSection(&builder, section.Title, section.Paragraphs)
	}

	writeSection(&builder, record.Conclusion.Title, record.Conclusion.Paragraphs)

	return strings.TrimSpace(builder.String())
}

func writeSection(builder *strings.Builder, title string, paragraphs []string) {
	if len(paragraphs) == 0 {
		return
	}

	sectionTitle := strings.TrimSpace(title)
	if sectionTitle == "" {
		sectionTitle = "Section"
	}

	if builder.Len() > 0 {
		builder.WriteString("## ")
	} else {
		builder.WriteString("## ")
	}
	builder.WriteString(sectionTitle)
	builder.WriteString("\n\n")

	for _, paragraph := range paragraphs {
		cleaned := normalizeParagraph(paragraph)
		if strings.TrimSpace(cleaned) == "" {
			continue
		}
		builder.WriteString(cleaned)
		builder.WriteString("\n\n")
	}
}

func parseMarkdown(markdown string) (title string, preamble []string, sections []parsedSection) {
	lines := strings.Split(markdown, "\n")
	var current *parsedSection
	inCodeBlock := false
	codeFence := ""

	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		trimmed := strings.TrimSpace(line)

		if fence, ok := isCodeFence(trimmed); ok {
			if !inCodeBlock {
				inCodeBlock = true
				codeFence = fence
			} else if fence == codeFence {
				inCodeBlock = false
			}
			appendLine(current, &preamble, line)
			continue
		}

		if !inCodeBlock {
			if strings.HasPrefix(trimmed, "# ") {
				if title == "" {
					title = strings.TrimSpace(trimmed[2:])
					current = nil
					continue
				}
			}

			if headingLevel(trimmed) >= 2 {
				sectionTitle := strings.TrimSpace(trimmed[headingPrefixLength(trimmed):])
				sections = append(sections, parsedSection{title: sectionTitle})
				current = &sections[len(sections)-1]
				continue
			}
		}

		appendLine(current, &preamble, line)
	}

	return title, preamble, sections
}

func appendLine(current *parsedSection, preamble *[]string, line string) {
	if current != nil {
		current.lines = append(current.lines, line)
		return
	}
	*preamble = append(*preamble, line)
}

func paragraphsFromLines(lines []string) []string {
	var (
		paragraphs []string
		buffer     []string
		inCode     bool
		codeFence  string
	)

	flush := func() {
		if len(buffer) == 0 {
			return
		}
		paragraph := normalizeParagraph(strings.Join(buffer, "\n"))
		if strings.TrimSpace(paragraph) != "" {
			paragraphs = append(paragraphs, paragraph)
		}
		buffer = nil
	}

	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		trimmed := strings.TrimSpace(line)

		if fence, ok := isCodeFence(trimmed); ok {
			if !inCode {
				inCode = true
				codeFence = fence
			} else if fence == codeFence {
				inCode = false
			}
			buffer = append(buffer, line)
			continue
		}

		if !inCode && trimmed == "" {
			flush()
			continue
		}

		buffer = append(buffer, line)
	}

	flush()

	return paragraphs
}

func normalizeParagraph(paragraph string) string {
	trimmed := strings.Trim(paragraph, "\n\r")
	lines := strings.Split(trimmed, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

func isIntroductionTitle(title string) bool {
	lowered := strings.ToLower(strings.TrimSpace(title))
	for _, keyword := range introKeywords {
		if lowered == keyword {
			return true
		}
	}
	return false
}

func isConclusionTitle(title string) bool {
	lowered := strings.ToLower(strings.TrimSpace(title))
	for _, keyword := range conclusionKeywords {
		if lowered == keyword {
			return true
		}
	}
	return false
}

func isCodeFence(line string) (string, bool) {
	if strings.HasPrefix(line, "```") {
		return "```", true
	}
	if strings.HasPrefix(line, "~~~") {
		return "~~~", true
	}
	return "", false
}

func headingLevel(line string) int {
	count := 0
	for i := 0; i < len(line); i++ {
		if line[i] == '#' {
			count++
		} else if line[i] == ' ' {
			return count
		} else {
			break
		}
	}
	return 0
}

func headingPrefixLength(line string) int {
	level := headingLevel(line)
	if level == 0 {
		return 0
	}
	return level + 1 // number of '#' plus trailing space
}
