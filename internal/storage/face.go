package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type PGer interface {
	DBSX() *sqlx.DB
	QueryTx(ctx context.Context, f func(tx *sqlx.Tx) error) error
}
