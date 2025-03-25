package handlers

import (
	"crud-echo/helper"
	"crud-echo/internal/models"
	uc "crud-echo/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BooksHandler struct {
	BUC *uc.BooksUseCase
}

func (h BooksHandler) CreateBook(c echo.Context) error {
	var b models.CreateBooksRequest

	if err := c.Bind(&b); err != nil {
		return c.JSON(http.StatusBadRequest, helper.CustomResponse(false, "Invalid request"))
	}

	if _, err := h.BUC.CreateBook(b); err != nil {
		// return echo.NewHTTPError(http.StatusInternalServerError, err)
		// return c.JSON(http.StatusInternalServerError, helper.CustomResponse(false, err.Error()))
		return err
	}

	resp := helper.CustomResponse(true, "Books has been created")
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) GetBook(c echo.Context) error {
	var b models.Books

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err = h.BUC.GetBook(&b, id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	resp := b.ToBooksResponse()
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) GetAllBooks(c echo.Context) error {
	var b models.BooksList

	// idk what's this for
	available, err := strconv.ParseBool(c.QueryParam("available"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameter (only true or false)")
	}

	if !available {
		return echo.NewHTTPError(http.StatusBadRequest, "Available is not true")
	}

	if err := h.BUC.GetAllBooks(&b); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	resp := b.ToBooksResponse()
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) UpdateBook(c echo.Context) error {
	var b models.UpdateBooksRequest

	if err := c.Bind(&b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.BUC.UpdateBook(b); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	resp := helper.CustomResponse(true, "Books has been updated")
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) DeleteBook(c echo.Context) error {
	var b models.Books

	if err := c.Bind(&b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.BUC.DeleteBook(&b); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	resp := helper.CustomResponse(true, "Books has been deleted")
	return c.JSON(http.StatusOK, resp)
}
