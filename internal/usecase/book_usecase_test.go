package usecase

import (
	"crud-echo/internal/models"
	"crud-echo/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/validator/v10"
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
		book    *models.CreateBooksRequest
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Success create book",
			book: &models.CreateBooksRequest{
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
		// idk if this is right or not
		{
			name: "Invalid Validation (Title required)",
			book: &models.CreateBooksRequest{
				Description: "Test Description",
				Qty:         10,
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Invalid Validation (Title minimum)",
			book: &models.CreateBooksRequest{
				Title:       gofakeit.LetterN(2),
				Description: "Test Description",
				Qty:         10,
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Invalid Validation (Title maximum)",
			book: &models.CreateBooksRequest{
				Title:       gofakeit.LetterN(101),
				Description: "Test Description",
				Qty:         10,
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Invalid Validation (Description required)",
			book: &models.CreateBooksRequest{
				Title: "Test Title",
				Qty:   10,
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Invalid Validation (Description minimum)",
			book: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: gofakeit.LetterN(2),
				Qty:         10,
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Invalid Validation (Description maximum)",
			book: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: gofakeit.LetterN(1001),
				Qty:         10,
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Invalid Validation (Qty less than)",
			book: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         gofakeit.Number(-255, -1),
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Invalid Validation (Qty greater than)",
			book: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         gofakeit.Number(101, 255),
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Repository error",
			book: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         3,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "books" (.+) VALUES (.+)`).
					WithArgs("Test Title", "Test Description", 3, sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			repo := repository.NewBooksRepository(gdb)
			validator := NewCustomValidator(validator.New())
			uc := NewBooksUseCase(repo, validator)

			got, err := uc.CreateBook(tt.book)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.book.Title, got.Title)
				assert.Equal(t, tt.book.Description, got.Description)
				assert.Equal(t, tt.book.Qty, got.Qty)
			}
		})
	}
}

func TestGetBook(t *testing.T) {
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
					AddRow(1, "Test Title", "Test Description", 10, time.Now(), time.Now())
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

			repo := repository.NewBooksRepository(gdb)
			validator := NewCustomValidator(validator.New())
			uc := NewBooksUseCase(repo, validator)

			got, err := uc.GetBook(tt.book, tt.id)

			if tt.wantErr {
				t.Log("testttt")
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got.ID, tt.book.ID)
				assert.Equal(t, got.Title, tt.book.Title)
				assert.Equal(t, got.Description, tt.book.Description)
				assert.Equal(t, got.Qty, tt.book.Qty)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	tests := []struct {
		name    string
		books   *models.BooksList
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name:  "Success get all books",
			books: &models.BooksList{},
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "qty", "created_at", "updated_at"}).
					AddRow(1, "Test Title", "Test Description", 10, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT \* FROM "books"`).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:  "No books found",
			books: &models.BooksList{},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "books"`).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
		{
			name:  "Database error",
			books: &models.BooksList{},
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

			repo := repository.NewBooksRepository(gdb)
			validator := NewCustomValidator(validator.New())
			uc := NewBooksUseCase(repo, validator)

			got, err := uc.GetAllBooks(tt.books)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)

				books := *tt.books
				got_books := *got
				assert.Equal(t, 1, len(books))
				assert.Equal(t, got_books[0].ID, books[0].ID)
				assert.Equal(t, got_books[0].Title, books[0].Title)
				assert.Equal(t, got_books[0].Description, books[0].Description)
				assert.Equal(t, got_books[0].Qty, books[0].Qty)
			}
		})
	}
}

// func TestUpdateBook(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		book    models.UpdateBooksRequest
// 		mock    func(mock sqlmock.Sqlmock)
// 		wantErr bool
// 	}{
// 		{
// 			name: "Success update book",
// 			book: models.UpdateBooksRequest{
// 				ID:          1,
// 				Title:       "Updated Book",
// 				Description: "Updated Description",
// 				Qty:         20,
// 			},
// 			mock: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec("UPDATE books").
// 					WithArgs("Updated Book", "Updated Description", 20, 1).
// 					WillReturnResult(sqlmock.NewResult(0, 1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Invalid Validation (Title required)",
// 			book: models.UpdateBooksRequest{
// 				ID:          1,
// 				Description: "Updated Description",
// 				Qty:         20,
// 			},
// 			mock:    func(mock sqlmock.Sqlmock) {},
// 			wantErr: true,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			db, mock, cleanup := setupTestDB(t)
// 			defer cleanup()
//
// 			tt.mock(mock)
//
// 			repo := repository.NewBooksRepository(db)
// 			validator := NewCustomValidator(validator.New())
// 			uc := NewBooksUseCase(repo, validator)
//
// 			err := uc.UpdateBook(tt.book)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}
//
// 			assert.NoError(t, err)
// 		})
// 	}
// }
//
// func TestDeleteBook(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		id      int
// 		mock    func(mock sqlmock.Sqlmock)
// 		wantErr bool
// 	}{
// 		{
// 			name: "Success delete book",
// 			id:   1,
// 			mock: func(mock sqlmock.Sqlmock) {
// 				go mock.ExpectExec("DELETE FROM books").
// 					WithArgs(1).
// 					WillReturnResult(sqlmock.NewResult(0, 1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Book not found",
// 			id:   999,
// 			mock: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec("DELETE FROM books").
// 					WithArgs(999).
// 					WillReturnError(gorm.ErrRecordNotFound)
// 			},
// 			wantErr: true,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			db, mock, cleanup := setupTestDB(t)
// 			defer cleanup()
//
// 			tt.mock(mock)
//
// 			repo := repository.NewBooksRepository(db)
// 			validator := NewCustomValidator(validator.New())
// 			uc := NewBooksUseCase(repo, validator)
//
// 			book := &models.Books{ID: tt.id}
// 			err := uc.DeleteBook(book)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}
//
// 			assert.NoError(t, err)
// 		})
// 	}
// }
