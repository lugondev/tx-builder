package gopg

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/lugondev/tx-builder/src/infra/postgres"
	"github.com/lugondev/wallet-signer-manager/pkg/errors"
)

type Client struct {
	db orm.DB
}

var _ postgres.Client = &Client{}

func New(appName string, cfg *Config) (*Client, error) {
	pgOptions, err := cfg.ToPGOptions(appName)
	if err != nil {
		return nil, err
	}

	return &Client{db: pg.Connect(pgOptions)}, nil
}

func (c *Client) ModelContext(ctx context.Context, models ...interface{}) postgres.Query {
	return NewQuery(c.db.ModelContext(ctx, models...))
}

func (c *Client) QueryOneContext(ctx context.Context, model, query interface{}, params ...interface{}) error {
	_, err := c.db.QueryOneContext(ctx, model, query, params...)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c Client) RunInTransaction(ctx context.Context, persist func(client postgres.Client) error) error {
	persistFunc := func(tx *pg.Tx) error {
		c.db = tx
		return persist(&c)
	}

	// CheckPermission whether we already are in a tx or not to allow for nested DB transactions
	dbtx, isTx := c.db.(*pg.Tx)
	if isTx {
		return persistFunc(dbtx)
	}

	return c.db.(*pg.DB).RunInTransaction(ctx, persistFunc)
}

func (c *Client) Exec(query interface{}, params ...interface{}) error {
	_, err := c.db.Exec(query, params)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *Client) Close() error {
	pgDB, ok := c.db.(*pg.DB)
	if !ok {
		return errors.PostgresError("cannot close connection on a PG transaction")
	}

	err := pgDB.Close()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}
