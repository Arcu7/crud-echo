package handlers

import (
	"crud-echo/internal/inbound/customvalidator"
	"crud-echo/internal/models"
	uc "crud-echo/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BooksHandler struct {
	buc *uc.BooksUseCase
	cv  *customvalidator.CustomValidator
}

func NewBooksHandler(buc *uc.BooksUseCase, validator *customvalidator.CustomValidator) *BooksHandler {
	return &BooksHandler{buc: buc, cv: validator}
}

func (h BooksHandler) CreateBook(c echo.Context) error {
	var b models.CreateBooksRequest

	if err := c.Bind(&b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	if err := h.cv.Validate(b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrValidationError.Error())
	}

	_, err := h.buc.CreateBook(&b)
	if err != nil {
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	return CustomResponse(c, http.StatusOK, true, "Book has been created", nil)
}

func (h BooksHandler) GetBookByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	resp, err := h.buc.GetBookByID(id)
	if err != nil {
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	// return c.JSON(http.StatusOK, resp)
	return CustomResponse(c, http.StatusOK, true, "Book retrieved successfully", resp)
}

func (h BooksHandler) GetAllBooks(c echo.Context) error {
	available, err := strconv.ParseBool(c.QueryParam("available"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	resp, err := h.buc.GetAllBooks(available)
	if err != nil {
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	// return c.JSON(http.StatusOK, resp)
	return CustomResponse(c, http.StatusOK, true, "Books retrieved successfully", resp)
}

func (h BooksHandler) UpdateBook(c echo.Context) error {
	var b models.UpdateBooksRequest

	if err := c.Bind(&b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	if err := h.cv.Validate(b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrValidationError.Error())
	}

	if err := h.buc.UpdateBook(&b); err != nil {
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	return CustomResponse(c, http.StatusOK, true, "Book with ID "+strconv.Itoa(b.ID)+" has been updated", nil)
}

func (h BooksHandler) DeleteBook(c echo.Context) error {
	var b models.DeleteBooksRequest

	if err := c.Bind(&b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	if err := h.cv.Validate(b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrValidationError.Error())
	}

	if err := h.buc.DeleteBook(&b); err != nil {
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	return CustomResponse(c, http.StatusOK, true, "Book with ID "+strconv.Itoa(b.ID)+" has been deleted", nil)
}
