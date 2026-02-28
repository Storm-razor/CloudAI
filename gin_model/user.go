package ginmodel

// UserRegisterReq 用户注册请求
type UserRegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"` // 用户名(3-20字符)
	Password string `json:"password" binding:"required,min=6,max=30"` // 密码(6-30字符)
	Email    string `json:"email" binding:"required,email"`           // 邮箱
	Phone    string `json:"phone" binding:"required,e164"`            // 手机号(E.164格式,如+8613800000000)
}

// LoginResponse 登录成功响应
type LoginResponse struct {
	AccessToken string `json:"access_token"` // JWT访问令牌
	ExpiresIn   int    `json:"expires_in"`   // 过期时间(秒)
	TokenType   string `json:"token_type"`   // 令牌类型(Bearer)
}

// UserNameLoginReq 用户名密码登录
type UserNameLoginReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"` // 用户名
	Password string `json:"password" binding:"required,min=8,max=30"` // 密码
}

// UserPhoneLogin 手机号密码登录
type UserPhoneLogin struct {
	Phone    string `json:"phone" binding:"required,e164"`            // 手机号
	Password string `json:"password" binding:"required,min=8,max=30"` // 密码
}
