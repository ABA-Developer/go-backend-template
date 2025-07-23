---
title: Clean Architecture
keywords: [clean, architecture, fiber, postgreSQL, go]
description: Implementing clean architecture in Go.
---

# Clean Architecture Example

[![Github](https://img.shields.io/static/v1?label=&message=Github&color=2ea44f&style=for-the-badge&logo=github)](https://github.com/gofiber/recipes/tree/master/clean-architecture) [![StackBlitz](https://img.shields.io/static/v1?label=&message=StackBlitz&color=2ea44f&style=for-the-badge&logo=StackBlitz)](https://stackblitz.com/github/gofiber/recipes/tree/master/clean-architecture)

This example demonstrates a Go Fiber application following the principles of Clean Architecture.

## Description

This project provides a starting point for building a web application with a clean architecture. It leverages Fiber for the web framework, PostgreSQL for the database, and follows the Clean Architecture principles to separate concerns and improve maintainability.

## Requirements

- [Go](https://golang.org/dl/) 1.18 or higher
- [PostgreSQL](https://www.postgresql.org/)
- [Git](https://git-scm.com/downloads)

## Project Structure

- `api/`: Contains the HTTP handlers, routes, and presenters.
- `build/`: Contains the built application.
- `cmd/`: Contains the main application entry point.
- `db/`: Contains the database migration.
- `internal/`: Contains internal application package.
- `pkg/`: Contains the core business logic and entities.

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/ABA-Developer/BE-dashboard-nba.git
   cd BE-dashboard-nba
   ```

2. Set the environment variables in a `.env` file:

   ```env
    DB_MIGRATOR_DRIVER="postgres"
    DB_USERNAME="auth_user"
    DB_PASSWORD="auth_user"
    DB_NAME="auth_user"
    DB_HOST="localhost"
    DB_PORT="5432"
    DB_SSLMODE="disable"
    DB_MAX_OPEN_CONNS=10
    DB_MAX_IDLE_CONNS=10
    DB_MAX_IDLE_TIME=5m
    DB_MAX_LIFETIME=5m
    DB_MAX_CONN_WAIT_TIME=5m
    DB_MAX_CONN_LIFETIME=5m
    DB_MAX_CONN_IDLE_TIME=5m
    DB_ADDR="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
   ```

3. Install the dependencies:

   ```bash
   go mod download
   ```

4. Run the application:
   ```bash
   make run
   ```

The API should now be running on `http://localhost:8080`.

## API Endpoints

The following endpoints are available in the API:

- **GET /user**: List all users.
- **POST /user**: Add a new user.
- **GET /user/:id**: Get a user by id.
- **PUT /user/:id**: Update an existing user by id.
- **DELETE /user/:id**: Remove a book by id.

## Clean Architecture Principles

Clean Architecture is a software design philosophy that emphasizes the separation of concerns, making the codebase more maintainable, testable, and scalable. In this example, the Go Fiber application follows Clean Architecture principles by organizing the code into distinct layers, each with its own responsibility.

### Layers in Clean Architecture

1. **Entities (Core Business Logic)**

- Located in the `pkg/entities` directory.
- Contains the core business logic and domain models, which are independent of any external frameworks or technologies.

2. **Use Cases (Application Logic)**

- Located in the `pkg/user` directory.
- Contains the application-specific business rules and use cases. This layer orchestrates the flow of data to and from the entities.

3. **Interface Adapters (Adapters and Presenters)**

- Located in the `api` directory.
- Contains the HTTP handlers, routes, and presenters. This layer is responsible for converting data from the use cases into a format suitable for the web framework (Fiber in this case).

4. **Frameworks and Drivers (External Interfaces)**

- Located in the `cmd` directory.
- Contains the main application entry point and any external dependencies like the web server setup.

### Example Breakdown

- **Entities**: The `entities.User` struct represents the core business model for a user.
- **Use Cases**: The `user.Service` interface defines the methods for interacting with user, such as `CreateUser`, `UpdateUser`, `DeleteUser`, and `GetUserById`.
- **Interface Adapters**: The `handlers` package contains the HTTP handlers that interact with the `user.Service` to process HTTP requests and responses.
- **Frameworks and Drivers**: The `cmd/main.go` file initializes the Fiber application and sets up the routes using the `routes.UserRouter` function.

By following Clean Architecture principles, this example ensures that each layer is independent and can be modified or replaced without affecting the other layers, leading to a more maintainable and scalable application.

## Conclusion

This example provides a basic setup for a Go Fiber application following Clean Architecture principles. It can be extended and customized further to fit the needs of more complex applications.

## References

- [Fiber Documentation](https://docs.gofiber.io)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html)
