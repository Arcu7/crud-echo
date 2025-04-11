// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	models "crud-echo/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// MockhandlerBookUsecase is an autogenerated mock type for the handlerBookUsecase type
type MockhandlerBookUsecase struct {
	mock.Mock
}

type MockhandlerBookUsecase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockhandlerBookUsecase) EXPECT() *MockhandlerBookUsecase_Expecter {
	return &MockhandlerBookUsecase_Expecter{mock: &_m.Mock}
}

// CreateBook provides a mock function with given fields: book
func (_m *MockhandlerBookUsecase) CreateBook(book *models.CreateBooksRequest) (*models.Books, error) {
	ret := _m.Called(book)

	if len(ret) == 0 {
		panic("no return value specified for CreateBook")
	}

	var r0 *models.Books
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.CreateBooksRequest) (*models.Books, error)); ok {
		return rf(book)
	}
	if rf, ok := ret.Get(0).(func(*models.CreateBooksRequest) *models.Books); ok {
		r0 = rf(book)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Books)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.CreateBooksRequest) error); ok {
		r1 = rf(book)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockhandlerBookUsecase_CreateBook_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBook'
type MockhandlerBookUsecase_CreateBook_Call struct {
	*mock.Call
}

// CreateBook is a helper method to define mock.On call
//   - book *models.CreateBooksRequest
func (_e *MockhandlerBookUsecase_Expecter) CreateBook(book interface{}) *MockhandlerBookUsecase_CreateBook_Call {
	return &MockhandlerBookUsecase_CreateBook_Call{Call: _e.mock.On("CreateBook", book)}
}

func (_c *MockhandlerBookUsecase_CreateBook_Call) Run(run func(book *models.CreateBooksRequest)) *MockhandlerBookUsecase_CreateBook_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.CreateBooksRequest))
	})
	return _c
}

func (_c *MockhandlerBookUsecase_CreateBook_Call) Return(_a0 *models.Books, _a1 error) *MockhandlerBookUsecase_CreateBook_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockhandlerBookUsecase_CreateBook_Call) RunAndReturn(run func(*models.CreateBooksRequest) (*models.Books, error)) *MockhandlerBookUsecase_CreateBook_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteBook provides a mock function with given fields: book
func (_m *MockhandlerBookUsecase) DeleteBook(book *models.DeleteBooksRequest) error {
	ret := _m.Called(book)

	if len(ret) == 0 {
		panic("no return value specified for DeleteBook")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.DeleteBooksRequest) error); ok {
		r0 = rf(book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockhandlerBookUsecase_DeleteBook_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteBook'
type MockhandlerBookUsecase_DeleteBook_Call struct {
	*mock.Call
}

// DeleteBook is a helper method to define mock.On call
//   - book *models.DeleteBooksRequest
func (_e *MockhandlerBookUsecase_Expecter) DeleteBook(book interface{}) *MockhandlerBookUsecase_DeleteBook_Call {
	return &MockhandlerBookUsecase_DeleteBook_Call{Call: _e.mock.On("DeleteBook", book)}
}

func (_c *MockhandlerBookUsecase_DeleteBook_Call) Run(run func(book *models.DeleteBooksRequest)) *MockhandlerBookUsecase_DeleteBook_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.DeleteBooksRequest))
	})
	return _c
}

func (_c *MockhandlerBookUsecase_DeleteBook_Call) Return(_a0 error) *MockhandlerBookUsecase_DeleteBook_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockhandlerBookUsecase_DeleteBook_Call) RunAndReturn(run func(*models.DeleteBooksRequest) error) *MockhandlerBookUsecase_DeleteBook_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllBooks provides a mock function with given fields: available
func (_m *MockhandlerBookUsecase) GetAllBooks(available bool) (*[]models.BooksSummary, error) {
	ret := _m.Called(available)

	if len(ret) == 0 {
		panic("no return value specified for GetAllBooks")
	}

	var r0 *[]models.BooksSummary
	var r1 error
	if rf, ok := ret.Get(0).(func(bool) (*[]models.BooksSummary, error)); ok {
		return rf(available)
	}
	if rf, ok := ret.Get(0).(func(bool) *[]models.BooksSummary); ok {
		r0 = rf(available)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]models.BooksSummary)
		}
	}

	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(available)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockhandlerBookUsecase_GetAllBooks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllBooks'
type MockhandlerBookUsecase_GetAllBooks_Call struct {
	*mock.Call
}

// GetAllBooks is a helper method to define mock.On call
//   - available bool
func (_e *MockhandlerBookUsecase_Expecter) GetAllBooks(available interface{}) *MockhandlerBookUsecase_GetAllBooks_Call {
	return &MockhandlerBookUsecase_GetAllBooks_Call{Call: _e.mock.On("GetAllBooks", available)}
}

func (_c *MockhandlerBookUsecase_GetAllBooks_Call) Run(run func(available bool)) *MockhandlerBookUsecase_GetAllBooks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *MockhandlerBookUsecase_GetAllBooks_Call) Return(_a0 *[]models.BooksSummary, _a1 error) *MockhandlerBookUsecase_GetAllBooks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockhandlerBookUsecase_GetAllBooks_Call) RunAndReturn(run func(bool) (*[]models.BooksSummary, error)) *MockhandlerBookUsecase_GetAllBooks_Call {
	_c.Call.Return(run)
	return _c
}

// GetBookByID provides a mock function with given fields: id
func (_m *MockhandlerBookUsecase) GetBookByID(id int) (*models.BooksSummary, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetBookByID")
	}

	var r0 *models.BooksSummary
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*models.BooksSummary, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *models.BooksSummary); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.BooksSummary)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockhandlerBookUsecase_GetBookByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBookByID'
type MockhandlerBookUsecase_GetBookByID_Call struct {
	*mock.Call
}

// GetBookByID is a helper method to define mock.On call
//   - id int
func (_e *MockhandlerBookUsecase_Expecter) GetBookByID(id interface{}) *MockhandlerBookUsecase_GetBookByID_Call {
	return &MockhandlerBookUsecase_GetBookByID_Call{Call: _e.mock.On("GetBookByID", id)}
}

func (_c *MockhandlerBookUsecase_GetBookByID_Call) Run(run func(id int)) *MockhandlerBookUsecase_GetBookByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockhandlerBookUsecase_GetBookByID_Call) Return(_a0 *models.BooksSummary, _a1 error) *MockhandlerBookUsecase_GetBookByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockhandlerBookUsecase_GetBookByID_Call) RunAndReturn(run func(int) (*models.BooksSummary, error)) *MockhandlerBookUsecase_GetBookByID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateBook provides a mock function with given fields: book
func (_m *MockhandlerBookUsecase) UpdateBook(book *models.UpdateBooksRequest) error {
	ret := _m.Called(book)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBook")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.UpdateBooksRequest) error); ok {
		r0 = rf(book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockhandlerBookUsecase_UpdateBook_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateBook'
type MockhandlerBookUsecase_UpdateBook_Call struct {
	*mock.Call
}

// UpdateBook is a helper method to define mock.On call
//   - book *models.UpdateBooksRequest
func (_e *MockhandlerBookUsecase_Expecter) UpdateBook(book interface{}) *MockhandlerBookUsecase_UpdateBook_Call {
	return &MockhandlerBookUsecase_UpdateBook_Call{Call: _e.mock.On("UpdateBook", book)}
}

func (_c *MockhandlerBookUsecase_UpdateBook_Call) Run(run func(book *models.UpdateBooksRequest)) *MockhandlerBookUsecase_UpdateBook_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.UpdateBooksRequest))
	})
	return _c
}

func (_c *MockhandlerBookUsecase_UpdateBook_Call) Return(_a0 error) *MockhandlerBookUsecase_UpdateBook_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockhandlerBookUsecase_UpdateBook_Call) RunAndReturn(run func(*models.UpdateBooksRequest) error) *MockhandlerBookUsecase_UpdateBook_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockhandlerBookUsecase creates a new instance of MockhandlerBookUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockhandlerBookUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockhandlerBookUsecase {
	mock := &MockhandlerBookUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
