package dto

import "echo-demo/internal/validator"

type CreateUpdateTheaterRequest struct {
	Name string `json:"name"`
	Rows int    `json:"rows"`
	Cols int    `json:"cols"`
}

type TheaterResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateSeatTypeRequest struct {
	SeatIDs    []int64 `json:"seat_ids"`
	SeatTypeID int64   `json:"seat_type_id"`
}

func (s CreateUpdateTheaterRequest) Validate(v *validator.Validator) {
	v.Check(s.Name != "", "name", "must be provided")
	v.Check(len(s.Name) <= 100, "name", "must not be more than 100 characters")

	v.Check(s.Rows > 0 && s.Rows < 20, "rows", "must not be more than 20 characters or empty")
	v.Check(s.Cols > 0 && s.Cols < 20, "cols", "must not be more than 20 characters or empty")
}

func (s UpdateSeatTypeRequest) Validate(v *validator.Validator) {
	v.Check(s.SeatIDs != nil, "seat_ids", "must be provided")
	v.Check(s.SeatTypeID > 0, "seat_type_id", "must be provided")
}
