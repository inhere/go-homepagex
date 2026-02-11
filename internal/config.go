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

// Config 应用主配置
type Config struct {
	Server      ServerConfig `yaml:"server"`
	Auth        AuthConfig   `yaml:"auth"`
	PagesDir    string       `yaml:"pages_dir"`
	FrontendDir string       `yaml:"frontend_dir"`
	Navs        []NavItem    `yaml:"navs" json:"navs"`
}

func newDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
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

	if err := yaml.Unmarshal(data, &config); err != nil {
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
