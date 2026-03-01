package internal

// LoginInfo 当前登录用户信息
type LoginInfo struct {
	Username string `json:"username"`
}

// PageDataResponse 页面数据接口响应 DTO
type PageDataResponse struct {
	*PageConfig
	// 覆盖 PageConfig.Navs，仅包含当前用户有权限访问的导航项
	Navs []NavItem `json:"navs"`
	// 当前登录用户信息（游客访问时为 null）
	UserInfo *LoginInfo `json:"user_info,omitempty"`
}
