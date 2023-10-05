# Audio Guide Bot: API service
API service is a dockerized Golang application that serves Guide Bot API and webhook for Telegram Bot API.

## Configuration
Required environment variables:
- `TELEGRAM_WEB_APP_URL` - URL to Guide UI service
- `TELEGRAM_BOT_TOKEN` - Telegram Bot token
- `TELEGRAM_PAYMENTS_TOKEN` - Telegram Payments token
- `JWT_SECRET` - Secret to sign and verify JWT tokens
- `DB_CONNECTION_STRING` - URL to DB
- `S3_CONNECTION_STRING` - URL to S3

Optional environment variables:
- `CORS_ALLOWED_ORIGINS` - list of allowed origins that may access the resource

## Service structure
Service is built on three abstractions:
- Providers - isolate interaction with other systems or packages
    - [auth](./auth/auth.go) - provides functionality for token-based authentication, implementations: [JWT](./auth/jwt.go)
    - [blob](./blob/blob.go) - provides I/O operations on immutable binary objects, implementations: [S3](./blob/s3.go)
    - [bot](./bot/bot.go) - provides interaction with Bot API, implementations: [Telegram API](./bot/bot.go)
    - [db](./db/db.go) - provides interaction with a database, implementations: [PostgreSQL](./db/postgres.go)
- Repositories - provide CRUD operations for data types, all interfaces are implemented as an aggregate [repository](./repository/repository.go) object
    - [repository/object](./repository/object.go) - implements CRUD operations for Object type
    - [repository/ticket](./repository/ticket.go) - implements CRUD operations for Ticket type
    - [repository/config](./repository/config.go) - implements CRUD operations for configuration variables
- Controllers - implement HTTP handlers with business logic, all handlers are implemented in compliance with [JSend](https://github.com/omniti-labs/jsend) specification
    - [controller/bot](./controller/bot.go) - implements logic to handle Telegram Bot API updates
    - [controller/objects](./controller/objects.go) - implements logic to interact with Object type
    - [controller/tickets](./controller/tickets.go) - implements logic to interact with Ticket type

All entities are constructed and injected in [main](main.go) and then HTTP handlers are served by [Fiber](https://github.com/gofiber/fiber).

## Database migrations
API service can be started in database migration mode. In this case, it will apply migrations from the implemented `DBProvider` and exit. To start the service in migration mode - specify `--migrate` execution argument.