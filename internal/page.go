package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/gookit/goutil/fsutil"
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
	Navs         []NavItem    `yaml:"navs" json:"navs"`

	// 内部设置，页面配置文件路径
	Pagefile string `yaml:"-" json:"-"`
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
	Tags     []string          `yaml:"tags" json:"tags"`
	Keywords string            `yaml:"keywords" json:"keywords"`
	URL      string            `yaml:"url" json:"url"`
	Target   string            `yaml:"target" json:"target"`
	Method   string            `yaml:"method" json:"method"`
	Headers  map[string]string `yaml:"headers" json:"headers"`
	Type     string            `yaml:"type" json:"type"`
}

// PageDataManager 页面数据管理器
type PageDataManager struct {
	Debug   bool
	PageDir string
	// 默认导航项
	Navs []NavItem

	// 页面配置缓存 key is page name TODO 支持缓存过期
	cacheMap map[string]*PageConfig
}

// PageDataMgr 页面数据管理器实例
var PageDataMgr *PageDataManager

// GetPageConfig 获取页面配置数据
func (m *PageDataManager) GetPageConfig(name string) (*PageConfig, error) {
	if m.cacheMap == nil {
		m.cacheMap = make(map[string]*PageConfig)
	}

	// 非调试模式下，从缓存中获取
	if page, ok := m.cacheMap[name]; ok && !m.Debug {
		return page, nil
	}

	page, err := m.LoadPageConfig(name)
	if err == nil {
		m.cacheMap[name] = page
	}

	return page, err
}

const (
	// DefaultPageName 默认页面名
	DefaultPageName = "home"
	// DefaultPageFile 默认页面文件名
	DefaultPageFile = "home.yaml"
)

// getFilename 生成页面配置文件名
func (m *PageDataManager) getFilename(name string) string {
	// 移除开头的无效字符 /.
	name = strings.TrimLeft(name, "/.")
	if name == "" {
		return DefaultPageName
	}
	return name
}

// LoadPageConfig 加载页面配置
func (m *PageDataManager) LoadPageConfig(name string) (*PageConfig, error) {
	filename := m.getFilename(name)
	pagefile := filepath.Join(m.PageDir, filename + ".yaml")

	var err error
	var data []byte

	// debug mode 下，优先使用 {name}.local.yaml
	if m.Debug {
		dotLocalFile := filepath.Join(m.PageDir, filename+".local.yaml")
		if fsutil.IsFile(dotLocalFile) {
			pagefile = dotLocalFile
			data, err = os.ReadFile(dotLocalFile)
		}
	}

	if len(data) == 0 {
		data, err = os.ReadFile(pagefile)
	}
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

	// 记录页面配置文件路径
	page.Pagefile = pagefile

	// 设置默认值
	if page.Style == "" {
		page.Style = "cards"
	}
	if page.Columns == "" {
		page.Columns = "3"
	}

	// 如果页面没有配置 navs，使用默认配置
	if len(page.Navs) == 0 {
		page.Navs = m.Navs
	}
	return &page, nil
}
