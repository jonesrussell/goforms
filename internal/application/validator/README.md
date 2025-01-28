# Validator Package

## Overview

The validator package provides a custom validator implementation for Echo that uses `go-playground/validator/v10` under the hood. This package bridges Echo's validator interface with go-playground/validator's powerful validation capabilities.

## Architecture

```
                                    implements
CustomValidator (our code) -----> echo.Validator (interface)
        |
        | uses
        v
go-playground/validator (library)
```

## Components

### CustomValidator

```go
type CustomValidator struct {
    validator *validator.Validate  // from go-playground/validator/v10
}
```

The `CustomValidator` wraps go-playground/validator's `Validate` instance and implements Echo's `Validator` interface.

### Interface Implementation

```go
// Echo's Validator interface requirement
type Validator interface {
    Validate(i interface{}) error
}

// Our implementation
func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}
```

## Usage

1. Create a new validator:
```go
v := validator.NewValidator()
```

2. Register with Echo:
```go
e := echo.New()
e.Validator = validator.NewValidator()
```

3. Use in your structs:
```go
type User struct {
    Email     string `validate:"required,email"`
    Password  string `validate:"required,min=8"`
    FirstName string `validate:"required"`
    LastName  string `validate:"required"`
}
```

## Validation Tags

The validator supports all tags from go-playground/validator/v10, including:

- `required`: Field must be set
- `email`: Must be valid email format
- `min`: Minimum length for strings/slices
- `max`: Maximum length for strings/slices
- `gte`: Greater than or equal to
- `lte`: Less than or equal to

For a complete list of validation tags, see [go-playground/validator documentation](https://pkg.go.dev/github.com/go-playground/validator/v10).

## Testing

The package includes comprehensive tests in `validator_test.go` that verify:
- Validator initialization
- Basic validation functionality
- Common validation scenarios
- Error handling

## Dependencies

- github.com/go-playground/validator/v10
- github.com/labstack/echo/v4 