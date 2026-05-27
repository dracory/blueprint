package file_manager

import (
	"testing"
)

func TestNormalizePath_EmptyDirWithFilename(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("", "file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "", "file.txt", err)
	}
	if got != "/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "", "file.txt", got, "/file.txt")
	}
}

func TestNormalizePath_RootDirWithFilename(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("/", "file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "/", "file.txt", err)
	}
	if got != "/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "/", "file.txt", got, "/file.txt")
	}
}

func TestNormalizePath_SubdirectoryWithFilename(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("documents", "file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "documents", "file.txt", err)
	}
	if got != "documents/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "documents", "file.txt", got, "documents/file.txt")
	}
}

func TestNormalizePath_NestedDirectoryWithFilename(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("documents/reports", "file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "documents/reports", "file.txt", err)
	}
	if got != "documents/reports/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "documents/reports", "file.txt", got, "documents/reports/file.txt")
	}
}

func TestNormalizePath_DirectoryWithTrailingSlash(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("documents/", "file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "documents/", "file.txt", err)
	}
	if got != "documents/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "documents/", "file.txt", got, "documents/file.txt")
	}
}

func TestNormalizePath_EmptyDirWithNestedFilename(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("", "documents/file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "", "documents/file.txt", err)
	}
	if got != "/documents/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "", "documents/file.txt", got, "/documents/file.txt")
	}
}

func TestNormalizePath_RootDirWithNestedFilename(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("/", "documents/file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "/", "documents/file.txt", err)
	}
	if got != "/documents/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "/", "documents/file.txt", got, "/documents/file.txt")
	}
}

func TestNormalizeDirPath_EmptyDirWithDirname(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("", "folder")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "", "folder", err)
	}
	if got != "/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "", "folder", got, "/folder")
	}
}

func TestNormalizeDirPath_RootDirWithDirname(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("/", "folder")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "/", "folder", err)
	}
	if got != "/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "/", "folder", got, "/folder")
	}
}

func TestNormalizeDirPath_SubdirectoryWithDirname(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("documents", "folder")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "documents", "folder", err)
	}
	if got != "documents/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "documents", "folder", got, "documents/folder")
	}
}

func TestNormalizeDirPath_NestedDirectoryWithDirname(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("documents/reports", "folder")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "documents/reports", "folder", err)
	}
	if got != "documents/reports/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "documents/reports", "folder", got, "documents/reports/folder")
	}
}

func TestNormalizeDirPath_DirectoryWithTrailingSlash(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("documents/", "folder")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "documents/", "folder", err)
	}
	if got != "documents/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "documents/", "folder", got, "documents/folder")
	}
}

func TestNormalizeDirPath_DirnameWithTrailingSlash(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("", "folder/")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "", "folder/", err)
	}
	if got != "/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "", "folder/", got, "/folder")
	}
}

func TestNormalizeDirPath_BothDirAndFilenameWithTrailingSlashes(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("documents/", "folder/")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "documents/", "folder/", err)
	}
	if got != "documents/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "documents/", "folder/", got, "documents/folder")
	}
}

func TestNormalizeDirPath_EmptyDirWithNestedDirname(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("", "documents/folder")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "", "documents/folder", err)
	}
	if got != "/documents/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "", "documents/folder", got, "/documents/folder")
	}
}

func TestNormalizeDirPath_RootDirWithNestedDirname(t *testing.T) {
	got, err := verifyAndNormalizeDirPath("/", "documents/folder")
	if err != nil {
		t.Errorf("normalizeDirPath(%q, %q) unexpected error: %v", "/", "documents/folder", err)
	}
	if got != "/documents/folder" {
		t.Errorf("normalizeDirPath(%q, %q) = %q, want %q", "/", "documents/folder", got, "/documents/folder")
	}
}

func TestNormalizePathSecurity_PathTraversalWithSingleDotDot(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("", "..")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "", "..")
	}
}

func TestNormalizePathSecurity_PathTraversalWithDoubleDotSlash(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("", "../")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "", "../")
	}
}

func TestNormalizePathSecurity_PathTraversalWithDoubleDotPrefix(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("", "../file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "", "../file.txt")
	}
}

func TestNormalizePathSecurity_PathTraversalWithMultipleDoubleDots(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("", "../../file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "", "../../file.txt")
	}
}

func TestNormalizePathSecurity_PathTraversalInMiddleOfPath(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("documents", "../file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "documents", "../file.txt")
	}
}

func TestNormalizePathSecurity_PathTraversalWithBackslash(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("", "..\\file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "", "..\\file.txt")
	}
}

func TestNormalizePathSecurity_PathTraversalWithMixedSeparators(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("", "..\\../file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "", "..\\../file.txt")
	}
}

func TestNormalizePathSecurity_PathTraversalWithEncodedDot(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("", "%2e%2e/file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "", "%2e%2e/file.txt", err)
	}
	if got != "/%2e%2e/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "", "%2e%2e/file.txt", got, "/%2e%2e/file.txt")
	}
}

func TestNormalizePathSecurity_CurrentDirectoryReference(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("", "./file.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "", "./file.txt", err)
	}
	if got != "/file.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "", "./file.txt", got, "/file.txt")
	}
}

func TestNormalizePathSecurity_PathTraversalWithDirectoryAndDoubleDot(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("documents", "reports/../file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "documents", "reports/../file.txt")
	}
}

func TestNormalizePathSecurity_ComplexPathTraversalAttempt(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("documents", "./reports/../../secret/file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "documents", "./reports/../../secret/file.txt")
	}
}

func TestNormalizePathSecurity_PathTraversalEscapingRoot(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("/", "../../etc/passwd")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "/", "../../etc/passwd")
	}
}

func TestNormalizePathSecurity_PathTraversalInDirParameter(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("../uploads", "file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "../uploads", "file.txt")
	}
}

func TestNormalizePathSecurity_PathTraversalInDirParameterWithMultipleLevels(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("../../uploads", "file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "../../uploads", "file.txt")
	}
}

func TestNormalizePathSecurity_ExactDotInDirParameter(t *testing.T) {
	_, err := verifyAndNormalizePathOrError(".", "file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", ".", "file.txt")
	}
}

func TestNormalizePathSecurity_ExactDotInFilename(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("", ".")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "", ".")
	}
}

func TestNormalizePathSecurity_PathStartingWithTildeInFilename(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("", "~/file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "", "~/file.txt")
	}
}

func TestNormalizePathSecurity_PathStartingWithTildeInDir(t *testing.T) {
	_, err := verifyAndNormalizePathOrError("~/uploads", "file.txt")
	if err == nil {
		t.Errorf("normalizePath(%q, %q) expected error but got none", "~/uploads", "file.txt")
	}
}

func TestNormalizePathSecurity_LegitimateFileWithDotInName(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("", ".hiddenfile")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "", ".hiddenfile", err)
	}
	if got != "/.hiddenfile" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "", ".hiddenfile", got, "/.hiddenfile")
	}
}

func TestNormalizePathSecurity_LegitimateFileWithMultipleDots(t *testing.T) {
	got, err := verifyAndNormalizePathOrError("", "file.name.with.dots.txt")
	if err != nil {
		t.Errorf("normalizePath(%q, %q) unexpected error: %v", "", "file.name.with.dots.txt", err)
	}
	if got != "/file.name.with.dots.txt" {
		t.Errorf("normalizePath(%q, %q) = %q, want %q", "", "file.name.with.dots.txt", got, "/file.name.with.dots.txt")
	}
}
