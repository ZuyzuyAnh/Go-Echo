package dto

import (
	"echo-demo/internal/validator"
)

type SignupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

type SignUpResponse struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 6, "password", "must be at least 6 characters")
}

func (s *SignupRequest) Validate(v *validator.Validator) error {
	v.Check(s.FullName != "", "name", "must be provided")
	v.Check(len(s.FullName) <= 100, "name", "must not be more than 100 characters")

	ValidateEmail(v, s.Email)
	ValidatePassword(v, s.Password)

	v.Check(s.PhoneNumber != "", "phone_number", "must be provided")
	v.Check(len(s.PhoneNumber) >= 10 && len(s.PhoneNumber) <= 11, "phone_number", "must be between 10 and 11 characters")

	return nil
}

func (l *LoginRequest) Validate(v *validator.Validator) error {
	ValidateEmail(v, l.Email)
	ValidatePassword(v, l.Password)

	return nil
}
