package gopg

import (
	"github.com/go-pg/pg/v10"

	"github.com/lugondev/tx-builder/pkg/errors"
)

func parseErrorResponse(err error) error {
	if pg.ErrNoRows == err {
		return errors.NotFoundError("resource not found")
	}
	if pg.ErrMultiRows == err {
		return errors.InvalidStateError("multiple resources found, only expected one")
	}

	pgErr, ok := err.(pg.Error)
	if !ok {
		return errors.PostgresConnectionError(err.Error())
	}

	switch {
	case pgErr.IntegrityViolation():
		return errors.ConstraintViolatedError(pgErr.Error())
	default:
		return errors.PostgresConnectionError(pgErr.Error())
	}
}
