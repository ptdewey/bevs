package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

type Page struct {
	Metadata map[string]interface{} `json:"metadata"`
	Content  string                 `json:"content"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category"`
}

func generateSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	return slug
}

func parseFrontMatter(content []byte) (map[string]interface{}, []byte, error) {
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---") {
		return nil, content, nil
	}
	parts := strings.SplitN(contentStr, "---", 3)
	if len(parts) < 3 {
		return nil, content, fmt.Errorf("invalid front-matter format")
	}
	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		return nil, nil, err
	}
	return metadata, []byte(parts[2]), nil
}

func getReadingTime(text string) int {
	words := strings.Fields(text)
	wordCount := len(words)

	// reading/speaking rate
	wordsPerMinute := 200.0
	return int(math.Round(float64(wordCount) / wordsPerMinute))
}

func processMarkdownFile(path string) (Page, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Page{}, err
	}

	metadata, markdownContent, err := parseFrontMatter(content)
	if err != nil {
		return Page{}, err
	}

	if _, ok := metadata["slug"]; !ok {
		fname := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		metadata["slug"] = generateSlug(fname)
		metadata["category"] = filepath.Base(filepath.Dir(path))
	}

	var htmlContent bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	if err := md.Convert(markdownContent, &htmlContent); err != nil {
		return Page{}, err
	}

	metadata["read_time"] = getReadingTime(string(markdownContent))

	return Page{
		Metadata: metadata,
		Content:  htmlContent.String(),
	}, nil
}

func processDirectory(dir string) ([]Page, error) {
	var pages []Page

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".md") {
			page, err := processMarkdownFile(path)
			if err != nil {
				return err
			}

			pages = append(pages, page)
		}
		return nil
	})

	return pages, err
}

func generateRSSFromJSON(inputPath string, outputPath string) error {
	file, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var pages []Page
	if err := json.Unmarshal(file, &pages); err != nil {
		return err
	}

	channel := Channel{
		Title:       "Patrick's Bevs",
		Link:        "https://bev.pdewey.com",
		Description: "RSS feed of Patrick Dewey's bev website",
		PubDate:     time.Now().Format(time.RFC1123Z),
	}

	for _, page := range pages {
		var categories []string
		if rawCategories, ok := page.Metadata["categories"].([]interface{}); ok {
			for _, category := range rawCategories {
				if strCategory, ok := category.(string); ok {
					categories = append(categories, strCategory)
				}
			}
		}

		date, ok := page.Metadata["date"].(string)
		if !ok {
			date = ""
		}

		item := Item{
			Title:       page.Metadata["title"].(string),
			Link:        fmt.Sprintf("https://bev.pdewey.com/%s", page.Metadata["slug"].(string)),
			Description: page.Content,
			PubDate:     date,
			Category:    strings.Join(categories, ", "),
		}
		channel.Items = append(channel.Items, item)
	}

	rss := RSS{
		Version: "2.0",
		Channel: channel,
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return err
	}

	rssHeader := []byte(xml.Header)
	output = append(rssHeader, output...)

	return os.WriteFile(outputPath, output, 0644)
}

func writeJSONFile(data interface{}, outputPath string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, jsonData, 0644)
}

func main() {
	pages, err := processDirectory("content")
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}

	// TODO: spread out output to different files based on subdir structure within content
	if err := writeJSONFile(pages, "static/data/pages.json"); err != nil {
		fmt.Printf("Error writing writing.json: %v\n", err)
		os.Exit(1)
	}

	if err := generateRSSFromJSON("static/data/pages.json", "static/rss.xml"); err != nil {
		fmt.Printf("Error writing rss.xml: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated writing.json and rss.xml")
}
