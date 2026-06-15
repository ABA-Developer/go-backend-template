# Go Backend Template

Backend template in Go with a clean architecture structure inspired by the NATA HR service pattern.

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
    shared/dto/
    auth/
    menu/
    role/
    user/
    ...
  domain/
    contract/
      external/
      repository/
      usecase/
    model/
  infrastructure/
    api/
    database/
    logger/
    repository/
    runtime/
      container/
      env/
    validator/
  presentation/
    handler/
    middleware/
    presenter/
  test/
    mock/
      external/
      repository/
      usecase/
scripts/
```

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
