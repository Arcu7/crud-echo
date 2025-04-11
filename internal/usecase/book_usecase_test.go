package usecase

import (
	"crud-echo/internal/mocks"
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
		mock        func(mock *mocks.MockusecaseBooksRepository)
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
			mock: func(mock *mocks.MockusecaseBooksRepository) {
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
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				mock.EXPECT().ExistsByTitle("Test Title").Return(true, nil)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", models.ErrResourceAlreadyExist),
		},
		{
			name: "Failed create book due to invalid DB",
			bookRequest: &models.CreateBooksRequest{
				Title:       "Test Title",
				Description: "Test Description",
				Qty:         10,
			},
			mock: func(mock *mocks.MockusecaseBooksRepository) {
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
			mock := mocks.NewMockusecaseBooksRepository(t)
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
		mock         func(mock *mocks.MockusecaseBooksRepository)
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
			mock: func(mock *mocks.MockusecaseBooksRepository) {
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
			mock: func(mock *mocks.MockusecaseBooksRepository) {
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
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				mock.EXPECT().GetByID(&models.Books{}, 1).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrInvalidDB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mocks.NewMockusecaseBooksRepository(t)
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
		booksRequest  *[]models.Books
		expectedBooks *[]models.BooksSummary
		available     bool
		mock          func(mock *mocks.MockusecaseBooksRepository)
		wantErr       bool
		errType       error
	}{
		{
			name:         "Success get all books",
			booksRequest: &[]models.Books{},
			expectedBooks: &[]models.BooksSummary{
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
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				var arg []models.Books
				mock.EXPECT().GetAll(&arg).RunAndReturn(func(books *[]models.Books) error {
					*books = []models.Books{
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
			booksRequest: &[]models.Books{},
			available:    false,
			mock:         func(mock *mocks.MockusecaseBooksRepository) {},
			wantErr:      true,
			errType:      fmt.Errorf("invalid parameter"),
		},
		{
			name:         "Failed get all books because no books found in database",
			booksRequest: &[]models.Books{},
			available:    true,
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				var arg []models.Books
				mock.EXPECT().GetAll(&arg).Return(models.ErrNotFound)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", models.ErrNotFound),
		},
		{
			name:         "Failed get all books because of database error",
			booksRequest: &[]models.Books{},
			available:    true,
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				var arg []models.Books
				mock.EXPECT().GetAll(&arg).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrInvalidDB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mocks.NewMockusecaseBooksRepository(t)
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

func TestUpdateBook(t *testing.T) {
	tests := []struct {
		name        string
		bookRequest *models.UpdateBooksRequest
		mock        func(mock *mocks.MockusecaseBooksRepository)
		wantErr     bool
		errType     error
	}{
		{
			name: "Success update book with ID of 1",
			bookRequest: &models.UpdateBooksRequest{
				ID:          1,
				Title:       "Updated Title",
				Description: "Updated Description",
				Qty:         15,
			},
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				mock.EXPECT().Update(&models.Books{
					ID:          1,
					Title:       "Updated Title",
					Description: "Updated Description",
					Qty:         15,
				}).Return(nil)
			},
		},
		{
			name: "Failed update book due to book not found",
			bookRequest: &models.UpdateBooksRequest{
				ID:          1,
				Title:       "Updated Title",
				Description: "Updated Description",
				Qty:         15,
			},
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				mock.EXPECT().Update(&models.Books{
					ID:          1,
					Title:       "Updated Title",
					Description: "Updated Description",
					Qty:         15,
				}).Return(models.ErrNotFound)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", models.ErrNotFound),
		},
		{
			name: "Failed update book due to invalid DB",
			bookRequest: &models.UpdateBooksRequest{
				ID:          1,
				Title:       "Updated Title",
				Description: "Updated Description",
				Qty:         15,
			},
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				mock.EXPECT().Update(&models.Books{
					ID:          1,
					Title:       "Updated Title",
					Description: "Updated Description",
					Qty:         15,
				}).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrInvalidDB),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mocks.NewMockusecaseBooksRepository(t)
			tt.mock(mock)

			uc := NewBooksUseCase(mock)

			err := uc.UpdateBook(tt.bookRequest)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	tests := []struct {
		name        string
		bookRequest *models.DeleteBooksRequest
		mock        func(mock *mocks.MockusecaseBooksRepository)
		wantErr     bool
		errType     error
	}{
		{
			name: "Success update book with ID of 1",
			bookRequest: &models.DeleteBooksRequest{
				ID: 1,
			},
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				mock.EXPECT().Delete(&models.Books{
					ID: 1,
				}).Return(nil)
			},
		},
		{
			name: "Failed update book due to book not found",
			bookRequest: &models.DeleteBooksRequest{
				ID: 99,
			},
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				mock.EXPECT().Delete(&models.Books{
					ID: 99,
				}).Return(models.ErrNotFound)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", models.ErrNotFound),
		},
		{
			name: "Failed update book due to invalid DB",
			bookRequest: &models.DeleteBooksRequest{
				ID: 1,
			},
			mock: func(mock *mocks.MockusecaseBooksRepository) {
				mock.EXPECT().Delete(&models.Books{
					ID: 1,
				}).Return(gorm.ErrInvalidDB)
			},
			wantErr: true,
			errType: fmt.Errorf("repository error: %w", gorm.ErrInvalidDB),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mocks.NewMockusecaseBooksRepository(t)
			tt.mock(mock)

			uc := NewBooksUseCase(mock)

			err := uc.DeleteBook(tt.bookRequest)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
