package bookmarks

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

type Type struct {
	OpenWith map[string]string `toml:"open_with"`
}

type Bookmark struct {
	Content     string   `toml:"content"`
	Description string   `toml:"description"`
	Tags        []string `toml:"tags"`
	Type        string   `toml:"type"`
}

type Config struct {
	Chooser   string          `toml:"chooser"`
	Types     map[string]Type `toml:"types"`
	Bookmarks []Bookmark      `toml:"bookmarks"`
}

type Query struct {
	Tags []string
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open bookmarks.toml: %w", err)
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read bookmarks.toml: %w", err)
	}

	err = toml.Unmarshal(content, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bookmarks.toml: %w", err)
	}

	return config, nil
}

func ParseQuery(query string) Query {
	tags := strings.Split(query, " ")
	return Query{Tags: tags}
}

func (c *Config) OpenBookmark(openWithType string, openWithFile string) error {
	return c.OpenBookmarkWithQuery(Query{}, openWithType, openWithFile)
}

func (c *Config) OpenBookmarkWithQuery(query Query, openWithType string, openWithFile string) error {
	bookmarks := c.Bookmarks

	if len(query.Tags) > 0 {
		bookmarks = filterByTags(bookmarks, query.Tags...)
	}

	b, err := chooseBookmark(c.Chooser, bookmarks)
	if err != nil {
		return err
	}

	openWith, err := getOpenerForBookmark(c.Types, b, openWithType)
	if err != nil {
		return fmt.Errorf("failed to get opener for bookmark: %w", err)
	}

	b = formatContentByType(b)

	return open(openWith, b, openWithType, openWithFile)
}

func (c *Config) AddBookmark(b Bookmark) error {
	c.Bookmarks = append(c.Bookmarks, b)
	return nil
}

func (c *Config) ListTags() ([]string, error) {
	var tags []string

	for _, bookmark := range c.Bookmarks {
		for _, tag := range bookmark.Tags {
			contains := false
			for _, t := range tags {
				if t == tag {
					contains = true
					break
				}
			}

			if !contains {
				tags = append(tags, tag)
			}
		}
	}

	return tags, nil
}

func (c *Config) Save(path string) error {
	content, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(path, content, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func open(openWith string, b Bookmark, openWithType string, openWithFile string) error {
	if !strings.Contains(openWith, "{}") {
		openWith = fmt.Sprintf("%s {}", openWith)
	}
	openWith = strings.Replace(openWith, "{}", b.Content, 1)

	switch openWithType {
	case "bash":
		err := os.WriteFile(openWithFile, []byte(openWith), 0600)
		if err != nil {
			return fmt.Errorf("failed to write openWithFile: %w", err)
		}
	default:
		args := strings.Split(openWith, " ")
		cmd := exec.Command(args[0], args[1:]...) // nolint:gosec
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to open bookmark: %w", err)
		}
	}
	return nil
}

func formatContentByType(b Bookmark) Bookmark {
	switch b.Type {
	case "folder":
		if b.Content[0] == '~' {
			b.Content = strings.Replace(b.Content, "~", os.Getenv("HOME"), 1)
		}
	}

	return b
}

func getOpenerForBookmark(t map[string]Type, b Bookmark, o string) (string, error) {
	bt, ok := t[b.Type]
	if !ok {
		return "", fmt.Errorf("no type %s found", b.Type)
	}

	openWith, ok := bt.OpenWith[o]
	if !ok {
		openWith = bt.OpenWith["default"]
	}

	return openWith, nil
}

func filterByTags(bookmarks []Bookmark, tags ...string) []Bookmark {
	var b []Bookmark

	for _, bookmark := range bookmarks {
		for _, tag := range tags {
			for _, t := range bookmark.Tags {
				if t == tag {
					b = append(b, bookmark)
				}
			}
		}
	}

	return b
}

func chooseBookmark(chooser string, bookmarks []Bookmark) (Bookmark, error) {
	args := strings.Split(chooser, " ")

	var cmd *exec.Cmd
	if len(args) == 1 {
		cmd = exec.Command(args[0]) // nolint:gosec
	} else {
		cmd = exec.Command(args[0], args[1:]...) // nolint:gosec
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return Bookmark{}, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	go func() {
		defer stdin.Close()
		for i, bookmark := range bookmarks {
			_, err := fmt.Fprintf(stdin, "%d: %s(%s)\n", i, bookmark.Description, bookmark.Content)
			if err != nil {
				return
			}
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return Bookmark{}, fmt.Errorf("failed to run chooser: %w", err)
	}

	index, err := strconv.Atoi(strings.Split(string(out), ":")[0])
	if err != nil {
		return Bookmark{}, fmt.Errorf("failed to convert index to int: %w", err)
	}

	return bookmarks[index], nil
}
