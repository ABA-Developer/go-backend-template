package dto

type CreateMenuPayload struct {
	ParentID    *int
	Name        string
	Description *string
	URL         *string
	Group       string
	Icon        *string
	Active      *bool
	Display     *bool
}
