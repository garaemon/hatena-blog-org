package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func convertOrgToMarkdown(orgFilePath string) (string, error) {
	if !fileExists(orgFilePath) {
		return "", fmt.Errorf("org file not found: %s", orgFilePath)
	}

	if !strings.HasSuffix(orgFilePath, ".org") {
		return "", fmt.Errorf("file is not an org file: %s", orgFilePath)
	}

	cmd := exec.Command("pandoc", "-f", "org", "-t", "markdown", orgFilePath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pandoc conversion failed: %v", err)
	}

	return string(output), nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getAbsPath(path string) (string, error) {
	return filepath.Abs(path)
}