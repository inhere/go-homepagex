package internal

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/goutil/strutil"
)

const (
	// IconLocalPrefix 图标缓存路径前缀
	IconLocalPrefix = "icons-local"
	PageApiPrefix = "/api/page"
)

// HealthHandler 健康检查
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	s.sendJSON(w, map[string]string{"status": "ok"})
}

// GetIconLocalHandler 图标缓存处理
// 当 icon 路径以 icons-local/ 开头时，从本地缓存读取，若不存在则下载并缓存
func (s *Server) GetIconLocalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取图标路径 (去掉 icons-local/ 前缀)
	iconPath := strings.TrimPrefix(r.URL.Path, "/icons-local/")
	if iconPath == "" {
		s.sendError(w, "Icon path required", http.StatusBadRequest)
		return
	}

	// 构建本地缓存路径
	cacheDir := filepath.Join(s.config.FrontendDir, IconLocalPrefix)
	localPath := filepath.Join(cacheDir, iconPath)

	// 检查本地缓存是否存在
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		// 本地不存在，尝试下载
		iconCdnKey := strutil.BeforeFirst(iconPath, "/")
		baseRemoteUrl, ok := s.config.IconsCDN[iconCdnKey]
		if !ok {
			log.Printf("Icon CDN key %q not found in config.icons_cdn", iconCdnKey)
			s.sendError(w, "Icon not found", http.StatusNotFound)
			return
		}

		remoteURL := baseRemoteUrl + iconPath
		log.Printf("Icon cache miss: %s, downloading from: %s", iconPath, remoteURL)

		// 下载文件
		if err := downloadIconFile(remoteURL, localPath); err != nil {
			log.Printf("Failed to download icon: %v", err)
			s.sendError(w, "Cache icon error", http.StatusInternalServerError)
			return
		}

		log.Printf("Icon cached: %s -> %s", iconPath, localPath)
	}

	// 读取并返回本地文件
	data, err := os.ReadFile(localPath)
	if err != nil {
		log.Printf("Failed to read cached icon: %v", err)
		s.sendError(w, "Icon read error", http.StatusInternalServerError)
		return
	}

	// 设置 Content-Type
	contentType := getContentType(localPath)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400") // 缓存 24 小时
	w.Write(data)
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
