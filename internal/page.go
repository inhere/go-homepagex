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
	Title        string       `yaml:"title" json:"title"`
	Subtitle     string       `yaml:"subtitle" json:"subtitle"`
	Logo         string       `yaml:"logo" json:"logo"`
	Header       string       `yaml:"header" json:"header"`
	Footer       string       `yaml:"footer" json:"footer"`
	Theme        string       `yaml:"theme" json:"theme"`
	Color        string       `yaml:"color" json:"color"`
	Style        string       `yaml:"style" json:"style"`
	Columns      string       `yaml:"columns" json:"columns"`
	Connectivity Connectivity `yaml:"connectivity" json:"connectivity"`
	Services     []Service    `yaml:"services" json:"services"`
}

// Connectivity 连接检查配置
type Connectivity struct {
	CheckInterval int    `yaml:"check_interval" json:"check_interval"`
	Mode          string `yaml:"mode" json:"mode"`
}

// Service 服务分组
type Service struct {
	Name  string `yaml:"name" json:"name"`
	Icon  string `yaml:"icon" json:"icon"`
	Items []Item `yaml:"items" json:"items"`
}

// Item 单个链接项
type Item struct {
	Name     string            `yaml:"name" json:"name"`
	Logo     string            `yaml:"logo" json:"logo"`
	Subtitle string            `yaml:"subtitle" json:"subtitle"`
	Tag      string            `yaml:"tag" json:"tag"`
	Keywords string            `yaml:"keywords" json:"keywords"`
	URL      string            `yaml:"url" json:"url"`
	Target   string            `yaml:"target" json:"target"`
	Method   string            `yaml:"method" json:"method"`
	Headers  map[string]string `yaml:"headers" json:"headers"`
	Type     string            `yaml:"type" json:"type"`
}

// LoadPageConfig 加载页面配置
func LoadPageConfig(pageDir, route string) (*PageConfig, error) {
	var filename string
	if route == "/" || route == "" {
		filename = "main.yaml"
	} else {
		// 移除开头的 /
		route = strings.TrimPrefix(route, "/")
		filename = fmt.Sprintf("page-%s.yaml", route)
	}

	path := filepath.Join(pageDir, filename)
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
