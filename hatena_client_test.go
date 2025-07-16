package main

import (
	"strings"
	"testing"
	"time"
)

func TestNewHatenaClient(t *testing.T) {
	client := NewHatenaClient("testuser", "testapi", "testblog.example.com")
	if client.HatenaID != "testuser" {
		t.Errorf("Expected HatenaID to be 'testuser', got '%s'", client.HatenaID)
	}
	if client.APIKey != "testapi" {
		t.Errorf("Expected APIKey to be 'testapi', got '%s'", client.APIKey)
	}
	if client.BlogDomain != "testblog.example.com" {
		t.Errorf("Expected BlogDomain to be 'testblog.example.com', got '%s'", client.BlogDomain)
	}
	expectedURL := "https://blog.hatena.ne.jp/testuser/testblog.example.com/atom"
	if client.BaseURL != expectedURL {
		t.Errorf("Expected BaseURL to be '%s', got '%s'", expectedURL, client.BaseURL)
	}
}

func TestCreateWSSEHeader(t *testing.T) {
	client := NewHatenaClient("testuser", "testapi", "testblog.example.com")
	header := client.createWSSEHeader()

	if !strings.Contains(header, "UsernameToken Username=\"testuser\"") {
		t.Error("WSSE header should contain username")
	}
	if !strings.Contains(header, "PasswordDigest=") {
		t.Error("WSSE header should contain password digest")
	}
	if !strings.Contains(header, "Nonce=") {
		t.Error("WSSE header should contain nonce")
	}
	if !strings.Contains(header, "Created=") {
		t.Error("WSSE header should contain created timestamp")
	}
}

func TestGenerateNonce(t *testing.T) {
	nonce1 := generateNonce()
	time.Sleep(1 * time.Millisecond)
	nonce2 := generateNonce()

	if nonce1 == nonce2 {
		t.Error("Generated nonces should be different")
	}
	if len(nonce1) == 0 || len(nonce2) == 0 {
		t.Error("Generated nonces should not be empty")
	}
}

func TestGenerateDigest(t *testing.T) {
	digest1 := generateDigest("nonce1", "created1", "password1")
	digest2 := generateDigest("nonce2", "created2", "password2")

	if digest1 == digest2 {
		t.Error("Different inputs should generate different digests")
	}
	if len(digest1) == 0 || len(digest2) == 0 {
		t.Error("Generated digests should not be empty")
	}
}

func TestCreateEntryXML(t *testing.T) {
	client := NewHatenaClient("testuser", "testapi", "testblog.example.com")
	entry := BlogEntry{
		Title:      "Test Title",
		Content:    "Test content",
		Categories: []string{"Test Category"},
		IsDraft:    true,
	}

	xml := client.createEntryXML(entry)

	if !strings.Contains(xml, "<title>Test Title</title>") {
		t.Error("XML should contain title")
	}
	if !strings.Contains(xml, "<content type=\"text/x-markdown\">Test content</content>") {
		t.Error("XML should contain content")
	}
	if !strings.Contains(xml, "<category term=\"Test Category\" />") {
		t.Error("XML should contain category")
	}
	if !strings.Contains(xml, "<app:draft>yes</app:draft>") {
		t.Error("XML should indicate draft status")
	}
	if !strings.Contains(xml, "<author><name>testuser</name></author>") {
		t.Error("XML should contain author")
	}
}

func TestCreateEntryXMLNoDraft(t *testing.T) {
	client := NewHatenaClient("testuser", "testapi", "testblog.example.com")
	entry := BlogEntry{
		Title:      "Test Title",
		Content:    "Test content",
		Categories: []string{},
		IsDraft:    false,
	}

	xml := client.createEntryXML(entry)

	if !strings.Contains(xml, "<app:draft>no</app:draft>") {
		t.Error("XML should indicate non-draft status")
	}
}

func TestExtractTitleFromMarkdown(t *testing.T) {
	markdown := `# Test Title

This is some content.

## Subsection

More content.`

	title := extractTitleFromMarkdown(markdown)
	if title != "Test Title" {
		t.Errorf("Expected title to be 'Test Title', got '%s'", title)
	}
}

func TestExtractTitleFromMarkdownNoTitle(t *testing.T) {
	markdown := `This is some content without a title.

## Subsection

More content.`

	title := extractTitleFromMarkdown(markdown)
	if title != "Untitled" {
		t.Errorf("Expected title to be 'Untitled', got '%s'", title)
	}
}

func TestRemoveTitleFromMarkdown(t *testing.T) {
	markdown := `# Test Title

This is some content.

## Subsection

More content.`

	content := removeTitleFromMarkdown(markdown)
	if content != markdown {
		t.Error("Content should remain unchanged since we no longer remove titles")
	}
}

func TestRemoveTitleFromMarkdownNoTitle(t *testing.T) {
	markdown := `This is some content without a title.

## Subsection

More content.`

	content := removeTitleFromMarkdown(markdown)
	if content != markdown {
		t.Error("Content should remain unchanged when no title is present")
	}
}

func TestCreateEntryXMLEscapesSpecialCharacters(t *testing.T) {
	client := NewHatenaClient("test<user>", "testapi", "testblog.example.com")
	entry := BlogEntry{
		Title:      "Test & Title with <tags>",
		Content:    "Test content with <script> & \"quotes\"",
		Categories: []string{"Test & Category with <tags>"},
		IsDraft:    false,
	}

	xml := client.createEntryXML(entry)

	expectedTitle := "<title>Test &amp; Title with &lt;tags&gt;</title>"
	if !strings.Contains(xml, expectedTitle) {
		t.Errorf("XML should escape title properly, expected: %s", expectedTitle)
	}

	expectedContent := "<content type=\"text/x-markdown\">Test content with &lt;script&gt; &amp; &#34;quotes&#34;</content>"
	if !strings.Contains(xml, expectedContent) {
		t.Errorf("XML should escape content properly, expected: %s", expectedContent)
	}

	expectedCategory := "<category term=\"Test &amp; Category with &lt;tags&gt;\" />"
	if !strings.Contains(xml, expectedCategory) {
		t.Errorf("XML should escape category properly, expected: %s", expectedCategory)
	}

	expectedAuthor := "<author><name>test&lt;user&gt;</name></author>"
	if !strings.Contains(xml, expectedAuthor) {
		t.Errorf("XML should escape author name properly, expected: %s", expectedAuthor)
	}
}

func TestPostEntryDebugMode(t *testing.T) {
	client := NewHatenaClient("testuser", "testapi", "testblog.example.com")
	entry := BlogEntry{
		Title:      "Test Title",
		Content:    "Test content",
		Categories: []string{"Test Category"},
		IsDraft:    false,
	}

	_, err := client.PostEntry(entry, true)
	if err == nil {
		t.Error("Expected HTTP error since we're not making a real request")
	}
}
