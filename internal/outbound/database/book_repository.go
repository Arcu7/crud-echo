package database

import (
	"crud-echo/internal/models"

	"gorm.io/gorm"
)

type RepositoryDBConn interface {
	Migrate() error
	GetDB() *gorm.DB
}

type BooksRepository struct {
	rdc RepositoryDBConn
}

func NewBooksRepository(repoDBConn RepositoryDBConn) *BooksRepository {
	return &BooksRepository{rdc: repoDBConn}
}

func (r *BooksRepository) Create(book *models.Books) error {
	result := r.rdc.GetDB().Create(&book)

	if result.Error != nil {
		return result.Error
	} else if book.ID == 0 {
		return models.ErrInternalServerError
	}

	return nil
}

func (r *BooksRepository) GetByID(book *models.Books, id int) error {
	result := r.rdc.GetDB().First(&book, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.ErrNotFound
		}
		return result.Error
	}

	return nil
}

func (r *BooksRepository) GetAll(books *[]models.Books) error {
	result := r.rdc.GetDB().Find(&books)

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected < 1 { // maybe only need to do one check
		return models.ErrEmptyTable
	}

	return nil
}

func (r *BooksRepository) Update(book *models.Books) error {
	result := r.rdc.GetDB().Model(&book).Updates(models.Books{
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
	result := r.rdc.GetDB().Delete(&book)

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected < 1 {
		return models.ErrNotFound
	}

	return nil
}

func (r *BooksRepository) ExistsByTitle(title string) (bool, error) {
	var count int64
	result := r.rdc.GetDB().Model(&models.Books{}).Where("title = ?", title).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
