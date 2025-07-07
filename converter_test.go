package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestConvertOrgToMarkdown(t *testing.T) {
	if !isPandocAvailable() {
		t.Skip("pandoc not available, skipping test")
	}

	tmpFile := filepath.Join(os.TempDir(), "test.org")
	orgContent := `* Test Title

This is a test paragraph.

** Subsection

- List item 1
- List item 2
`
	err := os.WriteFile(tmpFile, []byte(orgContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	markdown, err := convertOrgToMarkdown(tmpFile)
	if err != nil {
		t.Fatalf("convertOrgToMarkdown failed: %v", err)
	}

	if markdown == "" {
		t.Error("Expected non-empty markdown output")
	}

	if len(markdown) < 10 {
		t.Error("Expected longer markdown output")
	}
}

func TestConvertOrgToMarkdownFileNotFound(t *testing.T) {
	_, err := convertOrgToMarkdown("/nonexistent/file.org")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestConvertOrgToMarkdownWrongExtension(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "test.txt")
	err := os.WriteFile(tmpFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	_, err = convertOrgToMarkdown(tmpFile)
	if err == nil {
		t.Error("Expected error for wrong file extension")
	}
}

func TestFileExists(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "test_exists.org")
	err := os.WriteFile(tmpFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	if !fileExists(tmpFile) {
		t.Error("Expected file to exist")
	}

	if fileExists("/nonexistent/file.org") {
		t.Error("Expected file to not exist")
	}
}

func TestGetAbsPath(t *testing.T) {
	path, err := getAbsPath("./test.org")
	if err != nil {
		t.Fatalf("getAbsPath failed: %v", err)
	}

	if !filepath.IsAbs(path) {
		t.Error("Expected absolute path")
	}
}

func TestExtractTitleFromOrg(t *testing.T) {
	tests := []struct {
		name        string
		orgContent  string
		expectedTitle string
	}{
		{
			name:        "lowercase title",
			orgContent:  "#+title: Test Title\n\nContent here",
			expectedTitle: "Test Title",
		},
		{
			name:        "uppercase title",
			orgContent:  "#+TITLE: Uppercase Title\n\nContent here",
			expectedTitle: "Uppercase Title",
		},
		{
			name:        "mixed case title",
			orgContent:  "#+Title: Mixed Case Title\n\nContent here",
			expectedTitle: "Mixed Case Title",
		},
		{
			name:        "title with extra spaces",
			orgContent:  "#+title:   Spaced Title   \n\nContent here",
			expectedTitle: "Spaced Title",
		},
		{
			name:        "no title",
			orgContent:  "Just some content without title",
			expectedTitle: "Untitled",
		},
		{
			name:        "empty title",
			orgContent:  "#+title:\n\nContent here",
			expectedTitle: "Untitled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(os.TempDir(), "test_title.org")
			err := os.WriteFile(tmpFile, []byte(tt.orgContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile)

			title, err := extractTitleFromOrg(tmpFile)
			if err != nil {
				t.Fatalf("extractTitleFromOrg failed: %v", err)
			}

			if title != tt.expectedTitle {
				t.Errorf("Expected title %q, got %q", tt.expectedTitle, title)
			}
		})
	}
}

func TestExtractTitleFromOrgFileNotFound(t *testing.T) {
	_, err := extractTitleFromOrg("/nonexistent/file.org")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func isPandocAvailable() bool {
	_, err := exec.LookPath("pandoc")
	return err == nil
}
