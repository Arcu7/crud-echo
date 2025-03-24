package usecase

import (
	"crud-echo/internal/models"
	"crud-echo/internal/repository"
	vc "crud-echo/internal/usecase/validators_custom"
)

type BookUseCase struct {
	Repo      *repository.BookRepository
	Validator *vc.CustomValidator
}

func NewBookUseCase(repo *repository.BookRepository, validator *vc.CustomValidator) *BookUseCase {
	return &BookUseCase{Repo: repo, Validator: validator}
}

func (uc *BookUseCase) CreateBook(book *models.Books) error {
	if err := uc.Validator.Validate(book); err != nil {
		return err
	}

	return uc.Repo.Create(book)
}

func (uc *BookUseCase) GetBook(book *models.Books, id int) error {
	return uc.Repo.GetByID(book, id)
}

func (uc *BookUseCase) GetAllBooks(book *models.BooksList) error {
	return uc.Repo.GetAll(book)
}

func (uc *BookUseCase) UpdateBook(book *models.Books) error {
	if err := uc.Validator.Validate(book); err != nil {
		return err
	}
	return uc.Repo.Update(book)
}

func (uc *BookUseCase) DeleteBook(book *models.Books) error {
	return uc.Repo.Delete(book)
}
