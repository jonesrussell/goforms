# Module Hierarchy

## Domain Module
- **user.Module**
  - Provides:
    - `UserService`
      - Depends on:
        - `UserRepository`
        - `Logger`
    - `UserRepository`
      - Implemented by:
        - `NewStore`

- **contact.Module**
  - Provides:
    - `ContactService`
      - Depends on:
        - `ContactRepository`
        - `Logger`

---

## View Module
- **view.Module**
  - Provides:
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

---

## Presentation Module
- **presentation.Module**
  - Provides:
    - `Provider Name`
      - Depends on:
        - `Dependency1`
        - `Dependency2` 