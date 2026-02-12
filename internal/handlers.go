package internal

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/goutil/strutil"
)

// HealthHandler 健康检查
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	s.sendJSON(w, map[string]string{"status": "ok"})
}

// GetPageConfigHandler 获取页面配置数据
func (s *Server) GetPageConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从 URL 路径获取路由
	path := strings.TrimPrefix(r.URL.Path, "/api/page")
	// 检查 refresh 参数
	refresh := strutil.SafeBool(r.URL.Query().Get("refresh"))

	pageConfig, err := PageDataMgr.GetPageConfig(path, refresh)
	if err != nil {
		log.Printf("Error loading page data for %s: %v", r.URL.Path, err)
		s.sendError(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Request API GET %s, Pagefile: %s", r.RequestURI, pageConfig.Pagefile)
	s.sendJSON(w, pageConfig)
}

// StaticFileHandler 静态文件服务
func (s *Server) StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	// 清理路径防止目录遍历
	path := strings.TrimLeft(r.URL.Path, "/.")
	if path == "" {
		path = "index.html"
	}

	// 构建完整路径
	fullPath := filepath.Join(s.config.FrontendDir, path)

	// 检查文件是否存在
	info, err := os.Stat(fullPath)
	if err != nil {
		// 如果是目录，尝试 index.html
		if info != nil && info.IsDir() {
			fullPath = filepath.Join(fullPath, "index.html")
		} else {
			extName := filepath.Ext(path)
			if extName == "" {
				// 返回前端应用的 index.html（支持前端路由）
				fullPath = filepath.Join(s.config.FrontendDir, "index.html")
			} else {
				log.Printf("NOTICE File not found: %s", fullPath)
				// 否则返回 404
				s.sendError(w, "File not found", http.StatusNotFound)
				return
			}
		}
	}

	log.Printf("Request static: %s, Serving file: %s", r.URL.Path, fullPath)

	// 设置正确的 Content-Type
	contentType := getContentType(fullPath)
	w.Header().Set("Content-Type", contentType)

	http.ServeFile(w, r, fullPath)
}

// getContentType 根据文件扩展名获取 Content-Type
func getContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".html":
		return "text/html"
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
