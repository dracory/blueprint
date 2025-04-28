package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func toSnakeCase(name string) string {
	// Remove file extension
	base := strings.TrimSuffix(name, filepath.Ext(name))

	// Insert underscores before capitals and lowercase everything
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	base = re.ReplaceAllString(base, `${1}_${2}`)
	return strings.ToLower(base) + filepath.Ext(name)
}

func renameGoFiles(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) == ".go" {
			dir := filepath.Dir(path)
			snake := toSnakeCase(info.Name())
			if info.Name() != snake {
				newPath := filepath.Join(dir, snake)
				if _, err := os.Stat(newPath); err == nil {
					fmt.Fprintf(os.Stderr, "File already exists, skipping: %s\n", newPath)
					return nil
				}
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
	dir := flag.String("dir", ".", "Directory to recursively rename Go files to snake_case")
	flag.Parse()
	if err := renameGoFiles(*dir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
