package controller

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
