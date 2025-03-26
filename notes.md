Let me analyze the codebase to suggest appropriate unit tests based on the functionality present.

Search files...
Based on the codebase analysis, I can see this is a CRUD application for managing books using Echo framework and GORM. Let me suggest appropriate unit tests for each layer of the application following clean architecture principles.

1. **Repository Layer Tests** (`internal/repository/book_repository_test.go`):

```go
- TestCreateBook
  - Success case: Create a new book
  - Error case: Database error

- TestGetByID
  - Success case: Get existing book
  - Error case: Book not found
  - Error case: Database error

- TestGetAll
  - Success case: Get all books
  - Error case: Database error

- TestUpdate
  - Success case: Update existing book
  - Error case: Book not found
  - Error case: Database error

- TestDelete
  - Success case: Delete existing book
  - Error case: Book not found
  - Error case: Database error
```

2. **Use Case Layer Tests** (`internal/usecase/book_usecase_test.go`):

```go
- TestCreateBook
  - Success case: Valid book data
  - Error case: Invalid validation (title too short)
  - Error case: Invalid validation (quantity out of range)
  - Error case: Repository error

- TestGetBook
  - Success case: Get existing book
  - Error case: Book not found
  - Error case: Repository error

- TestGetAllBooks
  - Success case: Get all books
  - Error case: Repository error

- TestUpdateBook
  - Success case: Update with valid data
  - Error case: Invalid validation
  - Error case: Book not found
  - Error case: Repository error

- TestDeleteBook
  - Success case: Delete existing book
  - Error case: Invalid validation
  - Error case: Book not found
  - Error case: Repository error
```

3. **Handler Layer Tests** (`internal/inbound/handlers/book_handler_test.go`):

```go
- TestCreateBook
  - Success case: Valid JSON request
  - Error case: Invalid JSON
  - Error case: Invalid validation
  - Error case: Use case error

- TestGetBook
  - Success case: Get existing book
  - Error case: Invalid ID parameter
  - Error case: Book not found
  - Error case: Use case error

- TestGetAllBooks
  - Success case: Get all books
  - Error case: Invalid available parameter
  - Error case: Use case error

- TestUpdateBook
  - Success case: Update with valid data
  - Error case: Invalid JSON
  - Error case: Invalid validation
  - Error case: Use case error

- TestDeleteBook
  - Success case: Delete existing book
  - Error case: Invalid JSON
  - Error case: Invalid validation
  - Error case: Use case error
```

4. **Validator Tests** (`internal/usecase/validator_test.go`):

```go
- TestCustomValidator_Validate
  - Success case: Valid struct
  - Error case: Required field missing
  - Error case: Field too short
  - Error case: Field too long
  - Error case: Invalid quantity range
```

5. **Model Tests** (`internal/models/model_test.go`):

```go
- TestBooks_ToBooksSummary
  - Success case: Convert book to summary

- TestBooksList_ToBooksSummary
  - Success case: Convert list of books to summaries
```

For implementing these tests, you'll need to:

1. Use mocking for dependencies:
   - Use `go-sqlmock` for database mocking in repository tests
   - Create mock interfaces for repository in use case tests
   - Create mock interfaces for use case in handler tests

2. Use test helpers:
   - Create helper functions for common test setup
   - Create helper functions for creating test data
   - Create helper functions for assertions

3. Follow testing best practices:
   - Use table-driven tests where appropriate
   - Test both success and error cases
   - Test edge cases
   - Use meaningful test names
   - Keep tests focused and isolated

Would you like me to help you implement any specific test from this list? I can provide the detailed implementation for any of these test cases.
