package models

import (
	"time"
)

type BooksRepository interface {
	Create(book *Books) error
	GetByID(book *Books, id int) error
	GetAll(users *BooksList) error
	Update(book *Books) error
	Delete(book *Books) error
}

// type Books struct {
// 	ID int `json:"id" gorm:"primaryKey;autoIncrement;not null" validate:"required_without_all=Title Description Qty"`
// 	Title       string       `json:"title" gorm:"type:varchar(50);not null" validate:"required,min=3,max=50"`
// 	Description string       `json:"description" gorm:"type:varchar(255);not null" validate:"required,min=3,max=255"`
// 	Qty         int          `json:"qty" gorm:"not null" validate:"required,gte=0,lte=100"`
// 	Created_at  time.Time    `json:"created_at" gorm:"autoCreateTime;type:timestamptz;not null"`
// 	Updated_at  sql.NullTime `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
// }

type Books struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Title       string    `json:"title" gorm:"type:varchar(50);not null"`
	Description string    `json:"description" gorm:"type:varchar(255);not null"`
	Qty         int       `json:"qty" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

type CreateBooksRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=3,max=255"`
	Qty         int    `json:"qty" validate:"required,gte=0,lte=100"`
}

type UpdateBooksRequest struct {
	ID          int    `json:"id" validate:"required"`
	Title       string `json:"title" validate:"min=3,max=50"`
	Description string `json:"description" validate:"min=3,max=255"`
	Qty         int    `json:"qty" validate:"gte=0,lte=100"`
}

type DeleteBooksRequest struct {
	ID int `json:"id" validate:"required"`
}

type BooksList []Books

type BooksSummary struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Qty         int    `json:"qty"`
}

func (b Books) ToBooksSummary() *BooksSummary {
	return &BooksSummary{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Qty:         b.Qty,
	}
}

func (bl BooksList) ToBooksSummary() *[]BooksSummary {
	var resp []BooksSummary
	for _, b := range bl {
		resp = append(resp, BooksSummary{
			ID:          b.ID,
			Title:       b.Title,
			Description: b.Description,
			Qty:         b.Qty,
		})
	}
	return &resp
}
