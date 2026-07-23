# Artinux

A Linux artifact storage service that stores, retrieves, and manages artifacts on S3-compatible object storage (MinIO). Designed for CI/CD pipelines that produce binaries, packages, archives, and log files.

## Architecture

The project follows a hexagonal (ports and adapters) architecture with four layers:

- **`internal/domain/`** - Core business logic: entities, value objects, error types, interfaces (ports), and domain services (file validation).
- **`internal/application/usecases/`** - Application use cases that orchestrate domain logic. Each use case depends on interfaces, not concrete implementations.
- **`internal/infrastructure/`** - Concrete implementations of domain ports: MinIO storage client, SHA256 hasher, Viper-based config loader, Zap logger.
- **`internal/adapters/http/`** - HTTP handlers for the REST API, built with Echo. These adapters convert HTTP requests into use case inputs and use case outputs into HTTP responses.

Dependencies flow inward. The domain layer has no external dependencies. Infrastructure and adapters depend on domain interfaces through dependency injection.

## Project Structure

```
cmd/main.go                       Application entry point
internal/
  domain/
    entities/artifact.go          Artifact entity with naming and metadata logic
    ports/                        Interfaces for storage, logging, checksum
    services/file_validator.go    File validation rules
    errors.go                     Domain error types
  application/usecases/           Use cases for upload, download, list, delete
  infrastructure/
    config/                       Configuration loading and validation
    storage/minio.go              MinIO client wrapper
    checksum_hasher/sha256.go     SHA256 checksum implementation
    logger/zap_logger.go          Structured logger using zap
  adapters/http/                  Echo HTTP handlers and route setup
pkg/
  server/                         HTTP server with middleware
  response/                       Standardized API response helpers
  helper/                         OS/arch validation, unit conversion
```

## Prerequisites

- Go 1.23.7
- Docker and Docker Compose
- MinIO-compatible object storage (provided via Docker Compose)

## Quick Start

### The Simplest Way

You can copy/paste (& edit them to your own values) everything in [examples](https://github.com/dinno7/artinux/tree/main/examples) directory and then just run:

```bash
docker compose up

```

The service listens on `http://localhost:7000`. The Swagger UI is available at `http://localhost:7000/docs`.

## Configuration

Configuration is read from `config.yml` in the working directory.

### Config Reference

```yaml
env: dev # Runtime environment: dev or prod

logging:
  level: debug # debug, info, warn, error, fatal
  format: text # text or json

http_server:
  schema: http
  host: 0.0.0.0
  port: 7000

object_storage:
  endpoint: "minio:9000" # S3-compatible endpoint
  username: dinno
  password: strongpassword
  bucket_name: artinux
  region: us-east-1
  health_check_interval: 5s
  max_retries: 10

upload:
  max_size_mb: 100 # Maximum file size in MB
  allowed_file_exts: # Allowed file extensions (without dot)
    - deb
    - yml
    - md
```

MinIO credentials come from environment variables, you can set them in `.env` file or inside the `compoe.yaml` file:

- `MINIO_ROOT_USER`
- `MINIO_ROOT_PASSWORD`

## Build & Run

### Standalone(MinIO must be running)

```bash
go build -o ./artinux ./cmd/main.go

./artinux
```

### With Docker Compose

```bash
docker compose up --build
```

## API

| Method | Path                           | Description                                        |
| ------ | ------------------------------ | -------------------------------------------------- |
| GET    | `/api/v1/health`               | Health check                                       |
| POST   | `/api/v1/artifacts`            | Upload one or more files (multipart)               |
| GET    | `/api/v1/artifacts`            | List artifacts (supports `?prefix=` and `?limit=`) |
| GET    | `/api/v1/artifacts/download/*` | Download artifact by object key                    |
| DELETE | `/api/v1/artifacts/*`          | Delete single artifact by object key               |
| DELETE | `/api/v1/artifacts`            | Batch delete (JSON body with `object_keys`)        |

### Upload Parameters

All form parameters are required:

- `arch` - Target architecture (amd64, arm64, etc.)
- `os` - Target operating system (linux, darwin, windows, etc.)
- `username` - Uploader username
- `hostname` - Source hostname
- `artifacts` - File field, accepts multiple files

## Object Naming Strategy

Artifacts are stored using this key pattern:

```
{os}/{arch}/{YYYY}/{M}/{D}/{uuid}_{original_filename}
```

For example: `linux/amd64/2026/7/23/a1b2c3d4-package.deb`

This structure prevents name collisions through UUID prefixes, groups artifacts by platform and date for easy discovery, and scales well across millions of objects. The date-based hierarchy makes prefix-based listing efficient.

## Testing

```bash
# Run all unit tests
go test -v ./...

```

Unit tests cover validation logic, configuration loading and validation, checksum generation, object key generation, artifact metadata mapping, error handling paths, and all use case orchestration. Tests use table-driven patterns and mocked dependencies (gomock) for the storage and checksum interfaces.

## Assumptions

- The service runs in a trusted network. No authentication or authorization is implemented.
- Uploaded metadata (os, arch, username, hostname) is provided by the client and treated as informational.
- Files are read into memory for checksum computation before upload. Very large files may require streaming changes.
- MinIO is assumed to be the backend. Other S3-compatible stores may work but are not tested.
- AI-assisted code generation was used for a subset of unit tests (the use case and integration test scaffolding) and documenting. All architecture, domain logic, infrastructure code, and design decisions were implemented without AI assistance.

## Limitations

- No authentication or authorization layer.
- No pagination for the list endpoint beyond a simple limit parameter.
- No background cleanup for artifacts with integrity check failures (deletion is attempted synchronously).
- Compound file extensions (`.tar.gz`) are detected by `filepath.Ext`, which only returns the last component (`.gz`). The validator handles compound extensions separately as a workaround.

## Future Improvements

- Presigned URL support for direct upload/download.
- Concurrent processing with progress reporting.
- Retry logic with exponential backoff for transient storage failures.
- Prometheus metrics and structured health checks.
- Object versioning and lifecycle policies.
- Immutable artifact releases with checksum verification on download.
