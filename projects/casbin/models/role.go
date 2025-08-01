package models

// Role 角色
type Role struct {
	CommonField
	Name     string // e.g., admin, user, viewer
	TenantId uint64 // 👥 角色也与租户绑定
}
