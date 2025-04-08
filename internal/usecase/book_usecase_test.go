package usecase

import (
	"crud-echo/internal/models"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Helper function to create test create book request data
func createTestCreateBookRequest(title string, description string, qty int) *models.CreateBooksRequest {
	return &models.CreateBooksRequest{
		Title:       title,
		Description: description,
		Qty:         qty,
	}
}

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name        string
		bookRequest *models.CreateBooksRequest
		mock        func(mock *MockusecaseBooksRepository)
		wantErr     bool
		errType     error
		expectedID  int
	}{
		{
			name:        "Success create book with valid data",
			bookRequest: createTestCreateBookRequest("Test Title", "Test Description", 10),
			expectedID:  1,
			mock: func(mock *MockusecaseBooksRepository) {
				mock.EXPECT().Create(createTestCreateBookRequest("Test Title", "Test Description", 10)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Database error during create",
			bookRequest: createTestCreateBookRequest("Test Title", "Test Description", 10),
			mock: func(mock *MockusecaseBooksRepository) {
				mock.EXPECT().Create(createTestCreateBookRequest("Test Title", "Test Description", 10)).Return(nil)
			},
			wantErr: true,
			errType: gorm.ErrInvalidDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockusecaseBooksRepository(t)
			tt.mock(mock)

			ucValidator := NewCustomValidator(validator.New())

			uc := NewBooksUseCase(mock, ucValidator)

			book, err := uc.CreateBook(tt.bookRequest)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, book.ID) // Assuming the ID is set to 1 in the mock
				assert.Equal(t, tt.bookRequest.Title, book.Title)
				assert.Equal(t, tt.bookRequest.Description, book.Description)
				assert.Equal(t, tt.bookRequest.Qty, book.Qty)
			}
		})
	}
}

func TestGetBook(t *testing.T) {
}

func TestGetAllBooks(t *testing.T) {
}

func TestUpdateBook(t *testing.T) {
}

func TestDeleteBook(t *testing.T) {
}
