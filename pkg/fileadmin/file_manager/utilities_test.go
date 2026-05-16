package file_manager

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		filename string
		want     string
	}{
		{
			name:     "empty dir with filename",
			dir:      "",
			filename: "file.txt",
			want:     "/file.txt",
		},
		{
			name:     "root dir with filename",
			dir:      "/",
			filename: "file.txt",
			want:     "/file.txt",
		},
		{
			name:     "subdirectory with filename",
			dir:      "documents",
			filename: "file.txt",
			want:     "documents/file.txt",
		},
		{
			name:     "nested directory with filename",
			dir:      "documents/reports",
			filename: "file.txt",
			want:     "documents/reports/file.txt",
		},
		{
			name:     "directory with trailing slash",
			dir:      "documents/",
			filename: "file.txt",
			want:     "documents/file.txt",
		},
		{
			name:     "empty dir with nested filename",
			dir:      "",
			filename: "documents/file.txt",
			want:     "/documents/file.txt",
		},
		{
			name:     "root dir with nested filename",
			dir:      "/",
			filename: "documents/file.txt",
			want:     "/documents/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifyAndNormalizePathOrError(tt.dir, tt.filename)
			if err != nil {
				t.Errorf("normalizePath(%q, %q) unexpected error: %v", tt.dir, tt.filename, err)
			}
			if got != tt.want {
				t.Errorf("normalizePath(%q, %q) = %q, want %q", tt.dir, tt.filename, got, tt.want)
			}
		})
	}
}

func TestNormalizeDirPath(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		filename string
		want     string
	}{
		{
			name:     "empty dir with dirname",
			dir:      "",
			filename: "folder",
			want:     "/folder",
		},
		{
			name:     "root dir with dirname",
			dir:      "/",
			filename: "folder",
			want:     "/folder",
		},
		{
			name:     "subdirectory with dirname",
			dir:      "documents",
			filename: "folder",
			want:     "documents/folder",
		},
		{
			name:     "nested directory with dirname",
			dir:      "documents/reports",
			filename: "folder",
			want:     "documents/reports/folder",
		},
		{
			name:     "directory with trailing slash",
			dir:      "documents/",
			filename: "folder",
			want:     "documents/folder",
		},
		{
			name:     "dirname with trailing slash",
			dir:      "",
			filename: "folder/",
			want:     "/folder",
		},
		{
			name:     "both dir and filename with trailing slashes",
			dir:      "documents/",
			filename: "folder/",
			want:     "documents/folder",
		},
		{
			name:     "empty dir with nested dirname",
			dir:      "",
			filename: "documents/folder",
			want:     "/documents/folder",
		},
		{
			name:     "root dir with nested dirname",
			dir:      "/",
			filename: "documents/folder",
			want:     "/documents/folder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifyAndNormalizeDirPath(tt.dir, tt.filename)
			if err != nil {
				t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", tt.dir, tt.filename, err)
			}
			if got != tt.want {
				t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", tt.dir, tt.filename, got, tt.want)
			}
		})
	}
}

func TestNormalizePathSecurity(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		filename    string
		want        string
		expectError bool
	}{
		{
			name:        "path traversal with single dot dot",
			dir:         "",
			filename:    "..",
			expectError: true,
		},
		{
			name:        "path traversal with double dot slash",
			dir:         "",
			filename:    "../",
			expectError: true,
		},
		{
			name:        "path traversal with double dot prefix",
			dir:         "",
			filename:    "../file.txt",
			expectError: true,
		},
		{
			name:        "path traversal with multiple double dots",
			dir:         "",
			filename:    "../../file.txt",
			expectError: true,
		},
		{
			name:        "path traversal in middle of path",
			dir:         "documents",
			filename:    "../file.txt",
			expectError: true,
		},
		{
			name:        "path traversal with backslash (Windows)",
			dir:         "",
			filename:    "..\\file.txt",
			expectError: true,
		},
		{
			name:        "path traversal with mixed separators",
			dir:         "",
			filename:    "..\\../file.txt",
			expectError: true,
		},
		{
			name:     "path traversal with encoded dot",
			dir:      "",
			filename: "%2e%2e/file.txt",
			want:     "/%2e%2e/file.txt",
		},
		{
			name:     "current directory reference",
			dir:      "",
			filename: "./file.txt",
			want:     "/file.txt",
		},
		{
			name:        "path traversal with directory and double dot",
			dir:         "documents",
			filename:    "reports/../file.txt",
			expectError: true,
		},
		{
			name:        "complex path traversal attempt",
			dir:         "documents",
			filename:    "./reports/../../secret/file.txt",
			expectError: true,
		},
		{
			name:        "path traversal escaping root",
			dir:         "/",
			filename:    "../../etc/passwd",
			expectError: true,
		},
		{
			name:        "path traversal in dir parameter",
			dir:         "../uploads",
			filename:    "file.txt",
			expectError: true,
		},
		{
			name:        "path traversal in dir parameter with multiple levels",
			dir:         "../../uploads",
			filename:    "file.txt",
			expectError: true,
		},
		{
			name:        "exact dot in dir parameter",
			dir:         ".",
			filename:    "file.txt",
			expectError: true,
		},
		{
			name:        "exact dot in filename",
			dir:         "",
			filename:    ".",
			expectError: true,
		},
		{
			name:        "path starting with tilde in filename",
			dir:         "",
			filename:    "~/file.txt",
			expectError: true,
		},
		{
			name:        "path starting with tilde in dir",
			dir:         "~/uploads",
			filename:    "file.txt",
			expectError: true,
		},
		{
			name:     "legitimate file with dot in name",
			dir:      "",
			filename: ".hiddenfile",
			want:     "/.hiddenfile",
		},
		{
			name:     "legitimate file with multiple dots",
			dir:      "",
			filename: "file.name.with.dots.txt",
			want:     "/file.name.with.dots.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifyAndNormalizePathOrError(tt.dir, tt.filename)
			if tt.expectError {
				if err == nil {
					t.Errorf("normalizePath(%q, %q) expected error but got none", tt.dir, tt.filename)
				}
				return
			}
			if err != nil {
				t.Errorf("normalizePath(%q, %q) unexpected error: %v", tt.dir, tt.filename, err)
			}
			if got != tt.want {
				t.Errorf("normalizePath(%q, %q) = %q, want %q", tt.dir, tt.filename, got, tt.want)
			}
		})
	}
}
