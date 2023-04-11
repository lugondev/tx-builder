package gopg

import (
	"fmt"
	"github.com/lugondev/tx-builder/src/infra/postgres"
	"github.com/samber/lo"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
)

type Query struct {
	pgQuery *orm.Query
}

var _ postgres.Query = &Query{}

func NewQuery(pgQuery *orm.Query) *Query {
	return &Query{pgQuery: pgQuery}
}

func (q Query) WherePK() postgres.Query {
	q.pgQuery = q.pgQuery.WherePK()

	return &q
}

func (q Query) Where(condition string, params ...interface{}) postgres.Query {
	q.pgQuery = q.pgQuery.Where(condition, params...)

	return &q
}

func (q Query) Order(order string) postgres.Query {
	q.pgQuery = q.pgQuery.Order(order)

	return &q
}

func (q Query) OrderExpr(order string, params ...interface{}) postgres.Query {
	q.pgQuery = q.pgQuery.OrderExpr(order, params...)

	return &q
}

func (q Query) Column(columns ...string) postgres.Query {
	q.pgQuery = q.pgQuery.Column(columns...)

	return &q
}

func (q Query) ColumnExpr(expr string, params ...interface{}) postgres.Query {
	q.pgQuery = q.pgQuery.ColumnExpr(expr, params...)

	return &q
}

func (q Query) Join(join string, params ...interface{}) postgres.Query {
	q.pgQuery = q.pgQuery.Join(join, params...)

	return &q
}

func (q Query) OnConflict(s string, params ...interface{}) postgres.Query {
	q.pgQuery = q.pgQuery.OnConflict(s, params...)

	return &q
}

func (q Query) Set(set string, params ...interface{}) postgres.Query {
	q.pgQuery = q.pgQuery.Set(set, params...)

	return &q
}

func (q Query) Returning(s string, params ...interface{}) postgres.Query {
	q.pgQuery = q.pgQuery.Returning(s, params...)

	return &q
}

func (q Query) Relation(name string, apply ...func(query postgres.Query) (postgres.Query, error)) postgres.Query {
	if len(apply) > 0 {
		relationFunc := func(intQ *orm.Query) (*orm.Query, error) {
			newQ, err := apply[0](NewQuery(intQ))
			if err != nil {
				return nil, parseErrorResponse(err)
			}
			return newQ.(*Query).pgQuery, err
		}

		q.pgQuery = q.pgQuery.Relation(name, relationFunc)
	} else {
		q.pgQuery = q.pgQuery.Relation(name)
	}

	return &q
}

func (q Query) For(s string, params ...interface{}) postgres.Query {
	q.pgQuery = q.pgQuery.For(s, params...)

	return &q
}

func (q *Query) Insert() error {
	_, err := q.pgQuery.Insert()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (q *Query) UpdateNotZero() error {
	_, err := q.pgQuery.UpdateNotZero()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (q *Query) Update() error {
	_, err := q.pgQuery.Update()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (q *Query) Select() error {
	err := q.pgQuery.Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (q *Query) SelectOne() error {
	err := q.pgQuery.First()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (q *Query) SelectOrInsert() error {
	_, err := q.pgQuery.SelectOrInsert()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (q *Query) SelectColumn(result interface{}) error {
	err := q.pgQuery.Select(result)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (q *Query) Delete() error {
	_, err := q.pgQuery.Delete()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (q Query) WhereAllowedTenants(tenantIDLabel string, tenants []string) postgres.Query {
	if tenantIDLabel == "" {
		tenantIDLabel = "tenant_id"
	}

	return q.whereAllowedTenants(tenantIDLabel, tenants)
}

func (q Query) WhereAllowedOwner(ownerIDLabel, ownerID string) postgres.Query {
	if ownerIDLabel == "" {
		ownerIDLabel = "owner_id"
	}

	return q.whereAllowedOwner(ownerIDLabel, ownerID)
}

func (q Query) whereAllowedTenants(field string, tenants []string) postgres.Query {
	if len(tenants) == 0 || lo.Contains(tenants, multitenancy.WildcardTenant) {
		return &q
	}
	return q.Where(fmt.Sprintf("%s IN (?)", field), pg.In(tenants))
}

func (q Query) whereAllowedOwner(field, ownerID string) postgres.Query {
	if ownerID == multitenancy.WildcardOwner {
		return &q
	}

	if ownerID != "" {
		return q.Where(fmt.Sprintf("%s = ? OR %s IS NULL", field, field), ownerID)
	}

	return q.Where(fmt.Sprintf("%s is NULL", field))
}
