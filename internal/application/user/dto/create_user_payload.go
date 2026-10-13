package dto

type CreateUserPayload struct {
	Name     string
	FullName string
	Email    string
	Password string
	RoleID   int
	Phone    *string
	Active   *bool
}
