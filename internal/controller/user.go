package controller

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/repository"
	"echo-demo/internal/service"
	"echo-demo/internal/validator"
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	UserService *service.UserService
	Secret      string
	log         *zap.Logger
}

func NewUserController(userService *service.UserService, logger *zap.Logger, secret string) *UserController {
	return &UserController{UserService: userService, log: logger, Secret: secret}
}

func (uc *UserController) Register(c echo.Context) error {
	req := new(dto.SignupRequest)
	if err := c.Bind(req); err != nil {
		uc.log.Error("Bind error", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, MSG_INVALIDS_REQUEST_BODY)
	}

	v := validator.New()

	if req.Validate(v); !v.Valid() {
		return echo.NewHTTPError(http.StatusBadRequest, v.Errors)
	}

	res, err := uc.UserService.Register(req)
	if err != nil {
		uc.log.Error("Register error", zap.Error(err))

		var duplicateErr *repository.DuplicateError
		if errors.As(err, &duplicateErr) {
			return echo.NewHTTPError(http.StatusConflict, echo.Map{"error": duplicateErr.Error()})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": MSG_INTERNAL_SERVER_ERROR})
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *UserController) Login(c echo.Context) error {
	req := new(dto.LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MSG_INVALIDS_REQUEST_BODY)
	}

	v := validator.New()
	if req.Validate(v); !v.Valid() {
		return echo.NewHTTPError(http.StatusBadRequest, v.Errors)
	}

	res, err := uc.UserService.Login(req, uc.Secret)
	if err != nil {
		var notFoundErr *repository.NotFoundError
		if errors.As(err, &notFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, echo.Map{"error": notFoundErr.Error()})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": MSG_INTERNAL_SERVER_ERROR})
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *UserController) GetProfile(c echo.Context) error {
	id, err := readUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := uc.UserService.Profile(id)
	if err != nil {
		var notFoundErr *repository.NotFoundError
		if errors.As(err, &notFoundErr) {
			return echo.NewHTTPError(http.StatusNotFound, echo.Map{"error": notFoundErr.Error()})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": MSG_INTERNAL_SERVER_ERROR})
	}

	return c.JSON(http.StatusOK, res)
}
