# Project TODO List

## Completed: Email Subscription MVP ✅
- [x] 1. Project Setup
  - [x] Initialize Go module
  - [x] Create basic directory structure
  - [x] Add .gitignore
  - [x] Add .env.example
  - [x] Set up configuration management

- [x] Core Domain Implementation
  - [x] Contact Submissions
    - [x] CRUD Operations
    - [x] Status Management
    - [x] Input Validation
    - [x] Unit Tests
  - [x] Email Subscriptions
    - [x] CRUD Operations
    - [x] Status Management
    - [x] Input Validation
    - [x] Unit Tests
- [x] API Implementation
  - [x] RESTful Endpoints
  - [x] Standardized Response Format
  - [x] Error Handling
  - [x] Input Validation
  - [x] Unit Tests
- [x] Testing Infrastructure
  - [x] Mock Implementations
  - [x] Test Utilities
  - [x] Assertion Helpers
  - [x] Test Setup Utilities
- [x] Development Environment
  - [x] Dev Container Setup
  - [x] Task Automation
  - [x] Hot Reload
- [x] Security
  - [x] Authentication System
    - [x] User Model & Migration
    - [x] JWT Token Implementation
    - [x] Login Endpoint
    - [x] Signup Endpoint
    - [x] Middleware Protection

## In Progress

- [ ] Contact Form Demo
  - [x] Basic Form Implementation
  - [x] API Response Display
  - [x] Message History Display
  - [ ] Form Validation
  - [ ] Error Handling Improvements
  - [ ] Loading States
  - [ ] Success Feedback
  - [ ] Message Filtering
  - [ ] Message Sorting Options
  - [ ] Responsive Design
  - [ ] Accessibility Improvements

- [ ] Dependency Injection Improvements
  - [x] Logger Consistency
    - [x] Consolidate logger initialization
    - [x] Replace direct GetLogger calls with DI
    - [x] Add logger interface documentation
  - [x] Configuration Management
    - [x] Move server config to infrastructure layer
    - [x] Create unified config structure
    - [x] Add config validation
  - [x] Database Access
    - [x] Group store providers
    - [x] Ensure consistent store initialization
    - [x] Add store interfaces documentation
  - [ ] Handler Dependencies
    - [ ] Audit handler constructors
    - [ ] Ensure consistent DI usage
    - [ ] Document dependency requirements

## Current Sprint: MVP Stabilization 🚀
- [ ] 1. Deployment Readiness
  - [ ] Production Docker setup
  - [ ] Deployment documentation
  - [ ] Environment configuration guide
  - [ ] Logging setup review
  - [ ] Monitoring configuration

- [ ] 2. Code Quality
  - [ ] Complete test coverage
  - [ ] API documentation
  - [ ] Error handling review
  - [ ] Performance testing

## V2 Sprint: Multi-tenant Forms Platform 🎯
- [ ] 1. Multi-tenant Foundation
  - [ ] Create tenant table migration
    - [ ] ID, name, api_key fields
    - [ ] Domain whitelist
    - [ ] Rate limit configs
  - [ ] Tenant authentication system
  - [ ] Per-tenant rate limiting
  - [ ] Per-tenant CORS management

- [ ] 2. Form Management
  - [ ] Create forms table migration
    - [ ] ID, title, description, created_at, updated_at
    - [ ] Status (draft, published, archived)
    - [ ] Validation rules
  - [ ] Create form_fields table migration
    - [ ] Field types (text, number, email, etc.)
    - [ ] Validation rules
    - [ ] Required/optional status
  - [ ] Form CRUD endpoints
  - [ ] Field configuration
  - [ ] Input sanitization
  - [ ] API versioning (v1)
  - [ ] OpenAPI/Swagger annotations

- [ ] 3. Submission Management
  - [ ] Create form_submissions table migration
    - [ ] Submission data (JSON)
    - [ ] Metadata (IP, timestamp, etc.)
  - [ ] Form submission handling
  - [ ] Submission retrieval API
  - [ ] Export capabilities

- [ ] Form Management System
  - [ ] Form Schema Implementation
    - [ ] JSON Schema-based form definition
    - [ ] UI Schema for rendering configuration
    - [ ] Form settings and metadata
    - [ ] Version control for forms
  - [ ] Database Implementation
    - [ ] Forms table migration
    - [ ] Form submissions table migration
    - [ ] JSON validation functions
    - [ ] Form versioning system
  - [ ] Form Builder UI
    - [ ] Schema editor component
    - [ ] Live form preview
    - [ ] Settings configuration panel
    - [ ] Deployment instructions view
  - [ ] Form API Endpoints
    - [ ] Form CRUD operations
    - [ ] Form submission handling
    - [ ] Form analytics endpoints
    - [ ] Form embedding endpoints
  - [ ] JavaScript SDK
    - [ ] Form rendering library
    - [ ] Form submission handling
    - [ ] Validation implementation
    - [ ] Custom styling support
  - [ ] Security Features
    - [ ] Origin validation
    - [ ] Rate limiting per form
    - [ ] CAPTCHA integration
    - [ ] XSS protection
  - [ ] Integration System
    - [ ] Webhook support
    - [ ] Email notifications
    - [ ] Slack integration

- [ ] Testing Improvements
  - [ ] Integration Tests
  - [ ] Performance Tests
  - [ ] Load Tests
  - [ ] API Contract Tests
- [ ] Features
  - [ ] Form Versioning
  - [ ] Export/Import
  - [ ] Webhook Integration
  - [ ] Email Notifications
- [ ] Monitoring
  - [ ] Metrics Collection
  - [ ] Performance Monitoring
  - [ ] Error Tracking
  - [ ] Audit Logging
  - [ ] User Activity Tracking

- [ ] Infrastructure Separation of Concerns
  - [ ] Move domain-specific validation
    - [x] Move contact validation to domain/contact/validation
    - [x] Update validation imports and tests
  - [ ] Reorganize store layer
    - [ ] Move stores to respective domain packages
    - [ ] Update store interfaces and implementations
  - [ ] Clean up persistence layer
    - [ ] Merge or remove redundant persistence package
    - [ ] Update affected dependencies
  - [ ] Package Documentation
    - [x] Add README.md to each infrastructure package
    - [ ] Document package responsibilities
    - [ ] Document package interfaces
    - [ ] Add usage examples
