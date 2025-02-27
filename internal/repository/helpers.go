package repository

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

type NotFoundError struct {
	Message string
}

func (n *NotFoundError) Error() string {
	return n.Message
}

type DuplicateError struct {
	Message string
}

func (n *DuplicateError) Error() string {
	return n.Message
}

type ForeignKeyError struct {
	Message string
}

func (n *ForeignKeyError) Error() string {
	return n.Message
}

func isDuplicateError(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		parts := strings.Split(pgErr.ConstraintName, "_")

		if len(parts) >= 3 {
			return parts[1], true
		}

		return pgErr.ConstraintName, true
	}
	return "", false
}

func isForeignKeyError(err error) (table string, column string, ok bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		parts := strings.Split(pgErr.ConstraintName, "_")

		if len(parts) >= 2 {
			table = parts[0]
			column = parts[1]

			return
		}
	}

	return "", "", false
}
