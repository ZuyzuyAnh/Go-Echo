package controller

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/service"
	"echo-demo/internal/validator"
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
		return sendErrorResponse(c, err)
	}

	v := validator.New()

	if req.Validate(v); !v.Valid() {
		return sendValidationErrors(c, v)
	}

	res, err := uc.UserService.Register(req)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *UserController) Login(c echo.Context) error {
	req := new(dto.LoginRequest)
	if err := c.Bind(req); err != nil {
		return sendErrorResponse(c, err)
	}

	v := validator.New()
	if req.Validate(v); !v.Valid() {
		return sendValidationErrors(c, v)
	}

	res, err := uc.UserService.Login(req, uc.Secret)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (uc *UserController) GetProfile(c echo.Context) error {
	id, err := readUserID(c)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	res, err := uc.UserService.Profile(id)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, res)
}
