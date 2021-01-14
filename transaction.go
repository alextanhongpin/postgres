package postgres

import (
	"context"
	"database/sql"
)

type TxFn func(*sql.Tx) error

// WithTx creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTx(db *sql.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

type TxFnCtx func(context.Context, *sql.Tx) error

// WithTxContext creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTxContext(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn TxFnCtx) (err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)
	return err
}
