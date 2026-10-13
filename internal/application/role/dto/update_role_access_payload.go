package dto

type UpdateRoleAccessPayloadItem struct {
	AccessID  int
	HasAccess *bool
}

type UpdateRoleAccessPayload struct {
	AccessItem []UpdateRoleAccessPayloadItem
}
