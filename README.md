# Go Backend Template

Backend template in Go with the NATA HR transaction service layer structure and Fiber as the HTTP framework.

## Structure

```text
cmd/
constant/
database/
  migrations/
  seeders/
docs/
internal/
  application/
    auth/
      dto/
      mapper/
      usecase/
    menu/
      menu-permission/
    role/
    user/
    shared/dto/
  domain/
    contract/
      repository/
        params/
      usecase/
    model/
  infrastructure/
    api/
      routes/
    database/
    logger/
    repository/
    runtime/
      container/
      env/
    validator/
  presentation/
    auth/
      handler/
      mapper/
      presenter/
    menu/
    role/
    user/
    middleware/
    request/
    response/
  test/
    mock/
      repository/
scripts/
```

Handlers convert presenter requests to application DTOs through presentation mappers.
Use cases convert application DTOs to repository parameters through application mappers.
Repository contracts use types from `internal/domain/contract/repository/params`.

## Environment

Copy `.env.example` to `.env` and adjust the values for your environment.

Important values:

- `APP_NAME`
- `APP_HOST`
- `APP_PORT`
- `DB_DRIVER`
- `DB_USERNAME`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_HOST`
- `DB_PORT`
- `SWAGGER_ENABLED`
- `SWAGGER_SERVICES`

## Run

```bash
go mod download
make run
```

## Database seed

```bash
make migrate.seed
```

## Swagger

Swagger UI is enabled with `SWAGGER_ENABLED=true`.
Use `SWAGGER_SERVICES` to filter the services shown in the docs.

## Test

```bash
go test ./...
```
