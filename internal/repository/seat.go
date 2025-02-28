package repository

import (
	"context"
	"echo-demo/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type SeatRepository struct {
	DB *sqlx.DB
}

func NewSeatRepository(db *sqlx.DB) *SeatRepository {
	return &SeatRepository{
		DB: db,
	}
}

func (r *SeatRepository) CreateSeat(exec Executor, s *model.Seat) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		INSERT INTO seats (number, theater_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := exec.QueryRowxContext(ctx, query, s.Number, s.TheaterID).Scan(&s.ID)
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return &DuplicateError{Message: fmt.Sprintf("seat with %s already exists", field)}
		}
		return err
	}
	return nil
}

func (r *SeatRepository) ListSeatsByTheater(exec Executor, theaterID int64) ([]*model.Seat, error) {
	if exec == nil {
		exec = r.DB
	}

	query := `
		SELECT id, number, theater_id, seat_type_id
		FROM seats
		WHERE theater_id = $1
		ORDER BY number;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := exec.QueryxContext(ctx, query, theaterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []*model.Seat
	for rows.Next() {
		var s model.Seat
		if err := rows.Scan(&s.ID, &s.Number, &s.TheaterID, &s.SeatTypeID); err != nil {
			return nil, err
		}
		seats = append(seats, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seats, nil
}

func (r *SeatRepository) UpdateSeat(exec Executor, s *model.Seat) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		UPDATE seats
		SET number = $1, theater_id = $2, seat_type_id = $3
		WHERE id = $4;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exec.ExecContext(ctx, query, s.Number, s.TheaterID, s.SeatTypeID, s.ID)
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return &DuplicateError{Message: fmt.Sprintf("seat with %s already exists", field)}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &NotFoundError{Message: fmt.Sprintf("seat with id %d not found", s.ID)}
	}
	return nil
}

func (r *SeatRepository) BatchUpdateSeatType(exec Executor, seatTypeID int64, seatIDs []int64) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		UPDATE seats
		SET seat_type_id = $1
		WHERE id = ANY($2);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exec.ExecContext(ctx, query, seatTypeID, pq.Array(seatIDs))
	if err != nil {
		if field, ok := isDuplicateError(err); ok {
			return &DuplicateError{Message: fmt.Sprintf("seat with %s already exists", field)}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &NotFoundError{Message: "no seats found to update"}
	}
	return nil
}

func (r *SeatRepository) DeleteSeat(exec Executor, id int64) error {
	if exec == nil {
		exec = r.DB
	}

	query := `
		DELETE FROM seats
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
		return &NotFoundError{Message: fmt.Sprintf("seat with id %d not found", id)}
	}
	return nil
}
