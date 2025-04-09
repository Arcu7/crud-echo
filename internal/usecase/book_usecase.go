package usecase

import (
	"crud-echo/internal/models"
	"fmt"
)

type usecaseBooksRepository interface {
	Create(book *models.Books) error
	GetByID(book *models.Books, id int) error
	GetAll(book *[]models.Books) error
	Update(book *models.Books) error
	Delete(book *models.Books) error
	ExistsByTitle(title string) (bool, error)
}

type BooksUseCase struct {
	bookRepo usecaseBooksRepository
}

func NewBooksUseCase(repo usecaseBooksRepository) *BooksUseCase {
	return &BooksUseCase{bookRepo: repo}
}

func (uc *BooksUseCase) CreateBook(request *models.CreateBooksRequest) (*models.Books, error) {
	exists, err := uc.bookRepo.ExistsByTitle(request.Title)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("repository error: %w", models.ErrResourceExistAlready)
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

func (uc *BooksUseCase) GetBookByID(id int) (*models.BooksSummary, error) {
	var book models.Books
	if err := uc.bookRepo.GetByID(&book, id); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return book.ToBooksSummary(), nil
}

func (uc *BooksUseCase) GetAllBooks(available bool) (*[]models.BooksSummary, error) {
	var books []models.Books
	var booksList []models.BooksSummary
	if !available {
		return nil, models.ErrInvalidParam
	}

	if err := uc.bookRepo.GetAll(&books); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	for _, book := range books {
		booksList = append(booksList, *book.ToBooksSummary())
	}
	return &booksList, nil
}

func (uc *BooksUseCase) UpdateBook(book *models.UpdateBooksRequest) error {
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
	bookData := &models.Books{
		ID: book.ID,
	}

	if err := uc.bookRepo.Delete(bookData); err != nil {
		return fmt.Errorf("repository error: %w", err)
	}

	return nil
}
