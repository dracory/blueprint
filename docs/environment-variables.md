# Environment Variables Reference

Complete reference of all environment variables used by Blueprint.

## Quick Start

1. Copy `.env.example` to `.env`
2. Set required variables for your environment
3. Run `task env-dev` or `task env-prod` to configure encrypted variables

## Variable Categories

### Application

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| APP_ENV | Yes | - | Environment (local, development, staging, production, testing) |
| APP_NAME | No | Blueprint | Application name |
| APP_URL | Yes | - | Base URL for link generation |
| APP_HOST | Yes | - | Server host address |
| APP_PORT | Yes | - | Server port |
| APP_DEBUG | No | false | Enable debug mode |

### Database

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| DB_DRIVER | No | sqlite | Database type (sqlite, mysql, postgres) |
| DB_HOST | Conditional* | - | Database server host |
| DB_PORT | Conditional* | - | Database server port |
| DB_DATABASE | Yes | - | Database name or SQLite file path |
| DB_USERNAME | Conditional* | - | Database username |
| DB_PASSWORD | Conditional* | - | Database password |
| DB_SSL_MODE | No | disable | SSL mode (postgres only) |
| DB_CHARSET | No | utf8mb4 | Character set (mysql only) |
| DB_TIMEZONE | No | UTC | Database timezone |
| DB_MAX_OPEN_CONNS | No | varies | Max open connections |
| DB_MAX_IDLE_CONNS | No | varies | Max idle connections |
| DB_CONN_MAX_LIFETIME_SECONDS | No | varies | Connection lifetime |
| DB_CONN_MAX_IDLE_TIME_SECONDS | No | varies | Connection idle time |

*Required when DB_DRIVER is mysql or postgres

### Email

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| MAIL_DRIVER | No | smtp | Email driver |
| MAIL_HOST | Conditional* | - | SMTP host |
| MAIL_PORT | Conditional* | - | SMTP port |
| MAIL_USERNAME | Conditional* | - | SMTP username |
| MAIL_PASSWORD | Conditional* | - | SMTP password |
| MAIL_ENCRYPTION | No | null | SMTP encryption (null, tls, ssl) |
| MAIL_FROM_ADDRESS | Yes | - | Default from address |
| MAIL_FROM_NAME | Yes | - | Default from name |

*Required when MAIL_DRIVER=smtp

### Authentication

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| AUTH_REGISTRATION_ENABLED | No | yes | Allow user registration |
| AUTH_EMAILS_ALLOWED_ACCESS | No | - | Allowed email domains |

### LLM Providers

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| ANTHROPIC_API_KEY | No | - | Anthropic API key |
| ANTHROPIC_API_USED | No | no | Enable Anthropic API |
| ANTHROPIC_API_DEFAULT_MODEL | No | - | Default Anthropic model |
| GEMINI_API_KEY | No | - | Google Gemini API key |
| GEMINI_API_USED | No | no | Enable Gemini API |
| GEMINI_API_DEFAULT_MODEL | No | - | Default Gemini model |
| OPENAI_API_KEY | No | - | OpenAI API key |
| OPENAI_API_USED | No | no | Enable OpenAI API |
| OPENAI_API_DEFAULT_MODEL | No | - | Default OpenAI model |
| OPENROUTER_API_KEY | No | - | OpenRouter API key |
| OPENROUTER_API_USED | No | no | Enable OpenRouter API |
| OPENROUTER_API_DEFAULT_MODEL | No | - | Default OpenRouter model |
| VERTEX_AI_API_PROJECT_ID | No | - | Vertex AI project ID |
| VERTEX_AI_API_REGION_ID | No | - | Vertex AI region |
| VERTEX_AI_API_MODEL_ID | No | - | Vertex AI model ID |
| VERTEX_AI_API_USED | No | no | Enable Vertex AI API |

### Media Storage

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| MEDIA_DRIVER | No | local | Storage driver (local, s3, gcs, sql) |
| MEDIA_BUCKET | Conditional* | - | Storage bucket name |
| MEDIA_KEY | Conditional* | - | Storage access key |
| MEDIA_SECRET | Conditional* | - | Storage secret key |
| MEDIA_REGION | Conditional* | - | Storage region |
| MEDIA_ENDPOINT | No | - | Custom endpoint URL |
| MEDIA_URL | Yes | - | Public URL base |

*Required when MEDIA_DRIVER is s3 or gcs

### Payment

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| STRIPE_KEY_PRIVATE | Conditional* | - | Stripe private key |
| STRIPE_KEY_PUBLIC | Conditional* | - | Stripe public key |

*Required when using Stripe payments

### Security

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| ENVENC_USED | No | no | Enable environment encryption |
| ENVENC_KEY_PRIVATE | Conditional* | - | Encryption private key |
| SESSION_SECRET | Yes | - | Session secret key |

*Required when ENVENC_USED=yes

### CMS

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| MCP_API_KEY | No | - | CMS MCP API key |
| CMS_STORE_TEMPLATE_ID | Conditional* | - | CMS store template ID |
| VAULT_STORE_KEY | Conditional* | - | Vault store encryption key |

*Required when respective stores are enabled

## Environment-Specific Configuration

### Local Development

```bash
APP_ENV=local
APP_HOST=127.0.0.1
APP_PORT=32322
APP_URL=http://localhost:32322
APP_DEBUG=true
DB_DRIVER=sqlite
DB_DATABASE=./database.db
MAIL_DRIVER=smtp
MAIL_HOST=localhost
MAIL_PORT=1025
```

### Production

```bash
APP_ENV=production
APP_HOST=0.0.0.0
APP_PORT=8080
APP_URL=https://example.com
APP_DEBUG=false
DB_DRIVER=postgres
DB_HOST=your-db-host
DB_PORT=5432
DB_DATABASE=production_db
DB_USERNAME=production_user
DB_PASSWORD=secure_password
```

## Security Best Practices

1. **Never commit `.env` files** - Add to `.gitignore`
2. **Use encryption for secrets** - Enable ENVENC_USED for sensitive values
3. **Different values per environment** - Never use production values in development
4. **Rotate keys regularly** - Update API keys and passwords periodically
5. **Limit access** - Restrict who can access environment configuration
6. **Audit changes** - Track changes to environment variables
7. **Use strong session secrets** - Generate cryptographically secure random strings

## Troubleshooting

### Application won't start

1. Check all required variables are set
2. Verify variable values are correct for your environment
3. Check for typos in variable names
4. Review application logs for specific errors

### Database connection fails

1. Verify DB_DRIVER matches your setup
2. Check database host and port are accessible
3. Verify database credentials are correct
4. Ensure database exists and user has permissions
5. Check SSL mode settings for PostgreSQL

### Email not sending

1. Verify MAIL_DRIVER is set correctly
2. Check SMTP credentials are valid
3. Verify SMTP host and port are correct
4. Check firewall allows SMTP traffic
5. Verify MAIL_ENCRYPTION setting matches server requirements

### LLM API errors

1. Verify API key is correct and active
2. Check API provider's service status
3. Ensure the API_USED flag is set to "yes"
4. Verify model ID is valid for the provider
5. Check rate limits and quotas

## See Also

- [Package Documentation](docs/proposals/package-documentation.md)
- [Architecture Documentation](docs/proposals/architecture-documentation.md)
- [Environment Variable Documentation Proposal](docs/proposals/environment-variable-documentation.md)
