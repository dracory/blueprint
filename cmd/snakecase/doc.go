// Package snakecase provides utilities for converting filenames to snake_case format.
//
// This package is designed to help maintain consistent naming conventions in Go projects
// by automatically renaming Go files to follow snake_case naming standards.
//
// Command line usage:
//
//	go run ./cmd/snakecase -dir ./src -dry-run -verbose
//
// Flags:
//
//	-dir string
//	      Directory to recursively rename Go files to snake_case (default ".")
//	-dry-run
//	      Show what would be renamed without actually renaming files
//	-verbose
//	      Show detailed output including files that don't need renaming
//
// The conversion rules are:
//   - Consecutive capital letters followed by a capital+lowercase pattern are split
//     (e.g., "XMLParser" -> "xml_parser")
//   - A capital letter preceded by a lowercase letter or number gets an underscore
//     (e.g., "Parser" -> "parser", "file2Go" -> "file2_go")
//   - File extensions are preserved
//   - Already snake_case files are left unchanged
package main
