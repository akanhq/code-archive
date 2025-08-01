package models

// User 用户
type User struct {
	CommonField
	Username string `json:"username"`
	TenantId uint64 // 🔑 多租户关键字段
	RoleId   uint64
}
