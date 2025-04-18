// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	models "crud-echo/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// MockusecaseBooksRepository is an autogenerated mock type for the usecaseBooksRepository type
type MockusecaseBooksRepository struct {
	mock.Mock
}

type MockusecaseBooksRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockusecaseBooksRepository) EXPECT() *MockusecaseBooksRepository_Expecter {
	return &MockusecaseBooksRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: book
func (_m *MockusecaseBooksRepository) Create(book *models.Books) error {
	ret := _m.Called(book)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Books) error); ok {
		r0 = rf(book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockusecaseBooksRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockusecaseBooksRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - book *models.Books
func (_e *MockusecaseBooksRepository_Expecter) Create(book interface{}) *MockusecaseBooksRepository_Create_Call {
	return &MockusecaseBooksRepository_Create_Call{Call: _e.mock.On("Create", book)}
}

func (_c *MockusecaseBooksRepository_Create_Call) Run(run func(book *models.Books)) *MockusecaseBooksRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Books))
	})
	return _c
}

func (_c *MockusecaseBooksRepository_Create_Call) Return(_a0 error) *MockusecaseBooksRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockusecaseBooksRepository_Create_Call) RunAndReturn(run func(*models.Books) error) *MockusecaseBooksRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: book
func (_m *MockusecaseBooksRepository) Delete(book *models.Books) error {
	ret := _m.Called(book)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Books) error); ok {
		r0 = rf(book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockusecaseBooksRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockusecaseBooksRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - book *models.Books
func (_e *MockusecaseBooksRepository_Expecter) Delete(book interface{}) *MockusecaseBooksRepository_Delete_Call {
	return &MockusecaseBooksRepository_Delete_Call{Call: _e.mock.On("Delete", book)}
}

func (_c *MockusecaseBooksRepository_Delete_Call) Run(run func(book *models.Books)) *MockusecaseBooksRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Books))
	})
	return _c
}

func (_c *MockusecaseBooksRepository_Delete_Call) Return(_a0 error) *MockusecaseBooksRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockusecaseBooksRepository_Delete_Call) RunAndReturn(run func(*models.Books) error) *MockusecaseBooksRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// ExistsByTitle provides a mock function with given fields: title
func (_m *MockusecaseBooksRepository) ExistsByTitle(title string) (bool, error) {
	ret := _m.Called(title)

	if len(ret) == 0 {
		panic("no return value specified for ExistsByTitle")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(title)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(title)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockusecaseBooksRepository_ExistsByTitle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistsByTitle'
type MockusecaseBooksRepository_ExistsByTitle_Call struct {
	*mock.Call
}

// ExistsByTitle is a helper method to define mock.On call
//   - title string
func (_e *MockusecaseBooksRepository_Expecter) ExistsByTitle(title interface{}) *MockusecaseBooksRepository_ExistsByTitle_Call {
	return &MockusecaseBooksRepository_ExistsByTitle_Call{Call: _e.mock.On("ExistsByTitle", title)}
}

func (_c *MockusecaseBooksRepository_ExistsByTitle_Call) Run(run func(title string)) *MockusecaseBooksRepository_ExistsByTitle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockusecaseBooksRepository_ExistsByTitle_Call) Return(_a0 bool, _a1 error) *MockusecaseBooksRepository_ExistsByTitle_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockusecaseBooksRepository_ExistsByTitle_Call) RunAndReturn(run func(string) (bool, error)) *MockusecaseBooksRepository_ExistsByTitle_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: book
func (_m *MockusecaseBooksRepository) GetAll(book *[]models.Books) error {
	ret := _m.Called(book)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*[]models.Books) error); ok {
		r0 = rf(book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockusecaseBooksRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockusecaseBooksRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - book *[]models.Books
func (_e *MockusecaseBooksRepository_Expecter) GetAll(book interface{}) *MockusecaseBooksRepository_GetAll_Call {
	return &MockusecaseBooksRepository_GetAll_Call{Call: _e.mock.On("GetAll", book)}
}

func (_c *MockusecaseBooksRepository_GetAll_Call) Run(run func(book *[]models.Books)) *MockusecaseBooksRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*[]models.Books))
	})
	return _c
}

func (_c *MockusecaseBooksRepository_GetAll_Call) Return(_a0 error) *MockusecaseBooksRepository_GetAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockusecaseBooksRepository_GetAll_Call) RunAndReturn(run func(*[]models.Books) error) *MockusecaseBooksRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: book, id
func (_m *MockusecaseBooksRepository) GetByID(book *models.Books, id int) error {
	ret := _m.Called(book, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Books, int) error); ok {
		r0 = rf(book, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockusecaseBooksRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockusecaseBooksRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - book *models.Books
//   - id int
func (_e *MockusecaseBooksRepository_Expecter) GetByID(book interface{}, id interface{}) *MockusecaseBooksRepository_GetByID_Call {
	return &MockusecaseBooksRepository_GetByID_Call{Call: _e.mock.On("GetByID", book, id)}
}

func (_c *MockusecaseBooksRepository_GetByID_Call) Run(run func(book *models.Books, id int)) *MockusecaseBooksRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Books), args[1].(int))
	})
	return _c
}

func (_c *MockusecaseBooksRepository_GetByID_Call) Return(_a0 error) *MockusecaseBooksRepository_GetByID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockusecaseBooksRepository_GetByID_Call) RunAndReturn(run func(*models.Books, int) error) *MockusecaseBooksRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: book
func (_m *MockusecaseBooksRepository) Update(book *models.Books) error {
	ret := _m.Called(book)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Books) error); ok {
		r0 = rf(book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockusecaseBooksRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockusecaseBooksRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - book *models.Books
func (_e *MockusecaseBooksRepository_Expecter) Update(book interface{}) *MockusecaseBooksRepository_Update_Call {
	return &MockusecaseBooksRepository_Update_Call{Call: _e.mock.On("Update", book)}
}

func (_c *MockusecaseBooksRepository_Update_Call) Run(run func(book *models.Books)) *MockusecaseBooksRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*models.Books))
	})
	return _c
}

func (_c *MockusecaseBooksRepository_Update_Call) Return(_a0 error) *MockusecaseBooksRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockusecaseBooksRepository_Update_Call) RunAndReturn(run func(*models.Books) error) *MockusecaseBooksRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockusecaseBooksRepository creates a new instance of MockusecaseBooksRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockusecaseBooksRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockusecaseBooksRepository {
	mock := &MockusecaseBooksRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
