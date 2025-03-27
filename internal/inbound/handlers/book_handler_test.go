package handlers

import (
	"bytes"
	"crud-echo/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateMockDB(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Error in creating mock: ", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatal("GORM connection error: ", err)
	}

	return db, gdb, mock
}

func TestCreateBooks(t *testing.T) {
	e := echo.New()
	e.Validator = &vc.CustomValidator{Validator: validator.New()}

	t.Run("Success create user", func(t *testing.T) {
		sqldb, db, mock := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		user := models.Books{
			Name:  "testuser",
			Email: "test@gmail.com",
			Age:   20,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).
			WithArgs(user.Name, user.Email, user.Age).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(1))
		mock.ExpectCommit()

		expectedBooks := user
		expectedBooks.ID = 1

		body, err := json.Marshal(user)
		if err != nil {
			t.Fatal("Error in marshalling user: ", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, userHandler.CreateBooks(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var response models.Books

			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatal("Error in unmarshalling response: ", err)
			}

			assert.Equal(t, expectedBooks, response)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("Bad Request/Invalid input", func(t *testing.T) {
		sqldb, db, _ := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(`{test`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := userHandler.CreateBooks(c); assert.Error(t, err) {
			assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
		}
	})
}

func TestGetBooks(t *testing.T) {
	e := echo.New()
	e.Validator = &vc.CustomValidator{Validator: validator.New()}

	t.Run("Success get user", func(t *testing.T) {
		sqldb, db, mock := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		user := models.Books{
			ID:    1,
			Name:  "testuser",
			Email: "test@gmail.com",
		}

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(user.ID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(user.ID, user.Name, user.Email))

		req := httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(user.ID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(user.ID))

		if assert.NoError(t, userHandler.GetBooks(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response models.Books
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatal("Error in unmarshalling response: ", err)
			}

			assert.Equal(t, user, response)
		}
	})

	t.Run("Get but user not found", func(t *testing.T) {
		sqldb, db, mock := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		if err := userHandler.GetBooks(c); assert.Error(t, err) {
			assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
			assert.Equal(t, gorm.ErrRecordNotFound, err.(*echo.HTTPError).Message)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetAllBooks(t *testing.T) {
	e := echo.New()
	e.Validator = &vc.CustomValidator{Validator: validator.New()}

	t.Run("Success get all user", func(t *testing.T) {
		sqldb, db, mock := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		user := []models.Books{{ID: 1, Name: "testuser", Email: "tes@gmail.com"}}

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(user[0].ID, user[0].Name, user[0].Email))

		req := httptest.NewRequest(http.MethodGet, "/users/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, userHandler.GetAllBookss(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response []models.Books
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatal("Error in unmarshalling response: ", err)
			}

			assert.Equal(t, user, response)
		}
	})

	t.Run("Get all but no user found", func(t *testing.T) {
		sqldb, db, mock := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := userHandler.GetAllBookss(c); assert.Error(t, err) {
			assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
			assert.Equal(t, gorm.ErrRecordNotFound, err.(*echo.HTTPError).Message)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestUpdateBooks(t *testing.T) {
	e := echo.New()
	e.Validator = &vc.CustomValidator{Validator: validator.New()}

	t.Run("Success update user", func(t *testing.T) {
		sqldb, db, mock := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		user := models.Books{
			ID:    1,
			Name:  "testuserupdate",
			Email: "testupdate@gmail.com",
			Age:   20,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "name"=\$1,"email"=\$2,"age"=\$3 WHERE "id" = \$4`).
			WithArgs(user.Name, user.Email, user.Age, user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPut, "/users/"+strconv.Itoa(user.ID), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(user.ID))

		expectedMessage := "Books ID " + c.Param("id") + " has been updated\n"

		if assert.NoError(t, userHandler.UpdateBooks(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, rec.Body.String(), expectedMessage)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Bad Update Request/Invalid input", func(t *testing.T) {
		sqldb, db, _ := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewReader([]byte(`{test`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := userHandler.UpdateBooks(c); assert.Error(t, err) {
			assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
		}
	})
}

func TestDeleteBooks(t *testing.T) {
	e := echo.New()
	e.Validator = &vc.CustomValidator{Validator: validator.New()}

	t.Run("Success delete user", func(t *testing.T) {
		sqldb, db, mock := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		deletedBooksID := 999

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "users" WHERE "users"."id" = \$1`).
			WithArgs(deletedBooksID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest(http.MethodDelete, "/users/"+strconv.Itoa(deletedBooksID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(deletedBooksID))

		expectedMessage := "Books ID " + c.Param("id") + " has been deleted\n"

		if assert.NoError(t, userHandler.DeleteBooks(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, rec.Body.String(), expectedMessage)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("Delete but user not found", func(t *testing.T) {
		sqldb, db, mock := CreateMockDB(t)
		defer sqldb.Close()

		userHandler := BooksHandler{DB: db}

		deletedBooksID := 999

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "users" WHERE "users"."id" = \$1`).
			WithArgs(deletedBooksID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		req := httptest.NewRequest(http.MethodDelete, "/users/"+strconv.Itoa(deletedBooksID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(deletedBooksID))

		if err := userHandler.DeleteBooks(c); assert.Error(t, err) {
			assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
			assert.Equal(t, gorm.ErrRecordNotFound, err.(*echo.HTTPError).Message)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
