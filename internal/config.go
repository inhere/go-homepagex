package internal

import (
	"crypto/subtle"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode string `yaml:"mode" json:"mode"` // debug or release
	Port string `yaml:"port" json:"port"`
}

// NavItem 导航项
type NavItem struct {
	Name string `yaml:"name" json:"name"`
	Icon string `yaml:"icon" json:"icon"`
	URL  string `yaml:"url" json:"url"`
}

// PageDefaults 页面默认配置
type PageDefaults struct {
	Title    string `yaml:"title" json:"title"`
	Subtitle string `yaml:"subtitle" json:"subtitle"`
	Logo     string `yaml:"logo" json:"logo"`
	Header   string `yaml:"header" json:"header"`
	Footer   string `yaml:"footer" json:"footer"`
	Theme    string `yaml:"theme" json:"theme"`
	Color    string `yaml:"color" json:"color"`
	Style    string `yaml:"style" json:"style"`
	Columns  string `yaml:"columns" json:"columns"`
}

// Config 应用主配置
type Config struct {
	Server      ServerConfig `yaml:"server"`
	PagesDir    string       `yaml:"pages_dir"`
	FrontendDir string       `yaml:"frontend_dir"`
	// 图标 CDN 配置 see https://dashboardicons.com/ 搜索
	IconsCDN map[string]string `yaml:"icons_cdn"`
	// basic 认证配置
	Auths []string `yaml:"auths"`
	// 页面默认配置
	PageDefaults PageDefaults `yaml:"page_defaults"`
	PageNavs     []NavItem    `yaml:"page_navs"`
	// 解析后的认证配置，key为username
	parsedAuths map[string]*AuthConfig
}

// AuthConfig 一个解析后的认证配置
type AuthConfig struct {
	Username  string
	Password  string
	PathPerms []string
	// runtime 匹配后的权限：rw 读写, ro 只读, no 拒绝访问
	Permission string
}

// IsValid 是否有效
func (c *AuthConfig) IsValid() bool {
	return c.Username != "" || c.Password != "" || len(c.PathPerms) > 0
}

func newDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8090",
		},
		Auths:       []string{"@*"}, // 所有路径可访问，无需认证
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

	if err = yaml.Unmarshal(data, config); err != nil {
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

	err = config.parseAuths()
	if err != nil {
		return nil, fmt.Errorf("failed to parse auths: %w", err)
	}
	return config, nil
}

const (
	// 读写权限
	PermRW = "rw"
	// 只读权限
	PermRO = "ro"
	// 拒绝访问权限
	PermNO = "no"
)

func (c *Config) parseAuths() error {
	auths := make(map[string]*AuthConfig)
	for _, auth := range c.Auths {
		if auth == "" {
			continue
		}

		ac := &AuthConfig{
			PathPerms: []string{},
		}

		credStr, pathStr, ok := strings.Cut(auth, "@")
		if !ok {
			continue
		}

		if credStr != "" {
			colonIdx := strings.Index(credStr, ":")
			if colonIdx == -1 {
				ac.Username = credStr
			} else {
				ac.Username = credStr[:colonIdx]
				ac.Password = credStr[colonIdx+1:]
			}
		}

		// 裸 * 表示所有路径只读，直接归一化
		if pathStr == "" || pathStr == "*" {
			ac.PathPerms = append(ac.PathPerms, "/*:ro")
			auths[ac.Username] = ac
			continue
		}

		var noPerms []string
		var normalPerms []string

		for p := range strings.SplitSeq(pathStr, ",") {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}

			if after, ok := strings.CutPrefix(p, "!"); ok {
				// 排除路径：始终归一化为 /prefix:no，由 pathMatch 做通配匹配
				if !strings.HasPrefix(after, "/") {
					after = "/" + after
				}
				noPerms = append(noPerms, after+":no")
			} else {
				// 普通路径：已有 :perm 后缀则原样保存（pathMatch 负责归一化），
				// 否则补全 / 前缀和 :ro 默认权限
				hasPerm := strings.HasSuffix(p, ":rw") || strings.HasSuffix(p, ":ro")
				if !hasPerm {
					if !strings.HasPrefix(p, "/") {
						p = "/" + p
					}
					p = p + ":ro"
				}
				normalPerms = append(normalPerms, p)
			}
		}

		// 排除规则放在前面，优先匹配
		ac.PathPerms = append(noPerms, normalPerms...)

		if len(ac.PathPerms) > 0 {
			auths[ac.Username] = ac
		}
	}

	c.parsedAuths = auths
	return nil
}

func (c *Config) IsNeedAuth(reqPath string, isWrite bool) bool {
	var matchedPerm string
	var needCred bool

	if !strings.HasPrefix(reqPath, "/") {
		reqPath = "/" + reqPath
	}

	// 检查访客权限
	authCfg, exists := c.parsedAuths[""]
	if exists {
		for _, pathWithPerm := range authCfg.PathPerms {
			matched, perm := c.pathMatch(pathWithPerm, reqPath)
			if matched {
				if perm == PermNO {
					return true
				}
				if matchedPerm == "" {
					matchedPerm = perm
					needCred = authCfg.Username != "" || authCfg.Password != ""
				}
			}
		}
	}

	if matchedPerm == "" {
		return true
	}

	if matchedPerm == PermRW {
		return needCred
	}

	if matchedPerm == PermRO {
		if isWrite {
			return true
		}
		return needCred
	}

	return false
}

// MatchUserAuthConfig 根据用户名匹配认证配置
func (c *Config) MatchUserAuthConfig(username, reqPath string) (*AuthConfig, bool) {
	auth, exists := c.parsedAuths[username]
	if !exists {
		return nil, false
	}

	if !strings.HasPrefix(reqPath, "/") {
		reqPath = "/" + reqPath
	}

	for _, pathWithPerm := range auth.PathPerms {
		matched, perm := c.pathMatch(pathWithPerm, reqPath)
		if matched {
			result := *auth
			result.Permission = perm
			return &result, true
		}
	}

	return nil, false
}

// MatchAuthConfig 匹配认证配置（无用户名时使用）
func (c *Config) MatchAuthConfig(reqPath string) *AuthConfig {
	var matchedAuth *AuthConfig
	var noAuthMatched *AuthConfig

	if !strings.HasPrefix(reqPath, "/") {
		reqPath = "/" + reqPath
	}

	for _, auth := range c.parsedAuths {
		for _, pathWithPerm := range auth.PathPerms {
			matched, perm := c.pathMatch(pathWithPerm, reqPath)
			if matched {
				result := *auth
				result.Permission = perm

				if perm == PermNO {
					return &result
				}

				if auth.Username == "" && auth.Password == "" {
					if noAuthMatched == nil {
						noAuthMatched = &result
					}
				} else {
					if matchedAuth == nil {
						matchedAuth = &result
					}
				}
			}
		}
	}

	if noAuthMatched != nil {
		return noAuthMatched
	}
	return matchedAuth
}

// pathMatch 路径匹配
// pattern 格式: path:perm 或 path（默认 ro）
// 支持：* 或 /* 匹配所有路径；/prefix* 匹配前缀；/path 精确匹配或子路径匹配
func (c *Config) pathMatch(pathWithPerm string, reqPath string) (bool, string) {
	pattern, perm, found := strings.Cut(pathWithPerm, ":")
	if !found {
		pattern = pathWithPerm
		perm = PermRO
	}
	if pattern == "" {
		return false, ""
	}

	// 归一化：pattern 非裸 * 时补全 / 前缀
	if pattern != "*" && !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	// 归一化：reqPath 补全 / 前缀
	if !strings.HasPrefix(reqPath, "/") {
		reqPath = "/" + reqPath
	}

	// 通配后缀：*, /*, /prefix*, /prefix/* 均转为前缀匹配
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(reqPath, prefix), perm
	}

	// 精确匹配或子路径匹配（/path → /path 和 /path/...）
	return reqPath == pattern || strings.HasPrefix(reqPath, pattern+"/"), perm
}

// CheckCredentials 验证用户名和密码是否匹配任意已配置用户
func (c *Config) CheckCredentials(username, password string) bool {
	auth, exists := c.parsedAuths[username]
	if !exists || auth.Username == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(password), []byte(auth.Password)) == 1
}

// FilterNavsByPermission 过滤出当前用户有权访问的导航项
// - 公开可访问（无需认证）的路径对所有人显示
// - 需要认证的路径仅对有权限的已登录用户显示
func (c *Config) FilterNavsByPermission(navs []NavItem, username string) []NavItem {
	var filtered []NavItem
	for _, nav := range navs {
		navPath := nav.URL
		if navPath == "" {
			filtered = append(filtered, nav)
			continue
		}
		// 公开可访问 → 所有人显示
		if !c.IsNeedAuth(navPath, false) {
			filtered = append(filtered, nav)
			continue
		}
		// 需要认证 → 仅对有权限的已登录用户显示
		if username != "" {
			authConfig, exists := c.MatchUserAuthConfig(username, navPath)
			if exists && authConfig.Permission != PermNO {
				filtered = append(filtered, nav)
			}
		}
	}
	return filtered
}

// AuthEnabled 是否启用认证
func (c *Config) AuthEnabled() bool {
	return len(c.parsedAuths) > 0
}

// ParsedAuths 解析后的认证配置
func (c *Config) ParsedAuths() map[string]*AuthConfig {
	return c.parsedAuths
}
