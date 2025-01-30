# Dependency Documentation

## Domain Module

### Module Definition
```go
var Module = fx.Options(
    user.Module,
    contact.Module,
)
```

### Providers
- **user.Module**: Provides user-related services and repositories.
  - **Important Note**: Ensure that `user.Repository` is only provided once to avoid conflicts during application startup.

- **contact.Module**: Provides contact-related services and repositories.

### Dependencies
- **user.Module**:
  - `UserService`: Depends on `UserRepository`, `TokenRepository`, and `Logger`.
  - `UserRepository`: Implemented by `NewStore`, which interacts with the database.
- **contact.Module**:
  - `ContactService`: Depends on `ContactRepository` and `Logger`.
  - `ContactRepository`: Implemented by `NewStore`, which interacts with the database.

---

## Application Module

### Module Definition
```go
var Module = fx.Options(
    user.Module,
    contact.Module,
)
```

### Providers
- **application.Module**: Combines all application-level modules and providers.
  - **Important Note**: Ensure that handlers are registered correctly to avoid conflicts during application startup.

### Dependencies
- **application.Module**:
  - `NewWebHandler`: Creates a new instance of `WebHandler`.
  - `NewAuthHandler`: Creates a new instance of `AuthHandler`.
  - `NewRenderer`: Provides a new instance of `Renderer`.

---

## View Module

### Module Definition
```go
var Module = fx.Options(
    fx.Provide(
        NewRenderer,
    ),
)
```

### Providers
- **NewRenderer**: Creates a new instance of `Renderer`.

### Dependencies
- **NewRenderer**:
  - `Logger`: Required for logging within the renderer.

---

## Presentation Module

### Module Definition
```go
var Module = fx.Options(
    fx.Provide(
        NewPresenter,
    ),
)
```

### Providers
- **NewPresenter**: Creates a new instance of `Presenter`.

### Dependencies
- **NewPresenter**:
  - `Logger`: Required for logging within the presenter.

---

## Main Application Entry

### Module Initialization
```go
func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    loadEnvironment()
    versionInfo := createVersionInfo()
    app := createApp(versionInfo)
    return startApp(app)
}
```

### Application Modules
- **Modules Registered**:
  - `logging.Module`
  - `config.Module`
  - `database.Module`
  - `domain.Module`
  - `application.Module`
  - `user.Module`
  - `view.Module`

---

## Additional Modules

### Observations and Recommendations

1. **User Store Implementation**:
   - Ensure proper error handling and logging in database operations.
   - Verify that the `store` struct is instantiated correctly in the DI setup.

2. **User Service Implementation**:
   - Review error handling in `UpdateUser` and `DeleteUser` methods.
   - Ensure correct interaction with the `store`.

3. **Auth Handler**:
   - Validate user inputs in `handleSignup` and `handleLogin`.
   - Ensure proper token handling in `handleLogout`.

4. **Main Application Entry**:
   - Confirm that all modules are registered correctly in the `fx.New` call.
   - Check for missing dependencies that could cause startup failures.

5. **Mock Service for Testing**:
   - Ensure mock services simulate real behavior for effective testing.
   - Review tests for coverage of various scenarios, including error cases.

6. **Logging and Error Handling**:
   - Ensure that all modules utilize the logging framework consistently.
   - Implement structured error handling across all services.

7. **User Service Methods**:
   - Document the methods available in the `UserService` for clarity:
     - `DeleteUser(ctx context.Context, id uint) error`
     - `GetByEmail(email string) (*User, error)`
     - `GetUserByID(ctx context.Context, id uint) (*User, error)`
     - `ListUsers(ctx context.Context) ([]User, error)`
     - `Login(ctx context.Context, login *Login) (*TokenPair, error)`
     - `Logout(ctx context.Context, token string) error`
     - `SignUp(signup *Signup) (*User, error)`
     - `UpdateSubmissionStatus(ctx context.Context, id int64, status string) error`
     - `UpdateUser(ctx context.Context, user *User) error`
     - `IsTokenBlacklisted(token string) bool`
```