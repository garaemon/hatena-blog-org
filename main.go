package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		orgFile     = flag.String("file", "", "Path to the .org file")
		hatenaID    = flag.String("id", "", "Hatena ID")
		apiKey      = flag.String("key", "", "API Key")
		blogDomain  = flag.String("domain", "", "Blog domain")
		category    = flag.String("category", "", "Category for the blog post")
		isDraft     = flag.Bool("draft", false, "Post as draft")
		configFile  = flag.String("config", "", "Path to config file")
		interactive = flag.Bool("interactive", false, "Interactive mode")
	)
	flag.Parse()

	if *interactive {
		runInteractiveMode()
		return
	}

	if *orgFile == "" {
		fmt.Println("Error: -file is required")
		flag.Usage()
		os.Exit(1)
	}

	config, err := loadConfig(*configFile, *hatenaID, *apiKey, *blogDomain)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := validateConfig(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := postOrgFile(*orgFile, config, *category, *isDraft); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully posted to Hatena Blog!")
}

func runInteractiveMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter path to .org file: ")
	orgFile, _ := reader.ReadString('\n')
	orgFile = strings.TrimSpace(orgFile)

	fmt.Print("Enter Hatena ID: ")
	hatenaID, _ := reader.ReadString('\n')
	hatenaID = strings.TrimSpace(hatenaID)

	fmt.Print("Enter API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	fmt.Print("Enter blog domain: ")
	blogDomain, _ := reader.ReadString('\n')
	blogDomain = strings.TrimSpace(blogDomain)

	fmt.Print("Enter category (optional): ")
	category, _ := reader.ReadString('\n')
	category = strings.TrimSpace(category)

	fmt.Print("Post as draft? (y/n): ")
	draftInput, _ := reader.ReadString('\n')
	isDraft := strings.TrimSpace(strings.ToLower(draftInput)) == "y"

	config := &Config{
		HatenaID:   hatenaID,
		APIKey:     apiKey,
		BlogDomain: blogDomain,
	}

	if err := validateConfig(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := postOrgFile(orgFile, config, category, isDraft); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully posted to Hatena Blog!")
}

func postOrgFile(orgFile string, config *Config, category string, isDraft bool) error {
	absPath, err := getAbsPath(orgFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	title, err := extractTitleFromOrg(absPath)
	if err != nil {
		return fmt.Errorf("failed to extract title from org file: %v", err)
	}

	markdown, err := convertOrgToMarkdown(absPath)
	if err != nil {
		return fmt.Errorf("failed to convert org to markdown: %v", err)
	}

	content := removeTitleFromMarkdown(markdown)

	client := NewHatenaClient(config.HatenaID, config.APIKey, config.BlogDomain)
	entry := BlogEntry{
		Title:    title,
		Content:  content,
		Category: category,
		IsDraft:  isDraft,
	}

	return client.PostEntry(entry)
}

func validateConfig(config *Config) error {
	if config.HatenaID == "" {
		return fmt.Errorf("hatena ID is required")
	}
	if config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	if config.BlogDomain == "" {
		return fmt.Errorf("blog domain is required")
	}
	return nil
}
