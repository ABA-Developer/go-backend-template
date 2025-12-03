package permissions

import "be-dashboard-nba/constant"

var permissionHierarchy = map[constant.PermissionCode][]string{
	constant.ActionReadMenuPermission: {
		constant.ActionReadMenuPermission.String(),
		constant.ActionRead.String(),
	},
	constant.ActionCreateMenuPermission: {
		constant.ActionCreateMenuPermission.String(),
		constant.ActionCreate.String(),
	},
	constant.ActionUpdateMenuPermission: {
		constant.ActionUpdateMenuPermission.String(),
		constant.ActionUpdate.String(),
	},
	constant.ActionDeleteMenuPermission: {
		constant.ActionDeleteMenuPermission.String(),
		constant.ActionDelete.String(),
	},
}

func GetInheritedPermissions(code constant.PermissionCode) []string {
	if codesToTest, found := permissionHierarchy[code]; found {
		return codesToTest
	}
	return []string{code.String()}
}
