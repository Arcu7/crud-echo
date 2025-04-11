package handlers

import (
	"crud-echo/internal/inbound/customvalidator"
	"crud-echo/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handlerBookUsecase interface {
	CreateBook(book *models.CreateBooksRequest) (*models.Books, error)
	GetBookByID(id int) (*models.BooksSummary, error)
	GetAllBooks(available bool) (*[]models.BooksSummary, error)
	UpdateBook(book *models.UpdateBooksRequest) error
	DeleteBook(book *models.DeleteBooksRequest) error
}

type BooksHandler struct {
	buc handlerBookUsecase
	cv  *customvalidator.CustomValidator
}

func NewBooksHandler(buc handlerBookUsecase, validator *customvalidator.CustomValidator) *BooksHandler {
	return &BooksHandler{buc: buc, cv: validator}
}

func (h BooksHandler) CreateBook(c echo.Context) error {
	var b models.CreateBooksRequest

	if err := c.Bind(&b); err != nil {
		log.Printf("Error binding request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	if err := h.cv.Validate(b); err != nil {
		log.Printf("Error validating request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrValidationError.Error())
	}

	_, err := h.buc.CreateBook(&b)
	if err != nil {
		log.Printf("Error creating book: %v", err)
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	return CustomResponse(c, http.StatusOK, true, "Book has been created", nil)
}

func (h BooksHandler) GetBookByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Error converting id to integer: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	resp, err := h.buc.GetBookByID(id)
	if err != nil {
		log.Printf("Error retrieving book with ID %d: %v", id, err)
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	// return c.JSON(http.StatusOK, resp)
	return CustomResponse(c, http.StatusOK, true, "Book retrieved successfully", resp)
}

func (h BooksHandler) GetAllBooks(c echo.Context) error {
	available, err := strconv.ParseBool(c.QueryParam("available"))
	if err != nil {
		log.Printf("Error parsing available param: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	resp, err := h.buc.GetAllBooks(available)
	if err != nil {
		log.Printf("Error retrieving books: %v", err)
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	// return c.JSON(http.StatusOK, resp)
	return CustomResponse(c, http.StatusOK, true, "Books retrieved successfully", resp)
}

func (h BooksHandler) UpdateBook(c echo.Context) error {
	var b models.UpdateBooksRequest

	if err := c.Bind(&b); err != nil {
		log.Printf("Error binding request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	if err := h.cv.Validate(b); err != nil {
		log.Printf("Error validating request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrValidationError.Error())
	}

	if err := h.buc.UpdateBook(&b); err != nil {
		log.Printf("Error updating book: %v", err)
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	return CustomResponse(c, http.StatusOK, true, "Book with ID "+strconv.Itoa(b.ID)+" has been updated", nil)
}

func (h BooksHandler) DeleteBook(c echo.Context) error {
	var b models.DeleteBooksRequest

	if err := c.Bind(&b); err != nil {
		log.Printf("Error binding request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequest)
	}

	if err := h.cv.Validate(b); err != nil {
		log.Printf("Error validating request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrValidationError.Error())
	}

	if err := h.buc.DeleteBook(&b); err != nil {
		log.Printf("Error deleting book: %v", err)
		return echo.NewHTTPError(models.GetErrorHTTPStatusCode(err), models.GetErrorHTTPStatusMessage(err))
	}

	return CustomResponse(c, http.StatusOK, true, "Book with ID "+strconv.Itoa(b.ID)+" has been deleted", nil)
}
