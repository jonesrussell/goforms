# Contact Mocks

This package provides mock implementations for the contact domain interfaces.

## Mocks

### Service (`service.go`)
Implements `contact.Service` interface from `internal/domain/contact/service.go`

Methods:
- Create(ctx context.Context, sub *contact.Submission) error
- List(ctx context.Context) ([]contact.Submission, error)
- Get(ctx context.Context, id int64) (*contact.Submission, error)
- UpdateStatus(ctx context.Context, id int64, status contact.Status) error

### Store (`store.go`)
Implements `contact.Store` interface from `internal/domain/contact/store.go`

Methods:
- Create(ctx context.Context, sub *contact.Submission) error
- List(ctx context.Context) ([]contact.Submission, error)
- Get(ctx context.Context, id int64) (*contact.Submission, error)
- UpdateStatus(ctx context.Context, id int64, status contact.Status) error

## Usage

```go
// Example test setup
mockSvc := mockcontact.NewMockService()
mockStore := mockcontact.NewMockStore()

// Set expectations
mockSvc.ExpectCreate(ctx, submission, nil) // expect successful creation
mockStore.ExpectGet(ctx, 1, &submission, nil) // expect successful retrieval

// Verify expectations were met
if err := mockSvc.Verify(); err != nil {
    t.Error(err)
}
``` 