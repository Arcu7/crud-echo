package repository

import (
	"crud-echo/internal/models"

	"gorm.io/gorm"
)

type BooksRepository struct {
	DB *gorm.DB
}

func NewBooksRepository(db *gorm.DB) *BooksRepository {
	return &BooksRepository{DB: db}
}

func (r *BooksRepository) Create(book *models.Books) error {
	err := r.DB.Create(&book).Error
	return err
}

func (r *BooksRepository) GetByID(book *models.Books, id int) error {
	err := r.DB.First(&book, id).Error
	return err
}

func (r *BooksRepository) GetAll(books *models.BooksList) error {
	err := r.DB.Find(&books).Error
	return err
}

func (r *BooksRepository) Update(book *models.Books) error {
	result := r.DB.Model(&book).Updates(models.Books{
		Title:       book.Title,
		Description: book.Description,
		Qty:         book.Qty,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BooksRepository) Delete(book *models.Books) error {
	result := r.DB.Delete(&book).Error
	return result
}
