# Authentication Architecture Documentation

## Overview

The Blueprint authentication system implements a sophisticated privacy-first architecture that prioritizes user data protection while maintaining robust security controls. This document outlines the complete authentication flow, security measures, and architectural decisions.

## Login Methods

Blueprint ships with two login mechanisms. The active method is selected by the
`config.LOGIN_METHOD` constant in `internal/config/auth_config.go` — a
compile-time, one-time developer decision (not an environment variable).

- **`LOGIN_METHOD_OTP` (default)** — in-house email one-time-password login,
  fully self-contained, no external services. Two routes share the same path:
  - `GET auth/login` (HTML handler) renders a self-contained page
    (`app.html`/`app.js`/`app.css` embedded, Vue + Notiflix via CDN).
  - `POST auth/login` (JSON handler) dispatches on the `action` query param:
    - `otp-send-ajax` issues a 6-digit `crypto/rand` code plus a 128-bit
      nonce, stores `email:otp` in the memory cache under the nonce (15-min
      TTL, max 5 verify attempts, max 3 sends per email per window), and
      enqueues `EmailOTPTask`. The queued task receives **only the nonce** —
      the plaintext code is resolved from the memory cache at execution time
      and is never persisted in the task store. The send quota is consumed
      only after a successful enqueue, so transient failures don't lock the
      user out. Counter read-modify-writes are serialized with a mutex.
    - `otp-verify-ajax` checks the attempt counter before lookup
      (anti-enumeration), compares with `subtle.ConstantTimeCompare`, then
      calls the shared pipeline. The `return` param is accepted only if it is
      a relative path (open-redirect protection).
- **Registration completion** (`auth/register`) follows the same pattern:
  `GET` renders a Vue page, `POST` (JSON) exposes `action=save` (profile
  update, vault-aware) and `action=timezones` (country → timezone list).
- **`LOGIN_METHOD_AUTHKNIGHT`** — delegates to the external AuthKnight
  service: `auth/login` redirects to `authknight.com`, `auth/auth` exchanges
  the `once` token for the user email. Only in this mode is the `AUTH_AUTH`
  callback route registered. The echoed `backUrl` from the AuthKnight
  response is re-validated (must start with the app home URL) before being
  used as the post-login redirect.

An unrecognized `LOGIN_METHOD` value panics at startup in `routes.go`
rather than silently defaulting.

Both methods converge on the shared post-auth pipeline in
`internal/controllers/auth/shared` (`SessionLogin`): find-or-create user
(vault + blind index aware), session creation, auth cookie, redirect.

### Removing a Method

Because `LOGIN_METHOD` is a constant, deleting the unused method's package
produces a compile error pointing at the branch to remove in
`internal/controllers/auth/routes.go`.

- Remove OTP: set `LOGIN_METHOD = LOGIN_METHOD_AUTHKNIGHT`, then delete
  `internal/controllers/auth/login_otp/`, `internal/tasks/email_otp/`,
  `internal/emails/user_email_otp.go`, and the `EmailOTPTask`
  registration/alias.
- Remove AuthKnight: set `LOGIN_METHOD = LOGIN_METHOD_OTP`, then delete
  `internal/controllers/auth/login_authknight/`, `internal/controllers/auth/authentication_authknight/`,
  `links.authLinks.AuthKnightLogin`, and the `AUTH_AUTH` route/constant.

## Architecture Components

### 1. Authentication Controller (`internal/controllers/auth/`)

**Primary Responsibilities:**
- Route the selected login method (`config.LOGIN_METHOD`): in-house email
  OTP (`login_otp/`) or external AuthKnight (`login_authknight/` +
  `authentication_authknight/` callback)
- Manage user creation and session establishment via the shared pipeline
  (`auth/shared.SessionLogin`)
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
5. **OTP email not arriving**: The task queue runner must be running (it is
   started by `startBackgroundProcesses` when the task store is enabled).
   Check the `snv_tasks_task_queue` table — the task `details` column records
   failures (e.g. `smtp: server doesn't support AUTH` when `MAIL_USERNAME`
   is set against a server without AUTH, such as Mailpit — leave
   `MAIL_USERNAME`/`MAIL_PASSWORD` empty for local Mailpit).

### Debug Information
- Structured logging with correlation IDs
- Error context preservation
- Performance metrics collection
- Security event tracking

---

*This documentation should be updated whenever architectural changes are made to the authentication system.*
