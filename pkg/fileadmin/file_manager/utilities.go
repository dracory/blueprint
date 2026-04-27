package file_manager

import (
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
