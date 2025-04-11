package models

import (
	"time"
)

type Books struct {
	ID          int       `gorm:"primaryKey;autoIncrement;not null"`
	Title       string    `gorm:"type:varchar(50);not null"`
	Description string    `gorm:"type:varchar(255);not null"`
	Qty         int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime;type:timestamptz;not null"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;type:timestamptz"`
}

type CreateBooksRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=3,max=255"`
	Qty         int    `json:"qty" validate:"required,gte=0,lte=100"`
}

type UpdateBooksRequest struct {
	ID          int    `json:"id" validate:"required,gte=1"`
	Title       string `json:"title" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=3,max=255"`
	Qty         int    `json:"qty" validate:"required,gte=0,lte=100"`
}

type DeleteBooksRequest struct {
	ID int `json:"id" validate:"required,gte=1"`
}

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
