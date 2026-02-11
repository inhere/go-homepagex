# Home Dashboard

一个类似 Homer 的导航主页，使用 Go + Svelte 实现。

## 项目结构

```
home-dashboard/
├── internal/          # Go 后端服务
│   ├── config.go     # 配置加载
│   ├── page.go       # 页面配置解析
│   ├── auth.go       # Basic 认证
│   └── handlers.go   # HTTP 处理器
├── frontend/         # 前端应用
│   └── build/        # 构建输出
│       ├── index.html
│       └── app.js
├── pages/            # 页面 YAML 配置
│   ├── main.yaml     # 主页面配置
│   └── page-another.yaml
├── config.yaml   # 后端配置
├── main.go       # Go 入口文件
└── README.md
```

## 功能特性

- **YAML 配置**: 类似 Homer 的 YAML 格式定义页面链接
- **多页面支持**: `/` 对应 `main.yaml`，`/another` 对应 `page-another.yaml`
- **视图切换**: 支持卡片视图和列表视图切换
- **Basic 认证**: 可配置的用户名/密码验证
- **FontAwesome 图标**: 支持 FontAwesome 图标
- **响应式设计**: 适配桌面和移动设备

## 快速开始

### 1. 安装依赖

**Go (1.21+)**:
```bash
# 下载并安装 Go
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2. 编译后端

```bash
cd backend
go mod tidy
go build -o home-dashboard
```

### 3. 运行

```bash
# 使用默认配置
./backend/home-dashboard

# 或使用自定义配置文件
./backend/home-dashboard /path/to/config.yaml
```

### 4. 访问

打开浏览器访问: http://localhost:8080

## 配置说明

### 后端配置 (config/config.yaml)

```yaml
server:
  port: "8080"

# Basic 认证配置
auth:
  enabled: true
  username: "admin"
  password: "your-password"

# 页面配置文件存放目录
pages_dir: "./pages"

# 前端构建目录
frontend_dir: "./frontend/build"
```

### 页面配置 (pages/main.yaml)

```yaml
title: "Home Dashboard"
subtitle: "Welcome to your dashboard"
logo: "logo.png"

theme: "default"
color: "blue"
style: "cards"  # cards 或 list
columns: "3"

connectivity:
  check_interval: 30000
  mode: "ping"

services:
  - name: "Media"
    icon: "fas fa-play-circle"
    items:
      - name: "Plex"
        logo: "https://example.com/plex.png"
        subtitle: "Media server"
        tag: "app"
        url: "https://plex.example.com"
        target: "_blank"
```

## 页面配置规则

- `/` 路由对应 `pages/main.yaml`
- `/another` 路由对应 `pages/page-another.yaml`
- 依此类推: `{route}` -> `page-{route}.yaml`

## 开发

### 前端开发

前端使用原生 JavaScript，无需构建步骤。直接编辑 `frontend/build/` 目录下的文件即可。

### 后端开发

```bash
cd backend
go run . ../config/config.yaml
```

## 图标

支持 FontAwesome 图标:
- 使用 `fas fa-icon-name` 格式
- 完整图标列表: https://fontawesome.com/icons

## 许可证

MIT
