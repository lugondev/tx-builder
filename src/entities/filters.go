package entities

import "time"

type JobFilters struct {
	TxHashes      []string  `validate:"omitempty,unique,dive,isHash"`
	ChainUUID     string    `validate:"omitempty,uuid"`
	UpdatedAfter  time.Time `validate:"omitempty"`
	ParentJobUUID string    `validate:"omitempty"`
	OnlyParents   bool      `validate:"omitempty"`
	WithLogs      bool      `validate:"omitempty"`
}

type AccountFilters struct {
	Aliases  []string `validate:"omitempty,unique"`
	TenantID string   `validate:"omitempty"`
}
