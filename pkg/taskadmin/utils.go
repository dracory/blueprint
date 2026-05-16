package taskadmin

import (
	"net/http"
	"strings"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func sortableColumnLabel(columnName string, columnTableName string, request http.Request, sortBy string, sortOrder string) *hb.Tag {
	isSelected := strings.EqualFold(sortBy, columnTableName)

	direction := lo.If(isSelected && sortOrder == "asc", "desc").
		Else("asc")

	// link := links.Admin().Tasks(map[string]string{
	// 	"page": "0",
	// 	"by":   columnTableName,
	// 	"sort": direction,
	// })

	url := link(request, PathHome, map[string]string{
		"page": "0",
		"by":   columnTableName,
		"sort": direction,
	})
	return hb.NewHyperlink().
		HTML(columnName).
		Child(sortingIndicator(columnTableName, sortBy, sortOrder)).
		Href(url)
}

func sortingIndicator(columnName string, sortByColumnName string, sortOrder string) *hb.Tag {
	isSelected := strings.EqualFold(sortByColumnName, columnName)

	direction := lo.If(isSelected && sortOrder == "asc", "up").
		ElseIf(isSelected && sortOrder == "desc", "down").
		Else("none")

	sortingIndicator := hb.NewSpan().
		Class("sorting").
		HTMLIf(direction == "up", "&#8595;").
		HTMLIf(direction == "down", "&#8593;").
		HTMLIf(direction != "down" && direction != "up", "")

	return sortingIndicator
}
