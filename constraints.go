package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/lib/pq"
)

type CheckConstraintsOption struct {
	Column string
}

type CheckConstraintsOptionModifier func(o *CheckConstraintsOption)

func WithColumn(column string) CheckConstraintsOptionModifier {
	return func(o *CheckConstraintsOption) {
		o.Column = column
	}
}

// CheckConstraints checks if the postgres error matches the errcodes
// https://www.postgresql.org/docs/current/errcodes-appendix.html
func CheckConstraints(err error, constraint string, opts ...CheckConstraintsOptionModifier) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		constraint := pqErr.Code.Name() == constraint
		if len(opts) > 0 {
			var opt CheckConstraintsOption
			for _, modify := range opts {
				modify(&opt)
			}
			return constraint && match(pqErr.Detail, opt.Column)
		}
		return constraint
	}
	return false
}

func IsUniqueViolation(err error, opts ...CheckConstraintsOptionModifier) bool {
	return CheckConstraints(err, "unique_violation", opts...)
}

func IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}

func match(src, tgt string) bool {
	ok, _ := regexp.MatchString(fmt.Sprintf("\\b%s\\b", tgt), src)
	return ok
}
