package handlers

import (
	"crud-echo/internal/models"
	uc "crud-echo/internal/usecase"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BooksHandler struct {
	BUC *uc.BooksUseCase
}

func NewBooksHandler(buc *uc.BooksUseCase) *BooksHandler {
	return &BooksHandler{BUC: buc}
}

func (h BooksHandler) CreateBook(c echo.Context) error {
	var b models.CreateBooksRequest

	if err := c.Bind(&b); err != nil {
		return CustomResponse(c, http.StatusBadRequest, false, models.BadRequest)
	}

	if _, err := h.BUC.CreateBook(&b); err != nil {
		var ve *models.ValidationError
		if errors.As(err, &ve) {
			return ValidationErrorResponse(c, ve.Message, ve.Errors)
		} else if errors.Is(err, models.ErrResourceExistAlready) {
			return CustomResponse(c, http.StatusConflict, false, err.Error())
		} else {
			return CustomResponse(c, http.StatusInternalServerError, false, err.Error())
		}
	}

	resp := CustomResponse(c, http.StatusOK, true, "Book has been created")
	return resp
}

func (h BooksHandler) GetBookByID(c echo.Context) error {
	var b models.Books

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return CustomResponse(c, http.StatusBadRequest, false, models.BadRequest)
	}

	resp, err := h.BUC.GetBookByID(&b, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return CustomResponse(c, http.StatusNotFound, false, err.Error())
		} else {
			return CustomResponse(c, http.StatusInternalServerError, false, err.Error())
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) GetAllBooks(c echo.Context) error {
	var b models.BooksList

	available, err := strconv.ParseBool(c.QueryParam("available"))
	if err != nil {
		return CustomResponse(c, http.StatusBadRequest, false, models.BadRequest)
	}

	resp, err := h.BUC.GetAllBooks(&b, available)
	if err != nil {
		if errors.Is(err, models.ErrTableEmpty) {
			return CustomResponse(c, http.StatusNotFound, false, err.Error())
		} else if errors.Is(err, models.ErrInvalidParam) {
			return CustomResponse(c, http.StatusBadRequest, false, err.Error())
		} else {
			return CustomResponse(c, http.StatusInternalServerError, false, err.Error())
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) UpdateBook(c echo.Context) error {
	var b models.UpdateBooksRequest

	if err := c.Bind(&b); err != nil {
		return CustomResponse(c, http.StatusBadRequest, false, models.BadRequest)
	}

	if err := h.BUC.UpdateBook(&b); err != nil {
		var ve *models.ValidationError
		if errors.As(err, &ve) {
			return ValidationErrorResponse(c, ve.Message, ve.Errors)
		} else if errors.Is(err, models.ErrNotFound) {
			return CustomResponse(c, http.StatusNotFound, false, err.Error())
		} else {
			return CustomResponse(c, http.StatusInternalServerError, false, err.Error())
		}
	}

	resp := CustomResponse(c, http.StatusOK, true, "Book with ID "+strconv.Itoa(b.ID)+" has been updated")
	return resp
}

func (h BooksHandler) DeleteBook(c echo.Context) error {
	var b models.DeleteBooksRequest

	if err := c.Bind(&b); err != nil {
		return CustomResponse(c, http.StatusBadRequest, false, models.BadRequest)
	}

	if err := h.BUC.DeleteBook(&b); err != nil {
		var ve *models.ValidationError
		if errors.As(err, &ve) {
			return ValidationErrorResponse(c, ve.Message, ve.Errors)
		} else if errors.Is(err, models.ErrNotFound) {
			return CustomResponse(c, http.StatusNotFound, false, err.Error())
		} else {
			return CustomResponse(c, http.StatusInternalServerError, false, err.Error())
		}
	}

	resp := CustomResponse(c, http.StatusOK, true, "Book with ID "+strconv.Itoa(b.ID)+" has been deleted")
	return resp
}
