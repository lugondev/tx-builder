package postgres

import (
	"context"
)

//go:generate mockgen -source=postgres.go -destination=mocks/postgres.go -package=mocks

type Client interface {
	ModelContext(ctx context.Context, models ...interface{}) Query
	QueryOneContext(ctx context.Context, model interface{}, query interface{}, params ...interface{}) error
	RunInTransaction(ctx context.Context, persist func(client Client) error) error
	Exec(query interface{}, params ...interface{}) error
	Close() error
}

type Query interface {
	WherePK() Query
	Where(condition string, params ...interface{}) Query
	WhereAllowedOwner(ownerIDLabel, ownerID string) Query
	WhereAllowedTenants(tenantIDLabel string, tenants []string) Query
	Order(order string) Query
	OrderExpr(order string, params ...interface{}) Query
	Relation(name string, apply ...func(query Query) (Query, error)) Query
	Column(columns ...string) Query
	ColumnExpr(expr string, params ...interface{}) Query
	Join(join string, params ...interface{}) Query
	OnConflict(s string, params ...interface{}) Query
	Set(set string, params ...interface{}) Query
	Returning(s string, params ...interface{}) Query
	For(s string, params ...interface{}) Query
	Insert() error
	Update() error
	UpdateNotZero() error
	Select() error
	SelectOne() error
	SelectOrInsert() error
	SelectColumn(result interface{}) error
	Delete() error
}
