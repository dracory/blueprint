# Blueprint 

![tests](https://github.com/dracory/blueprint/workflows/tests/badge.svg)

## URLS

- https://YOURAPPURL
- https://YOURAPPURL.a.run.app (Dev)

## Description

**Build production-ready Go apps in minutes, not days and months.**

Blueprint is a rapid application development (RAD) starter template built on the Dracory framework. Think of it as "Laravel on steroids" for Go - a batteries-included foundation for building production-ready web applications and APIs with pre-configured infrastructure and best practices. The philosophy is simple: it's easier to remove what you don't need than to add missing infrastructure later.

**What's Included:**

**Core Infrastructure:**
- Database connection setup (SQLite, MySQL, PostgreSQL support)
- High-performance router with middleware chains
- Background task queue system
- Cron-like job scheduler
- Configuration management with validation
- Environment-specific configuration (dev, staging, prod)
- Encrypted environment variable support (Data Vault)
- Graceful shutdown handling
- Logging and structured logging

**Authentication & Security:**
- Complete authentication system
- Session management
- API authentication middleware
- Security headers middleware
- Jail bots middleware (IP-based bot protection)
- Email allowlist middleware
- HTTPS redirect middleware
- Blind Index for searchable encrypted data
- Vault to securely store secrets

**Admin & CMS:**
- Full admin interface with Vue.js
- CMS blocks system
- Blog management (blogadmin)
- File management (fileadmin)
- User management (useradmin)
- Shop/e-commerce (shopadmin)
- Log viewing (logadmin)
- Task management (task admin)
- CMS management (cmsadmin)
- Social media integration
- Testimonials system

**Development Tools:**
- CLI command dispatcher
- Deployment utilities
- Environment encryption tool (envenc)
- Load testing utility
- Task runner (taskfile.yml)
- Gitpod & GitHub CodeSpaces ready
- Docker support
- CI/CD pipelines

**Testing & Quality:**
- Comprehensive test utilities
- Integration test setup
- Coverage reporting
- Test fixtures and helpers

**Email & Notifications:**
- Email system with SMTP
- Email templates
- Admin email notifications
- User registration emails
- Contact form emails

**Additional Features:**
- Theme system
- Widget components
- Layout templates
- Helper functions
- Resource management
- Link management

## Installation

```bash
git clone https://github.com/dracory/blueprint
```

## Environment Variables

- Copy the `.env.example` file to `.env`

```bash
cp .env.example .env
```

- Set the dev vault values

```bash
task env-dev
```

- Set the prod vault values

```bash
task env-prod
```

For a complete reference of all available environment variables, see [Environment Variables Documentation](docs/environment-variables.md).


## Local Development

- Just starting
```bash
task dev:init
```

- Run in development mode
```bash
task dev
```

## Development on Gitpod

Use the link on the top of this README

## Testing

Running all tests

```bash
task test
```

-Running individual test

```
go test -run ^TestGuestFunnelTestSuite$
```

## Coverage Report

```bash
task cover
```

## CLI Commands

Deploy Live:

```bash
task deploy:live
```

Deploy Staging:

```bash
task deploy:staging
```

List Routes:

```bash
go run ./cmd/server routes list
```

Run task:

```bash
go run ./cmd/server task run ...
```

Run job:

```bash
go run ./cmd/server job run ...
```

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt).

For commercial use, please use the [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.
