package internal

import (
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

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
						next(w, r)
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
		next(w, r)
	}
}

func (s *Server) requestAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="HomePagex"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
