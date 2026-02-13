package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// downloadIconFile 下载图标文件到本地缓存
func downloadIconFile(remoteURL, localPath string) error {
	// 下载文件
	resp, err := client.Get(remoteURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download icon: %s, status: %d", remoteURL, resp.StatusCode)
	}

	// 创建目录
	if err = os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}

	// 写入文件
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}