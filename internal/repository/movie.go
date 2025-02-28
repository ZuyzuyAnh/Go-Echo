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

type MovieRepository struct {
	DB *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) *MovieRepository {
	return &MovieRepository{
		DB: db,
	}
}

func (r *MovieRepository) InsertMovie(exec Executor, m *model.Movie) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		INSERT INTO movies(title, description, duration, cover_url, background_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{m.Title, m.Description, m.Duration, m.CoverURL, m.BackgroundURL}

	err := exec.QueryRowxContext(ctx, query, args...).Scan(&m.ID)
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return &DuplicateError{Message: fmt.Sprintf("movie with %s already exists", field)}
		}

		return err
	}

	return nil
}

func (r *MovieRepository) GetMovieByID(exec Executor, id int64) (*model.Movie, error) {
	if exec == nil {
		exec = r.DB
	}

	var m model.Movie

	query := `
		SELECT id, title, description, duration, cover_url, background_url
		FROM movies WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := exec.QueryRowxContext(ctx, query, id).Scan(
		&m.ID,
		&m.Title,
		&m.Description,
		&m.Duration,
		&m.CoverURL,
		&m.BackgroundURL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &NotFoundError{Message: fmt.Sprintf("movie with id %d not found", id)}
		}

		return nil, err
	}

	return &m, nil
}

func (r *MovieRepository) UpdateMovie(exec Executor, m *model.Movie) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		UPDATE movies
		SET title = $1, description = $2, duration = $3, cover_url = $4, background_url = $5
		WHERE id = $6
	`

	args := []any{m.Title, m.Description, m.Duration, m.CoverURL, m.BackgroundURL}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return &DuplicateError{fmt.Sprintf("user with %s already exists", field)}
		}

		return err
	}

	return nil
}

func (r *MovieRepository) DeleteMovie(exec Executor, id int64) error {
	if exec == nil {
		exec = r.DB
	}

	query := `DELETE FROM movies WHERE id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return &NotFoundError{Message: fmt.Sprintf("movie with id %d not found", id)}
	}

	return nil
}
