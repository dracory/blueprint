package aipostcontentupdate

import (
	"strings"
	"testing"

	"project/pkg/blogai"
)

func TestMarkdownToRecordPost_BasicStructure(t *testing.T) {
	markdown := `# Contract Basics

Intro paragraph one.

## Section One
Section paragraph.

## Conclusion
Closing paragraph.`

	record := MarkdownToRecordPost(markdown, "Fallback Title")

	if record.Title != "Contract Basics" {
		t.Fatalf("expected title to be 'Contract Basics', got %q", record.Title)
	}

	if len(record.Introduction.Paragraphs) == 0 {
		t.Fatal("expected introduction paragraphs")
	}

	if len(record.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(record.Sections))
	}

	if !strings.EqualFold(record.Conclusion.Title, "Conclusion") {
		t.Fatalf("expected conclusion title to remain 'Conclusion', got %q", record.Conclusion.Title)
	}

	if len(record.Conclusion.Paragraphs) != 1 {
		t.Fatalf("expected conclusion paragraphs to be 1, got %d", len(record.Conclusion.Paragraphs))
	}
}

func TestMarkdownToRecordPost_NoHeadings(t *testing.T) {
	markdown := `This is a paragraph without headings.`

	record := MarkdownToRecordPost(markdown, "Fallback Title")

	if record.Title != "Fallback Title" {
		t.Fatalf("expected fallback title, got %q", record.Title)
	}

	if len(record.Introduction.Paragraphs) != 1 {
		t.Fatalf("expected introduction paragraph from body, got %d", len(record.Introduction.Paragraphs))
	}

	if len(record.Sections) != 0 {
		t.Fatalf("expected no sections, got %d", len(record.Sections))
	}
}

func TestRecordPostToMarkdown_RoundTrip(t *testing.T) {
	original := `# Contract Basics

Intro paragraph one.

## Section One
Section paragraph.

## Conclusion
Closing paragraph.`

	record := MarkdownToRecordPost(original, "Fallback Title")
	rebuilt := RecordPostToMarkdown(record)

	if strings.TrimSpace(rebuilt) == "" {
		t.Fatal("expected rebuilt markdown to have content")
	}

	if !strings.Contains(rebuilt, "## Section One") {
		t.Fatalf("expected rebuilt markdown to contain section heading, got %q", rebuilt)
	}

	roundTripRecord := MarkdownToRecordPost(rebuilt, "Fallback Title")
	if len(roundTripRecord.Sections) != len(record.Sections) {
		t.Fatalf("expected same number of sections after roundtrip, got %d vs %d", len(roundTripRecord.Sections), len(record.Sections))
	}
}

func TestMarkdownToRecordPost_IntroductionFromFirstSection(t *testing.T) {
	markdown := `## Introduction
Intro paragraph one.

## Section One
Section paragraph.`

	record := MarkdownToRecordPost(markdown, "Fallback Title")

	if !strings.EqualFold(record.Introduction.Title, "Introduction") {
		t.Fatalf("expected introduction title to be 'Introduction', got %q", record.Introduction.Title)
	}

	if len(record.Introduction.Paragraphs) != 1 {
		t.Fatalf("expected introduction paragraph to be captured, got %d", len(record.Introduction.Paragraphs))
	}

	if len(record.Sections) != 1 {
		t.Fatalf("expected one main section after extracting intro, got %d", len(record.Sections))
	}
}

func TestMarkdownToRecordPost_ConclusionExtraction(t *testing.T) {
	markdown := `## Section One
Section paragraph.

## Conclusion
Closing paragraph.`

	record := MarkdownToRecordPost(markdown, "Fallback Title")

	if len(record.Conclusion.Paragraphs) != 1 {
		t.Fatalf("expected conclusion paragraph extracted, got %d", len(record.Conclusion.Paragraphs))
	}

	if len(record.Sections) != 1 {
		t.Fatalf("expected remaining sections to be 1, got %d", len(record.Sections))
	}
}

func TestRecordPostToMarkdown_EmptySections(t *testing.T) {
	record := blogai.RecordPost{
		Title: "Sample",
		Introduction: blogai.PostContentIntroduction{
			Title:      "Intro",
			Paragraphs: []string{""},
		},
		Sections: []blogai.PostContentSection{
			{
				Title:      "",
				Paragraphs: []string{""},
			},
		},
	}

	result := RecordPostToMarkdown(record)
	if strings.Contains(result, "## Section") {
		t.Fatalf("expected empty sections to be skipped, got %q", result)
	}
}
