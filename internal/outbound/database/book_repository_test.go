package database

import (
	"crud-echo/internal/models"
	psgr "crud-echo/pkg/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*psgr.PostgresDB, sqlmock.Sqlmock, func()) {
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

	psqlgdb := &psgr.PostgresDB{
		DB: gdb,
	}

	cleanup := func() {
		db.Close()
	}

	return psqlgdb, mock, cleanup
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		bookRequest *models.Books
		expectedID  int
		mock        func(mock sqlmock.Sqlmock)
		wantErr     bool
		errType     error
	}{
		{
			name: "Success create book with valid data",
			bookRequest: &models.Books{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			expectedID: 1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "books" (.+) VALUES (.+)`).
					WithArgs("Test Title", "Test Description", 10, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Database error during create",
			bookRequest: &models.Books{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "books" (.+) VALUES (.+)`).
					WithArgs("Test Title", "Test Description", 10, sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			err := repo.Create(tt.bookRequest)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, tt.bookRequest.ID)
				assert.NotEmpty(t, tt.bookRequest.CreatedAt)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name         string
		bookRequest  *models.Books
		expectedBook *models.Books
		id           int
		mock         func(mock sqlmock.Sqlmock)
		wantErr      bool
		errType      error
	}{
		{
			name:        "Success get book by ID",
			bookRequest: &models.Books{},
			expectedBook: &models.Books{
				ID:          1,
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			id: 1,
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
			name:        "Book not found with non-existent ID",
			bookRequest: &models.Books{},
			id:          999,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "books" WHERE "books"."id" = (.+) ORDER BY "books"."id" LIMIT (.+)`).
					WithArgs(999, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: models.ErrNotFound,
		},
		{
			name:        "Database error during get",
			bookRequest: &models.Books{},
			id:          1,
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

			err := repo.GetByID(tt.bookRequest, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBook.ID, tt.bookRequest.ID)
				assert.Equal(t, tt.expectedBook.Title, tt.bookRequest.Title)
				assert.Equal(t, tt.expectedBook.Description, tt.bookRequest.Description)
				assert.Equal(t, tt.expectedBook.Qty, tt.bookRequest.Qty)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name          string
		books         *[]models.Books
		expectedBooks *[]models.Books
		mock          func(mock sqlmock.Sqlmock)
		wantErr       bool
		errType       error
	}{
		{
			name:  "Success get all books",
			books: &[]models.Books{},
			expectedBooks: &[]models.Books{
				{
					ID:          1,
					Title:       "Test Title 1",
					Description: "Test Description 1",
					Qty:         10,
				},
				{
					ID:          2,
					Title:       "Test Title 2",
					Description: "Test Description 2",
					Qty:         20,
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "qty", "created_at", "updated_at"}).
					AddRow(1, "Test Title 1", "Test Description 1", 10, time.Now(), time.Now()).
					AddRow(2, "Test Title 2", "Test Description 2", 20, time.Now(), time.Now())
				mock.ExpectQuery(`SELECT \* FROM "books"`).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:  "No books found in database",
			books: &[]models.Books{},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "books"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "qty", "created_at", "updated_at"}))
			},
			wantErr: true,
			errType: models.ErrEmptyTable,
		},
		{
			name:  "Database error during get all",
			books: &[]models.Books{},
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

			err := repo.GetAll(tt.books)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(*tt.expectedBooks), len(*tt.books))
				for i, expected := range *tt.expectedBooks {
					actual := (*tt.books)[i]
					assert.Equal(t, expected.ID, actual.ID)
					assert.Equal(t, expected.Title, actual.Title)
					assert.Equal(t, expected.Description, actual.Description)
					assert.Equal(t, expected.Qty, actual.Qty)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name string
		book *models.Books
		// expectedRowsAffected int64
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name: "Success update book with valid data",
			book: &models.Books{
				ID:          1,
				Title:       "Updated Title",
				Description: "Updated Description",
				Qty:         15,
			},
			// expectedRowsAffected: 1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"description"=\$2,"qty"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
					WithArgs("Updated Title", "Updated Description", 15, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Book not found during update",
			book: &models.Books{
				ID:          99,
				Title:       "Updated Title",
				Description: "Updated Description",
				Qty:         15,
			},
			// expectedRowsAffected: 1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"description"=\$2,"qty"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
					WithArgs("Updated Title", "Updated Description", 15, sqlmock.AnyArg(), 99).
					WillReturnResult(sqlmock.NewResult(99, 0))
				mock.ExpectCommit()
			},
			wantErr: true,
			errType: models.ErrNotFound,
		},
		{
			name: "Database error during update",
			book: &models.Books{
				ID:          1,
				Title:       "Updated Title",
				Description: "Updated Description",
				Qty:         15,
			},
			// expectedRowsAffected: 1,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"description"=\$2,"qty"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
					WithArgs("Updated Title", "Updated Description", 15, sqlmock.AnyArg(), 1).
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
				// assert.Equal(t, tt.expectedRowsAffected, int64(1))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		book    *models.Books
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name: "Success delete existing book",
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
			name: "Book not found during delete",
			book: &models.Books{
				ID: 99,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "books" WHERE "books"."id" = \$1`).
					WithArgs(99).
					WillReturnResult(sqlmock.NewResult(99, 0))
				mock.ExpectCommit()
			},
			wantErr: true,
			errType: models.ErrNotFound,
		},
		{
			name: "Database error during delete",
			book: &models.Books{
				ID: 1,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "books" WHERE "books"."id" = \$1`).
					WithArgs(1).
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

func TestExistsByTitle(t *testing.T) {
	tests := []struct {
		name       string
		title      string
		mock       func(mock sqlmock.Sqlmock)
		wantExists bool
		wantErr    bool
		errType    error
	}{
		{
			name:       "Book exists with given title",
			title:      "Test Title",
			wantExists: true,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "books" WHERE title = \$1`).
					WithArgs("Test Title").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:       "Book does not exist with given title",
			title:      "Non-existent Title",
			wantExists: false,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "books" WHERE title = \$1`).
					WithArgs("Non-existent Title").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:       "Database error during check",
			title:      "Test Title",
			wantExists: false,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "books" WHERE title = \$1`).
					WithArgs("Test Title").
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

			exists, err := repo.ExistsByTitle(tt.title)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantExists, exists)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
