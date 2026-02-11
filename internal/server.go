package internal

import (
	"encoding/json"
	"net/http"
)

// APIResponse API 响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

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

// sendJSON 发送 JSON 响应
func (s *Server) sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

// sendError 发送错误响应
func (s *Server) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   message,
	})
}
