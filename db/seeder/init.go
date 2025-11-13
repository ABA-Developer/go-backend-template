package seeder

import "github.com/kristijorgji/goseeder"

func init() {
	goseeder.Register(userSeeder)
	goseeder.Register(roleSeeder)
	goseeder.Register(menuSeeder)
	goseeder.Register(userRoleSeeder)
	goseeder.Register(menuPermissionSeeder)
	goseeder.Register(roleAccessSeeder)
}
