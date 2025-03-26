package handlers

import (
	"crud-echo/internal/models"
	uc "crud-echo/internal/usecase"
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
		return c.JSON(http.StatusBadRequest, err)
	}

	if _, err := h.BUC.CreateBook(b); err != nil {
		return err
	}

	resp := CustomResponse(c, http.StatusOK, true, "Book has been created")
	return resp
}

func (h BooksHandler) GetBook(c echo.Context) error {
	var b models.Books

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err = h.BUC.GetBook(&b, id); err != nil {
		return err
	}

	resp := b.ToBooksResponse()
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) GetAllBooks(c echo.Context) error {
	var b models.BooksList

	// idk what's this for
	available, err := strconv.ParseBool(c.QueryParam("available"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid query parameter (only true or false)")
	}

	if !available {
		return c.JSON(http.StatusBadRequest, "Available is not true")
	}

	if err := h.BUC.GetAllBooks(&b); err != nil {
		return err
	}

	resp := b.ToBooksResponse()
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) UpdateBook(c echo.Context) error {
	var b models.UpdateBooksRequest

	if err := c.Bind(&b); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.BUC.UpdateBook(b); err != nil {
		return err
	}

	resp := CustomResponse(c, http.StatusOK, true, "Book with ID "+strconv.Itoa(b.ID)+" has been updated")
	return resp
}

func (h BooksHandler) DeleteBook(c echo.Context) error {
	var b models.Books

	if err := c.Bind(&b); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.BUC.DeleteBook(&b); err != nil {
		return err
	}

	resp := CustomResponse(c, http.StatusOK, true, "Book with ID "+strconv.Itoa(b.ID)+" has been deleted")
	return resp
}
