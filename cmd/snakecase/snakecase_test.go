package main

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple.go", "simple.go"},
		{"camelCase.go", "camel_case.go"},
		{"PascalCase.go", "pascal_case.go"},
		{"HTTPClient.go", "http_client.go"},
		{"XMLParser.go", "xml_parser.go"},
		{"MyXMLParser.go", "my_xml_parser.go"},
		{"file2Go.go", "file2_go.go"},
		{"already_snake.go", "already_snake.go"},
		{"IDs.go", "ids.go"},           // Updated expectation
		{"UserIDs.go", "user_ids.go"}, // Updated expectation
		{"Server123.go", "server123.go"},
		{"JSONToXML.go", "json_to_xml.go"},
		{"with-dash.go", "with-dash.go"},
	}

	for _, test := range tests {
		result := toSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("toSnakeCase(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}