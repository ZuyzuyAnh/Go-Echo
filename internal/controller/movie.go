package controller

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/service"
	"echo-demo/internal/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type MovieController struct {
	MovieService *service.MovieService
	log          *zap.Logger
}

func NewMovieController(movieService *service.MovieService, logger *zap.Logger) *MovieController {
	return &MovieController{
		MovieService: movieService,
		log:          logger,
	}
}

func (m *MovieController) CreateMovie(c echo.Context) error {
	req := new(dto.CreateMovieReq)
	if err := c.Bind(req); err != nil {
		return sendErrorResponse(c, err)
	}

	v := validator.New()
	if req.Validate(v); !v.Valid() {
		return sendValidationErrors(c, v)
	}

	res, err := m.MovieService.CreateMovie(req)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (m *MovieController) GetMovieByID(c echo.Context) error {
	id, er := readIntID(c)
	if er != nil {
		return sendErrorResponse(c, er)
	}

	res, err := m.MovieService.GetMovieByID(id)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (m *MovieController) UpdateMovie(c echo.Context) error {
	id, err := readIntID(c)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	req := new(dto.CreateMovieReq)
	if err := c.Bind(req); err != nil {
		return sendErrorResponse(c, err)
	}

	v := validator.New()
	if req.ValidateUpdate(v); !v.Valid() {
		return sendValidationErrors(c, v)
	}

	err = m.MovieService.UpdateMovie(req, id)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, id)
}

func (m *MovieController) DeleteMovie(c echo.Context) error {
	id, err := readIntID(c)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	err = m.MovieService.DeleteMovie(id)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
