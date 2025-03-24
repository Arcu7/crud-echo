package models

import (
	"database/sql"
	"time"
)

type Books struct {
	ID          int          `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Title       string       `json:"title" gorm:"type:varchar(50);not null" validate:"required,min=3,max=50"`
	Description string       `json:"description" gorm:"type:varchar(255);not null" validate:"required,min=3,max=255"`
	Qty         int          `json:"qty" gorm:"not null" validate:"required,gte=0,lte=100"`
	Created_at  time.Time    `json:"created_at" gorm:"autoCreateTime;type:timestamptz;not null"`
	Updated_at  sql.NullTime `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

type BooksList []Books

type BooksResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Qty         int    `json:"qty"`
}

func (b Books) ToBooksResponse() BooksResponse {
	return BooksResponse{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Qty:         b.Qty,
	}
}

// is there a better way?
func (bl BooksList) ToBooksResponse() []BooksResponse {
	var resp []BooksResponse
	for _, b := range bl {
		resp = append(resp, BooksResponse{
			ID:          b.ID,
			Title:       b.Title,
			Description: b.Description,
			Qty:         b.Qty,
		})
	}
	return resp
}
