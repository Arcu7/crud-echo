package usecase

import (
	"crud-echo/internal/models"
	"fmt"
)

type usecaseBooksRepository interface {
	Create(book *models.Books) error
	GetByID(book *models.Books, id int) error
	GetAll(users *models.BooksList) error
	Update(book *models.Books) error
	Delete(book *models.Books) error
	ExistsByTitle(title string) (bool, error)
}

type BooksUseCase struct {
	bookRepo  usecaseBooksRepository
	Validator *CustomValidator
}

func NewBooksUseCase(repo usecaseBooksRepository, validator *CustomValidator) *BooksUseCase {
	return &BooksUseCase{bookRepo: repo, Validator: validator}
}

func (uc *BooksUseCase) CreateBook(request *models.CreateBooksRequest) (*models.Books, error) {
	if err := uc.Validator.Validate(request); err != nil {
		return nil, fmt.Errorf("validaiton error: %w", err)
	}

	exists, err := uc.bookRepo.ExistsByTitle(request.Title)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if exists {
		return nil, models.ErrResourceExistAlready
	}

	bookData := &models.Books{
		Title:       request.Title,
		Description: request.Description,
		Qty:         request.Qty,
	}

	if err := uc.bookRepo.Create(bookData); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	return bookData, nil
}

func (uc *BooksUseCase) GetBookByID(book *models.Books, id int) (*models.BooksSummary, error) {
	if err := uc.bookRepo.GetByID(book, id); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return book.ToBooksSummary(), nil
}

func (uc *BooksUseCase) GetAllBooks(book *models.BooksList, available bool) (*[]models.BooksSummary, error) {
	if !available {
		return nil, models.ErrInvalidParam
	}

	if err := uc.bookRepo.GetAll(book); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return book.ToBooksSummary(), nil
}

func (uc *BooksUseCase) UpdateBook(book *models.UpdateBooksRequest) error {
	if err := uc.Validator.Validate(book); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	bookData := &models.Books{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Qty:         book.Qty,
	}

	if err := uc.bookRepo.Update(bookData); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	return nil
}

func (uc *BooksUseCase) DeleteBook(book *models.DeleteBooksRequest) error {
	if err := uc.Validator.Validate(book); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	bookData := &models.Books{
		ID: book.ID,
	}

	if err := uc.bookRepo.GetByID(bookData, bookData.ID); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	if err := uc.bookRepo.Delete(bookData); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	return nil
}
