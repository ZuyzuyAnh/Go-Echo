package controller

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/service"
	"echo-demo/internal/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type TheaterController struct {
	TheaterService *service.TheaterService
	log            *zap.Logger
}

func NewTheaterController(theaterService *service.TheaterService, logger *zap.Logger) *TheaterController {
	return &TheaterController{
		TheaterService: theaterService,
		log:            logger,
	}
}

func (t *TheaterController) CreateTheater(c echo.Context) error {
	req := new(dto.CreateUpdateTheaterRequest)
	if err := c.Bind(req); err != nil {
		return sendErrorResponse(c, err)
	}

	v := validator.New()
	if req.Validate(v); !v.Valid() {
		return sendValidationErrors(c, v)
	}

	res, err := t.TheaterService.CreateTheater(req)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (t *TheaterController) GetTheaterByID(c echo.Context) error {
	id, err := readIntID(c)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	res, err := t.TheaterService.GetTheaterByID(id)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (t *TheaterController) ListTheaters(c echo.Context) error {
	res, err := t.TheaterService.ListTheaters()
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (t *TheaterController) UpdateTheater(c echo.Context) error {
	id, err := readIntID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req := new(dto.CreateUpdateTheaterRequest)
	if err := c.Bind(req); err != nil {
		return sendErrorResponse(c, err)
	}

	v := validator.New()
	if req.Validate(v); !v.Valid() {
		return sendValidationErrors(c, v)
	}

	err = t.TheaterService.UpdateTheater(req, id)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, id)
}

func (t *TheaterController) DeleteTheater(c echo.Context) error {
	id, err := readIntID(c)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	err = t.TheaterService.DeleteTheater(id)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (t *TheaterController) UpdateSeatType(c echo.Context) error {
	req := new(dto.UpdateSeatTypeRequest)
	if err := c.Bind(req); err != nil {
		return sendErrorResponse(c, err)
	}

	v := validator.New()
	if req.Validate(v); !v.Valid() {
		return sendValidationErrors(c, v)
	}

	err := t.TheaterService.UpdateSeatType(req)
	if err != nil {
		return sendErrorResponse(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
