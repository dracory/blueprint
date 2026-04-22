package search

import (
	"fmt"
	"net/url"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
)

// renderSearchHTML renders the search block HTML
func renderSearchHTML(query, placeholder string, results []SearchResult, totalResults, pageNum, resultsPerPage int) (string, error) {
	var content []hb.TagInterface

	// Search box section
	searchBox := renderSearchBox(query, placeholder)
	content = append(content, searchBox)

	// Results section (if query exists)
	if query != "" {
		resultsSection := renderResultsSection(query, results, totalResults, pageNum, resultsPerPage)
		content = append(content, resultsSection)
	}

	// Wrap in section
	section := hb.Section().
		ID("SectionSearch").
		Style("background:#fff; padding:50px 0px 80px 0px;").
		Child(bs.Container().Children(content))

	return section.ToHTML(), nil
}

// renderSearchBox renders the search input form
func renderSearchBox(query, placeholder string) hb.TagInterface {
	searchInput := hb.Input().
		Type("search").
		Name("q").
		Class("form-control form-control-lg").
		Style("border-radius: 8px; border: 2px solid #794FC6; padding: 15px 20px; font-size: 18px;").
		Placeholder(placeholder).
		Value(query)

	searchButton := hb.Button().
		Type("submit").
		Class("btn btn-lg").
		Style("background:#794FC6; color:#fff; border-radius: 8px; padding: 15px 30px; margin-left: 10px;").
		HTML("Search")

	searchForm := hb.Form().
		Method("GET").
		Action("").
		Style("display: flex; align-items: center; max-width: 800px; margin: 0 auto 40px auto;").
		Children([]hb.TagInterface{
			searchInput,
			searchButton,
		})

	return hb.Div().
		Style("text-align: center; margin-bottom: 30px;").
		Child(searchForm)
}

// renderResultsSection renders the search results and pagination
func renderResultsSection(query string, results []SearchResult, totalResults, pageNum, resultsPerPage int) hb.TagInterface {
	var children []hb.TagInterface

	// Results header
	resultsCount := len(results)
	if resultsCount == 0 {
		noResults := hb.Div().
			Style("text-align: center; padding: 40px 0; color: #6c757d;").
			Children([]hb.TagInterface{
				hb.Heading4().HTML("No results found"),
				hb.Paragraph().HTML(fmt.Sprintf("No results found for \"<strong>%s</strong>\"", escapeHTML(query))),
			})
		children = append(children, noResults)
	} else {
		// Results count
		resultText := fmt.Sprintf("Found %d result(s) for \"<strong>%s</strong>\"", resultsCount, escapeHTML(query))
		if totalResults > resultsCount {
			resultText = fmt.Sprintf("Showing %d of %d+ results for \"<strong>%s</strong>\"", resultsCount, totalResults, escapeHTML(query))
		}

		resultsHeader := hb.Div().
			Style("margin-bottom: 30px; color: #6c757d; font-size: 16px;").
			HTML(resultText)
		children = append(children, resultsHeader)

		// Results list
		resultsList := hb.Div().Class("search-results")
		for _, result := range results {
			resultItem := renderResultItem(result)
			resultsList.Child(resultItem)
		}
		children = append(children, resultsList)

		// Pagination
		totalPages := (totalResults + resultsPerPage - 1) / resultsPerPage
		if totalPages > 1 {
			pagination := renderPaginationControls(query, pageNum, totalPages)
			children = append(children, pagination)
		}
	}

	return hb.Div().
		Style("max-width: 900px; margin: 0 auto;").
		Children(children)
}

// renderResultItem renders a single search result
func renderResultItem(result SearchResult) hb.TagInterface {
	// Type badge
	var badgeColor string
	var badgeText string
	switch result.Type {
	case "page":
		badgeColor = "#1ba1b6"
		badgeText = "Page"
	case "post":
		badgeColor = "#794FC6"
		badgeText = "Blog Post"
	default:
		badgeColor = "#6c757d"
		badgeText = "Content"
	}

	badge := hb.Span().
		Style(fmt.Sprintf("display: inline-block; background: %s; color: #fff; padding: 3px 10px; border-radius: 4px; font-size: 12px; margin-bottom: 8px;", badgeColor)).
		Text(badgeText)

	// Title link
	titleLink := hb.Hyperlink().
		Href(result.URL).
		Style("text-decoration: none; color: #794FC6; font-size: 22px; font-weight: 600; display: block; margin-bottom: 8px;").
		HTML(escapeHTML(result.Title))

	var metaParts []string
	if result.Date != "" {
		metaParts = append(metaParts, result.Date)
	}

	var metaInfo hb.TagInterface
	if len(metaParts) > 0 {
		metaText := ""
		for i, part := range metaParts {
			if i > 0 {
				metaText += " | "
			}
			metaText += part
		}
		metaInfo = hb.Paragraph().
			Style("font-size: 14px; color: #6c757d; margin-bottom: 8px;").
			Text(metaText)
	}

	// Summary
	summary := hb.Paragraph().
		Style("font-size: 16px; color: #333; line-height: 1.5; margin-bottom: 0;").
		Text(escapeHTML(result.Summary))

	resultContent := []hb.TagInterface{badge, titleLink}
	if metaInfo != nil {
		resultContent = append(resultContent, metaInfo)
	}
	resultContent = append(resultContent, summary)

	return hb.Div().
		Style("padding: 25px; border-bottom: 1px solid #e9ecef;").
		Children(resultContent)
}

// renderPaginationControls renders pagination buttons
func renderPaginationControls(query string, currentPage, totalPages int) hb.TagInterface {
	var pages []hb.TagInterface

	// Previous button
	if currentPage > 0 {
		prevURL := buildPageURL(query, currentPage-1)
		prevBtn := hb.Hyperlink().
			Href(prevURL).
			Class("page-link").
			HTML("&laquo; Previous")
		prevLi := hb.Li().Class("page-item").Child(prevBtn)
		pages = append(pages, prevLi)
	}

	// Page numbers (show up to 5 pages around current)
	startPage := currentPage - 2
	if startPage < 0 {
		startPage = 0
	}
	endPage := startPage + 5
	if endPage > totalPages {
		endPage = totalPages
		startPage = endPage - 5
		if startPage < 0 {
			startPage = 0
		}
	}

	for i := startPage; i < endPage && i < totalPages; i++ {
		pageNum := i + 1 // Display as 1-based
		pageURL := buildPageURL(query, i)

		pageLink := hb.Hyperlink().
			Href(pageURL).
			Class("page-link").
			Text(fmt.Sprintf("%d", pageNum))

		pageItem := hb.Li().Class("page-item").Child(pageLink)

		// Highlight current page
		if i == currentPage {
			pageItem.Class("page-item active")
			pageLink = hb.Span().Class("page-link").Text(fmt.Sprintf("%d", pageNum))
			pageItem = hb.Li().Class("page-item active").Child(pageLink)
		}

		pages = append(pages, pageItem)
	}

	// Next button
	if currentPage < totalPages-1 {
		nextURL := buildPageURL(query, currentPage+1)
		nextBtn := hb.Hyperlink().
			Href(nextURL).
			Class("page-link").
			HTML("Next &raquo;")
		nextLi := hb.Li().Class("page-item").Child(nextBtn)
		pages = append(pages, nextLi)
	}

	pagination := hb.Ul().
		Class("pagination justify-content-center mt-4").
		Children(pages)

	return hb.Nav().
		Style("margin-top: 30px;").
		Child(pagination)
}

// buildPageURL builds a URL for a specific page with the search query
func buildPageURL(query string, pageNum int) string {
	params := url.Values{}
	if query != "" {
		params.Set("q", query)
	}
	// Only include page parameter if not on first page (1-based for URL)
	if pageNum > 0 {
		params.Set("page", fmt.Sprintf("%d", pageNum+1))
	}

	if len(params) == 0 {
		return ""
	}
	return "?" + params.Encode()
}

// escapeHTML escapes HTML special characters
func escapeHTML(text string) string {
	replacer := map[string]string{
		"&":  "&amp;",
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
	}

	result := text
	// Process & first to avoid double-escaping
	result = replaceAll(result, "&", "&amp;")
	for old, new := range replacer {
		if old != "&" {
			result = replaceAll(result, old, new)
		}
	}
	return result
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := 0
		for i := 0; i <= len(s)-len(old); i++ {
			if s[i:i+len(old)] == old {
				idx = i
				break
			}
		}
		if idx == 0 && (len(s) < len(old) || s[:len(old)] != old) {
			result += s
			break
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	return result
}
