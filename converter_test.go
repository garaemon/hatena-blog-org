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
		name          string
		orgContent    string
		expectedTitle string
	}{
		{
			name:          "lowercase title",
			orgContent:    "#+title: Test Title\n\nContent here",
			expectedTitle: "Test Title",
		},
		{
			name:          "uppercase title",
			orgContent:    "#+TITLE: Uppercase Title\n\nContent here",
			expectedTitle: "Uppercase Title",
		},
		{
			name:          "mixed case title",
			orgContent:    "#+Title: Mixed Case Title\n\nContent here",
			expectedTitle: "Mixed Case Title",
		},
		{
			name:          "title with extra spaces",
			orgContent:    "#+title:   Spaced Title   \n\nContent here",
			expectedTitle: "Spaced Title",
		},
		{
			name:          "no title",
			orgContent:    "Just some content without title",
			expectedTitle: "Untitled",
		},
		{
			name:          "empty title",
			orgContent:    "#+title:\n\nContent here",
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

func TestExtractCategoriesFromOrg(t *testing.T) {
	tests := []struct {
		name               string
		orgContent         string
		expectedCategories []string
	}{
		{
			name:               "lowercase filetags",
			orgContent:         "#+filetags: tag1 tag2 tag3\n\nContent here",
			expectedCategories: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:               "uppercase filetags",
			orgContent:         "#+FILETAGS: TAG1 TAG2\n\nContent here",
			expectedCategories: []string{"TAG1", "TAG2"},
		},
		{
			name:               "mixed case filetags",
			orgContent:         "#+FileTags: MixedCase Another\n\nContent here",
			expectedCategories: []string{"MixedCase", "Another"},
		},
		{
			name:               "filetags with colons",
			orgContent:         "#+filetags: :tag1: :tag2: :tag3:\n\nContent here",
			expectedCategories: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:               "filetags colon-separated",
			orgContent:         "#+filetags: tag1:tag2:tag3\n\nContent here",
			expectedCategories: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:               "filetags mixed separators",
			orgContent:         "#+filetags: tag1:tag2 tag3:tag4\n\nContent here",
			expectedCategories: []string{"tag1", "tag2", "tag3", "tag4"},
		},
		{
			name:               "filetags with extra spaces",
			orgContent:         "#+filetags:   tag1   tag2   tag3   \n\nContent here",
			expectedCategories: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:               "no filetags",
			orgContent:         "Just some content without filetags",
			expectedCategories: []string{},
		},
		{
			name:               "empty filetags",
			orgContent:         "#+filetags:\n\nContent here",
			expectedCategories: []string{},
		},
		{
			name:               "single tag",
			orgContent:         "#+filetags: singletag\n\nContent here",
			expectedCategories: []string{"singletag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(os.TempDir(), "test_categories.org")
			err := os.WriteFile(tmpFile, []byte(tt.orgContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile)

			categories, err := extractCategoriesFromOrg(tmpFile)
			if err != nil {
				t.Fatalf("extractCategoriesFromOrg failed: %v", err)
			}

			if len(categories) != len(tt.expectedCategories) {
				t.Errorf("Expected %d categories, got %d", len(tt.expectedCategories), len(categories))
			}

			for i, expectedCategory := range tt.expectedCategories {
				if i >= len(categories) || categories[i] != expectedCategory {
					t.Errorf("Expected category %q at index %d, got %q", expectedCategory, i, categories[i])
				}
			}
		})
	}
}

func TestExtractCategoriesFromOrgFileNotFound(t *testing.T) {
	_, err := extractCategoriesFromOrg("/nonexistent/file.org")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestFilterOrgMetadata(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "filter org filetags block",
			input:    "```{=org}\n#+filetags: :emacs:go:claude:\n```\n\n# Title\n\nContent here",
			expected: "\n# Title\n\nContent here",
		},
		{
			name:     "no org metadata",
			input:    "# Title\n\nContent here",
			expected: "# Title\n\nContent here",
		},
		{
			name:     "multiple org blocks",
			input:    "```{=org}\n#+title: Test\n```\n\n# Title\n\n```{=org}\n#+filetags: :test:\n```\n\nContent",
			expected: "\n# Title\n\n\nContent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterOrgMetadata(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
			}
		})
	}
}

func TestExtractImageLinks(t *testing.T) {
	tmpDir := t.TempDir()
	orgFile := filepath.Join(tmpDir, "test.org")
	
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "single image with file: prefix",
			content: `#+TITLE: Test Article

This is an image:
[[file:./images/test.jpg]]

End of article.`,
			expected: []string{filepath.Join(tmpDir, "images", "test.jpg")},
		},
		{
			name: "single image without file: prefix",
			content: `#+TITLE: Test Article

This is an image:
[[./images/test.png]]

End of article.`,
			expected: []string{filepath.Join(tmpDir, "images", "test.png")},
		},
		{
			name: "multiple images",
			content: `#+TITLE: Test Article

First image: [[file:./img1.jpg]]
Second image: [[./img2.png]]
Third image: [[file:/absolute/path/img3.gif]]

End of article.`,
			expected: []string{
				filepath.Join(tmpDir, "img1.jpg"),
				filepath.Join(tmpDir, "img2.png"),
				"/absolute/path/img3.gif",
			},
		},
		{
			name: "no images",
			content: `#+TITLE: Test Article

No images here.
[[https://example.com/link]]

End of article.`,
			expected: []string{},
		},
		{
			name: "mixed image formats",
			content: `#+TITLE: Test Article

JPG: [[file:./test.jpg]]
JPEG: [[./test.jpeg]]
PNG: [[file:./test.png]]
GIF: [[./test.gif]]
BMP: [[file:./test.bmp]]
WEBP: [[./test.webp]]
SVG: [[file:./test.svg]]

End of article.`,
			expected: []string{
				filepath.Join(tmpDir, "test.jpg"),
				filepath.Join(tmpDir, "test.jpeg"),
				filepath.Join(tmpDir, "test.png"),
				filepath.Join(tmpDir, "test.gif"),
				filepath.Join(tmpDir, "test.bmp"),
				filepath.Join(tmpDir, "test.webp"),
				filepath.Join(tmpDir, "test.svg"),
			},
		},
		{
			name: "same line multiple images",
			content: `#+TITLE: Test Article

Multiple on same line: [[./img1.jpg]] and [[file:./img2.png]]

End of article.`,
			expected: []string{
				filepath.Join(tmpDir, "img1.jpg"),
				filepath.Join(tmpDir, "img2.png"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile(orgFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result, err := extractImageLinks(orgFile)
			if err != nil {
				t.Fatalf("extractImageLinks() error = %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("extractImageLinks() returned %d links, want %d", len(result), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("extractImageLinks()[%d] = %q, want %q", i, result[i], expected)
				}
			}
		})
	}
}

func TestExtractImageLinksFileNotFound(t *testing.T) {
	_, err := extractImageLinks("/nonexistent/file.org")
	if err == nil {
		t.Error("extractImageLinks() should return error for nonexistent file")
	}
}

func isPandocAvailable() bool {
	_, err := exec.LookPath("pandoc")
	return err == nil
}
