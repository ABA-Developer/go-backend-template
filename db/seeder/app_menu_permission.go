package seeder

import (
	"log"

	"github.com/kristijorgji/goseeder"
)

type MenuPermission struct {
	ID         int
	MenuID     int
	Code       string
	ActionName string
	CreatedBy  string
}

var menuPermissions = []MenuPermission{
	{ID: 101, MenuID: 1, Code: "R", ActionName: "read", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},

	{ID: 301, MenuID: 3, Code: "C", ActionName: "create", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},
	{ID: 302, MenuID: 3, Code: "R", ActionName: "read", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},
	{ID: 303, MenuID: 3, Code: "U", ActionName: "update", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},
	{ID: 304, MenuID: 3, Code: "D", ActionName: "delete", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},

	{ID: 401, MenuID: 4, Code: "C", ActionName: "create", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},
	{ID: 402, MenuID: 4, Code: "R", ActionName: "read", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},
	{ID: 403, MenuID: 4, Code: "U", ActionName: "update", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},
	{ID: 404, MenuID: 4, Code: "D", ActionName: "delete", CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2"},
}

func menuPermissionSeeder(s goseeder.Seeder) {
	const statement = `
	INSERT INTO app_menu_permission(
		id, menu_id, code, action_name, created_by, created_at
	)
	VALUES($1, $2, $3, $4, $5, NOW())
	ON CONFLICT (id) DO UPDATE
	SET code = EXCLUDED.code, action_name = EXCLUDED.action_name, menu_id = EXCLUDED.menu_id;
	`
	for i := range menuPermissions {
		_, err := s.DB.Exec(statement,
			menuPermissions[i].ID,
			menuPermissions[i].MenuID,
			menuPermissions[i].Code,
			menuPermissions[i].ActionName,
			menuPermissions[i].CreatedBy,
		)
		if err != nil {
			log.Fatalf("❌ ERROR execute menu_permission seeder : %v", err.Error())
		}
	}
}
