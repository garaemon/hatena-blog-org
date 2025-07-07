package main

import (
	"bufio"
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

func extractTitleFromOrg(orgFilePath string) (string, error) {
	file, err := os.Open(orgFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open org file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(strings.ToLower(line), "#+title:") {
			title := strings.TrimSpace(line[8:])
			if title != "" {
				return title, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read org file: %v", err)
	}

	return "Untitled", nil
}

func getAbsPath(path string) (string, error) {
	return filepath.Abs(path)
}
