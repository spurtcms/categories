package categories

import "errors"

var (
	ErrorAuth         = errors.New("auth enabled not initialised")
	ErrorPermission   = errors.New("permissions enabled not initialised")
	ErrorCategoryName = errors.New("given some values is empty")
)

func AuthandPermission(category *Categories) error {

	//check auth enable if enabled, use auth pkg otherwise it will return error
	if category.AuthEnable && !category.Auth.AuthFlg {

		return ErrorAuth
	}
	//check permission enable if enabled, use team-role pkg otherwise it will return error
	if category.PermissionEnable && !category.Permissions.PermissionFlg {

		return ErrorPermission

	}

	return nil
}
