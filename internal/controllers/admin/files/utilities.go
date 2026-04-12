package admin

import (
	"fmt"
)

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
		if exp >= 5 { // Prevent index out of bounds for "kMGTPE"[exp]
			break
		}
	}
	return fmt.Sprintf("%.1f %cB",
		float64(size)/float64(div), "kMGTPE"[exp])
}
