package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Pre-compiled regex patterns for better performance
	// specific fix for common "IDs" pattern to ensure "IDs" becomes "ids" or "user_ids" instead of "i_ds"
	commonAcronyms       = regexp.MustCompile(`IDs([A-Z]|$)`)
	consecutiveCapitals  = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	capitalFollowedLower = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

// toSnakeCase converts a filename to snake_case format while preserving the file extension.
func toSnakeCase(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	// Pre-process common acronyms like "IDs" -> "Ids" so they are handled correctly by subsequent rules
	// "IDs" -> "Ids" -> "ids"
	// "UserIDs" -> "UserIds" -> "user_ids"
	base = commonAcronyms.ReplaceAllString(base, `Ids${1}`)

	// Handle consecutive capitals (e.g., "XMLParser" -> "XML_Parser")
	base = consecutiveCapitals.ReplaceAllString(base, `${1}_${2}`)

	// Handle capital followed by lowercase (e.g., "Parser" -> "parser")
	// Also handles numbers followed by letters (e.g., "file2Go" -> "file2_go")
	base = capitalFollowedLower.ReplaceAllString(base, `${1}_${2}`)

	return strings.ToLower(base) + ext
}

// shouldIgnore checks if the directory or file should be ignored
func shouldIgnore(path string) bool {
	base := filepath.Base(path)
	// Ignore hidden files/directories (starting with .) and common dependency folders
	if strings.HasPrefix(base, ".") || base == "vendor" || base == "node_modules" {
		return true
	}
	return false
}

// renameGoFiles recursively renames Go files in the specified directory to snake_case format.
func renameGoFiles(root string, dryRun bool, verbose bool) error {
	// Validate directory exists
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("directory %s does not exist: %w", root, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path %s is not a directory", root)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if info.IsDir() {
			if shouldIgnore(path) && path != root {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(info.Name()) == ".go" {
			dir := filepath.Dir(path)
			snake := toSnakeCase(info.Name())

			// Skip if name hasn't changed
			if info.Name() == snake {
				if verbose {
					fmt.Printf("Already in snake_case: %s\n", path)
				}
				return nil
			}

			newPath := filepath.Join(dir, snake)

			// Check for collision
			// On case-insensitive filesystems (Windows/macOS), we need to be careful.
			// If strings.EqualFold(path, newPath) is true, it's a case-only rename.
			isCaseRename := strings.EqualFold(path, newPath)

			if !isCaseRename {
				// If it's NOT a case rename, check if the target exists
				if _, err := os.Stat(newPath); !os.IsNotExist(err) {
					if verbose {
						fmt.Fprintf(os.Stderr, "File already exists, skipping: %s\n", newPath)
					}
					return nil
				}
			}

			if dryRun {
				fmt.Printf("[DRY RUN] Would rename: %s -> %s\n", path, newPath)
			} else {
				fmt.Printf("Renaming %s -> %s\n", path, newPath)
				if err := os.Rename(path, newPath); err != nil {
					return fmt.Errorf("failed to rename %s to %s: %w", path, newPath, err)
				}
			}
		}
		return nil
	})
}

func main() {
	var (
		dir     = flag.String("dir", ".", "Directory to recursively rename Go files to snake_case")
		dryRun  = flag.Bool("dry-run", false, "Show what would be renamed without actually renaming files")
		verbose = flag.Bool("verbose", false, "Show detailed output including files that don't need renaming")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Recursively renames Go files in the specified directory to snake_case format.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                    # Rename files in current directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -dir ./src         # Rename files in src directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -dry-run           # Preview what would be renamed\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -verbose           # Show detailed output\n", os.Args[0])
	}

	flag.Parse()

	if *dryRun {
		fmt.Println("DRY RUN MODE - No files will be actually renamed")
	}

	if err := renameGoFiles(*dir, *dryRun, *verbose); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *dryRun {
		fmt.Println("Dry run completed. Use without -dry-run flag to actually rename files.")
	} else {
		fmt.Println("File renaming completed successfully.")
	}
}