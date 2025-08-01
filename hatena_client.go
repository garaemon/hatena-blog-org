package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"path"
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
	Title      string
	Content    string
	Categories []string
	IsDraft    bool
}

type AtomEntry struct {
	XMLName xml.Name `xml:"entry"`
	ID      string   `xml:"id"`
	Links   []struct {
		Rel  string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
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

	for _, category := range entry.Categories {
		if category != "" {
			xml += fmt.Sprintf(`
  <category term="%s" />`, html.EscapeString(category))
		}
	}

	xml += fmt.Sprintf(`
  <app:control>
    <app:draft>%s</app:draft>
  </app:control>
</entry>`, draftStatus)

	return fmt.Sprintf(xml, html.EscapeString(entry.Title), html.EscapeString(c.HatenaID), html.EscapeString(entry.Content), time.Now().Format(time.RFC3339))
}

func (c *HatenaClient) PostEntry(entry BlogEntry, debug bool) (string, error) {
	entryXML := c.createEntryXML(entry)
	if debug {
		fmt.Println("Generated XML:")
		fmt.Println(entryXML)
	}
	req, err := http.NewRequest("POST", c.BaseURL+"/entry", bytes.NewBufferString(entryXML))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("X-WSSE", c.createWSSEHeader())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var atomEntry AtomEntry
	if err := xml.Unmarshal(body, &atomEntry); err != nil {
		return "", fmt.Errorf("failed to parse response XML: %v", err)
	}

	var editURL string
	for _, link := range atomEntry.Links {
		if link.Rel == "edit" {
			editURL = link.Href
			break
		}
	}

	if editURL == "" {
		return "", fmt.Errorf("edit link not found in API response")
	}

	editPageURL := fmt.Sprintf("https://blog.hatena.ne.jp/%s/%s/edit?entry=%s", c.HatenaID, c.BlogDomain, extractEntryIDFromURL(editURL))
	return editPageURL, nil
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
	return content
}

func extractEntryIDFromURL(editURL string) string {
	return path.Base(editURL)
}
