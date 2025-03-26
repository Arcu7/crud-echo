package repository

import (
	"crud-echo/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return gdb, mock, cleanup
}

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name    string
		book    *models.Books
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Success create book",
			book: &models.Books{
				Title:       "Test Book",
				Description: "Test Description",
				Qty:         10,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "books" (.+) VALUES (.+)`).
					WithArgs("Test Book", "Test Description", 10, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Database error",
			book: &models.Books{
				Title:       "Test Book",
				Description: "Test Description",
				Qty:         10,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "books" (.+) VALUES (.+)`).
					WithArgs("Test Book", "Test Description", 10, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.mock(mock)

			repo := NewBooksRepository(gdb)

			err := repo.Create(tt.book)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrInvalidDB, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, tt.book.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetBookByID(t *testing.T) {
	tests := []struct {
		name    string
		book    *models.Books
		id      int
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name: "Success get book by ID",
			book: &models.Books{},
			id:   1,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "qty", "created_at", "updated_at"}).
					AddRow(1, "Test Book", "Test Description", 10, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT (.+) FROM "books" WHERE "books"."id" = (.+) ORDER BY "books"."id" LIMIT (.+)`).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "Book not found",
			book: &models.Books{},
			id:   999,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "books" WHERE "books"."id" = (.+) ORDER BY "books"."id" LIMIT (.+)`).
					WithArgs(999, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
		{
			name: "Database error",
			book: &models.Books{},
			id:   1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "books" WHERE "books"."id" = (.+) ORDER BY "books"."id" LIMIT (.+)`).
					WithArgs(1, 1).
					WillReturnError(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: gorm.ErrInvalidDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.mock(mock)

			repo := NewBooksRepository(gdb)

			err := repo.GetByID(tt.book, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, tt.book.ID)
				assert.Equal(t, "Test Book", tt.book.Title)
				assert.Equal(t, "Test Description", tt.book.Description)
				assert.Equal(t, 10, tt.book.Qty)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAllBook(t *testing.T) {
	tests := []struct {
		name    string
		book    *models.BooksList
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name: "Success get all books",
			book: &models.BooksList{},
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "qty", "created_at", "updated_at"}).
					AddRow(1, "Test Book", "Test Description", 10, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT \* FROM "books"`).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "No books found",
			book: &models.BooksList{},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "books"`).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
		{
			name: "Database error",
			book: &models.BooksList{},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "books"`).
					WillReturnError(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: gorm.ErrInvalidDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.mock(mock)

			repo := NewBooksRepository(gdb)

			err := repo.GetAll(tt.book)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)

				books := *tt.book
				assert.Equal(t, 1, len(books))
				assert.Equal(t, 1, books[0].ID)
				assert.Equal(t, "Test Book", books[0].Title)
				assert.Equal(t, "Test Description", books[0].Description)
				assert.Equal(t, 10, books[0].Qty)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	tests := []struct {
		name    string
		book    *models.Books
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name: "Success update book",
			book: &models.Books{
				ID:          1,
				Title:       "Test Update Book",
				Description: "Test Update Description",
				Qty:         99,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"description"=\$2,"qty"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
					WithArgs("Test Update Book", "Test Update Description", 99, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Book not found",
			book: &models.Books{
				ID:          999,
				Title:       "Test Update Book",
				Description: "Test Update Description",
				Qty:         99,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"description"=\$2,"qty"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
					WithArgs("Test Update Book", "Test Update Description", 99, sqlmock.AnyArg(), 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
		{
			name: "Database error",
			book: &models.Books{
				ID:          1,
				Title:       "Test Update Book",
				Description: "Test Update Description",
				Qty:         99,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"description"=\$2,"qty"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			wantErr: true,
			errType: gorm.ErrInvalidDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.mock(mock)

			repo := NewBooksRepository(gdb)

			err := repo.Update(tt.book)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	tests := []struct {
		name    string
		book    *models.Books
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name: "Success delete book",
			book: &models.Books{
				ID: 1,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "books" WHERE "books"."id" = \$1`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Book not found",
			book: &models.Books{
				ID: 999,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "books" WHERE "books"."id" = \$1`).
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
		{
			name: "Database error",
			book: &models.Books{
				ID: 1,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "books" WHERE "books"."id" = \$1`).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
			},
			wantErr: true,
			errType: gorm.ErrInvalidDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.mock(mock)

			repo := NewBooksRepository(gdb)

			err := repo.Delete(tt.book)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
