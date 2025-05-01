package v1

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

// book service

type LibroServiceServer interface {
	CreateLibro(echo.Context) error
	GetLibro(echo.Context) error
	UpdateLibro(echo.Context) error
	DeleteLibro(echo.Context) error
	
	CreateLibroReview(echo.Context) error
	GetLibroReview(echo.Context) error
	UpdateLibroReview(echo.Context) error
	DeleteLibroReview(echo.Context) error

	ReadLibroReview(echo.Context) error
	ReportLibroReview(echo.Context) error
}

func (s *APIV1Service) CreateLibro(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) GetLibro(c echo.Context) error {
	id := c.Param("id")
	slog.Debug("GetLibro: ", "id", id)
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) UpdateLibro(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) DeleteLibro(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) CreateLibroReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) GetLibroReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) UpdateLibroReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) DeleteLibroReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) ReadLibroReview(c echo.Context) error {
	year := c.QueryParam("year")
	slog.Debug("ReadReviewLibro: ", "year", year)
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) ReportLibroReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}
