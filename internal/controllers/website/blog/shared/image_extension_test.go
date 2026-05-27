package shared

import (
	"testing"
)

func TestImageExtension_JPG(t *testing.T) {
	result := ImageExtension("image.jpg")
	if result != "jpg" {
		t.Errorf("ImageExtension(\"image.jpg\") = %q, want \"jpg\"", result)
	}
}

func TestImageExtension_JPEG(t *testing.T) {
	result := ImageExtension("image.jpeg")
	if result != "jpg" {
		t.Errorf("ImageExtension(\"image.jpeg\") = %q, want \"jpg\"", result)
	}
}

func TestImageExtension_WEBP(t *testing.T) {
	result := ImageExtension("image.webp")
	if result != "webp" {
		t.Errorf("ImageExtension(\"image.webp\") = %q, want \"webp\"", result)
	}
}

func TestImageExtension_PNG(t *testing.T) {
	result := ImageExtension("image.png")
	if result != "png" {
		t.Errorf("ImageExtension(\"image.png\") = %q, want \"png\"", result)
	}
}

func TestImageExtension_NoExtension(t *testing.T) {
	result := ImageExtension("image")
	if result != "png" {
		t.Errorf("ImageExtension(\"image\") = %q, want \"png\"", result)
	}
}

func TestImageExtension_UppercaseJPG(t *testing.T) {
	result := ImageExtension("image.JPG")
	if result != "jpg" {
		t.Errorf("ImageExtension(\"image.JPG\") = %q, want \"jpg\"", result)
	}
}

func TestImageExtension_UppercaseJPEG(t *testing.T) {
	result := ImageExtension("image.JPEG")
	if result != "jpg" {
		t.Errorf("ImageExtension(\"image.JPEG\") = %q, want \"jpg\"", result)
	}
}

func TestImageExtension_UppercaseWEBP(t *testing.T) {
	result := ImageExtension("image.WEBP")
	if result != "webp" {
		t.Errorf("ImageExtension(\"image.WEBP\") = %q, want \"webp\"", result)
	}
}

func TestImageExtension_UppercasePNG(t *testing.T) {
	result := ImageExtension("image.PNG")
	if result != "png" {
		t.Errorf("ImageExtension(\"image.PNG\") = %q, want \"png\"", result)
	}
}

func TestImageExtension_MixedCase(t *testing.T) {
	result := ImageExtension("image.JpG")
	if result != "jpg" {
		t.Errorf("ImageExtension(\"image.JpG\") = %q, want \"jpg\"", result)
	}
}
