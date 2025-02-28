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

type TheaterRepository struct {
	DB *sqlx.DB
}

func NewTheaterRepository(db *sqlx.DB) *TheaterRepository {
	return &TheaterRepository{
		DB: db,
	}
}

func (r *TheaterRepository) CreateTheater(exec Executor, t *model.Theater) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		INSERT INTO theaters(name)
		VALUES($1)
		RETURNING id;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := exec.QueryRowxContext(ctx, query, t.Name).Scan(&t.ID)
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return &DuplicateError{Message: fmt.Sprintf("theater with %s already exists", field)}
		}
		return err
	}
	return nil
}

func (r *TheaterRepository) GetTheaterByID(exec Executor, id int64) (*model.Theater, error) {
	if exec == nil {
		exec = r.DB
	}

	var t model.Theater
	query := `
		SELECT id, name
		FROM theaters
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := exec.QueryRowxContext(ctx, query, id).Scan(&t.ID, &t.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &NotFoundError{Message: fmt.Sprintf("theater with id %d not found", id)}
		}
		return nil, err
	}

	return &t, nil
}

func (r *TheaterRepository) ListTheaters(exec Executor) ([]*model.Theater, error) {
	if exec == nil {
		exec = r.DB
	}

	query := `
		SELECT id, name
		FROM theaters;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := exec.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var theaters []*model.Theater
	for rows.Next() {
		var t model.Theater
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		theaters = append(theaters, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return theaters, nil
}

func (r *TheaterRepository) UpdateTheater(exec Executor, t *model.Theater) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		UPDATE theaters
		SET name = $1
		WHERE id = $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exec.ExecContext(ctx, query, t.Name, t.ID)
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return &DuplicateError{Message: fmt.Sprintf("theater with %s already exists", field)}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &NotFoundError{Message: fmt.Sprintf("theater with id %d not found", t.ID)}
	}

	return nil
}

func (r *TheaterRepository) DeleteTheater(exec Executor, id int64) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		DELETE FROM theaters
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &NotFoundError{Message: fmt.Sprintf("theater with id %d not found", id)}
	}

	return nil
}
