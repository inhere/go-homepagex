package internal

import "net/http"

// Server HTTP 服务器
type Server struct {
	config *Config
	fs     http.FileSystem
}

// NewServer 创建新的 HTTP 服务器
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		fs:     http.Dir(config.FrontendDir),
	}
}
