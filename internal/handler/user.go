package handler

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/model"
	"echo-demo/internal/validator"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (h *Handler) Signup(c echo.Context) error {
	req := new(dto.SignupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	v := validator.NewValidator()
	if req.Validate(v); !v.Valid() {
		return c.JSON(http.StatusBadRequest, v.Errors)
	}

	tx := h.DB.Begin()
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	var existingUser model.User
	err := tx.Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).First(&existingUser).Error
	if err == nil {
		if existingUser.Email == req.Email {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email already exists"})
		}
		if existingUser.PhoneNumber == req.PhoneNumber {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Phone number already exists"})
		}
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	u := &model.User{
		Email:       req.Email,
		Password:    genPassword(req.Password),
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user"})
	}

	userRole := model.UserRole{
		UserID: u.ID,
		RoleID: 3,
	}
	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to assign role"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Transaction commit failed"})
	}

	res := &dto.SignUpResponse{
		Email:       u.Email,
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
	}

	return c.JSON(http.StatusCreated, echo.Map{"user": res})
}

func (h *Handler) Login(c echo.Context) (err error) {
	req := new(dto.LoginRequest)
	if err = c.Bind(req); err != nil {
		return
	}

	u := new(model.User)

	db := h.DB
	if err = db.Where("email = ?", req.Email).First(u).Error; err != nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: "user not found",
		}
	}

	if !matchPassword(u.Password, req.Password) {

	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenStr, err := token.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		return err
	}

	u.Password = ""
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": tokenStr,
	})
}

func genPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func matchPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
