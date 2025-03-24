package repository

import (
	"crud-echo/internal/models"

	"gorm.io/gorm"
)

type BookRepository struct {
	DB *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) Create(book *models.Books) error {
	return r.DB.Create(book).Error
}

func (r *BookRepository) GetByID(book *models.Books, id int) error {
	err := r.DB.First(&book, id).Error
	return err
}

func (r *BookRepository) GetAll(books *models.BooksList) error {
	err := r.DB.Find(&books).Error
	return err
}

func (r *BookRepository) Update(book *models.Books) error {
	return r.DB.Save(&book).Error
}

func (r *BookRepository) Delete(book *models.Books) error {
	return r.DB.Delete(&book).Error
}
