package dto

import "echo-demo/internal/validator"

type CreateUpdateTheaterRequest struct {
	Name string `json:"name"`
}

type TheaterResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (s *CreateUpdateTheaterRequest) Validate(v *validator.Validator) error {
	v.Check(s.Name != "", "name", "must be provided")
	v.Check(len(s.Name) <= 100, "name", "must not be more than 100 characters")
	return nil
}
