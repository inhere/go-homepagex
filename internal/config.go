package internal

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode string `yaml:"mode" json:"mode"` // debug or release
	Port string `yaml:"port" json:"port"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// NavItem 导航项
type NavItem struct {
	Name string `yaml:"name" json:"name"`
	Icon string `yaml:"icon" json:"icon"`
	URL  string `yaml:"url" json:"url"`
}

// PageDefaults 页面默认配置
type PageDefaults struct {
	Title        string       `yaml:"title" json:"title"`
	Subtitle     string       `yaml:"subtitle" json:"subtitle"`
	Logo         string       `yaml:"logo" json:"logo"`
	Header       string       `yaml:"header" json:"header"`
	Footer       string       `yaml:"footer" json:"footer"`
	Theme        string       `yaml:"theme" json:"theme"`
	Color        string       `yaml:"color" json:"color"`
	Style        string       `yaml:"style" json:"style"`
	Columns      string       `yaml:"columns" json:"columns"`
}

// Config 应用主配置
type Config struct {
	Server      ServerConfig `yaml:"server"`
	PagesDir    string       `yaml:"pages_dir"`
	FrontendDir string       `yaml:"frontend_dir"`
	// 图标 CDN 配置 see https://dashboardicons.com/ 搜索
	IconsCDN map[string]string `yaml:"icons_cdn" json:"icons_cdn"`
	Auth     AuthConfig        `yaml:"auth"`
	Navs     []NavItem         `yaml:"navs" json:"navs"`
}

func newDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8090",
		},
		Auth: AuthConfig{
			Enabled: false,
		},
		PagesDir:    "./pages",
		FrontendDir: "./frontend/build",
	}
}

// LoadConfig 从 YAML 文件加载配置
func LoadConfig(path string) (*Config, error) {
	// 加载默认配置
	config := newDefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值
	if config.Server.Port == "" {
		config.Server.Port = "8090"
	}
	if config.PagesDir == "" {
		config.PagesDir = "./pages"
	}
	if config.FrontendDir == "" {
		config.FrontendDir = "./frontend/build"
	}
	return config, nil
}
