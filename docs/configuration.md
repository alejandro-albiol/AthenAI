# Configuration Guide

AthenAI uses environment variables for configuration management. This guide explains all available configuration options.

## Quick Start

1. **Development**: Copy `.env.development` to `.env`
2. **Production**: Copy `.env.production` to `.env` and update with real values

## Configuration Categories

### Database Configuration

| Variable  | Description                | Default    | Required |
| --------- | -------------------------- | ---------- | -------- |
| `DB_TYPE` | Database type (postgres)   | `postgres` | ✅       |
| `DB_DSN`  | Database connection string | -          | ✅       |

**Example DSN**: `host=localhost port=5432 user=athenai password=password dbname=athenai sslmode=disable`

### Application Configuration

| Variable      | Description                     | Default     | Required |
| ------------- | ------------------------------- | ----------- | -------- |
| `APP_ENV`     | Environment mode (`dev`/`prod`) | `prod`      | ✅       |
| `SERVER_PORT` | Server port                     | `8080`      | ❌       |
| `SERVER_HOST` | Server host                     | `localhost` | ❌       |

**Security Note**:

- `APP_ENV=dev`: Returns detailed error information including SQL errors
- `APP_ENV=prod`: Returns only error codes for security

### JWT Configuration

| Variable                   | Description                       | Default | Required |
| -------------------------- | --------------------------------- | ------- | -------- |
| `JWT_SECRET`               | JWT signing secret (min 256 bits) | -       | ✅       |
| `JWT_ACCESS_TOKEN_EXPIRY`  | Access token lifetime             | `15m`   | ❌       |
| `JWT_REFRESH_TOKEN_EXPIRY` | Refresh token lifetime            | `7d`    | ❌       |

### CORS Configuration

| Variable               | Description                     | Default                                       |
| ---------------------- | ------------------------------- | --------------------------------------------- |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins | `http://localhost:3000,http://localhost:8080` |
| `CORS_ALLOWED_METHODS` | Comma-separated allowed methods | `GET,POST,PUT,DELETE,OPTIONS`                 |
| `CORS_ALLOWED_HEADERS` | Comma-separated allowed headers | `Content-Type,Authorization,X-Gym-ID`         |

### Domain Configuration

| Variable          | Description                       | Default       |
| ----------------- | --------------------------------- | ------------- |
| `BASE_DOMAIN`     | Base domain for tenant subdomains | `athenai.com` |
| `PLATFORM_DOMAIN` | Platform domain                   | `athenai.com` |

**Multi-tenancy**:

- Platform: `athenai.com` or `localhost`
- Tenants: `{gym_uuid}.athenai.com`

### Email Configuration

| Variable          | Description      | Default               | Required |
| ----------------- | ---------------- | --------------------- | -------- |
| `SMTP_HOST`       | SMTP server host | -                     | ✅       |
| `SMTP_PORT`       | SMTP server port | `587`                 | ❌       |
| `SMTP_USERNAME`   | SMTP username    | -                     | ✅       |
| `SMTP_PASSWORD`   | SMTP password    | -                     | ✅       |
| `SMTP_FROM_NAME`  | From name        | `AthenAI`             | ❌       |
| `SMTP_FROM_EMAIL` | From email       | `noreply@athenai.com` | ❌       |

### File Upload Configuration

| Variable               | Description                | Default                                          |
| ---------------------- | -------------------------- | ------------------------------------------------ |
| `UPLOAD_MAX_SIZE`      | Maximum file size          | `10MB`                                           |
| `UPLOAD_ALLOWED_TYPES` | Comma-separated MIME types | `image/jpeg,image/png,image/gif,application/pdf` |
| `UPLOAD_STORAGE_PATH`  | Upload directory           | `./uploads`                                      |

### Security Configuration

| Variable             | Description                | Default |
| -------------------- | -------------------------- | ------- |
| `BCRYPT_COST`        | BCrypt hashing cost (4-31) | `12`    |
| `SESSION_TIMEOUT`    | Session timeout duration   | `24h`   |
| `MAX_LOGIN_ATTEMPTS` | Max failed login attempts  | `5`     |
| `LOCKOUT_DURATION`   | Account lockout duration   | `15m`   |

### Rate Limiting

| Variable                         | Description                | Default |
| -------------------------------- | -------------------------- | ------- |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Requests per minute per IP | `60`    |
| `RATE_LIMIT_BURST`               | Burst capacity             | `100`   |

### Logging Configuration

| Variable        | Description   | Default          | Options                          |
| --------------- | ------------- | ---------------- | -------------------------------- |
| `LOG_LEVEL`     | Logging level | `info`           | `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT`    | Log format    | `json`           | `json`, `text`                   |
| `LOG_FILE_PATH` | Log file path | `./logs/app.log` | -                                |

### Redis Configuration (Optional)

| Variable         | Description           | Default     |
| ---------------- | --------------------- | ----------- |
| `REDIS_HOST`     | Redis host            | `localhost` |
| `REDIS_PORT`     | Redis port            | `6379`      |
| `REDIS_PASSWORD` | Redis password        | -           |
| `REDIS_DB`       | Redis database number | `0`         |

### Payment Configuration

| Variable                | Description            | Required |
| ----------------------- | ---------------------- | -------- |
| `STRIPE_PUBLIC_KEY`     | Stripe publishable key | ✅       |
| `STRIPE_SECRET_KEY`     | Stripe secret key      | ✅       |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook secret  | ✅       |

### External APIs

| Variable              | Description         | Purpose                     |
| --------------------- | ------------------- | --------------------------- |
| `GOOGLE_MAPS_API_KEY` | Google Maps API key | Gym location features       |
| `WEATHER_API_KEY`     | Weather API key     | Weather-based notifications |

### Feature Flags

| Variable                   | Description                  | Default                        |
| -------------------------- | ---------------------------- | ------------------------------ |
| `ENABLE_SWAGGER`           | Enable Swagger UI            | `true` in dev, `false` in prod |
| `ENABLE_METRICS`           | Enable metrics endpoint      | `true`                         |
| `ENABLE_HEALTH_CHECK`      | Enable health check endpoint | `true`                         |
| `ENABLE_USER_REGISTRATION` | Allow new user registration  | `true`                         |

## Environment-Specific Recommendations

### Development

- Use `APP_ENV=dev` for detailed error messages
- Enable Swagger with `ENABLE_SWAGGER=true`
- Use test payment keys
- Lower security settings for convenience

### Production

- Use `APP_ENV=prod` for security
- Disable Swagger with `ENABLE_SWAGGER=false`
- Use live payment keys
- Higher security settings
- Enable monitoring and logging

## Security Best Practices

1. **Never commit real .env files**: Add `.env` to `.gitignore`
2. **Use strong JWT secrets**: Minimum 256 bits, cryptographically random
3. **Use environment-specific configurations**: Different keys for dev/prod
4. **Rotate secrets regularly**: Especially JWT secrets and API keys
5. **Use TLS in production**: Set `sslmode=require` for database
6. **Monitor configuration**: Use proper logging and monitoring

## Usage in Code

```go
package main

import "github.com/alejandro-albiol/athenai/config"

func main() {
    cfg := config.Load()

    // Check environment
    if cfg.IsDevelopment() {
        // Development-specific logic
    }

    // Use configuration
    fmt.Printf("Server running on %s:%s", cfg.ServerHost, cfg.ServerPort)
}
```
