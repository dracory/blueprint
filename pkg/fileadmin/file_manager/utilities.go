package file_manager

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileEntry represents a file or directory entry
type FileEntry struct {
	IsDir             bool
	Path              string
	URL               string
	Name              string
	Size              int64
	SizeHuman         string
	LastModified      time.Time
	LastModifiedHuman string
	Depth             int
}

// HumanFilesize converts bytes to human readable format
func (c *FileManagerController) HumanFilesize(size int64) string {
	const unit = 1000
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(size)/float64(div), "kMGTPE"[exp])
}

// normalizePath normalizes a directory path by handling root directory representation, removing double slashes, and preventing path traversal
func verifyAndNormalizePathOrError(dir, filename string) (string, error) {
	if dir == "/" {
		dir = ""
	}

	// Check for path traversal attempts BEFORE cleaning
	// This prevents any attempt to use ".." to escape the directory structure
	if dir == "." || dir == ".." {
		return "", errors.New("path traversal detected")
	}
	if filename == "." || filename == ".." {
		return "", errors.New("path traversal detected")
	}
	if strings.HasPrefix(dir, "~") || strings.HasPrefix(filename, "~") {
		return "", errors.New("path traversal detected")
	}
	if strings.Contains(dir, "..") {
		return "", errors.New("path traversal detected")
	}
	if strings.Contains(filename, "..") {
		return "", errors.New("path traversal detected")
	}

	// Clean the filename to resolve any other relative components
	filename = filepath.Clean(filename)

	path := dir + "/" + filename
	path = strings.ReplaceAll(path, "//", "/")

	// Clean the final path to resolve any remaining relative components
	path = filepath.Clean(path)

	// Ensure forward slashes are used consistently (Windows compatibility)
	path = strings.ReplaceAll(path, "\\", "/")

	return path, nil
}

// normalizeDirPath normalizes a directory path by handling root directory representation, removing double slashes, and trimming trailing slashes
func verifyAndNormalizeDirPath(dir, filename string) (string, error) {
	path, err := verifyAndNormalizePathOrError(dir, filename)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(path, "/"), nil
}

// allDirectories recursively lists all directories starting from the given path
func (c *FileManagerController) allDirectories(rootPath string) ([]FileEntry, error) {
	return c.allDirectoriesDepth(rootPath, 0)
}

func (c *FileManagerController) allDirectoriesDepth(rootPath string, depth int) ([]FileEntry, error) {
	result := []FileEntry{}
	dirs, err := c.storage.Directories(rootPath)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		if dir == "." || dir == ".." {
			continue
		}
		size, _ := c.storage.Size(dir)
		hSize := "-"
		if size > 0 {
			hSize = c.HumanFilesize(size)
		}
		modified, _ := c.storage.LastModified(dir)
		result = append(result, FileEntry{
			Path:              dir,
			Name:              filepath.Base(dir),
			Size:              size,
			SizeHuman:         hSize,
			LastModified:      modified,
			LastModifiedHuman: "",
			Depth:             depth,
		})

		subDirs, err := c.allDirectoriesDepth(dir, depth+1)
		if err != nil {
			return nil, err
		}
		result = append(result, subDirs...)
	}

	// Sort by path for consistent ordering
	sort.Slice(result, func(i, j int) bool {
		return strings.ToLower(result[i].Path) < strings.ToLower(result[j].Path)
	})

	return result, nil
}
