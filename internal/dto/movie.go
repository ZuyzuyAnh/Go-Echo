package dto

import "echo-demo/internal/validator"

type CreateMovieReq struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      int64  `json:"duration"`
	CoverURL      string `json:"cover_url"`
	BackgroundURL string `json:"background_url"`
}

type MovieResp struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      int64  `json:"duration"`
	CoverURL      string `json:"cover_url"`
	BackgroundURL string `json:"background_url"`
}

func (r CreateMovieReq) Validate(v *validator.Validator) {
	v.Check(r.Title != "", "title", "must be provided")
	v.Check(r.Description != "", "description", "must be provided")
	v.Check(r.Duration != 0, "duration", "must be provided")

	v.Check(len(r.Title) < 100, "title", "must not be more than 100 characters")
	v.Check(len(r.Description) < 255, "description", "must not be more than 255 characters")
}

func (r CreateMovieReq) ValidateUpdate(v *validator.Validator) {
	v.Check(len(r.Title) < 100, "title", "must not be more than 100 characters")
	v.Check(len(r.Description) < 255, "description", "must not be more than 255 characters")
}
