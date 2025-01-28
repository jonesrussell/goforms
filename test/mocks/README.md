# Mocks Package

This package contains all mock implementations used in testing.

## Structure

Mocks are organized by domain:

```
mocks/
├── contact/           # Contact domain mocks
│   ├── service.go    # Contact service mock
│   └── store.go      # Contact store mock
├── subscription/      # Subscription domain mocks
│   ├── service.go    # Subscription service mock
│   └── store.go      # Subscription store mock
└── user/             # User domain mocks
    ├── service.go    # User service mock
    └── store.go      # User store mock
```

## Standards

1. Directory Structure
   - One directory per domain
   - Flat structure within domain directories
   - No nested mock directories

2. File Naming
   - No "mock_" prefix
   - Named after the interface being mocked
   - Use .go extension

3. Implementation Pattern
   - Use sync.Mutex for thread safety
   - Implement expectations system
   - Provide Verify() method
   - Provide Reset() method
   - Include interface compliance check

4. Documentation
   - Each mock package must have README.md
   - Each mock must document the interface it implements
   - Include usage examples in tests 