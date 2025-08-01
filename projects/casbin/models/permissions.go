package models

// Permission 权限
type Permission struct {
	CommonField
	Name string // e.g., "product:read", "user:create"
}

// RolePermission 角色权限关联
type RolePermission struct {
	CommonField
	RoleId       uint64
	PermissionId uint64
}
