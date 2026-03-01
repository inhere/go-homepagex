package internal

import (
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil/assert"
)

func TestParseAuths(t *testing.T) {
	tests := []struct {
		name    string
		auths   []string
		wantLen int
		checkFn func(*testing.T, map[string]*AuthConfig)
	}{
		{
			name:    "空配置",
			auths:   []string{},
			wantLen: 0,
		},
		{
			name:    "公开访问 @*",
			auths:   []string{"@*"},
			wantLen: 1,
			checkFn: func(t *testing.T, auths map[string]*AuthConfig) {
				var auth *AuthConfig
				for _, a := range auths {
					auth = a
					break
				}
				if auth.Username != "" {
					t.Errorf("期望用户名为空，实际为 %q", auth.Username)
				}
				if auth.Password != "" {
					t.Errorf("期望密码为空，实际为 %q", auth.Password)
				}
				if len(auth.PathPerms) != 1 || auth.PathPerms[0] != "/*:ro" {
					t.Errorf("期望路径权限为 [/*:ro]，实际为 %v", auth.PathPerms)
				}
			},
		},
		{
			name:    "带用户名密码的完整配置",
			auths:   []string{"admin:admin123@*:rw"},
			wantLen: 1,
			checkFn: func(t *testing.T, auths map[string]*AuthConfig) {
				auth, exists := auths["admin"]
				if !exists {
					t.Errorf("期望找到用户 admin")
					return
				}
				if auth.Username != "admin" {
					t.Errorf("期望用户名 admin，实际为 %q", auth.Username)
				}
				if auth.Password != "admin123" {
					t.Errorf("期望密码 admin123，实际为 %q", auth.Password)
				}
				if len(auth.PathPerms) != 1 || auth.PathPerms[0] != "*:rw" {
					t.Errorf("期望路径权限为 [*:rw]，实际为 %v", auth.PathPerms)
				}
			},
		},
		{
			name:    "多路径配置",
			auths:   []string{"user:pass@/api:rw,/static:ro"},
			wantLen: 1,
			checkFn: func(t *testing.T, auths map[string]*AuthConfig) {
				auth, exists := auths["user"]
				if !exists {
					t.Errorf("期望找到用户 user")
					return
				}
				if auth.Username != "user" {
					t.Errorf("期望用户名 user，实际为 %q", auth.Username)
				}
				if auth.Password != "pass" {
					t.Errorf("期望密码 pass，实际为 %q", auth.Password)
				}
				if len(auth.PathPerms) != 2 {
					t.Errorf("期望 2 个路径权限，实际为 %d", len(auth.PathPerms))
				}
			},
		},
		{
			name:    "匿名用户多路径",
			auths:   []string{"@/public:ro,/api"},
			wantLen: 1,
			checkFn: func(t *testing.T, auths map[string]*AuthConfig) {
				var auth *AuthConfig
				for _, a := range auths {
					auth = a
					break
				}
				if auth.Username != "" {
					t.Errorf("期望用户名为空，实际为 %q", auth.Username)
				}
				if len(auth.PathPerms) != 2 {
					t.Errorf("期望 2 个路径权限，实际为 %d", len(auth.PathPerms))
				}
			},
		},
		{
			name:    "排除路径配置",
			auths:   []string{"@*,!/inner"},
			wantLen: 1,
			checkFn: func(t *testing.T, auths map[string]*AuthConfig) {
				var auth *AuthConfig
				for _, a := range auths {
					auth = a
					break
				}
				if len(auth.PathPerms) != 2 {
					t.Errorf("期望 2 个路径权限，实际为 %d", len(auth.PathPerms))
				}
			},
		},
		{
			name:    "仅用户名无密码",
			auths:   []string{"admin@*:rw"},
			wantLen: 1,
			checkFn: func(t *testing.T, auths map[string]*AuthConfig) {
				auth, exists := auths["admin"]
				if !exists {
					t.Errorf("期望找到用户 admin")
					return
				}
				if auth.Username != "admin" {
					t.Errorf("期望用户名 admin，实际为 %q", auth.Username)
				}
				if auth.Password != "" {
					t.Errorf("期望密码为空，实际为 %q", auth.Password)
				}
			},
		},
		{
			name:    "仅通配符路径",
			auths:   []string{"user:pass@*"},
			wantLen: 1,
			checkFn: func(t *testing.T, auths map[string]*AuthConfig) {
				auth, exists := auths["user"]
				if !exists {
					t.Errorf("期望找到用户 user")
					return
				}
				if len(auth.PathPerms) != 1 || auth.PathPerms[0] != "/*:ro" {
					t.Errorf("期望路径权限为 [/*:ro]，实际为 %v", auth.PathPerms)
				}
			},
		},
		{
			name:    "多用户配置",
			auths:   []string{"admin:admin123@*:rw", "user1:user123@/tools:rw"},
			wantLen: 2,
		},
		{
			name:    "空字符串过滤",
			auths:   []string{"", "admin:pass@*", ""},
			wantLen: 1,
		},
		{
			name:    "/inner 需要认证",
			auths:   []string{"@*,!/inner*"},
			wantLen: 1,
			checkFn: func(t *testing.T, auths map[string]*AuthConfig) {
				var auth *AuthConfig
				for _, a := range auths {
					auth = a
					break
				}
				if auth.Username != "" {
					t.Errorf("期望用户名为空，实际为 %q", auth.Username)
				}
				if auth.Password != "" {
					t.Errorf("期望密码为空，实际为 %q", auth.Password)
				}
				if len(auth.PathPerms) != 2 || auth.PathPerms[0] != "/inner*:no" || auth.PathPerms[1] != "/*:ro" {
					t.Errorf("期望路径权限为 [/inner*:no /*:ro]，实际为 %v", auth.PathPerms)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Auths: tt.auths,
			}
			err := c.parseAuths()
			if err != nil {
				t.Fatalf("parseAuths 失败: %v", err)
			}

			if len(c.parsedAuths) != tt.wantLen {
				t.Errorf("期望解析出 %d 个认证配置，实际为 %d", tt.wantLen, len(c.parsedAuths))
			}

			if tt.checkFn != nil && len(c.parsedAuths) > 0 {
				tt.checkFn(t, c.parsedAuths)
			}
		})
	}
}

func TestMatchAuthConfig(t *testing.T) {
	tests := []struct {
		name      string
		auths     []string
		reqPath   string
		wantMatch bool
		wantUser  string
		wantPerm  string
	}{
		{
			name:      "公开路径匹配",
			auths:     []string{"@*"},
			reqPath:   "/anything",
			wantMatch: true,
			wantUser:  "",
			wantPerm:  "ro",
		},
		{
			name:      "通配符路径匹配",
			auths:     []string{"admin:pass@*:rw"},
			reqPath:   "/api/test",
			wantMatch: true,
			wantUser:  "admin",
			wantPerm:  "rw",
		},
		{
			name:      "精确路径匹配",
			auths:     []string{"user:pass@/api:rw"},
			reqPath:   "/api",
			wantMatch: true,
			wantUser:  "user",
			wantPerm:  "rw",
		},
		{
			name:      "子路径匹配",
			auths:     []string{"user:pass@/api/*"},
			reqPath:   "/api/users/123",
			wantMatch: true,
			wantUser:  "user",
			wantPerm:  "ro",
		},
		{
			name:      "排除路径匹配",
			auths:     []string{"@*,!/inner"},
			reqPath:   "/inner",
			wantMatch: true,
			wantPerm:  "no",
		},
		{
			name:      "排除路径匹配2",
			auths:     []string{"@*,!/inner/*"},
			reqPath:   "/inner/tools",
			wantMatch: true,
			wantPerm:  "no",
		},
		{
			name:      "多用户配置-第一个匹配",
			auths:     []string{"admin:admin123@*:rw", "user1:user123@/tools:rw"},
			reqPath:   "/home",
			wantMatch: true,
			wantUser:  "admin",
			wantPerm:  "rw",
		},
		{
			name:      "多用户配置-匹配到第一个",
			auths:     []string{"admin:admin123@*:rw", "user1:user123@/tools:rw"},
			reqPath:   "/tools",
			wantMatch: true,
			wantUser:  "admin",
			wantPerm:  "rw",
		},
		{
			name:      "无配置-不匹配",
			auths:     []string{},
			reqPath:   "/api",
			wantMatch: false,
		},
		{
			name:      "无匹配路径",
			auths:     []string{"admin:pass@/admin/*"},
			reqPath:   "/public",
			wantMatch: false,
		},
		{
			name:      "多路径配置-第一个匹配",
			auths:     []string{"user:pass@/api:rw,/static:ro"},
			reqPath:   "/api",
			wantMatch: true,
			wantPerm:  "rw",
		},
		{
			name:      "多路径配置-第二个匹配",
			auths:     []string{"user:pass@/api:rw,/static:ro"},
			reqPath:   "/static",
			wantMatch: true,
			wantPerm:  "ro",
		},
		{
			name:      "不带权限后缀默认为ro",
			auths:     []string{"user:pass@/api"},
			reqPath:   "/api",
			wantMatch: true,
			wantPerm:  "ro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Auths: tt.auths,
			}
			err := c.parseAuths()
			if err != nil {
				t.Fatalf("parseAuths 失败: %v", err)
			}

			result := c.MatchAuthConfig(tt.reqPath)

			if tt.wantMatch {
				if result == nil {
					t.Errorf("期望匹配到认证配置，实际为 nil")
					return
				}
				if tt.wantUser != "" && result.Username != tt.wantUser {
					t.Errorf("期望用户名 %q，实际为 %q", tt.wantUser, result.Username)
				}
				if result.Permission != tt.wantPerm {
					t.Errorf("期望权限 %q，实际为 %q", tt.wantPerm, result.Permission)
				}
			} else {
				if result != nil {
					t.Errorf("期望不匹配，实际匹配到了 %v", result)
				}
			}
		})
	}
}

func TestIsNeedAuth(t *testing.T) {
	c := &Config{
		Auths: []string{"admin:admin123@*:rw", "user1:user123@/tools:rw", "@*,!/inner*"},
	}
	err := c.parseAuths()
	if err != nil {
		t.Fatalf("parseAuths 失败: %v", err)
	}

	dump.Config(dump.WithoutColor())
	dump.Clear(c.parsedAuths)

	assert.False(t, c.IsNeedAuth("/", false))
	assert.True(t, c.IsNeedAuth("/inner-tools", false))
}

func TestPathMatch(t *testing.T) {
	tests := []struct {
		name      string
		pattern   string
		reqPath   string
		wantMatch bool
		wantPerm  string
	}{
		{
			name:      "通配符匹配",
			pattern:   "/*:ro",
			reqPath:   "/anything",
			wantMatch: true,
			wantPerm:  "ro",
		},
		{
			name:      "星号匹配",
			pattern:   "*:rw",
			reqPath:   "/test",
			wantMatch: true,
			wantPerm:  "rw",
		},
		{
			name:      "精确路径匹配",
			pattern:   "/api:rw",
			reqPath:   "/api",
			wantMatch: true,
			wantPerm:  "rw",
		},
		{
			name:      "子路径匹配",
			pattern:   "/api:ro",
			reqPath:   "/api/users",
			wantMatch: true,
			wantPerm:  "ro",
		},
		{
			name:      "通配符前缀匹配",
			pattern:   "/api/*:rw",
			reqPath:   "/api/users/123",
			wantMatch: true,
			wantPerm:  "rw",
		},
		{
			name:      "通配符前缀不匹配",
			pattern:   "/api/*:rw",
			reqPath:   "/other",
			wantMatch: false,
			wantPerm:  "rw",
		},
		{
			name:      "精确路径不匹配",
			pattern:   "/api:rw",
			reqPath:   "/other",
			wantMatch: false,
			wantPerm:  "rw",
		},
		{
			name:      "无权限前缀",
			pattern:   "/test",
			reqPath:   "/test",
			wantMatch: true,
			wantPerm:  "ro",
		},
		{
			name:      "无斜杠前缀自动添加",
			pattern:   "api/*:rw",
			reqPath:   "/api/users",
			wantMatch: true,
			wantPerm:  "rw",
		},
		{
			name:      "无斜杠请求路径自动添加",
			pattern:   "/api/*:rw",
			reqPath:   "api/users",
			wantMatch: true,
			wantPerm:  "rw",
		},
		{
			name:      "no权限匹配",
			pattern:   "/inner:no",
			reqPath:   "/inner",
			wantMatch: true,
			wantPerm:  "no",
		},
		{
			name:      "空pattern不匹配",
			pattern:   ":ro",
			reqPath:   "/test",
			wantMatch: false,
			wantPerm:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			matched, perm := c.pathMatch(tt.pattern, tt.reqPath)

			if matched != tt.wantMatch {
				t.Errorf("期望匹配 %v，实际为 %v", tt.wantMatch, matched)
			}
			if perm != tt.wantPerm {
				t.Errorf("期望权限 %q，实际为 %q", tt.wantPerm, perm)
			}
		})
	}
}
