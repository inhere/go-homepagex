package internal

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"
)

// BasicAuthMiddleware Basic 认证中间件
func (s *Server) BasicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.config.Auth.Enabled {
			next(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.requestAuth(w)
			return
		}

		// 解析 Basic Auth
		const prefix = "Basic "
		if !strings.HasPrefix(authHeader, prefix) {
			s.requestAuth(w)
			return
		}

		encoded := strings.TrimPrefix(authHeader, prefix)
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			s.requestAuth(w)
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			s.requestAuth(w)
			return
		}

		username, password := parts[0], parts[1]

		// 使用 constant-time 比较防止时序攻击
		usernameMatch := subtle.ConstantTimeCompare(
			[]byte(username),
			[]byte(s.config.Auth.Username),
		) == 1
		passwordMatch := subtle.ConstantTimeCompare(
			[]byte(password),
			[]byte(s.config.Auth.Password),
		) == 1

		if !usernameMatch || !passwordMatch {
			s.requestAuth(w)
			return
		}

		next(w, r)
	}
}

// requestAuth 请求认证
func (s *Server) requestAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="HomePagex"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
