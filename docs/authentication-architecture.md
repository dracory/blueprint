# Authentication Architecture Documentation

## Overview

The Blueprint authentication system implements a sophisticated privacy-first architecture that prioritizes user data protection while maintaining robust security controls. This document outlines the complete authentication flow, security measures, and architectural decisions.

## Architecture Components

### 1. Authentication Controller (`internal/controllers/auth/`)

**Primary Responsibilities:**
- Handle external authentication via AuthKnight service
- Manage user creation and session establishment
- Implement privacy-first email encryption
- Coordinate between multiple storage systems

**Key Features:**
- Vault-based email encryption
- Blind indexing for privacy protection
- Context-aware HTTP requests
- Comprehensive error handling

### 2. External Authentication Integration

**AuthKnight Service Integration:**
- **Endpoint**: `https://authknight.com/api/who`
- **Method**: POST with form-encoded data
- **Authentication**: One-time "once" tokens
- **Timeout**: 10 seconds with context cancellation
- **Testing**: Predefined responses for test environment

**Flow:**
1. User redirected to AuthKnight with once parameter
2. AuthKnight validates authentication
3. Controller receives email and user data
4. Session created and user redirected

### 3. Privacy-First Data Architecture

#### Vault Store Integration

**Purpose**: Encrypt sensitive user data (email addresses)

**Implementation:**
```go
// Email encryption process
emailToken, err := vaultStore.TokenCreate(ctx, email, vaultKey, 20)
user.SetEmail(emailToken) // Store encrypted token, not plaintext
```

**Benefits:**
- Email addresses never stored in plaintext
- Encryption keys managed separately
- Configurable encryption strength

#### Blind Index System

**Purpose**: Enable email-based lookups without storing searchable plaintext

**Implementation:**
```go
// Blind index creation
searchValue := blindindexstore.NewSearchValue().
    SetSourceReferenceID(userID).
    SetSearchValue(email) // Original email for indexing
err := blindIndexStore.SearchValueCreate(ctx, searchValue)
```

**Benefits:**
- Prevents email enumeration attacks
- Enables efficient user lookup
- Maintains privacy while supporting functionality

### 4. Session Management

**Session Creation Process:**
1. Generate unique session key
2. Set session metadata (user ID, IP, User-Agent)
3. Configure expiration (2 hours production, 4 hours development)
4. Store in session store
5. Set secure authentication cookie

**Security Features:**
- Secure cookie configuration
- IP address tracking
- User-Agent validation
- Automatic expiration

### 5. Rate Limiting Strategy

**Global Rate Limits:**
- 20 requests per second
- 180 requests per minute  
- 12,000 requests per hour

**Authentication-Specific Limits:**
- Auth/Login: 5 requests per minute per IP
- Registration: 3 requests per minute per IP
- Logout: No additional limits (relies on global limits)

## Security Controls

### 1. Input Validation
- Once parameter validation and sanitization
- Email format verification
- Request context propagation

### 2. Error Handling
- Structured error messages
- No information leakage in responses
- Comprehensive logging with appropriate levels

### 3. Data Protection
- Email encryption via vault store
- Blind indexing for privacy
- Secure session management
- HTTPS-only external communications

### 4. Attack Prevention
- Rate limiting against brute force
- CSRF protection via secure cookies
- XSS prevention via HttpOnly cookies
- Session hijacking protection

## Configuration Options

### Environment Variables
```bash
# Authentication
AUTH_REGISTRATION_ENABLED="yes"

# Vault Store (for email encryption)
VAULT_STORE_KEY="your-long-vault-key"
USER_STORE_VAULT_ENABLED="yes"

# Session Management
SESSION_SECRET="your-secure-random-string"
```

### Feature Flags
- Registration enable/disable
- Vault store encryption toggle
- Development vs production configurations

## Database Schema Considerations

### Users Table
```sql
-- When vault enabled: stores encrypted email token
-- When vault disabled: stores plaintext email
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255), -- Encrypted token or plaintext
    status VARCHAR(50),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Blind Index Table
```sql
-- Enables email lookups without storing searchable plaintext
CREATE TABLE blind_index_email (
    id VARCHAR(255) PRIMARY KEY,
    search_value VARCHAR(255), -- Hashed email
    source_reference_id VARCHAR(255), -- User ID
    created_at TIMESTAMP
);
```

## Testing Strategy

### Unit Tests
- Controller logic testing
- Error scenario coverage
- Context cancellation testing
- Mock external service responses

### Integration Tests
- Complete authentication flow
- Database transaction testing
- Session management validation
- Rate limiting verification

### Test Environment
- Predefined AuthKnight responses
- Isolated test database
- Mock vault and blind index stores

## Performance Considerations

### Optimizations
- HTTP client reuse (future enhancement)
- Blind index result caching (future enhancement)
- Database connection pooling
- Efficient session storage

### Monitoring
- Authentication success/failure rates
- Session creation metrics
- External service response times
- Rate limiting effectiveness

## Security Best Practices Implemented

1. **Defense in Depth**: Multiple layers of security controls
2. **Privacy by Design**: Email encryption and blind indexing
3. **Secure Defaults**: Secure cookie configurations
4. **Fail Safe**: Proper error handling and logging
5. **Least Privilege**: Minimal data exposure in responses

## Future Enhancements

### Planned Improvements
- HTTP client pooling for performance
- Blind index caching for frequently accessed emails
- Enhanced monitoring and alerting
- Additional authentication providers

### Security Roadmap
- Multi-factor authentication support
- Advanced session management
- Enhanced audit logging
- Automated security testing

## Troubleshooting

### Common Issues
1. **Vault Store Not Available**: Check VAULT_STORE_KEY configuration
2. **Session Creation Fails**: Verify session store initialization
3. **Rate Limiting Issues**: Review global and specific rate limits
4. **External Service Timeouts**: Check AuthKnight service availability

### Debug Information
- Structured logging with correlation IDs
- Error context preservation
- Performance metrics collection
- Security event tracking

---

*This documentation should be updated whenever architectural changes are made to the authentication system.*
