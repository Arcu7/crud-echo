package database

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
	result := r.DB.Create(&book)

	if result.Error != nil {
		return result.Error
	} else if book.ID == 0 {
		return models.ErrInternalServerError
	}

	return nil
}

func (r *BooksRepository) GetByID(book *models.Books, id int) error {
	result := r.DB.First(&book, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.ErrNotFound
		}
		return result.Error
	}

	return nil
}

func (r *BooksRepository) GetAll(books *[]models.Books) error {
	result := r.DB.Find(&books)

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected < 1 { // maybe only need to do one check
		return models.ErrTableEmpty
	}

	return nil
}

func (r *BooksRepository) Update(book *models.Books) error {
	result := r.DB.Model(&book).Updates(models.Books{
		Title:       book.Title,
		Description: book.Description,
		Qty:         book.Qty,
	})

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected < 1 {
		return models.ErrNotFound
	}

	return nil
}

func (r *BooksRepository) Delete(book *models.Books) error {
	result := r.DB.Delete(&book)

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected < 1 {
		return models.ErrNotFound
	}

	return nil
}

func (r *BooksRepository) ExistsByTitle(title string) (bool, error) {
	var count int64
	result := r.DB.Model(&models.Books{}).Where("title = ?", title).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
