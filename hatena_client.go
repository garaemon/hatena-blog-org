package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HatenaClient struct {
	HatenaID   string
	APIKey     string
	BlogDomain string
	BaseURL    string
}

type BlogEntry struct {
	Title    string
	Content  string
	Category string
	IsDraft  bool
}

func NewHatenaClient(hatenaID, apiKey, blogDomain string) *HatenaClient {
	return &HatenaClient{
		HatenaID:   hatenaID,
		APIKey:     apiKey,
		BlogDomain: blogDomain,
		BaseURL:    fmt.Sprintf("https://blog.hatena.ne.jp/%s/%s/atom", hatenaID, blogDomain),
	}
}

func (c *HatenaClient) createWSSEHeader() string {
	nonce := generateNonce()
	created := time.Now().Format(time.RFC3339)
	digest := generateDigest(nonce, created, c.APIKey)
	
	return fmt.Sprintf(`UsernameToken Username="%s", PasswordDigest="%s", Nonce="%s", Created="%s"`,
		c.HatenaID, digest, nonce, created)
}

func generateNonce() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
}

func generateDigest(nonce, created, password string) string {
	h := sha1.New()
	h.Write([]byte(nonce + created + password))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c *HatenaClient) createEntryXML(entry BlogEntry) string {
	draftStatus := "no"
	if entry.IsDraft {
		draftStatus = "yes"
	}

	xml := `<?xml version="1.0" encoding="utf-8"?>
<entry xmlns="http://www.w3.org/2005/Atom"
       xmlns:app="http://www.w3.org/2007/app">
  <title>%s</title>
  <author><name>%s</name></author>
  <content type="text/x-markdown">%s</content>
  <updated>%s</updated>`

	if entry.Category != "" {
		xml += fmt.Sprintf(`
  <category term="%s" />`, entry.Category)
	}

	xml += fmt.Sprintf(`
  <app:control>
    <app:draft>%s</app:draft>
  </app:control>
</entry>`, draftStatus)

	return fmt.Sprintf(xml, entry.Title, c.HatenaID, entry.Content, time.Now().Format(time.RFC3339))
}

func (c *HatenaClient) PostEntry(entry BlogEntry) error {
	entryXML := c.createEntryXML(entry)
	
	req, err := http.NewRequest("POST", c.BaseURL+"/entry", bytes.NewBufferString(entryXML))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("X-WSSE", c.createWSSEHeader())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func extractTitleFromMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}
	return "Untitled"
}

func removeTitleFromMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			if i+1 < len(lines) && lines[i+1] == "" {
				return strings.Join(lines[i+2:], "\n")
			}
			return strings.Join(lines[i+1:], "\n")
		}
	}
	return content
}