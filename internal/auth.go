package internal

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

type contextKey string

// ContextKeyUsername request context 中存储已认证用户名的 key
const ContextKeyUsername contextKey = "username"

func (s *Server) BasicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path
		// 权限按页面配置，有页面的权限就有对应 api 的权限（api 访问去除 api/page 前缀后检查）
		if after, ok := strings.CutPrefix(reqPath, PageApiPrefix); ok {
			reqPath = after
		}

		isWrite := r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// 有认证信息：解析并验证
			after, ok := strings.CutPrefix(authHeader, "Basic ")
			if !ok {
				s.requestAuth(w)
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(after)
			if err != nil {
				log.Printf("BasicAuthMiddleware: base64 decode error: %v", err)
				s.requestAuth(w)
				return
			}

			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) == 2 {
				username, password := parts[0], parts[1]

				authConfig, exists := s.config.MatchUserAuthConfig(username, reqPath)
				if exists {
					if authConfig.Permission == PermNO {
						s.requestAuth(w)
						return
					}

					usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(authConfig.Username)) == 1
					passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(authConfig.Password)) == 1

					if usernameMatch && passwordMatch {
						// 写操作需要 rw 权限
						if isWrite && authConfig.Permission != PermRW {
							http.Error(w, "Forbidden", http.StatusForbidden)
							return
						}
						// 将认证用户名注入 request context
						ctx := context.WithValue(r.Context(), ContextKeyUsername, username)
						next(w, r.WithContext(ctx))
						return
					}
				} else {
					// 用户名不在配置中：回退到公开访问检查
					// 场景：前端退出登录后发送空凭据（Basic Og==），应视为游客
					if !s.config.IsNeedAuth(reqPath, isWrite) {
						ctx := context.WithValue(r.Context(), ContextKeyUsername, "")
						next(w, r.WithContext(ctx))
						return
					}
				}
			}

			s.requestAuth(w)
			return
		}

		// 没有认证信息：检查 path 是否允许匿名访问
		if s.config.IsNeedAuth(reqPath, isWrite) {
			s.requestAuth(w)
			return
		}
		// 公开访问，注入空用户名（游客）
		ctx := context.WithValue(r.Context(), ContextKeyUsername, "")
		next(w, r.WithContext(ctx))
	}
}

func (s *Server) requestAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="HomePagex"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

// AuthHandler 触发 Basic Auth 认证，成功后重定向
// 用于前端"登录"按钮：浏览器导航到此 URL 时弹出认证对话框，认证成功后跳回原页面
func (s *Server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		s.requestAuth(w)
		return
	}

	after, ok := strings.CutPrefix(authHeader, "Basic ")
	if !ok {
		s.requestAuth(w)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(after)
	if err != nil {
		s.requestAuth(w)
		return
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) == 2 {
		username, password := parts[0], parts[1]
		if s.config.CheckCredentials(username, password) {
			returnURL := r.URL.Query().Get("return")
			if returnURL == "" || !strings.HasPrefix(returnURL, "/") {
				returnURL = "/"
			}
			http.Redirect(w, r, returnURL, http.StatusFound)
			return
		}
	}

	s.requestAuth(w)
}

// LogoutHandler 退出登录：始终返回 401，触发浏览器清除缓存的 Basic Auth 凭据
func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="HomePagex"`)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusUnauthorized)
}
