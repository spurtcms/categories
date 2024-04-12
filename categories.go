package categories

func CategoriesSetup(config Config) *Categories {

	MigrateTables(config.DB)

	return &Categories{
		DB:               config.DB,
		AuthEnable:       config.AuthEnable,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
		Permissions:      config.Permissions,
	}

}
