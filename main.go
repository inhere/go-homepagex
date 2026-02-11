package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/inhere/homepagex/internal"
)

func main() {
	// 默认配置文件路径
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// 加载配置
	config, err := internal.LoadConfig(configPath)
	if err != nil {
		log.Printf("Warning: %v, using defaults", err)
		config = &internal.Config{
			Server: internal.ServerConfig{
				Port: "8080",
			},
			Auth: internal.AuthConfig{
				Enabled: false,
			},
			PagesDir:    "./pages",
			FrontendDir: "./frontend/build",
		}
	}

	server := internal.NewServer(config)

	// 设置路由
	mux := http.NewServeMux()

	// API 路由
	mux.HandleFunc("/api/page", server.BasicAuthMiddleware(server.GetPageConfigHandler))
	mux.HandleFunc("/api/page/", server.BasicAuthMiddleware(server.GetPageConfigHandler))
	mux.HandleFunc("/api/health", server.HealthHandler)

	// 静态文件路由（前端应用）
	mux.HandleFunc("/", server.BasicAuthMiddleware(server.StaticFileHandler))

	// 启动服务器
	addr := fmt.Sprintf(":%s", config.Server.Port)
	log.Printf("🚀 Starting server on http://localhost%s", addr)
	log.Printf("Page data directory: %s", config.PagesDir)
	log.Printf("Frontend directory: %s", config.FrontendDir)
	fmt.Println()
	if config.Auth.Enabled {
		log.Printf("Basic auth enabled for user: %s", config.Auth.Username)
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
