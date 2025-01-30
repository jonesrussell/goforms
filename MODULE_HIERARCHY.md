# Module Hierarchy

## Domain Module
- **user.Module**
  - Provides:
    - `UserService`
      - Depends on:
        - `UserRepository`
        - `TokenRepository`
        - `Logger`
      - Methods:
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
    - `UserRepository`
      - Implemented by:
        - `NewStore`
      - **Important Note**: Ensure that `user.Repository` is only provided once to avoid conflicts during application startup.

- **contact.Module**
  - Provides:
    - `ContactService`
      - Depends on:
        - `ContactRepository`
        - `Logger`
      - Methods:
        - `Submit(ctx context.Context, sub *common.Submission) error`
        - `ListSubmissions(ctx context.Context) ([]common.Submission, error)`
        - `GetSubmission(ctx context.Context, id int64) (*common.Submission, error)`
        - `UpdateSubmissionStatus(ctx context.Context, id int64, status string) error`
    - `ContactRepository`
      - Implemented by:
        - `NewStore`

---

## View Module
- **view.Module**
  - Provides:
    - `NewRenderer`
      - Depends on:
        - `Logger`

---

## Application Module
- **application.Module**
  - Provides:
    - `NewWebHandler`
      - Depends on:
        - `Logger`
        - `Renderer`
        - `ContactService`
    - `NewAuthHandler`
      - Depends on:
        - `Logger`
        - `UserService`
    - `NewRenderer`
      - Depends on:
        - `Logger`

---

## Main Application
- **Main Entry Point**: `cmd/goforms/main.go`
  - Initializes the application with the following modules:
    - `logging.Module`
    - `config.Module`
    - `database.Module`
    - `domain.Module`
    - `application.Module`
    - `user.Module`
    - `view.Module`
