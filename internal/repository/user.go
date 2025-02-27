package repository

import (
	"context"
	"database/sql"
	"echo-demo/internal/model"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) InsertUser(exec Executor, u *model.User) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		INSERT INTO users (email, name, password, phone_number)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{u.Email, u.Name, u.Password, u.Phone}

	err := exec.QueryRowxContext(ctx, query, args...).Scan(&u.ID)
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return &DuplicateError{fmt.Sprintf("user with %s already exists", field)}
		}

		return err
	}

	return nil
}

func (r *UserRepository) GetUserByID(exec Executor, id int64) (*model.User, error) {
	if exec == nil {
		exec = r.DB
	}

	query := `
		SELECT id, email, name, password, phone_number
		FROM users
		WHERE id = $1
	`

	var u model.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := exec.QueryRowxContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Password,
		&u.Phone,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, &NotFoundError{"user not found"}
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (r *UserRepository) GetUserByEmail(exec Executor, email string) (*model.User, error) {
	if exec == nil {
		exec = r.DB
	}

	query := `
		SELECT id, email, name, password, phone_number
		FROM users
		WHERE email = $1
	`

	var u model.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := exec.QueryRowxContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Password,
		&u.Phone,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, &NotFoundError{"user not found"}
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (r *UserRepository) UpdateUser(exec Executor, u *model.User) (*model.User, error) {
	if exec == nil {
		exec = r.DB
	}

	query := `
		UPDATE users
		SET email = $2, name = $3, password = $4, phone = $5
		WHERE id = $6
	`

	args := []any{
		u.Email,
		u.Name,
		u.Password,
		u.Phone,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := exec.QueryRowxContext(ctx, query, args...).Scan(&u.ID)
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return nil, &DuplicateError{fmt.Sprintf("user with %s already exists", field)}
		}

		return nil, err
	}

	return u, nil
}
