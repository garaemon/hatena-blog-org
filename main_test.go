package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)


type mockHatenaClient struct {
	uploadedImages map[string]string
	shouldError    bool
}

func (m *mockHatenaClient) uploadImage(imagePath string) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("mock upload error")
	}
	
	if m.uploadedImages == nil {
		m.uploadedImages = make(map[string]string)
	}
	
	filename := filepath.Base(imagePath)
	mockURL := "https://cdn.hatena.ne.jp/fotolife/user/" + filename
	m.uploadedImages[imagePath] = mockURL
	
	return mockURL, nil
}

func TestProcessImageUploads(t *testing.T) {
	tests := []struct {
		name         string
		markdown     string
		imageLinks   []string
		shouldError  bool
		expectedMD   string
		setupFiles   map[string]string
	}{
		{
			name: "single image upload success",
			markdown: `# Test Article

Here is an image:
![test.jpg](./test.jpg)

End of article.`,
			imageLinks: []string{"./test.jpg"},
			shouldError: false,
			expectedMD: `# Test Article

Here is an image:
![test.jpg](https://cdn.hatena.ne.jp/fotolife/user/test.jpg)

End of article.`,
			setupFiles: map[string]string{
				"./test.jpg": "fake image content",
			},
		},
		{
			name: "multiple images upload success",
			markdown: `# Test Article

First image: ![img1.jpg](./img1.jpg)
Second image: ![img2.png](./img2.png)

End of article.`,
			imageLinks: []string{"./img1.jpg", "./img2.png"},
			shouldError: false,
			expectedMD: `# Test Article

First image: ![img1.jpg](https://cdn.hatena.ne.jp/fotolife/user/img1.jpg)
Second image: ![img2.png](https://cdn.hatena.ne.jp/fotolife/user/img2.png)

End of article.`,
			setupFiles: map[string]string{
				"./img1.jpg": "fake image 1",
				"./img2.png": "fake image 2",
			},
		},
		{
			name: "no images",
			markdown: `# Test Article

No images here.

End of article.`,
			imageLinks: []string{},
			shouldError: false,
			expectedMD: `# Test Article

No images here.

End of article.`,
			setupFiles: map[string]string{},
		},
		{
			name: "image file not found",
			markdown: `# Test Article

Here is an image:
![missing.jpg](./missing.jpg)

End of article.`,
			imageLinks: []string{"./missing.jpg"},
			shouldError: false,
			expectedMD: `# Test Article

Here is an image:
![missing.jpg](./missing.jpg)

End of article.`,
			setupFiles: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			
			err := os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}

			for filePath, content := range tt.setupFiles {
				dir := filepath.Dir(filePath)
				if dir != "." && dir != "/" {
					err := os.MkdirAll(dir, 0755)
					if err != nil {
						t.Fatalf("Failed to create directory %s: %v", dir, err)
					}
				}
				
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", filePath, err)
				}
			}

			mockClient := &mockHatenaClient{shouldError: tt.shouldError}
			
			result, err := processImageUploads(tt.markdown, tt.imageLinks, mockClient)
			
			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
				return
			}
			
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expectedMD {
				t.Errorf("Expected markdown:\n%q\nGot:\n%q", tt.expectedMD, result)
			}
		})
	}
}

func TestProcessImageUploadsUploadError(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	err := os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	err = os.WriteFile("test.jpg", []byte("fake image"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	markdown := "![test.jpg](./test.jpg)"
	imageLinks := []string{"./test.jpg"}
	
	mockClient := &mockHatenaClient{shouldError: true}
	
	_, err = processImageUploads(markdown, imageLinks, mockClient)
	if err == nil {
		t.Error("Expected error for upload failure")
	}

	if !strings.Contains(err.Error(), "failed to upload image") {
		t.Errorf("Expected upload error message, got: %v", err)
	}
}