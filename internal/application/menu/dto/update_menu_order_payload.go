package dto

type UpdateMenuOrderPayload struct {
	Group     string
	ParentID  *int
	SortedIDs []int
}
