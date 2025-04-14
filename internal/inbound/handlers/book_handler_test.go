package handlers

import (
	"bytes"
	vc "crud-echo/internal/inbound/customvalidator"
	"crud-echo/internal/mocks"
	"crud-echo/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestContext struct {
	Echo    *echo.Echo
	Handler *BooksHandler
	Mock    *mocks.MockhandlerBookUsecase
}

func initialSetup(t *testing.T) *TestContext {
	e := echo.New()
	e.HTTPErrorHandler = CustomHTTPErrorHandler

	mockUsecase := mocks.NewMockhandlerBookUsecase(t)
	testValidator := &vc.CustomValidator{Validator: validator.New()}

	handler := NewBooksHandler(mockUsecase, testValidator)

	return &TestContext{
		Echo:    e,
		Handler: handler,
		Mock:    mockUsecase,
	}
}

func (tc *TestContext) executeRequest(method string, path string, body string, handler echo.HandlerFunc) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	rec := httptest.NewRecorder()
	c := tc.Echo.NewContext(req, rec)

	err := handler(c)
	if err != nil {
		tc.Echo.HTTPErrorHandler(err, c)
	}

	return rec
}

// maybe there is a better way...
func (tc *TestContext) executeRequestWithParam(method string, path string, paramName string, paramValue string, body string, handler echo.HandlerFunc) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	rec := httptest.NewRecorder()
	c := tc.Echo.NewContext(req, rec)
	if paramName != "" && paramValue != "" {
		c.SetPath(path)
		c.SetParamNames(paramName)
		c.SetParamValues(paramValue)
	}

	err := handler(c)
	if err != nil {
		tc.Echo.HTTPErrorHandler(err, c)
	}

	return rec
}

func (tc *TestContext) executeRequestWithQuery(method string, path string, query string, body string, handler echo.HandlerFunc) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if query != "" {
		req.URL.RawQuery = query
	}

	rec := httptest.NewRecorder()
	c := tc.Echo.NewContext(req, rec)

	err := handler(c)
	if err != nil {
		tc.Echo.HTTPErrorHandler(err, c)
	}

	return rec
}

// maybe shouldn't be a method...
func (tc *TestContext) unmarshalJSONResponse(t *testing.T, body string) Response {
	var response Response
	err := json.Unmarshal([]byte(body), &response)
	t.Logf("response: %#v", response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	return response
}

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      string
		m                func(mockuc *mocks.MockhandlerBookUsecase)
		expectedStatus   int
		expectedResponse Response
	}{
		{
			name:        "Success create book",
			requestBody: `{"title":"Test Book","description":"Test Description","qty":10}`,
			m: func(mockuc *mocks.MockhandlerBookUsecase) {
				mockuc.EXPECT().CreateBook(&models.CreateBooksRequest{
					Title:       "Test Book",
					Description: "Test Description",
					Qty:         10,
				}).Return(&models.Books{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: Response{
				Status:  true,
				Message: "Book has been created",
				Data:    nil,
			},
		},
		{
			name:           "Failed create book due to bind error",
			requestBody:    `{,,,}`,
			m:              func(mockuc *mocks.MockhandlerBookUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: Response{
				Status:  false,
				Message: models.BadRequest,
				Data:    nil,
			},
		},
		{
			name:           "Failed create book due to validation error",
			requestBody:    `{"title":"xx","description":"xx","qty":101}`,
			m:              func(mockuc *mocks.MockhandlerBookUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: Response{
				Status:  false,
				Message: models.ValidationError,
				Data:    nil,
			},
		},
		{
			name:        "Failed create book due to book exist already",
			requestBody: `{"title":"Test Book","description":"Test Description","qty":10}`,
			m: func(mockuc *mocks.MockhandlerBookUsecase) {
				mockuc.EXPECT().CreateBook(&models.CreateBooksRequest{
					Title:       "Test Book",
					Description: "Test Description",
					Qty:         10,
				}).Return(nil, fmt.Errorf("repository error: %w", models.ErrResourceAlreadyExist))
			},
			expectedStatus: http.StatusConflict,
			expectedResponse: Response{
				Status:  false,
				Message: models.ResourceAlreadyExist,
				Data:    nil,
			},
		},
		{
			name:        "Failed create book due to 0 ID", // don't know about this
			requestBody: `{"title":"Test Book","description":"Test Description","qty":10}`,
			m: func(mockuc *mocks.MockhandlerBookUsecase) {
				mockuc.EXPECT().CreateBook(&models.CreateBooksRequest{
					Title:       "Test Book",
					Description: "Test Description",
					Qty:         10,
				}).Return(nil, fmt.Errorf("repository error: %w", models.ErrInternalServerError))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: Response{
				Status:  false,
				Message: "internal server error",
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := initialSetup(t)
			tt.m(tc.Mock)

			rec := tc.executeRequest(http.MethodPost, "/", tt.requestBody, tc.Handler.CreateBook)
			actualResponse := tc.unmarshalJSONResponse(t, rec.Body.String())

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedResponse.Status, actualResponse.Status)
			assert.Equal(t, tt.expectedResponse.Message, actualResponse.Message)
			assert.Equal(t, tt.expectedResponse.Data, actualResponse.Data)
		})
	}
}

func TestGetBookByID(t *testing.T) {
	tests := []struct {
		name             string
		param            string
		m                func(mockuc *mocks.MockhandlerBookUsecase)
		expectedStatus   int
		expectedResponse Response
	}{
		{
			name:  "Success get book by ID",
			param: "1",
			m: func(mockuc *mocks.MockhandlerBookUsecase) {
				mockuc.EXPECT().GetBookByID(1).Return(&models.BooksSummary{
					ID:          1,
					Title:       "Test Book",
					Description: "Test Description",
					Qty:         10,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: Response{
				Status:  true,
				Message: "Book retrieved successfully",
				Data: &models.BooksSummary{
					ID:          1,
					Title:       "Test Book",
					Description: "Test Description",
					Qty:         10,
				},
			},
		},
		{
			name:           "Failed get book by ID due to error converting ID param",
			param:          "!@#$%^&*()",
			m:              func(mock *mocks.MockhandlerBookUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: Response{
				Status:  false,
				Message: "invalid parameter",
				Data:    nil,
			},
		},
		{
			name:  "Failed get book by ID due to book not found",
			param: "99",
			m: func(mock *mocks.MockhandlerBookUsecase) {
				mock.EXPECT().GetBookByID(99).
					Return(nil, fmt.Errorf("repository error: %w", models.ErrNotFound))
			},
			expectedStatus: http.StatusNotFound,
			expectedResponse: Response{
				Status:  false,
				Message: "record not found",
				Data:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := initialSetup(t)
			tt.m(tc.Mock)

			rec := tc.executeRequestWithParam(http.MethodPost, "/books/:id", "id", tt.param, "", tc.Handler.GetBookByID)
			t.Logf("rec.Body.String(): %#v", rec.Body.String())
			actualResponse := tc.unmarshalJSONResponse(t, rec.Body.String())

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedResponse.Status, actualResponse.Status)
			assert.Equal(t, tt.expectedResponse.Message, actualResponse.Message)

			// hackaround?
			if tt.expectedResponse.Data != nil {
				if expectedData, ok := tt.expectedResponse.Data.(*models.BooksSummary); ok {
					if actualData, ok := actualResponse.Data.(*models.BooksSummary); ok {
						assert.Equal(t, expectedData.ID, actualData.ID)
						assert.Equal(t, expectedData.Title, actualData.Title)
						assert.Equal(t, expectedData.Description, actualData.Description)
						assert.Equal(t, expectedData.Qty, actualData.Qty)
					}
				}
			} else {
				assert.Nil(t, actualResponse.Data)
			}
		})
	}
}

// func TestGetAllBooks(t *testing.T) {
// 	tests := []struct {
// 		name             string
// 		param            string
// 		m                func(mockuc *mocks.MockhandlerBookUsecase)
// 		expectedStatus   int
// 		expectedResponse Response
// 	}{
// 		{
// 			name:  "Success get all books",
// 			param: "true",
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 				mockuc.EXPECT().GetAllBooks(true).Return(&[]models.BooksSummary{
// 					{
// 						ID:          1,
// 						Title:       "Test Book",
// 						Description: "Test Description",
// 						Qty:         10,
// 					},
// 					{
// 						ID:          2,
// 						Title:       "Test Book 2",
// 						Description: "Test Description 2",
// 						Qty:         20,
// 					},
// 				}, nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedResponse: Response{
// 				Status:  true,
// 				Message: "Books retrieved successfully",
// 				Data: &[]models.BooksSummary{
// 					{
// 						ID:          1,
// 						Title:       "Test Book",
// 						Description: "Test Description",
// 						Qty:         10,
// 					},
// 					{
// 						ID:          2,
// 						Title:       "Test Book 2",
// 						Description: "Test Description 2",
// 						Qty:         20,
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:           "Failed get all books due to error parsing available param",
// 			param:          "!@#$%^&*()",
// 			m:              func(mock *mocks.MockhandlerBookUsecase) {},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: "bad request",
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name:  "Failed get all books due to false available param",
// 			param: "false",
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 				mockuc.EXPECT().GetAllBooks(false).Return(nil, fmt.Errorf("repository error: %w", models.ErrInvalidParam))
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: "invalid parameter",
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name:  "Failed get all books due to empty table", // hmm maybe not fail?
// 			param: "true",
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 				mockuc.EXPECT().GetAllBooks(true).Return(nil, fmt.Errorf("repository error: %w", models.ErrEmptyTable))
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: "table is empty",
// 				Data:    nil,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e, handler, mock := initialSetup(t)
// 			tt.m(mock)
//
// 			expectedResponseBytes, err := json.Marshal(tt.expectedResponse)
// 			if err != nil {
// 				t.Fatalf("Failed to marshal expected response: %v", err)
// 			}
//
// 			req := httptest.NewRequest(http.MethodPost, "/", nil)
// 			req.URL.RawQuery = fmt.Sprintf("available=%s", tt.param)
// 			rec := httptest.NewRecorder()
// 			c := e.NewContext(req, rec)
//
// 			err = handler.GetAllBooks(c)
// 			if err != nil {
// 				e.HTTPErrorHandler(err, c)
// 			}
//
// 			assert.Equal(t, tt.expectedStatus, rec.Code)
// 			assert.JSONEq(t, string(expectedResponseBytes), rec.Body.String())
// 		})
// 	}
// }
//
// func TestUpdateBook(t *testing.T) {
// 	tests := []struct {
// 		name             string
// 		requestBody      string
// 		m                func(mockuc *mocks.MockhandlerBookUsecase)
// 		expectedStatus   int
// 		expectedResponse Response
// 	}{
// 		{
// 			name:        "Success update book",
// 			requestBody: `{"id": 1,"title":"Test Book","description":"Test Description","qty":10}`,
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 				mockuc.EXPECT().UpdateBook(&models.UpdateBooksRequest{
// 					ID:          1,
// 					Title:       "Test Book",
// 					Description: "Test Description",
// 					Qty:         10,
// 				}).Return(nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedResponse: Response{
// 				Status:  true,
// 				Message: "Book with ID 1 has been updated",
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name:           "Failed update book due to bind error",
// 			requestBody:    `{,,,}`,
// 			m:              func(mockuc *mocks.MockhandlerBookUsecase) {},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: models.BadRequest,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name:           "Failed update book due to validation error",
// 			requestBody:    `{"id": -1,"title":"xx","description":"xx","qty":101}`,
// 			m:              func(mockuc *mocks.MockhandlerBookUsecase) {},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: models.ValidationError,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name:        "Failed update book due to book not found",
// 			requestBody: `{"id": 1,"title":"Test Book","description":"Test Description","qty":10}`,
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 				mockuc.EXPECT().UpdateBook(&models.UpdateBooksRequest{
// 					ID:          1,
// 					Title:       "Test Book",
// 					Description: "Test Description",
// 					Qty:         10,
// 				}).Return(fmt.Errorf("repository error: %w", models.ErrNotFound))
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: models.NotFound,
// 				Data:    nil,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e, handler, mock := initialSetup(t)
// 			tt.m(mock)
//
// 			expectedResponseBytes, err := json.Marshal(tt.expectedResponse)
// 			if err != nil {
// 				t.Fatalf("Failed to marshal expected response: %v", err)
// 			}
//
// 			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 			rec := httptest.NewRecorder()
// 			c := e.NewContext(req, rec)
//
// 			err = handler.UpdateBook(c)
// 			if err != nil {
// 				e.HTTPErrorHandler(err, c)
// 			}
//
// 			assert.Equal(t, tt.expectedStatus, rec.Code)
// 			assert.JSONEq(t, string(expectedResponseBytes), rec.Body.String())
// 		})
// 	}
// }
//
// func TestDeleteBook(t *testing.T) {
// 	tests := []struct {
// 		name             string
// 		requestBody      string
// 		m                func(mock *mocks.MockhandlerBookUsecase)
// 		expectedStatus   int
// 		expectedResponse Response
// 	}{
// 		{
// 			name:        "Success update book",
// 			requestBody: `{"id": 1}`,
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 				mockuc.EXPECT().DeleteBook(&models.DeleteBooksRequest{
// 					ID: 1,
// 				}).Return(nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedResponse: Response{
// 				Status:  true,
// 				Message: "Book with ID 1 has been deleted",
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name:        "Failed delete book due to bind error",
// 			requestBody: `{,,,}`,
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: models.BadRequest,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name:        "Failed delete book due to validation error",
// 			requestBody: `{"id": -1}`,
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: models.ValidationError,
// 				Data:    nil,
// 			},
// 		},
// 		{
// 			name:        "Failed delete book due to book not found",
// 			requestBody: `{"id": 1}`,
// 			m: func(mockuc *mocks.MockhandlerBookUsecase) {
// 				mockuc.EXPECT().DeleteBook(&models.DeleteBooksRequest{
// 					ID: 1,
// 				}).Return(fmt.Errorf("repository error: %w", models.ErrNotFound))
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			expectedResponse: Response{
// 				Status:  false,
// 				Message: models.NotFound,
// 				Data:    nil,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e, handler, mock := initialSetup(t)
// 			tt.m(mock)
//
// 			expectedResponseBytes, err := json.Marshal(tt.expectedResponse)
// 			if err != nil {
// 				t.Fatalf("Failed to marshal expected response: %v", err)
// 			}
//
// 			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.requestBody))
// 			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 			rec := httptest.NewRecorder()
// 			c := e.NewContext(req, rec)
//
// 			err = handler.DeleteBook(c)
// 			if err != nil {
// 				e.HTTPErrorHandler(err, c)
// 			}
//
// 			assert.Equal(t, tt.expectedStatus, rec.Code)
// 			assert.JSONEq(t, string(expectedResponseBytes), rec.Body.String())
// 		})
// 	}
// }
