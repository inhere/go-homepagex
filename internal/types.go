package internal

type LoginInfo struct {
	Username string `json:"username,omitempty"`
	Perm     string `json:"perm,omitempty"`
}

type PageDataResponse struct {
	*PageConfig
	// （当已经登录后）登录用户信息 TODO
	LoginInfo *LoginInfo `json:"login_info,omitempty"`
}
