package post_update

import (
	"context"
	"testing"

	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
)

// TestCreatePostVersioning_NilRegistry tests with nil registry
func TestCreatePostVersioning_NilRegistry(t *testing.T) {
	err := createPostVersioning(context.TODO(), nil, nil)
	if err == nil {
		t.Error("createPostVersioning() with nil registry should return error")
	}
	if err.Error() != "blog store not available" {
		t.Errorf("Expected 'blog store not available', got: %v", err)
	}
}

// TestCreatePostVersioning_NilPost tests with nil post
func TestCreatePostVersioning_NilPost(t *testing.T) {
	// Note: When registry doesn't have blog store, it returns "blog store not available"
	// before checking if post is nil. This test documents current behavior.
	registry := testutils.Setup()
	err := createPostVersioning(context.TODO(), registry, nil)
	if err == nil {
		t.Error("createPostVersioning() should return error")
	}
	// Registry without blog store returns this error first
	if err.Error() != "blog store not available" {
		t.Errorf("Expected 'blog store not available', got: %v", err)
	}
}

// TestCreatePostVersioning_RegistryWithoutBlogStore tests with registry that has no blog store
func TestCreatePostVersioning_RegistryWithoutBlogStore(t *testing.T) {
	registry := testutils.Setup()
	err := createPostVersioning(context.TODO(), registry, nil)
	if err == nil {
		t.Error("createPostVersioning() without blog store should return error")
	}
}

// TestCreatePostVersioning_StructFields tests error message constants
func TestCreatePostVersioning_ErrorMessages(t *testing.T) {
	// Test the error messages are correct
	tests := []struct {
		name    string
		setup   func() (context.Context, registry.RegistryInterface, blogstore.PostInterface)
		wantErr string
	}{
		{
			name: "nil registry",
			setup: func() (context.Context, registry.RegistryInterface, blogstore.PostInterface) {
				return context.TODO(), nil, nil
			},
			wantErr: "blog store not available",
		},
		{
			name: "nil post",
			setup: func() (context.Context, registry.RegistryInterface, blogstore.PostInterface) {
				return context.TODO(), testutils.Setup(), nil
			},
			wantErr: "blog store not available", // Registry checked before post
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, reg, post := tt.setup()
			err := createPostVersioning(ctx, reg, post)
			if err == nil {
				t.Error("Expected error")
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Expected %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// TestConstants tests the package constant
func TestConstants(t *testing.T) {
	if PostEditorMarkdownEasyMDE != "markdown_easymde" {
		t.Errorf("PostEditorMarkdownEasyMDE = %q, want %q", PostEditorMarkdownEasyMDE, "markdown_easymde")
	}
}
