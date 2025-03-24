package handlers

import (
	"crud-echo/helper"
	"crud-echo/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type BooksHandler struct {
	DB *gorm.DB
}

func (h BooksHandler) CreateBook(c echo.Context) error {
	b := new(models.Books)
	if err := c.Bind(b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(b); err != nil {
		return err
	}

	if err := h.DB.Create(&b).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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

	if err := h.DB.First(&b, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	resp := b.ToBooksResponse()
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) GetAllBooks(c echo.Context) error {
	var b models.BooksList

	// idk what's this for
	available := c.QueryParam("available")
	if available != "true" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameter (only true or false)")
	}

	if err := h.DB.Find(&b).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	resp := b.ToBooksResponse()
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) UpdateBook(c echo.Context) error {
	b := new(models.Books)
	if err := c.Bind(b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(b); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	b.ID = id

	if err := h.DB.Save(&b).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// old
	// return c.String(http.StatusOK, "Books ID "+c.Param("id")+" has been updated\n")
	resp := helper.CustomResponse(true, "Books has been updated")
	return c.JSON(http.StatusOK, resp)
}

func (h BooksHandler) DeleteBook(c echo.Context) error {
	var b models.Books

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := h.DB.Delete(&b, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// ol
	// return c.String(http.StatusOK, "Books ID "+c.Param("id")+" has been deleted\n")
	resp := helper.CustomResponse(true, "Books has been deleted")
	return c.JSON(http.StatusOK, resp)
}
