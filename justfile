
# Get help message
help:
  just --list


# Generate Swagger documentation
@gen-docs:
  swag init -d "pkg/server,pkg/response,internal/domain/entities,internal/adapters/http" -g server.go -o "./docs"
  echo "Formatiing..."
  swag fmt -d "pkg/server,pkg/response,internal/domain/entities,internal/adapters/http" -g server.go

@gen-mocks:
  mockgen -package ports internal/domain/ports/storage.go ObjectStorage  >  internal/domain/ports/storage_mock.go
  mockgen -package ports internal/domain/ports/checksum_hasher.go ChecksumHasher >  internal/domain/ports/checksum_hasher_mock.go
