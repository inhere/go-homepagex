package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// PageConfig 页面配置（类似 Homer 的格式）
type PageConfig struct {
	Title        string       `yaml:"title"`
	Subtitle     string       `yaml:"subtitle"`
	Logo         string       `yaml:"logo"`
	Header       string       `yaml:"header"`
	Footer       string       `yaml:"footer"`
	Theme        string       `yaml:"theme"`
	Color        string       `yaml:"color"`
	Style        string       `yaml:"style"` // "cards" 或 "list"
	Columns      string       `yaml:"columns"`
	Connectivity Connectivity `yaml:"connectivity"`
	Services     []Service    `yaml:"services"`
}

// Connectivity 连接检查配置
type Connectivity struct {
	CheckInterval int    `yaml:"check_interval"`
	Mode          string `yaml:"mode"`
}

// Service 服务分组
type Service struct {
	Name  string `yaml:"name"`
	Icon  string `yaml:"icon"`
	Items []Item `yaml:"items"`
}

// Item 单个链接项
type Item struct {
	Name     string            `yaml:"name"`
	Logo     string            `yaml:"logo"`
	Subtitle string            `yaml:"subtitle"`
	Tag      string            `yaml:"tag"`
	Keywords string            `yaml:"keywords"`
	URL      string            `yaml:"url"`
	Target   string            `yaml:"target"`
	Method   string            `yaml:"method"`
	Headers  map[string]string `yaml:"headers"`
	Type     string            `yaml:"type"`
}

// LoadPageConfig 加载页面配置
func (s *Server) LoadPageConfig(route string) (*PageConfig, error) {
	var filename string
	if route == "/" || route == "" {
		filename = "main.yaml"
	} else {
		// 移除开头的 /
		route = strings.TrimPrefix(route, "/")
		filename = fmt.Sprintf("page-%s.yaml", route)
	}

	path := filepath.Join(s.config.PagesDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("page config not found: %s", filename)
		}
		return nil, fmt.Errorf("failed to read page config: %w", err)
	}

	var page PageConfig
	if err := yaml.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("failed to parse page config: %w", err)
	}

	// 设置默认值
	if page.Style == "" {
		page.Style = "cards"
	}
	if page.Columns == "" {
		page.Columns = "3"
	}

	return &page, nil
}
