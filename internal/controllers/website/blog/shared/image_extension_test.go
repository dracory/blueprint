package shared

import (
	"testing"
)

func TestImageExtension(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{"jpg extension", "image.jpg", "jpg"},
		{"jpeg extension", "image.jpeg", "jpg"},
		{"webp extension", "image.webp", "webp"},
		{"png extension", "image.png", "png"},
		{"no extension", "image", "png"},
		{"uppercase JPG", "image.JPG", "jpg"},
		{"uppercase JPEG", "image.JPEG", "jpg"},
		{"uppercase WEBP", "image.WEBP", "webp"},
		{"uppercase PNG", "image.PNG", "png"},
		{"mixed case", "image.JpG", "jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ImageExtension(tt.url)
			if result != tt.expected {
				t.Errorf("ImageExtension(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}
