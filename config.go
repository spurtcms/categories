package categories

import (
	"github.com/spurtcms/auth"
	role "github.com/spurtcms/team-roles"
	"gorm.io/gorm"
)

type Config struct {
	DB               *gorm.DB
	AuthEnable       bool
	PermissionEnable bool
	Auth             *auth.Auth
	Permissions      *role.PermissionConfig
}

type Categories struct {
	DB               *gorm.DB
	AuthEnable       bool
	PermissionEnable bool
	Auth             *auth.Auth
	Permissions      *role.PermissionConfig
}
