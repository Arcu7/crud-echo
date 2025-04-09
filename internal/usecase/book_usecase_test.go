package usecase

import (
	"crud-echo/internal/mock"
	"crud-echo/internal/models"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var timeNow = time.Now

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name        string
		bookRequest *models.CreateBooksRequest
		expectedID  int
		mock        func(mock *mock.MockusecaseBooksRepository)
		wantErr     bool
		errType     error
	}{
		{
			name: "Success create book with valid data",
			bookRequest: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			expectedID: 1,
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().ExistsByTitle("Test Title").Return(false, nil)
				mock.EXPECT().Create(&models.Books{
					Title:       "Test Title",
					Description: "Test Description",
					Qty:         10,
				}).RunAndReturn(func(book *models.Books) error {
					book.ID = 1 // hackaround? maybe not the right approach
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "Failed create book due to title already exists",
			bookRequest: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().ExistsByTitle("Test Title").Return(true, nil)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", models.ErrResourceExistAlready),
		},
		{
			name: "Failed create book due to invalid DB",
			bookRequest: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().ExistsByTitle("Test Title").Return(false, nil)
				mock.EXPECT().Create(&models.Books{
					Title:       "Test Title",
					Description: "Test Description",
					Qty:         10,
				}).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrInvalidDB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mock.NewMockusecaseBooksRepository(t)
			tt.mock(mock)

			uc := NewBooksUseCase(mock)

			book, err := uc.CreateBook(tt.bookRequest)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, book.ID)
				assert.Equal(t, tt.bookRequest.Title, book.Title)
				assert.Equal(t, tt.bookRequest.Description, book.Description)
				assert.Equal(t, tt.bookRequest.Qty, book.Qty)
			}
		})
	}
}

func TestGetBookByID(t *testing.T) {
	tests := []struct {
		name         string
		bookRequest  *models.Books
		expectedBook *models.BooksSummary
		id           int
		mock         func(mock *mock.MockusecaseBooksRepository)
		wantErr      bool
		errType      error
	}{
		{
			name:        "Success get book with ID of 1",
			bookRequest: &models.Books{},
			expectedBook: &models.BooksSummary{
				ID:          1,
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			id: 1,
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().GetByID(&models.Books{}, 1).RunAndReturn(func(book *models.Books, id int) error {
					book.ID = id
					book.Title = "Test Title"
					book.Description = "Test Description"
					book.Qty = 10
					return nil
				})
			},
			wantErr: false,
		},
		{
			name:        "Failed get book with ID of 99",
			bookRequest: &models.Books{},
			expectedBook: &models.BooksSummary{
				ID:          99,
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			id: 99,
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().GetByID(&models.Books{}, 99).Return(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrRecordNotFound),
		},
		{
			name:        "Failed get book with ID due to invalid DB",
			bookRequest: &models.Books{},
			expectedBook: &models.BooksSummary{
				ID:          99,
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			id: 1,
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().GetByID(&models.Books{}, 1).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrInvalidDB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mock.NewMockusecaseBooksRepository(t)
			tt.mock(mock)

			uc := NewBooksUseCase(mock)

			book, err := uc.GetBookByID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBook.ID, book.ID)
				assert.Equal(t, tt.expectedBook.Title, book.Title)
				assert.Equal(t, tt.expectedBook.Description, book.Description)
				assert.Equal(t, tt.expectedBook.Qty, book.Qty)
			}
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	tests := []struct {
		name          string
		booksRequest  *models.BooksList
		expectedBooks *models.BooksList
		available     bool
		mock          func(mock *mock.MockusecaseBooksRepository)
		wantErr       bool
		errType       error
	}{
		{
			name:         "Success get all books",
			booksRequest: &models.BooksList{},
			expectedBooks: &models.BooksList{
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
			available: true,
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().GetAll(&models.BooksList{}).RunAndReturn(func(books *models.BooksList) error {
					*books = models.BooksList{
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
					}
					return nil
				})
			},
			wantErr: false,
		},
		{
			name:         "Failed get all books because of invalid parameter",
			booksRequest: &models.BooksList{},
			available:    false,
			mock: func(mock *mock.MockusecaseBooksRepository) {
			},
			wantErr: true,
			errType: fmt.Errorf("invalid parameter"),
		},
		{
			name:         "Failed get all books because no books found in database",
			booksRequest: &models.BooksList{},
			available:    true,
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().GetAll(&models.BooksList{}).Return(gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrRecordNotFound),
		},
		{
			name:         "Failed get all books because of database error",
			booksRequest: &models.BooksList{},
			available:    true,
			mock: func(mock *mock.MockusecaseBooksRepository) {
				mock.EXPECT().GetAll(&models.BooksList{}).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrInvalidDB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mock.NewMockusecaseBooksRepository(t)
			tt.mock(mock)

			uc := NewBooksUseCase(mock)

			books, err := uc.GetAllBooks(tt.available)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(*tt.expectedBooks), len(*books))
				for i, expected := range *tt.expectedBooks {
					actual := (*books)[i]
					assert.Equal(t, expected.ID, actual.ID)
					assert.Equal(t, expected.Title, actual.Title)
					assert.Equal(t, expected.Description, actual.Description)
					assert.Equal(t, expected.Qty, actual.Qty)
				}
			}
		})
	}
}

// func TestUpdateBook(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		bookRequest  *models.Books
// 		expectedBook *models.BooksSummary
// 		mock         func(mock *mock.MockusecaseBooksRepository)
// 		wantErr      bool
// 		errType      error
// 	}{
// 		{
// 			name: "Success update book with ID of 1",
// 			bookRequest: &models.Books{
// 				ID:          1,
// 				Title:       "Updated Title",
// 				Description: "Updated Description",
// 				Qty:         15,
// 			},
// 			expectedBook: &models.BooksSummary{
// 				ID:          1,
// 				Title:       "Updated Title",
// 				Description: "Updated Description",
// 				Qty:         15,
// 			},
// 			mock: func(mock *mock.MockusecaseBooksRepository) {
// 				mock.EXPECT().Update(&models.Books{
// 					ID:          1,
// 					Title:       "Updated Title",
// 					Description: "Updated Description",
// 					Qty:         15,
// 				}).Return(nil)
// 			},
// 		},
// 	}
// }

// func TestDeleteBook(t *testing.T) {
// }
