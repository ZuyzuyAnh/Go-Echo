package repository

import (
	"context"
	"echo-demo/internal/model"
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
