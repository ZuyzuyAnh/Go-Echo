package controller

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/repository"
	"echo-demo/internal/service"
	"errors"
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
		return echo.NewHTTPError(http.StatusBadRequest, MSG_INVALIDS_REQUEST_BODY)
	}

	res, err := m.MovieService.CreateMovie(req)
	if err != nil {
		var duplicateErr *repository.DuplicateError
		if errors.As(err, &duplicateErr) {
			return echo.NewHTTPError(http.StatusConflict, echo.Map{"error": duplicateErr.Error()})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": MSG_INTERNAL_SERVER_ERROR})
	}

	return c.JSON(http.StatusOK, res)
}
