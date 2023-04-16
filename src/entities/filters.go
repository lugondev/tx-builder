package entities

type AccountFilters struct {
	ID       int    `validate:"omitempty"`
	TenantID string `validate:"omitempty"`
	OwnerID  string `validate:"omitempty"`
}
