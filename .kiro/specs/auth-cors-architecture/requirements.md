# Requirements Document: Auth & CORS Architecture

## Introduction

This document defines the requirements for documenting the authentication and CORS architecture in the Main Service of ZooPlatforma. The system uses a microservices architecture with an API Gateway as the single entry point, handling authentication and CORS for all services.

## Glossary

- **Gateway**: API Gateway service running on port 80, acts as the single entry point for all requests
- **Main_Backend**: Main backend service running on port 8000
- **Main_Frontend**: Next.js frontend application running on port 3000
- **JWT**: JSON Web Token used for authentication
- **CORS**: Cross-Origin Resource Sharing mechanism for allowing frontend-backend communication
- **DevAuthMiddleware**: Middleware used in development to validate JWT tokens locally
- **AuthMiddleware**: Middleware used in production to read user data from Gateway headers
- **apiClient**: TypeScript client for making API requests with automatic URL configuration
- **X-User-Headers**: HTTP headers (X-User-ID, X-User-Email, X-User-Role) added by Gateway

## Requirements

### Requirement 1: Authentication Architecture Documentation

**User Story:** As a developer, I want to understand how authentication works in development vs production, so that I can correctly implement and debug authentication flows.

#### Acceptance Criteria

1. THE Documentation SHALL describe the development authentication flow using DevAuthMiddleware
2. THE Documentation SHALL describe the production authentication flow using Gateway
3. WHEN a developer reads the documentation, THE Documentation SHALL explain how JWT tokens are validated in each environment
4. THE Documentation SHALL describe how user context is passed to handlers in both environments
5. THE Documentation SHALL include diagrams showing authentication flow in development and production

### Requirement 2: CORS Architecture Documentation

**User Story:** As a developer, I want to understand why Gateway manages CORS and how to avoid duplicate headers, so that I can prevent CORS errors.

#### Acceptance Criteria

1. THE Documentation SHALL explain why Gateway is the single source of CORS headers
2. THE Documentation SHALL describe how Gateway filters CORS headers from backend responses
3. WHEN a developer needs to add a new frontend origin, THE Documentation SHALL provide clear instructions
4. THE Documentation SHALL explain the consequences of backend services setting CORS headers
5. THE Documentation SHALL include examples of correct and incorrect CORS configurations

### Requirement 3: API Client Documentation

**User Story:** As a developer, I want to understand how apiClient works and why it should be used, so that I can make API requests correctly in all environments.

#### Acceptance Criteria

1. THE Documentation SHALL describe how apiClient determines the correct API URL
2. THE Documentation SHALL explain the difference between development and production URL configuration
3. WHEN a developer makes API requests, THE Documentation SHALL explain why apiClient should be used instead of direct fetch
4. THE Documentation SHALL describe how apiClient handles credentials and authorization headers
5. THE Documentation SHALL provide examples of correct API client usage

### Requirement 4: Environment Configuration Documentation

**User Story:** As a developer, I want to understand how to configure environment variables for development and production, so that I can deploy the application correctly.

#### Acceptance Criteria

1. THE Documentation SHALL list all required environment variables for development
2. THE Documentation SHALL list all required environment variables for production
3. WHEN switching between environments, THE Documentation SHALL explain which configuration changes are needed
4. THE Documentation SHALL describe the role of NEXT_PUBLIC_API_URL in different environments
5. THE Documentation SHALL explain why JWT_SECRET must match between Gateway and Backend

### Requirement 5: Troubleshooting Guide

**User Story:** As a developer, I want a guide for common authentication and CORS errors, so that I can quickly resolve issues.

#### Acceptance Criteria

1. WHEN a 401 Unauthorized error occurs, THE Documentation SHALL provide diagnostic steps
2. WHEN CORS errors occur, THE Documentation SHALL explain how to identify the root cause
3. WHEN "Failed to fetch" errors occur, THE Documentation SHALL provide troubleshooting steps
4. THE Documentation SHALL include examples of error messages and their solutions
5. THE Documentation SHALL provide commands for verifying correct configuration

### Requirement 6: Architecture Diagrams

**User Story:** As a developer, I want visual diagrams of the authentication and CORS architecture, so that I can quickly understand the system flow.

#### Acceptance Criteria

1. THE Documentation SHALL include a diagram showing development architecture
2. THE Documentation SHALL include a diagram showing production architecture
3. THE Documentation SHALL include a diagram showing CORS header flow
4. THE Documentation SHALL include a diagram showing authentication flow
5. WHEN viewing diagrams, THE Documentation SHALL use consistent notation and clear labels

### Requirement 7: Best Practices

**User Story:** As a developer, I want documented best practices for authentication and CORS, so that I can follow established patterns.

#### Acceptance Criteria

1. THE Documentation SHALL list best practices for using apiClient
2. THE Documentation SHALL list best practices for handling authentication in components
3. THE Documentation SHALL list best practices for CORS configuration
4. THE Documentation SHALL explain when to use DevAuthMiddleware vs AuthMiddleware
5. THE Documentation SHALL provide examples of correct implementation patterns

### Requirement 8: Migration Guide

**User Story:** As a developer, I want a guide for migrating from direct fetch to apiClient, so that I can update existing code correctly.

#### Acceptance Criteria

1. WHEN migrating existing code, THE Documentation SHALL provide before/after examples
2. THE Documentation SHALL explain how to identify code that needs migration
3. THE Documentation SHALL provide a checklist for verifying correct migration
4. THE Documentation SHALL explain common migration pitfalls
5. THE Documentation SHALL include examples for different types of API requests (GET, POST, PUT, DELETE)
