package dto

type UpdateUserPayload struct {
	Name     string
	FullName string
	Email    string
	RoleID   int
	Phone    *string
	Active   *bool
}
