package usecase

import (
	"crud-echo/internal/models"
)

type BooksUseCase struct {
	bookRepo  models.BooksRepository
	Validator *CustomValidator
}

func NewBooksUseCase(repo models.BooksRepository, validator *CustomValidator) *BooksUseCase {
	return &BooksUseCase{bookRepo: repo, Validator: validator}
}

func (uc *BooksUseCase) CreateBook(request models.CreateBooksRequest) (*models.Books, error) {
	if err := uc.Validator.Validate(request); err != nil {
		return nil, err
	}

	bookData := &models.Books{
		Title:       request.Title,
		Description: request.Description,
		Qty:         request.Qty,
	}

	if err := uc.bookRepo.Create(bookData); err != nil {
		return nil, err
	}
	return bookData, nil
}

func (uc *BooksUseCase) GetBook(book *models.Books, id int) (*models.BooksSummary, error) {
	// return uc.bookRepo.GetByID(book, id)
	if err := uc.bookRepo.GetByID(book, id); err != nil {
		return nil, err
	}

	return book.ToBooksSummary(), nil
}

func (uc *BooksUseCase) GetAllBooks(book *models.BooksList) (*[]models.BooksSummary, error) {
	// return uc.bookRepo.GetAll(book)
	if err := uc.bookRepo.GetAll(book); err != nil {
		return nil, err
	}

	return book.ToBooksSummary(), nil
}

func (uc *BooksUseCase) UpdateBook(book models.UpdateBooksRequest) error {
	if err := uc.Validator.Validate(book); err != nil {
		return err
	}

	bookData := &models.Books{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Qty:         book.Qty,
	}

	if err := uc.bookRepo.Update(bookData); err != nil {
		return err
	}

	return nil
}

func (uc *BooksUseCase) DeleteBook(book *models.Books) error {
	if err := uc.Validator.Validate(book); err != nil {
		return err
	}

	bookData := &models.Books{
		ID: book.ID,
	}

	if err := uc.bookRepo.Delete(bookData); err != nil {
		return err
	}

	return nil
}
