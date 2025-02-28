package dto

import (
	"echo-demo/internal/validator"
)

type SignupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type SignUpResponse struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 6, "password", "must be at least 6 characters")
}

func (s SignupRequest) Validate(v *validator.Validator) {
	v.Check(s.Name != "", "name", "must be provided")
	v.Check(len(s.Name) <= 100, "name", "must not be more than 100 characters")

	ValidateEmail(v, s.Email)
	ValidatePassword(v, s.Password)

	v.Check(s.PhoneNumber != "", "phone_number", "must be provided")
	v.Check(len(s.PhoneNumber) >= 10 && len(s.PhoneNumber) <= 11, "phone_number", "must be between 10 and 11 characters")
}

func (l LoginRequest) Validate(v *validator.Validator) {
	ValidateEmail(v, l.Email)
	ValidatePassword(v, l.Password)

}
