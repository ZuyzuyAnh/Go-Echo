package controller

import (
	"echo-demo/internal/repository"
	"echo-demo/internal/validator"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

func readIntID(c echo.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, errors.New("id must be integer")
	}

	return id, nil
}

func readUserID(c echo.Context) (int64, error) {
	user := c.Get("user")
	if user == nil {
		return 0, errors.New("user not found")
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	log.Error("claims: ", claims)

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("id must be integer")
	}

	return int64(id), nil
}

var (
	errBadRequest     = errors.New("bad request")
	errInternalServer = errors.New("internal server error")
	errIntParam       = errors.New("param must be integer")
)

func errorMessage(err error) echo.Map {
	return echo.Map{"message": err.Error()}
}

func sendErrorResponse(e echo.Context, err error) error {
	var notFoundError *repository.NotFoundError
	var duplicateError *repository.DuplicateError
	var foreignKeyError *repository.ForeignKeyError
	var echoError *echo.HTTPError

	switch {
	case errors.As(err, &echoError):
		return echo.NewHTTPError(http.StatusBadRequest, errorMessage(errBadRequest))
	case errors.As(err, &notFoundError):
		return echo.NewHTTPError(http.StatusNotFound, errorMessage(err))
	case errors.As(err, &duplicateError):
		return echo.NewHTTPError(http.StatusBadRequest, errorMessage(err))
	case errors.As(err, &foreignKeyError):
		return echo.NewHTTPError(http.StatusBadRequest, errorMessage(err))
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage(errInternalServer))
	}
}

func sendValidationErrors(e echo.Context, v *validator.Validator) error {
	return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
		"error": v.Errors,
	})
}
